package service

import (
	"context"
	"log"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
)

type TripService struct {
	storyGenerator FallbackStoryGenerator
	executor       WorkerOne
}

func (s *TripService) Generate(ctx context.Context, userIDStr string) error {
	return s.executor.Execute(
		func(ctx context.Context) error {
			input := domain.JourneyGenerationInput{

				UserID: userIDStr,
			}
			story, err := s.storyGenerator.Generate(ctx, userIDStr)
			if err != nil {
				log.Printf(
					"failed to generate trip for user %s: %v",
					userIDStr,
					err,
				)
			}
			_ = story
			return nil
		},
	)
}
