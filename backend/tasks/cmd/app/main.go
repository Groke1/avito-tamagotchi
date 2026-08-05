package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
	taskhttp "github.com/cayman444/avito-gamification-hackathon.tasks/internal/http"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is required")
	}

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&postgres.TaskModel{}, &postgres.UserTaskModel{}); err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewTaskRepository(db)
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
