package reward

import (
	"context"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Service interface {
	GetAllRewards(ctx context.Context, userID string) ([]entity.UserReward, error)
	GetActiveRewards(ctx context.Context, userID string) ([]entity.UserReward, error)
	GetReward(ctx context.Context, userID, rewardID string) (*entity.UserReward, error)
	RedeemReward(ctx context.Context, userID, promoCode string) error
	GetDefinition(ctx context.Context, code string) (*entity.RewardDefinition, error)
	GrantReward(ctx context.Context, userID, code, earnedReason string) (*entity.UserReward, error)
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
	internal *mux.Router,
	authMiddleware func(http.Handler) http.Handler,
) {
	rewards := api.PathPrefix("/rewards").Subrouter()
	rewards.Use(authMiddleware)
	rewards.HandleFunc("", c.GetAllRewards).Methods(http.MethodGet)
	rewards.HandleFunc("/active", c.GetActiveRewards).Methods(http.MethodGet)
	rewards.HandleFunc("/{reward_id}", c.GetReward).Methods(http.MethodGet)
	rewards.HandleFunc("/redeem", c.RedeemReward).Methods(http.MethodPost)

	internalRewards := internal.PathPrefix("/rewards").Subrouter()
	internalRewards.HandleFunc("", c.GrantReward).Methods(http.MethodPost)
	internalRewards.HandleFunc("/{code}", c.GetDefinition).Methods(http.MethodGet)
}
