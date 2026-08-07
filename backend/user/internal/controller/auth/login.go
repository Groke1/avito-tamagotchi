package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if !isValidEmail(req.Email) || req.Password == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	tokens, err := c.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrInvalidCredentials)
			return
		}
		c.logger.Error("Login", zap.Error(err), zap.String("email", req.Email))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, tokensResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
