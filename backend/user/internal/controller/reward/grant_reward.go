package reward

import (
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"go.uber.org/zap"
)

func (c *controller) GrantReward(w http.ResponseWriter, r *http.Request) {
	var req grantRewardRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	reward, err := c.service.GrantReward(r.Context(), req.UserID, req.Code)

	if err != nil {
		c.logger.Error("failed to grant reward", zap.String("user_id", req.UserID),
			zap.String("code", req.Code), zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, toUserRewardResponse(*reward))
}
