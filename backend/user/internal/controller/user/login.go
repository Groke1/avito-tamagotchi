package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, errValidationError)
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if !isValidEmail(req.Email) || req.Password == "" {
		writeError(w, http.StatusUnprocessableEntity, errValidationError)
		return
	}

	tokens, err := c.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, errInvalidCredentials)
			return
		}
		c.logger.Error("Login", zap.Error(err), zap.String("email", req.Email))
		writeError(w, http.StatusInternalServerError, errInternalError)
		return
	}

	writeJSON(w, http.StatusOK, authResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
