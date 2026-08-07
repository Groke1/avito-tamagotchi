package main

import (
	"log"
	"net/http"

	cors "github.com/cayman444/avito-gamification-hackathon.pkg/middleware"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/api"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/config"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env not found: %v", err)
	}

	cfg := config.NewConfig()

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	repository := repository.NewPetRepository(db)
	service := service.NewPetService(repository, cfg.UserServiceURL)
	handler := api.NewPetHandler(service)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(cors.CorsHandler)
		r.Use(api.JwtMiddleware(cfg.JwtSecret))

		r.Route("/pet", func(r chi.Router) {
			r.Get("/", handler.GetPet)
			r.Post("/", handler.CreatePet)
			r.Post("/feed", handler.FeedPet)
			r.Post("/stroke", handler.StrokePet)
		})

		r.Get("/leaderboard", handler.GetLeaderboard)
	})

	r.Route("/internal", func(r chi.Router) {
		r.Post("/daily-bonus", handler.DailyBonus)
	})

	log.Println("[SERVICE STARTED]")
	err = http.ListenAndServe(":"+cfg.HttpPort, r)
	if err != nil {
		log.Fatal(err)
	}
}
