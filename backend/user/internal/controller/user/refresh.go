package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := decodeJSON(w, r, &req); err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		writeError(w, http.StatusUnauthorized, errInvalidRefreshToken)
		return
	}

	tokens, err := c.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidRefreshToken) || errors.Is(err, entity.ErrRefreshTokenNotFound) {
			writeError(w, http.StatusUnauthorized, errInvalidRefreshToken)
			return
		}
		c.logger.Error("Refresh", zap.Error(err))
		writeError(w, http.StatusInternalServerError, errInternalError)
		return
	}

	writeJSON(w, http.StatusOK, authResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
