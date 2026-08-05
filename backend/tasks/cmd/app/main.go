package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
	taskhttp "github.com/cayman444/avito-gamification-hackathon.tasks/internal/http"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres"
	"github.com/cayman444/avito-gamification-hackathon.tasks/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

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

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is required")
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	if err := runMigrations(dsn); err != nil {
		log.Fatalf("migration error: %v", err)
	}
	cnx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(cnx, dsn)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(cnx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	repo := postgres.NewTaskRepository(pool)
	handlers := taskhttp.NewHandlers(
		controller.NewGetTaskHandler(repo),
		controller.NewGetTodayTasksHandler(repo),
		controller.NewCompleteTaskHandler(repo),
	)
	router := taskhttp.NewRouter(handlers)

	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("tasks service listening on %s", addr)
	log.Fatal(server.ListenAndServe())
}
