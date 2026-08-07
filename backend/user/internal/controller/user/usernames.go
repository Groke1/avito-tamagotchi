package user

import (
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"go.uber.org/zap"
)

const maxUserIDs = 1000

func (c *controller) Usernames(w http.ResponseWriter, r *http.Request) {
	var req usernamesRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	if len(req.UserIDs) == 0 || len(req.UserIDs) > maxUserIDs {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	users, err := c.service.GetUsers(r.Context(), req.UserIDs)
	if err != nil {
		c.logger.Error("failed to get usernames", zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	responseUsers := make([]usernameResponse, len(users))
	for i, foundUser := range users {
		responseUsers[i] = usernameResponse{ID: foundUser.ID, Username: foundUser.Username}
	}

	httpx.WriteJSON(w, http.StatusOK, usernamesResponse{Users: responseUsers})
}
