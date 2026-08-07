package reward

import (
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (c *controller) GetReward(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rewardID := vars["reward_id"]

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrUnauthorized)
		return
	}
	userReward, err := c.service.GetReward(r.Context(), userID, rewardID)

	if err != nil {
		if errors.Is(err, entity.ErrRewardNotFound) {
			httpx.WriteError(w, http.StatusNotFound, httpx.ErrRewardNotFound)
			return
		}

		c.logger.Error("failed to get reward", zap.String("user_id", userID), zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, toUserRewardResponse(*userReward))
}
