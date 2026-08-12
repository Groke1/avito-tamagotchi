package service

import (
	"context"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
)

// JourneyStoryGenerator — абстракция над конкретной нейросетью, чтобы
// бизнес-логика Journey Worker не зависела напрямую от GigaChat, OpenAI
// и т.д. Реализации: gigachat.Client (LLM) и TemplateStoryGenerator (fallback).
type JourneyStoryGenerator interface {
	Generate(ctx context.Context, input domain.JourneyGenerationInput) (domain.JourneyStory, error)
}

// FallbackStoryGenerator оборачивает основной генератор (LLM) и подстраховывает
// его шаблонным вариантом. Путешествие не должно зависеть от того, жив ли
// внешний AI-провайдер: пользователь в любом случае получит историю.
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

// TemplateStoryGenerator — примитивный, но всегда доступный генератор.
// Не ходит в сеть, собирает историю из шаблонных фраз на основе
// JourneyResult. Используется, когда LLM недоступна.
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
