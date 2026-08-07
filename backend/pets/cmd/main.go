package main

import (
	"log"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/api"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/config"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/service"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/websocket"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("[MAIN] .env not found: %v", err)
	}

	cfg := config.NewConfig()

	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[MAIN] Database connection failed: %v", err)
	}
	defer db.Close()

	repository := repository.NewPetRepository(db)
	wsClientManager := websocket.NewClient()
	wsTicketManager := websocket.NewTicketManager()
	service := service.NewPetService(repository, cfg.UserServiceURL, wsClientManager)
	petHandler := api.NewPetHandler(service)
	wsHandler := api.NewWSHandler(wsClientManager, wsTicketManager, service)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(api.CorsMiddleware)

		r.Get("/ws", wsHandler.OpenWSConn)

		r.Group(func(r chi.Router) {
			r.Use(api.JwtMiddleware(cfg.JwtSecret))

			r.Route("/pet", func(r chi.Router) {
				r.Get("/", petHandler.GetPet)
				r.Post("/", petHandler.CreatePet)
				r.Post("/feed", petHandler.FeedPet)
				r.Post("/stroke", petHandler.StrokePet)
				r.Post("/ws-ticket", wsHandler.CreateTicket)
			})

			r.Get("/leaderboard", petHandler.GetLeaderboard)
		})
	})

	r.Route("/internal", func(r chi.Router) {
		r.Post("/daily-bonus", petHandler.DailyBonus)
	})

	log.Printf("[MAIN] Service starting on port %s", ":"+cfg.HttpPort)
	err = http.ListenAndServe(":"+cfg.HttpPort, r)
	if err != nil {
		log.Fatalf("[MAIN] http server shut down with error: %v", err)
	}
}
