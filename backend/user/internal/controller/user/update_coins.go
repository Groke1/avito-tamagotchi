package user

import (
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) UpdateCoins(w http.ResponseWriter, r *http.Request) {
	var req updateCoinsRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	updatedUser, err := c.service.UpdateCoins(r.Context(), req.UserID, req.Delta)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			httpx.WriteError(w, http.StatusNotFound, httpx.ErrUserNotFound)
		case errors.Is(err, entity.ErrInsufficientCoins):
			httpx.WriteError(w, http.StatusConflict, httpx.ErrInsufficientCoins)
		default:
			c.logger.Error("failed to update user coins",
				zap.String("user_id", req.UserID),
				zap.Int64("delta", req.Delta),
				zap.Error(err),
			)
			httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		}
		return
	}

	httpx.WriteJSON(w, http.StatusOK, updateCoinsResponse{
		UserID: updatedUser.ID,
		Coins:  updatedUser.Coins,
	})
}
