package controller

import (
	"net/http"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/auth"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/event"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/reward"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/user"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type controllerInterface interface {
	InitRoutes(router *mux.Router, internal *mux.Router, authMiddleware func(http.Handler) http.Handler)
}

type controller struct {
	auth   controllerInterface
	user   controllerInterface
	reward controllerInterface
	event  controllerInterface

	tokenValidator middleware.AccessTokenValidator
}

func NewController(
	logger *zap.Logger,
	authService auth.Service,
	userService user.Service,
	rewardService reward.Service,
	eventService event.Service,
	tokenValidator middleware.AccessTokenValidator,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *controller {
	return &controller{
		auth:           auth.NewController(logger, authService, accessTokenTTL, refreshTokenTTL),
		user:           user.NewController(logger, userService),
		reward:         reward.NewController(logger, rewardService),
		event:          event.NewController(logger, eventService),
		tokenValidator: tokenValidator,
	}
}

func (c *controller) InitRoutes(router *mux.Router) {
	api := router.PathPrefix("/api/v1").Subrouter()
	internal := router.PathPrefix("/internal").Subrouter()
	requireAccessToken := middleware.RequireAccessToken(c.tokenValidator)

	c.auth.InitRoutes(api, internal, requireAccessToken)
	c.user.InitRoutes(api, internal, requireAccessToken)
	c.reward.InitRoutes(api, internal, requireAccessToken)
	c.event.InitRoutes(api, internal, requireAccessToken)
}
