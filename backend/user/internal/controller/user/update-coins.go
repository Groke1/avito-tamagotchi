package user

import (
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) UpdateCoins(w http.ResponseWriter, r *http.Request) {
	var req updateCoinsRequest

	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusUnprocessableEntity, errValidationError)
		return
	}

	err := c.userService.UpdateCoins(r.Context(), req.UserID, req.Delta)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			writeError(w, http.StatusNotFound, errUserNotFound)

		case errors.Is(err, entity.ErrInsufficientCoins):
			writeError(w, http.StatusConflict, errInsufficientCoins)

		default:
			c.logger.Error("failed to update user coins", zap.String("user_id", req.UserID),
				zap.Int64("delta", req.Delta), zap.Error(err))

			writeError(w, http.StatusInternalServerError, errInternalError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
