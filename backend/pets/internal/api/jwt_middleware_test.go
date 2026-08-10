package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken(t *testing.T, secret string, method jwt.SigningMethod, claims jwt.MapClaims) string {
	t.Helper()

	token := jwt.NewWithClaims(method, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return signed
}

func TestJwtMiddleware(t *testing.T) {
	secret := "test-secret"

	var handlerCalled bool
	var capturedUserID string

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		userID, _ := GetUserIDFromContext(r.Context())
		capturedUserID = userID
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		authHeader     string
		wantStatus     int
		wantNextCalled bool
	}{
		{
			name:           "no authorization header",
			authHeader:     "",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "missing bearer prefix",
			authHeader:     generateToken(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "user-1"}),
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "too many parts",
			authHeader:     "Bearer extra " + generateToken(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "user-1"}),
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "wrong scheme",
			authHeader:     "Basic " + generateToken(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "user-1"}),
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "invalid signature",
			authHeader:     "Bearer " + generateToken(t, "wrong-secret", jwt.SigningMethodHS256, jwt.MapClaims{"sub": "user-1"}),
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "unexpected signing method",
			authHeader:     "Bearer " + generateToken(t, secret, jwt.SigningMethodHS384, jwt.MapClaims{"sub": "user-1"}),
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "malformed token string",
			authHeader:     "Bearer not-a-jwt",
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "missing sub claim",
			authHeader:     "Bearer " + generateToken(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{}),
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "sub claim wrong type",
			authHeader:     "Bearer " + generateToken(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{"sub": 123}),
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "expired token",
			authHeader:     "Bearer " + generateToken(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "user-1", "exp": time.Now().Add(-time.Hour).Unix()}),
			wantStatus:     http.StatusUnauthorized,
			wantNextCalled: false,
		},
		{
			name:           "valid token",
			authHeader:     "Bearer " + generateToken(t, secret, jwt.SigningMethodHS256, jwt.MapClaims{"sub": "user-1"}),
			wantStatus:     http.StatusOK,
			wantNextCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerCalled = false
			capturedUserID = ""

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()

			JwtMiddleware(secret)(next).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if handlerCalled != tt.wantNextCalled {
				t.Errorf("next called = %v, want %v", handlerCalled, tt.wantNextCalled)
			}
			if tt.wantNextCalled && capturedUserID != "user-1" {
				t.Errorf("userID in context = %q, want %q", capturedUserID, "user-1")
			}
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	t.Run("value present", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, "user-42")

		userID, err := GetUserIDFromContext(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if userID != "user-42" {
			t.Errorf("userID = %q, want %q", userID, "user-42")
		}
	})

	t.Run("value missing", func(t *testing.T) {
		_, err := GetUserIDFromContext(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("value wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, 123)

		_, err := GetUserIDFromContext(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
