package main

import (
	"context"
	"log"
	"net/http"

	cors "github.com/cayman444/avito-gamification-hackathon.pkg/middleware"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/api"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/clients/gigachat"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/config"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/service"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/websocket"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/worker"

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

	petRepository := repository.NewPetRepository(db)

	wsClientManager := websocket.NewUserManager()
	wsTicketManager := websocket.NewTicketManager()
	levelPolicy := domain.NewLevelPolicy()

	gigaConfig, err := gigachat.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("[MAIN] Failed to create GigaChat config: %v", err)
	}
	gigaClient := gigachat.NewClient(gigaConfig)
	templateClient := service.NewTemplateStoryGenerator()
	clientOrchestrator := service.NewFallbackStoryGenerator(gigaClient, templateClient)
	executor := worker.NewGoroutineWorker()

	trip_repo := repository.NewTripRepository(db)
	tripRepository := repository.NewTripRepository(db)

	tripService := service.NewTripService(clientOrchestrator, executor, trip_repo, petRepository, wsClientManager)
	petService := service.NewPetService(petRepository, tripRepository, cfg.UserServiceURL, wsClientManager, levelPolicy)
	petHandler := api.NewPetHandler(petService, tripService)
	wsHandler := api.NewWSHandler(wsClientManager, wsTicketManager, petService)

	tripWorker := worker.NewTripWorker(tripRepository, tripService)
	go tripWorker.Run(context.Background())

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(cors.CorsHandler)

		r.Get("/ws", wsHandler.OpenWSConn)

		r.Group(func(r chi.Router) {
			r.Use(api.JwtMiddleware(cfg.JwtSecret))

			r.Route("/pet", func(r chi.Router) {
				r.Get("/", petHandler.GetPet)
				r.Post("/", petHandler.CreatePet)
				r.Post("/feed", petHandler.FeedPet)
				r.Post("/stroke", petHandler.StrokePet)
				r.Post("/ws-ticket", wsHandler.CreateWSTicket)
				r.Get("/trip/last", petHandler.LastTrip)
				r.Post("/trip/{pet_id}", petHandler.MakeTrip)
			})

			r.Route("/leaderboard", func(r chi.Router) {
				r.Get("/", petHandler.GetLeaderboard)
				r.Get("/weekly", petHandler.GetWeeklyLeaderboard)
			})

			r.Get("/next-reward-description", petHandler.GetNextRewardDescription)
		})
	})

	r.Route("/internal", func(r chi.Router) {
		r.Post("/daily-bonus", petHandler.DailyBonusForStreak)
		r.Get("/daily-gained-xp", petHandler.DailyGainedXP)
		r.Put("/update-xp", petHandler.UpdateXP)
	})

	log.Printf("[MAIN] Service starting on port %s", ":"+cfg.HTTPPort)
	err = http.ListenAndServe(":"+cfg.HTTPPort, r)
	if err != nil {
		log.Printf("[MAIN] http server shut down with error: %v", err)
	}
}
