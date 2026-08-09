package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.pkg/middleware"
	httppets "github.com/cayman444/avito-gamification-hackathon.user/internal/adapter/pets/http"
	httptasks "github.com/cayman444/avito-gamification-hackathon.user/internal/adapter/tasks/http"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/config"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller"
	rewardrepo "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/reward"
	streakrepo "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/streak"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/repository/token/postgres"
	tokenrepo "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/token/redis"
	userrepo "github.com/cayman444/avito-gamification-hackathon.user/internal/repository/user"
	authserv "github.com/cayman444/avito-gamification-hackathon.user/internal/usecase/auth"
	rewardserv "github.com/cayman444/avito-gamification-hackathon.user/internal/usecase/reward"
	userserv "github.com/cayman444/avito-gamification-hackathon.user/internal/usecase/user"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/worker"
	"github.com/cayman444/avito-gamification-hackathon.user/migrations"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	redis "github.com/redis/go-redis/v9"
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

	httpClient := &http.Client{
		Timeout: cfg.Settings.HTTPClientTimeout,
	}

	petClient, err := httppets.NewPetClient(cfg.Clients.PetsHTTPAddr, httpClient)
	if err != nil {
		logger.Error("can not create pet client", zap.Error(err))
	}

	tasksClient, err := httptasks.NewTasksClient(cfg.Clients.TasksHTTPAddr, httpClient)
	if err != nil {
		logger.Error("can not create tasks client", zap.Error(err))
	}

	userRepo := userrepo.NewUserRepository(dbPool)
	streakRepo := streakrepo.NewStreakRepository(dbPool)
	rewardRepo := rewardrepo.NewRewardRepository(dbPool)
	transactor := db.NewTransactor(dbPool)

	var (
		tokenRepo    authserv.TokenRepository
		tokenCleaner worker.Cleaner
		redisClient  *redis.Client
	)

	switch cfg.Settings.SessionStore {
	case config.SessionStorePostgres:
		pgTokenRepo := postgres.NewTokenRepository(dbPool)
		tokenRepo = pgTokenRepo
		tokenCleaner = pgTokenRepo
	case config.SessionStoreRedis:
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.ConstructRedisAddr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			logger.Error("can not connect to redis", zap.Error(err))
			return err
		}
		defer func() {
			if err := redisClient.Close(); err != nil {
				logger.Warn("can not close redis client", zap.Error(err))
			}
		}()

		tokenRepo = tokenrepo.NewRedisRepository(redisClient)
	default:
		return errors.New("unsupported session store")
	}

	rewardService := rewardserv.NewRewardService(rewardRepo)

	authService := authserv.NewAuthService(userRepo, tokenRepo, transactor,
		rewardRepo, streakRepo, rewardService, authserv.Config{
			JWTSecret:              []byte(cfg.Settings.JWTSecret),
			AccessTokenTTL:         cfg.Settings.AccessTokenTTL,
			RefreshTokenTTL:        cfg.Settings.RefreshTokenTTL,
			RegistrationBonusCoins: cfg.Settings.RegistrationBonusCoins,
		})
	userService := userserv.NewUserService(userRepo, streakRepo, rewardRepo, transactor, petClient, tasksClient)

	router := mux.NewRouter()
	contr := controller.NewController(logger, authService,
		userService, rewardService, authService)
	contr.InitRoutes(router)

	handler := middleware.CorsHandler(router)
	server := &http.Server{
		Addr:              ":" + cfg.HTTP.Port,
		Handler:           handler,
		ReadHeaderTimeout: cfg.Settings.ServerReadHeaderTimeout,
	}

	if tokenCleaner != nil {
		go worker.StartTokenCleaner(ctx, logger, tokenCleaner, cfg.Settings.TokenCleanupInterval)
	}

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
