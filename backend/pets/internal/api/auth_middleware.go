package api

import (
	"context"
	"net/http"
	"strings"
)

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// TODO
			}

			jwtToken := strings.Split(authHeader, " ")[1]
			_ = jwtToken

			userID := int64(1) // TODO

			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
