package user

import (
	"net/http"

	"go.uber.org/zap"
)

const maxUserIDs = 1000

func (c *controller) Usernames(w http.ResponseWriter, r *http.Request) {
	var req usernamesRequest

	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, errValidationError)
		return
	}

	if len(req.UserIDs) == 0 || len(req.UserIDs) > maxUserIDs {
		writeError(w, http.StatusUnprocessableEntity, errValidationError)
		return
	}

	users, err := c.userService.GetUsers(r.Context(), req.UserIDs)
	if err != nil {
		c.logger.Error("failed to get usernames", zap.Error(err))

		writeError(w, http.StatusInternalServerError, errInternalError)
		return
	}

	responseUsers := make([]usernameResponse, len(users))
	for i, user := range users {
		responseUsers[i] = usernameResponse{
			ID:       user.ID,
			Username: user.Username,
		}
	}

	writeJSON(w, http.StatusOK,
		usernamesResponse{
			Users: responseUsers,
		},
	)
}
