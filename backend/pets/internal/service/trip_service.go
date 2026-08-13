package service

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/clients"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
)

const (
	defaultLocation    = "Авито-Сити"
	recentStoriesLimit = 1
	eventsLimit        = 2
)

type TripDTO struct {
	PetID  int64  `json:"pet_id"`
	UserID string `json:"user_id"`
	Coins  int32  `json:"coins"`
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
	tripRepository repository.TripRepositoryInterface
	petRepository  *repository.PetRepository
	eventNotifier  EventNotifier
	userClient     *clients.UserClient
	rewardPolicy   *domain.RewardPolicy
}

func NewTripService(
	storyGenerator TripStoryGeneratorInterface,
	executor Executor,
	tripRepository repository.TripRepositoryInterface,
	petRepository *repository.PetRepository,
	userServiceURL string,
	eventNotifier EventNotifier,
) *TripService {
	return &TripService{
		storyGenerator: storyGenerator,
		executor:       executor,
		tripRepository: tripRepository,
		userClient:     clients.NewUserClient(userServiceURL + "/internal"),
		petRepository:  petRepository,
		eventNotifier:  eventNotifier,
		rewardPolicy:   domain.NewRewardPolicy(),
	}
}

// Generate запускает фоновую генерацию истории и сразу возвращает управление
// хендлеру (для ответа 202 Accepted). Вся логика — в generateAndSave,
// выполняется внутри executor.Execute в отдельной горутине.
func (s *TripService) Generate(ctx context.Context, dto TripDTO) error {
	IsNegativeInt := rand.IntN(2)
	coins := 100

	// проверить что нет активных путешествий
	petTripReturned, err := s.tripRepository.GetLastTripByPetID(ctx, dto.PetID)
	if err != nil {
		return err
	}
	if petTripReturned.EndedAt.After(time.Now()) && petTripReturned.Status != domain.Delivered {
		fmt.Println("ErrPetAlreadyTravelling")
		return domain.ErrPetAlreadyTravelling
	}
	fmt.Println("[trip_service] активных путешествий нет")

	// спросить есть ли деньги
	err = s.userClient.WithdrawCoins(ctx, dto.UserID, coins)
	if err != nil {
		return err
	}
	fmt.Println("[trip_service] деньги есть")

	events, err := s.tripRepository.GetTripEvents(ctx)
	if err != nil {
		fmt.Println("91: " + err.Error())
		return err
	}
	fmt.Println("events: " + strconv.Itoa(len(events)))
	if len(events) == 0 {
		fmt.Println("no events")
		return domain.ErrTripEventsNotExist
	}

	shuffleEvents := getRandomDescriptions(&events, IsNegativeInt == 0, eventsLimit)
	fmt.Println("events: " + strconv.Itoa(len(shuffleEvents)))

	journey := domain.JourneyResult{
		Location: defaultLocation,
		Events:   shuffleEvents,
		Reward: domain.JourneyReward{
			RewardXP:    int32(IsNegativeInt * rand.IntN(100)),
			RewardCoins: int32(IsNegativeInt*rand.IntN(50) + 30),
		},
	}

	lastPetTrips, err := s.tripRepository.GetLastDeliveredTripsByPetID(ctx, dto.PetID, recentStoriesLimit)
	if err != nil {
		return fmt.Errorf("get pet memory: %w", err)
	}
	memory := make([]string, len(lastPetTrips))
	for i := 0; i < len(lastPetTrips); i++ {
		memory[i] = lastPetTrips[i].Story
	}
	input := domain.JourneyGenerationInput{
		Journey: journey,
		Memory:  memory,
	}

	story, err := s.storyGenerator.Generate(ctx, input)
	if err != nil {
		return fmt.Errorf("generate story: %w", err)
	}
	fmt.Println(story.Title + " :: " + story.Story + " ::in service")

	now := time.Now().UTC()
	tripFinal := domain.PetTrip{
		PetID:       dto.PetID,
		UserID:      dto.UserID,
		Location:    journey.Location,
		RewardXP:    &journey.Reward.RewardXP,
		RewardCoins: &journey.Reward.RewardCoins,
		RewardCode:  journey.Reward.RewardCode,
		Story:       story.Story,
		StartedAt:   now,
		EndedAt:     now.Add(60 * time.Second), // TODO: заменить на реальное время окончания путешествия
		Status:      domain.PendingDelivery,
	}
	err = s.tripRepository.CreateTrip(ctx, tripFinal)
	if err != nil {
		return fmt.Errorf("save story: %w", err)
	}
	fmt.Println("[trip_service] trip created")
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

func (ts *TripService) CompleteTrip(ctx context.Context, trip *domain.PetTrip) (*domain.PetTrip, error) {
	ok := ts.eventNotifier.SendToClient(trip.UserID, domain.EventTripCompleted, struct{}{})
	if !ok {
		ts.tripRepository.MarkTripPendingDelivery(ctx, trip.ID)
		return trip, nil
	}

	ts.tripRepository.MarkTripDelivered(ctx, trip.ID)

	return trip, nil
}

func (ts *TripService) GetLastTrip(ctx context.Context, userID string) (*domain.PetTrip, *domain.Reward, error) {
	pet, err := ts.petRepository.GetPet(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	trip, err := ts.tripRepository.GetLastTripByPetID(ctx, pet.ID)
	if err != nil {
		return nil, nil, err
	}

	if trip.Status == domain.PendingDelivery {
		var reward *domain.Reward
		if trip.RewardCode != nil {
			reward, err = ts.userClient.ClaimReward(ctx, trip.UserID, *trip.RewardCode)
			if err != nil {
				log.Printf("[LAST TRIP] lost reward trip `%d`", trip.ID)
			}
		}

		return trip, reward, nil
	}

	return nil, nil, domain.ErrNotPendingTrip
}
