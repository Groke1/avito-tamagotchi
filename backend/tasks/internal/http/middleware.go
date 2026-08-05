package http

import (
	"net/http"
	"strings"

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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			sendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Требуется повторная авторизация")
			c.Abort()
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] == "" {
			sendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Невалидный токен авторизации")
			c.Abort()
			return
		}

		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			sendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Идентификатор пользователя не найден")
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
