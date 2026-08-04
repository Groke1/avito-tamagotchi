package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.pkg/middleware"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/config"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/user"
	tokenrepo "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/token"
	userrepo "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/user"
	userserv "github.com/cayman444/avito-gamification-hackathon.user/internal/usecase/user"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/worker"
	"github.com/cayman444/avito-gamification-hackathon.user/migrations"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger, cfg *config.Config) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pgxcfg, err := pgxpool.ParseConfig(cfg.ConstructPostgresURL())

	if err != nil {
		logger.Error("can not create pgxpool cfg", zap.Error(err))
		return err
	}

	pgxcfg.MaxConns = cfg.PG.MaxConns
	pgxcfg.MinConns = cfg.PG.MinConns
	pgxcfg.HealthCheckPeriod = cfg.Settings.PGHealthCheckPeriod
	pgxcfg.MaxConnLifetime = cfg.Settings.PGMaxConnLifetime
	pgxcfg.MaxConnIdleTime = cfg.Settings.PGMaxConnIdleTime

	dbPool, err := pgxpool.NewWithConfig(ctx, pgxcfg)

	if err != nil {
		logger.Error("can not create pgxpool", zap.Error(err))
		return err
	}
	defer dbPool.Close()

	migrations.SetupPostgres(dbPool, logger)

	userRepo := userrepo.NewUserRepository(dbPool)
	tokenRepo := tokenrepo.NewTokenRepository(dbPool)
	transactor := db.NewTransactor(dbPool)

	userService := userserv.NewUserService(userRepo, tokenRepo, transactor, userserv.Config{
		JWTSecret:              []byte(cfg.Settings.JWTSecret),
		AccessTokenTTL:         cfg.Settings.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Settings.RefreshTokenTTL,
		RegistrationBonusCoins: cfg.Settings.RegistrationBonusCoins,
	})

	router := mux.NewRouter()
	controller := user.NewController(logger, userService)
	controller.InitRoutes(
		router,
		user.RequireAccessToken(userService),
	)

	handler := middleware.CorsHandler(router)
	server := &http.Server{
		Addr:              ":" + cfg.HTTP.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go worker.StartTokenCleaner(ctx, logger, tokenRepo, cfg.Settings.TokenCleanupInterval)

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting user service", zap.String("addr", server.Addr))
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", zap.Error(err))
			return err
		}
	case <-ctx.Done():
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Settings.ShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("graceful shutdown failed", zap.Error(err))
		return err
	}

	return nil
}
