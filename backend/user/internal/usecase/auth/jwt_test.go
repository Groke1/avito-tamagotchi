package auth

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

func newTestService(ttl time.Duration) *authService {
	return &authService{cfg: &Config{JWTSecret: []byte("test-secret"), AccessTokenTTL: ttl}}
}

func buildRawToken(t *testing.T, secret []byte, claims accessTokenClaims) string {
	t.Helper()
	header, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("marshal claims: %v", err)
	}
	signingInput := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	return signingInput + "." + signJWT(signingInput, secret)
}

func TestNewAccessToken_ValidatesSuccessfully(t *testing.T) {
	s := newTestService(time.Hour)

	token, err := s.newAccessToken("user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sub, err := s.ValidateAccessToken(nil, token)
	if err != nil {
		t.Fatalf("expected token to validate, got error: %v", err)
	}
	if sub != "user-1" {
		t.Fatalf("expected sub user-1, got %q", sub)
	}
}

func TestValidateAccessToken_Expired(t *testing.T) {
	s := newTestService(time.Hour)
	now := time.Now().UTC()
	token := buildRawToken(t, s.cfg.JWTSecret, accessTokenClaims{Sub: "user-1", Iat: now.Add(-2 * time.Hour).Unix(), Exp: now.Add(-time.Minute).Unix()})

	_, err := s.ValidateAccessToken(nil, token)
	if err != entity.ErrInvalidAccessToken {
		t.Fatalf("expected ErrInvalidAccessToken, got %v", err)
	}
}

func TestValidateAccessToken_MalformedParts(t *testing.T) {
	s := newTestService(time.Hour)

	_, err := s.ValidateAccessToken(nil, "only.two")
	if err != entity.ErrInvalidAccessToken {
		t.Fatalf("expected ErrInvalidAccessToken, got %v", err)
	}
}

func TestValidateAccessToken_TamperedSignature(t *testing.T) {
	s := newTestService(time.Hour)
	now := time.Now().UTC()
	token := buildRawToken(t, s.cfg.JWTSecret, accessTokenClaims{Sub: "user-1", Iat: now.Unix(), Exp: now.Add(time.Hour).Unix()})
	tampered := token[:len(token)-1] + "x"

	_, err := s.ValidateAccessToken(nil, tampered)
	if err != entity.ErrInvalidAccessToken {
		t.Fatalf("expected ErrInvalidAccessToken, got %v", err)
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	s := newTestService(time.Hour)
	now := time.Now().UTC()
	token := buildRawToken(t, []byte("other-secret"), accessTokenClaims{Sub: "user-1", Iat: now.Unix(), Exp: now.Add(time.Hour).Unix()})

	_, err := s.ValidateAccessToken(nil, token)
	if err != entity.ErrInvalidAccessToken {
		t.Fatalf("expected ErrInvalidAccessToken, got %v", err)
	}
}

func TestValidateAccessToken_MissingSub(t *testing.T) {
	s := newTestService(time.Hour)
	now := time.Now().UTC()
	token := buildRawToken(t, s.cfg.JWTSecret, accessTokenClaims{Sub: "", Iat: now.Unix(), Exp: now.Add(time.Hour).Unix()})

	_, err := s.ValidateAccessToken(nil, token)
	if err != entity.ErrInvalidAccessToken {
		t.Fatalf("expected ErrInvalidAccessToken, got %v", err)
	}
}

func TestValidateAccessToken_InvalidBase64Payload(t *testing.T) {
	s := newTestService(time.Hour)
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	signingInput := header + ".not-valid-base64!!!"
	token := signingInput + "." + signJWT(signingInput, s.cfg.JWTSecret)

	_, err := s.ValidateAccessToken(nil, token)
	if err != entity.ErrInvalidAccessToken {
		t.Fatalf("expected ErrInvalidAccessToken, got %v", err)
	}
}
