package user

import (
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"go.uber.org/zap"
)

func (c *controller) DailyStat(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrUnauthorized)
		return
	}

	stat, err := c.service.GetDailyStat(r.Context(), userID)
	if err != nil {
		c.logger.Error("Stat", zap.Error(err), zap.String("userID", userID))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, stat)
}
