package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func sendError(c *gin.Context, status int, code string, message string) {
	c.JSON(status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

type accessTokenClaims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

var (
	errInvalidToken = errors.New("invalid token")
	errExpiredToken = errors.New("token expired")
)

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := bearerToken(c.GetHeader("Authorization"))
		if err != nil {
			sendError(
				c,
				http.StatusUnauthorized,
				"UNAUTHORIZED",
				"Невалидный токен авторизации",
			)
			c.Abort()
			return
		}

		claims, err := parseAccessToken(token, jwtSecret)
		if err != nil {
			message := "Невалидный токен авторизации"
			if errors.Is(err, errExpiredToken) {
				message = "Срок действия токена истёк"
			}

			sendError(
				c,
				http.StatusUnauthorized,
				"UNAUTHORIZED",
				message,
			)
			c.Abort()
			return
		}

		c.Set("userID", claims.Sub)
		c.Next()
	}
}

func bearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errInvalidToken
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 ||
		!strings.EqualFold(parts[0], "Bearer") ||
		parts[1] == "" {
		return "", errInvalidToken
	}

	return parts[1], nil
}

func parseAccessToken(token string, secret []byte) (*accessTokenClaims, error) {
	if len(secret) == 0 {
		return nil, errInvalidToken
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errInvalidToken
	}

	encodedHeader := parts[0]
	encodedPayload := parts[1]
	encodedSignature := parts[2]

	headerBytes, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil {
		return nil, errInvalidToken
	}

	var header jwtHeader
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, errInvalidToken
	}

	if header.Alg != "HS256" || header.Typ != "JWT" {
		return nil, errInvalidToken
	}

	signingInput := encodedHeader + "." + encodedPayload
	expectedSignature := signJWT(signingInput, secret)
	actualSignature, err := base64.RawURLEncoding.DecodeString(encodedSignature)
	if err != nil {
		return nil, errInvalidToken
	}

	expectedSignatureBytes, err := base64.RawURLEncoding.DecodeString(expectedSignature)
	if err != nil {
		return nil, errInvalidToken
	}

	if !hmac.Equal(actualSignature, expectedSignatureBytes) {
		return nil, errInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return nil, errInvalidToken
	}

	var claims accessTokenClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errInvalidToken
	}

	if claims.Sub == "" || claims.Exp == 0 || claims.Iat == 0 {
		return nil, errInvalidToken
	}

	now := time.Now().UTC().Unix()

	if claims.Exp <= now {
		return nil, errExpiredToken
	}

	const allowedClockSkew = 30 * time.Second
	if claims.Iat > now+int64(allowedClockSkew.Seconds()) {
		return nil, errInvalidToken
	}

	if claims.Exp <= claims.Iat {
		return nil, errInvalidToken
	}
	return &claims, nil
}

func signJWT(signingInput string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(signingInput))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
