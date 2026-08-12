package ai

import (
	"context"
	"log"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
)

type SafeStoryGenerator struct {
	primary  domain.JourneyStoryGenerator
	fallback domain.JourneyStoryGenerator
}

func NewSafeStoryGenerator(primary, fallback domain.JourneyStoryGenerator) *SafeStoryGenerator {
	return &SafeStoryGenerator{primary: primary, fallback: fallback}
}

func (s *SafeStoryGenerator) Generate(ctx context.Context, journey domain.JourneyResult) (domain.JourneyStory, error) {
	if s.primary == nil {
		return s.fallback.Generate(ctx, journey)
	}

	story, err := s.primary.Generate(ctx, journey)
	if err != nil {
		log.Printf("[STORY GENERATOR] primary (gigachat) failed, falling back to template: %v", err)
		return s.fallback.Generate(ctx, journey)
	}

	return story, nil
}
