package main

import (
	"log"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/app"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatalf("can not initialize logger: %s", err)
	}

	cfg, err := config.New()

	if err != nil {
		log.Fatalf("can not initialize config: %s", err)
	}

	if err := app.Run(logger, cfg); err != nil {
		log.Fatalf("can not initialize app: %s", err)
	}
}
