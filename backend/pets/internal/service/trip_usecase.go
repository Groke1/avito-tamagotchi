package service

import (
	"context"
	"fmt"

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

	// сюда стоит добавить лог/метрику: "gigachat unavailable, using fallback"
	return g.fallback.Generate(ctx, input)
}

type TemplateStoryGenerator struct{}

func NewTemplateStoryGenerator() *TemplateStoryGenerator {
	return &TemplateStoryGenerator{}
}

func (g *TemplateStoryGenerator) Generate(_ context.Context, input domain.JourneyGenerationInput) (domain.JourneyStory, error) {
	j := input.Journey

	story := fmt.Sprintf("Я вернулся из путешествия в «%s»!", j.Location)
	if len(j.Events) > 0 {
		story += " Там произошло вот что: "
		for i, e := range j.Events {
			if i > 0 {
				story += ", "
			}
			story += e
		}
		story += "."
	}
	if j.Reward.Item != "" {
		story += fmt.Sprintf(" А ещё я нашёл кое-что интересное: %s.", j.Reward.Item)
	}

	return domain.JourneyStory{
		Title:  fmt.Sprintf("Возвращение из «%s»", j.Location),
		Story:  story,
		Teaser: "Кажется, рядом осталось ещё что-то неизведанное 👀",
	}, nil
}
