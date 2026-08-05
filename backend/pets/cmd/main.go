package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/api"
	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/repository"
	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/service"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env not found: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	authSecret := os.Getenv("AUTH_SECRET")
	if authSecret == "" {
		log.Fatal("AUTH_SECRET is not set")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	repository := repository.NewPetRepository(db)
	service := service.NewPetService(repository)
	handler := api.NewPetHandler(service)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1/pet", func(r chi.Router) {
		r.Use(api.AuthMiddleware(authSecret))
		r.Get("/", handler.GetPet)
		r.Post("/", handler.CreatePet)
		r.Post("/feed", handler.FeedPet)
		r.Post("/stroke", handler.StrokePet)
	})

	log.Println("[SERVICE STARTED]")
	http.ListenAndServe(":"+port, r)
}
