package service

import (
	"context"
	"fmt"
	"log"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
)

type JourneyStoryGenerator interface {
	Generate(ctx context.Context, input domain.JourneyGenerationInput) (domain.JourneyStory, error)
}

type FallbackStoryGenerator struct {
	primary  JourneyStoryGenerator
	fallback JourneyStoryGenerator
}

func NewFallbackStoryGenerator(primary, fallback JourneyStoryGenerator) *FallbackStoryGenerator {
	return &FallbackStoryGenerator{primary: primary, fallback: fallback}
}

func (g *FallbackStoryGenerator) Generate(ctx context.Context, input domain.JourneyGenerationInput) (domain.JourneyStory, error) {
	story, err := g.primary.Generate(ctx, input)
	if err == nil {
		return story, nil
	}

	log.Printf("primary story generator failed, falling back to template: %v", err)
	return g.fallback.Generate(ctx, input)
}

type TemplateStoryGenerator struct{}

func NewTemplateStoryGenerator() *TemplateStoryGenerator {
	return &TemplateStoryGenerator{}
}

func (g *TemplateStoryGenerator) Generate(_ context.Context, input domain.JourneyGenerationInput) (domain.JourneyStory, error) {
	j := input.Journey

	story := fmt.Sprintf("Я вернулся из путешествия в «%s»! "+
		"В этот раз ничего интересного не произошло.", j.Location)

	return domain.JourneyStory{
		Title: fmt.Sprintf("Возвращение из «%s»", j.Location),
		Story: story,
		// Teaser: "Кажется, рядом осталось ещё что-то неизведанное 👀",
	}, nil
}
