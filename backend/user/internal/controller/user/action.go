package user

import (
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) Action(w http.ResponseWriter, r *http.Request) {
	var req actionRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	err := c.service.UpdateStreak(r.Context(), req.UserID, req.OccurredAt)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			httpx.WriteError(w, http.StatusNotFound, httpx.ErrUserNotFound)
			return
		}

		c.logger.Error("failed to update streak", zap.String("user_id", req.UserID), zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *controller) ProtectedAction(w http.ResponseWriter, r *http.Request) {
	
}
