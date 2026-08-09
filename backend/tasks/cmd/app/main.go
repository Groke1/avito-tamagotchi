package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	"github.com/cayman444/avito-gamification-hackathon.pkg/middleware"
	client "github.com/cayman444/avito-gamification-hackathon.tasks/internal/clients"
	config "github.com/cayman444/avito-gamification-hackathon.tasks/internal/config"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
	taskhttp "github.com/cayman444/avito-gamification-hackathon.tasks/internal/http"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres"
	"github.com/cayman444/avito-gamification-hackathon.tasks/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	dsn := cfg.ConstructPostgresURL()
	err = runMigrations(dsn)
	if err != nil {
		return fmt.Errorf("migration error: %v", err)
	}

	jwtSecret := cfg.JWT.Secret
	if jwtSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	cnx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(cnx, dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	err = pool.Ping(cnx)
	if err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}

	repo := postgres.NewTaskRepository(pool)
	transactor := db.NewTransactor(pool)
	userClient := client.NewUserServiceClient(cfg.Client.CoinsServiceURL)
	petClient := client.NewPetServiceClient(cfg.Pet.PetServiceURL)
	handlers := taskhttp.NewHandlers(
		controller.NewGetTaskHandler(repo),
		controller.NewGetTodayTasksHandler(repo),
		controller.NewCompleteTaskHandler(repo, userClient, petClient, transactor),
		controller.NewGetCompletedTasksHandler(repo),
	)
	router := taskhttp.NewRouter(handlers, []byte(jwtSecret))
	handler := middleware.CorsHandler(router)
	server := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      handler,
		ReadTimeout:  config.ServerReadTimeout,
		WriteTimeout: config.HTTPClientTimeout,
	}

	log.Printf("tasks service listening on %s", cfg.HTTP.Addr)

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server stopped: %w", err)
	}

	return nil
}

func runMigrations(dsn string) error {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open db for migrations: %w", err)
	}
	defer db.Close()
	log.Println("applying database migrations...")
	if err := migrations.RunMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("migrations applied successfully")
	return nil
}
