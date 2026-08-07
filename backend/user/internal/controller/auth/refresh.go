package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil || strings.TrimSpace(req.RefreshToken) == "" {
		httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrInvalidRefreshToken)
		return
	}

	tokens, err := c.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidRefreshToken) || errors.Is(err, entity.ErrRefreshTokenNotFound) {
			httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrInvalidRefreshToken)
			return
		}
		c.logger.Error("Refresh", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, tokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
