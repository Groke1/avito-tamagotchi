package auth

import (
	"context"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

//go:generate mockgen -source=controller.go -destination=mocks/controller_mocks.go -package=mocks

type Service interface {
	Register(ctx context.Context, user entity.User) (*entity.JWT, error)
	Login(ctx context.Context, email, password string) (*entity.JWT, error)
	Refresh(ctx context.Context, refreshToken string) (*entity.JWT, error)
	Logout(ctx context.Context, userID, refreshToken string) error
}

type controller struct {
	logger  *zap.Logger
	service Service
}

func NewController(logger *zap.Logger, service Service) *controller {
	return &controller{logger: logger, service: service}
}

func (c *controller) InitRoutes(api *mux.Router, _ *mux.Router, authMiddleware func(http.Handler) http.Handler) {
	router := api.PathPrefix("/auth").Subrouter()
	router.HandleFunc("/register", c.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", c.Login).Methods(http.MethodPost)
	router.HandleFunc("/refresh", c.Refresh).Methods(http.MethodPost)
	logout := router.PathPrefix("/logout").Subrouter()
	logout.Use(authMiddleware)
	logout.HandleFunc("", c.Logout).Methods(http.MethodPost)
}
