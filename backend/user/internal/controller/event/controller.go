package event

import (
	"context"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Service interface {
	GetNewEvents(ctx context.Context, userID string) ([]entity.UserEventDetails, error)
}

type controller struct {
	logger  *zap.Logger
	service Service
}

func NewController(logger *zap.Logger, service Service) *controller {
	return &controller{
		logger:  logger,
		service: service,
	}
}

func (c controller) InitRoutes(
	api *mux.Router,
	_ *mux.Router,
	authMiddleware func(http.Handler) http.Handler,
) {
	events := api.PathPrefix("/events").Subrouter()
	events.Use(authMiddleware)
	events.HandleFunc("", c.GetNewEvents).Methods(http.MethodGet)
}
