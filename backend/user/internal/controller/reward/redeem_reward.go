package reward

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) RedeemReward(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrUnauthorized)
		return
	}

	var req redeemRewardRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	req.PromoCode = strings.TrimSpace(req.PromoCode)
	if req.PromoCode == "" {
		httpx.WriteError(w, http.StatusUnprocessableEntity, httpx.ErrValidation)
		return
	}

	err := c.service.RedeemReward(r.Context(), userID, req.PromoCode)

	if err != nil {
		if errors.Is(err, entity.ErrRewardUnavailable) {
			httpx.WriteError(w, http.StatusNotFound, httpx.ErrRewardUnavailable)
			return
		}

		c.logger.Error("failed to redeem reward", zap.String("user_id", userID), zap.Error(err))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
