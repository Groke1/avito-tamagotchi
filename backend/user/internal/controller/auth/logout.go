package auth

import (
	"net/http"
	"strings"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"go.uber.org/zap"
)

func (c *controller) Logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	req.RefreshToken = strings.TrimSpace(req.RefreshToken)
	if req.RefreshToken == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrUnauthorized)
		return
	}

	err := c.service.Logout(r.Context(), userID, req.RefreshToken)

	if err != nil {
		c.logger.Error("failed to logout user", zap.String("user_id", userID), zap.Error(err))

		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	middleware.ClearAuthCookies(w, r)
	w.WriteHeader(http.StatusNoContent)
}
