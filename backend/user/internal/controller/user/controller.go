package user

import (
	"context"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Service interface {
	Profile(ctx context.Context, userID string) (*entity.User, error)
	GetUsers(ctx context.Context, userIDs []string) ([]entity.User, error)
	UpdateCoins(ctx context.Context, userID string, coins int64) (*entity.User, error)
	UpdateStreak(ctx context.Context, userID, occurredAt string) error
}

type controller struct {
	logger  *zap.Logger
	service Service
}

func NewController(logger *zap.Logger, service Service) *controller {
	return &controller{logger: logger, service: service}
}

func (c *controller) InitRoutes(
	api *mux.Router,
	internal *mux.Router,
	authMiddleware func(http.Handler) http.Handler,
) {
	profile := api.PathPrefix("/profile").Subrouter()
	profile.Use(authMiddleware)
	profile.HandleFunc("", c.Profile).Methods(http.MethodGet)

	internal.HandleFunc("/usernames", c.Usernames).Methods(http.MethodPost)
	internal.HandleFunc("/update-coins", c.UpdateCoins).Methods(http.MethodPut)
	internal.HandleFunc("/action", c.Action).Methods(http.MethodPost)
}
