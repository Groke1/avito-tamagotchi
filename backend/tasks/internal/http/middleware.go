package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"strings"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
	"github.com/gin-gonic/gin"
)

const (
	ExpectedTokenParts = 3
)

type accessTokenClaims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := bearerToken(c.GetHeader("Authorization"))
		if err != nil {
			SendError(c, err)
			c.Abort()
			return
		}

		claims, err := parseAccessToken(token, jwtSecret)
		if err != nil {
			SendError(c, err)
			c.Abort()
			return
		}

		c.Set("userID", claims.Sub)
		c.Next()
	}
}
func bearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", controller.ErrInvalidToken
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 ||
		!strings.EqualFold(parts[0], "Bearer") ||
		parts[1] == "" {
		return "", controller.ErrInvalidToken
	}

	return parts[1], nil
}

func parseAccessToken(token string, secret []byte) (*accessTokenClaims, error) {
	if len(secret) == 0 {
		return nil, controller.ErrInvalidToken
	}

	parts := strings.Split(token, ".")
	if len(parts) != ExpectedTokenParts {
		return nil, controller.ErrInvalidToken
	}

	encodedHeader := parts[0]
	encodedPayload := parts[1]
	encodedSignature := parts[2]

	headerBytes, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil {
		return nil, controller.ErrInvalidToken
	}

	var header jwtHeader
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return nil, controller.ErrInvalidToken
	}

	if header.Alg != "HS256" || header.Typ != "JWT" {
		return nil, controller.ErrInvalidToken
	}

	signingInput := encodedHeader + "." + encodedPayload
	expectedSignature := signJWT(signingInput, secret)
	actualSignature, err := base64.RawURLEncoding.DecodeString(encodedSignature)
	if err != nil {
		return nil, controller.ErrInvalidToken
	}

	expectedSignatureBytes, err := base64.RawURLEncoding.DecodeString(expectedSignature)
	if err != nil {
		return nil, controller.ErrInvalidToken
	}

	if !hmac.Equal(actualSignature, expectedSignatureBytes) {
		return nil, controller.ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return nil, controller.ErrInvalidToken
	}

	var claims accessTokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, controller.ErrInvalidToken
	}

	if claims.Sub == "" || claims.Exp == 0 || claims.Iat == 0 {
		return nil, controller.ErrInvalidToken
	}

	now := time.Now().UTC().Unix()

	if claims.Exp <= now {
		return nil, controller.ErrExpiredToken
	}

	const allowedClockSkew = 30 * time.Second
	if claims.Iat > now+int64(allowedClockSkew.Seconds()) {
		return nil, controller.ErrInvalidToken
	}

	if claims.Exp <= claims.Iat {
		return nil, controller.ErrInvalidToken
	}
	return &claims, nil
}

func signJWT(signingInput string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(signingInput))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
