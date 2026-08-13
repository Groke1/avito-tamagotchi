package service

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
)

const (
	defaultLocation    = "AvitoLand"
	recentStoriesLimit = 3
)

type TripDTO struct {
	PetID  int64  `json:"pet_id"`
	UserID string `json:"user_id"`
}

type Executor interface {
	Execute(job func(context.Context) error) error
}

type TripStoryGeneratorInterface interface {
	Generate(
		ctx context.Context,
		input domain.JourneyGenerationInput,
	) (domain.JourneyStory, error)
}

type TripService struct {
	storyGenerator TripStoryGeneratorInterface
	executor       Executor
	repository     repository.TripRepositoryInterface
}

func NewTripService(
	storyGenerator TripStoryGeneratorInterface,
	executor Executor,
	repository repository.TripRepositoryInterface,
) *TripService {
	return &TripService{
		storyGenerator: storyGenerator,
		executor:       executor,
		repository:     repository,
	}
}

// Generate запускает фоновую генерацию истории и сразу возвращает управление
// хендлеру (для ответа 202 Accepted). Вся логика — в generateAndSave,
// выполняется внутри executor.Execute в отдельной горутине.
func (s *TripService) Generate(ctx context.Context, dto TripDTO) error {
	return s.executor.Execute(func(ctx context.Context) error {
		if err := s.generateAndSave(ctx, dto); err != nil {
			log.Printf("failed to generate trip for user %s: %v", dto.UserID, err)
		}
		return nil
	})
}

func (s *TripService) generateAndSave(ctx context.Context, dto TripDTO) error {
	IsNegativeInt := rand.IntN(2)

	events, err := s.repository.GetTripEvents(ctx)
	shuffleEvents := getRandomDescriptions(&events, IsNegativeInt == 1, recentStoriesLimit)
	if err != nil {
		return fmt.Errorf("get recent stories: %w", err)
	}
	journey := domain.JourneyResult{
		Location: defaultLocation,
		Events:   shuffleEvents,
		Reward: domain.JourneyReward{
			RewardXP:    int32(IsNegativeInt * rand.IntN(100)),
			RewardCoins: int32(IsNegativeInt*rand.IntN(50) + 30),
		},
	}

	memory, err := s.repository.GetLastDeliveredTripsByPetID(ctx, dto.PetID, recentStoriesLimit)
	if err != nil {
		return fmt.Errorf("get pet memory: %w", err)
	}

	input := domain.JourneyGenerationInput{
		Journey: journey,
		Memory:  memory,
	}

	story, err := s.storyGenerator.Generate(ctx, input)
	if err != nil {
		return fmt.Errorf("generate story: %w", err)
	}

	tripFinal := domain.PetTrip{
		PetID:       dto.PetID,
		UserID:      dto.UserID,
		Location:    journey.Location,
		RewardXP:    &journey.Reward.RewardXP,
		RewardCoins: &journey.Reward.RewardCoins,
		Story:       story.Story,
	}
	if err := s.repository.CreateTrip(ctx, tripFinal); err != nil {
		return fmt.Errorf("save story: %w", err)
	}

	return nil
}

func getRandomDescriptions(events *[]domain.TripEvent, neg bool, count int) []string {
	if events == nil {
		return nil
	}
	var indices []int
	for i, event := range *events {
		if event.IsNegative == neg {
			indices = append(indices, i)
		}
	}
	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})
	if len(indices) < count {
		count = len(indices)
	}
	result := make([]string, count)
	for i := 0; i < count; i++ {
		result[i] = (*events)[indices[i]].Description
	}

	return result
}
