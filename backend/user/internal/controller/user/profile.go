package user

import (
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"go.uber.org/zap"
)

func (c *controller) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, errUnauthorized)
		return
	}

	profile, err := c.userService.Profile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			writeError(w, http.StatusUnauthorized, errUnauthorized)
			return
		}
		c.logger.Error("Profile", zap.Error(err), zap.String("userID", userID))
		writeError(w, http.StatusInternalServerError, errInternalError)
		return
	}

	writeJSON(w, http.StatusOK, profileResponse{
		UserID:   profile.ID,
		Username: profile.Username,
		Email:    profile.Email,
		Coins:    profile.Coins,
	})
}
