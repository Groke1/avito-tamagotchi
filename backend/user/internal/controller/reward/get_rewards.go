package reward

import (
	"context"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) GetAllRewards(w http.ResponseWriter, r *http.Request) {
	c.getRewardsImpl(w, r, "GetAllRewards", c.service.GetAllRewards)
}

func (c *controller) GetActiveRewards(w http.ResponseWriter, r *http.Request) {
	c.getRewardsImpl(w, r, "GetActiveRewards", c.service.GetActiveRewards)
}

func (c *controller) getRewardsImpl(
	w http.ResponseWriter, r *http.Request,
	logMethod string, getRewards func(context.Context, string) ([]entity.UserReward, error),
) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrUnauthorized)
		return
	}

	rewards, err := getRewards(r.Context(), userID)
	if err != nil {
		c.logger.Error(logMethod, zap.Error(err), zap.String("userID", userID))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	items := make([]userRewardResponse, 0, len(rewards))
	for _, reward := range rewards {
		items = append(items, toUserRewardResponse(reward))
	}

	httpx.WriteJSON(w, http.StatusOK, userListRewardResponse{
		Items: items,
	})
}
