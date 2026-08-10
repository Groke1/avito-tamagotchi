package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func buildToken(t *testing.T, secret []byte, claims accessTokenClaims, alg, typ string) string {
	t.Helper()

	header := jwtHeader{Alg: alg, Typ: typ}
	headerBytes, err := json.Marshal(header)
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("marshal claims: %v", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedPayload := base64.RawURLEncoding.EncodeToString(claimsBytes)
	signingInput := encodedHeader + "." + encodedPayload
	signature := signJWT(signingInput, secret)

	return encodedHeader + "." + encodedPayload + "." + signature
}

func validClaims(now time.Time) accessTokenClaims {
	return accessTokenClaims{
		Sub: "user-1",
		Iat: now.Unix(),
		Exp: now.Add(time.Hour).Unix(),
	}
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		wantToken  string
		wantErr    error
	}{
		{name: "valid header", authHeader: "Bearer abc.def.ghi", wantToken: "abc.def.ghi", wantErr: nil},
		{name: "case-insensitive scheme", authHeader: "bearer abc.def.ghi", wantToken: "abc.def.ghi", wantErr: nil},
		{name: "empty header", authHeader: "", wantToken: "", wantErr: controller.ErrInvalidToken},
		{name: "missing token", authHeader: "Bearer", wantToken: "", wantErr: controller.ErrInvalidToken},
		{name: "wrong scheme", authHeader: "Basic abc.def.ghi", wantToken: "", wantErr: controller.ErrInvalidToken},
		{name: "too many parts", authHeader: "Bearer abc def", wantToken: "", wantErr: controller.ErrInvalidToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := bearerToken(tt.authHeader)
			if token != tt.wantToken {
				t.Errorf("expected token %q, got %q", tt.wantToken, token)
			}
			if tt.wantErr == nil && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if tt.wantErr != nil && err != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestParseAccessToken_Valid(t *testing.T) {
	secret := []byte("test-secret")
	now := time.Now().UTC()
	token := buildToken(t, secret, validClaims(now), "HS256", "JWT")

	claims, err := parseAccessToken(token, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claims.Sub != "user-1" {
		t.Fatalf("expected sub user-1, got %q", claims.Sub)
	}
}

func TestParseAccessToken_InvalidCases(t *testing.T) {
	secret := []byte("test-secret")
	now := time.Now().UTC()

	t.Run("empty secret", func(t *testing.T) {
		token := buildToken(t, secret, validClaims(now), "HS256", "JWT")
		_, err := parseAccessToken(token, []byte{})
		if err != controller.ErrInvalidToken {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("malformed token (wrong number of parts)", func(t *testing.T) {
		_, err := parseAccessToken("only.two", secret)
		if err != controller.ErrInvalidToken {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("wrong alg", func(t *testing.T) {
		token := buildToken(t, secret, validClaims(now), "HS512", "JWT")
		_, err := parseAccessToken(token, secret)
		if err != controller.ErrInvalidToken {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("wrong typ", func(t *testing.T) {
		token := buildToken(t, secret, validClaims(now), "HS256", "JWS")
		_, err := parseAccessToken(token, secret)
		if err != controller.ErrInvalidToken {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("wrong signature (signed with different secret)", func(t *testing.T) {
		token := buildToken(t, []byte("other-secret"), validClaims(now), "HS256", "JWT")
		_, err := parseAccessToken(token, secret)
		if err != controller.ErrInvalidToken {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		claims := accessTokenClaims{Sub: "user-1", Iat: now.Add(-2 * time.Hour).Unix(), Exp: now.Add(-time.Hour).Unix()}
		token := buildToken(t, secret, claims, "HS256", "JWT")
		_, err := parseAccessToken(token, secret)
		if err != controller.ErrExpiredToken {
			t.Fatalf("expected ErrExpiredToken, got %v", err)
		}
	})

	t.Run("iat too far in the future", func(t *testing.T) {
		claims := accessTokenClaims{Sub: "user-1", Iat: now.Add(time.Hour).Unix(), Exp: now.Add(2 * time.Hour).Unix()}
		token := buildToken(t, secret, claims, "HS256", "JWT")
		_, err := parseAccessToken(token, secret)
		if err != controller.ErrInvalidToken {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("missing sub", func(t *testing.T) {
		claims := accessTokenClaims{Sub: "", Iat: now.Unix(), Exp: now.Add(time.Hour).Unix()}
		token := buildToken(t, secret, claims, "HS256", "JWT")
		_, err := parseAccessToken(token, secret)
		if err != controller.ErrInvalidToken {
			t.Fatalf("expected ErrInvalidToken, got %v", err)
		}
	})
}

func TestAuthMiddleware_Integration(t *testing.T) {
	secret := []byte("test-secret")
	now := time.Now().UTC()

	router := gin.New()
	router.Use(AuthMiddleware(secret))
	router.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})

	t.Run("valid token passes through and sets userID", func(t *testing.T) {
		token := buildToken(t, secret, validClaims(now), "HS256", "JWT")
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d, body: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("missing header is rejected with 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})

	t.Run("expired token is rejected with 401", func(t *testing.T) {
		claims := accessTokenClaims{Sub: "user-1", Iat: now.Add(-2 * time.Hour).Unix(), Exp: now.Add(-time.Hour).Unix()}
		token := buildToken(t, secret, claims, "HS256", "JWT")
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", rec.Code)
		}
	})
}
