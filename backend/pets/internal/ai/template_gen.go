package ai

import (
	"context"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
)

type TemplateStoryGenerator struct{}

func NewTemplateStoryGenerator() *TemplateStoryGenerator {
	return &TemplateStoryGenerator{}
}

func (t *TemplateStoryGenerator) Generate(_ context.Context, journey domain.JourneyResult) (domain.JourneyStory, error) {
	story := "Я вернулся из путешествия! В этот раз ничего особенного не произошло"

	return domain.JourneyStory{
		Title:  fmt.Sprintf("Путешествие в %s", journey.Location),
		Text:   story,
		Teaser: "Кажется, рядом есть ещё неизведанные места",
	}, nil
}
