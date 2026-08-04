package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/api"
	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/repository"
	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://postgres:postgres@localhost:5432/pet_db?sslmode=disable"
	driver := "postgres"

	db, err := sqlx.Connect(driver, dsn)
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

	r.Route("/pet", func(r chi.Router) {
		r.Use(api.AuthMiddleware("jwt-secret")) // TODO
		r.Get("/", handler.GetPet)
		r.Post("/", handler.CreatePet)
	})

	fmt.Println("RUN")
	http.ListenAndServe(":8080", r)
}
