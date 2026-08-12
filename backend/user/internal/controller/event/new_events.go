package event

import (
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"go.uber.org/zap"
)

func (c *controller) GetNewEvents(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, httpx.ErrUnauthorized)
		return
	}

	events, err := c.service.GetNewEvents(r.Context(), userID)
	if err != nil {
		c.logger.Error("failed to get new rewards", zap.Error(err),
			zap.String("userID", userID))
		httpx.WriteError(w, http.StatusInternalServerError, httpx.ErrInternal)
		return
	}

	items := make([]eventResponse, 0, len(events))
	for _, event := range events {
		items = append(items, toEventResponse(&event))
	}

	httpx.WriteJSON(w, http.StatusOK, listEventsResponse{
		UserID: userID,
		Items:  items,
	})
}
