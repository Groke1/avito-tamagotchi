package user

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

const refreshTokenSize = 32

func (s *userService) generateTokens(ctx context.Context, userID string) (*entity.JWT, error) {
	accessToken, err := s.newAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenHash, refreshExpiresAt, err := s.newRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.tokenRepository.AddToken(ctx, userID, entity.RefreshToken{
		TokenHash: refreshTokenHash,
		ExpiresAt: refreshExpiresAt,
	}); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &entity.JWT{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

type accessTokenClaims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

func (s *userService) newAccessToken(userID string) (string, error) {
	now := time.Now().UTC()
	header, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("marshal jwt header: %w", err)
	}

	payload, err := json.Marshal(accessTokenClaims{
		Sub: userID,
		Exp: now.Add(s.cfg.AccessTokenTTL).Unix(),
		Iat: now.Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("marshal jwt payload: %w", err)
	}

	signingInput := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	signature := signJWT(signingInput, s.cfg.JWTSecret)

	return signingInput + "." + signature, nil
}

func (s *userService) newRefreshToken() (token string, tokenHash string, expiresAt time.Time, err error) {
	buf := make([]byte, refreshTokenSize)
	if _, err = rand.Read(buf); err != nil {
		return "", "", time.Time{}, fmt.Errorf("generate refresh token: %w", err)
	}

	token = base64.RawURLEncoding.EncodeToString(buf)
	tokenHash = hashToken(token)
	expiresAt = time.Now().UTC().Add(s.cfg.RefreshTokenTTL)

	return token, tokenHash, expiresAt, nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func signJWT(signingInput string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
