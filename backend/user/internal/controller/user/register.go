package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, errValidationError)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if !isValidUsername(req.Username) || !isValidEmail(req.Email) || !isValidPassword(req.Password) {
		writeError(w, http.StatusUnprocessableEntity, errValidationError)
		return
	}

	jwt, err := c.authService.Register(r.Context(), entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, entity.ErrUsernameAlreadyExists) || errors.Is(err, entity.ErrEmailAlreadyExists) {
			writeError(w, http.StatusConflict, errUserAlreadyExists)
			return
		}
		c.logger.Error("Register: ", zap.Error(err), zap.String("username", req.Username))
		writeError(w, http.StatusInternalServerError, errInternalError)
		return
	}

	writeJSON(w, http.StatusCreated, authResponse{
		AccessToken:  jwt.AccessToken,
		RefreshToken: jwt.RefreshToken,
	})
}
