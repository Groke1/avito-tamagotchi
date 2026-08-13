package user

import (
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrUnauthorized)
		return
	}

	profile, err := c.service.Profile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrUnauthorized)
			return
		}
		c.logger.Error("Profile", zap.Error(err), zap.String("userID", userID))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, profileResponse{
		UserID:   profile.ID,
		Username: profile.Username,
		Email:    profile.Email,
		Coins:    profile.Coins,
		Streak:   profile.CurrentStreak,
	})
}
