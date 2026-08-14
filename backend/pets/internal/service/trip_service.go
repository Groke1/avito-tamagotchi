package service

import (
	"context"
	"errors"
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
	TripDuration       = 30 * time.Second
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
	petService     *PetService
	eventNotifier  EventNotifier
	userClient     *clients.UserClient
	rewardPolicy   *domain.RewardPolicy
}

func NewTripService(
	storyGenerator TripStoryGeneratorInterface,
	executor Executor,
	tripRepository repository.TripRepositoryInterface,
	petRepository *repository.PetRepository,
	petService *PetService,
	userServiceURL string,
	eventNotifier EventNotifier,
) *TripService {
	return &TripService{
		storyGenerator: storyGenerator,
		executor:       executor,
		tripRepository: tripRepository,
		userClient:     clients.NewUserClient(userServiceURL + "/internal"),
		petRepository:  petRepository,
		petService:     petService,
		eventNotifier:  eventNotifier,
		rewardPolicy:   domain.NewRewardPolicy(),
	}
}

// Generate запускает фоновую генерацию истории и сразу возвращает управление
// хендлеру (для ответа 202 Accepted). Вся логика — в generateAndSave,
// выполняется внутри executor.Execute в отдельной горутине.
func (ts *TripService) Generate(ctx context.Context, dto TripDTO) error {
	IsNegativeInt := rand.IntN(2)
	coins := 100

	// проверить что нет активных путешествий
	petTripReturned, err := ts.tripRepository.GetLastTripByPetID(ctx, dto.PetID)
	switch {
	case errors.Is(err, domain.ErrTripNotFound):

	case err != nil:
		return err

	case petTripReturned != nil &&
		petTripReturned.EndedAt.After(time.Now()) &&
		petTripReturned.Status != domain.Delivered:
		return domain.ErrPetAlreadyTravelling
	}
	fmt.Println("[trip_service] активных путешествий нет")

	// спросить есть ли деньги
	err = ts.userClient.AdjustCoins(ctx, dto.UserID, coins)
	if err != nil {
		return err
	}
	fmt.Println("[trip_service] деньги есть")

	events, err := ts.tripRepository.GetTripEvents(ctx)
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
		Reward:   ts.rewardPolicy.GenerateReward(IsNegativeInt == 1),
	}

	lastPetTrips, err := ts.tripRepository.GetLastDeliveredTripsByPetID(ctx, dto.PetID, recentStoriesLimit)
	if err != nil {
		return fmt.Errorf("get pet memory: %w", err)
	}
	memory := make([]string, len(lastPetTrips))
	for i := 0; i < len(lastPetTrips); i++ {
		memory[i] = lastPetTrips[i].Story
	}
	input := domain.JourneyGenerationInput{
		Location: journey.Location,
		Events:   journey.Events,
		Memory:   memory,
	}

	story, err := ts.storyGenerator.Generate(ctx, input)
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
		EndedAt:     now.Add(TripDuration), // TODO: заменить на реальное время окончания путешествия
		Status:      domain.PendingDelivery,
	}
	err = ts.tripRepository.CreateTrip(ctx, tripFinal)
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

func (ts *TripService) SendToClient(trip *domain.PetTrip) {
	_ = ts.eventNotifier.SendToClient(trip.UserID, domain.EventTripCompleted, struct{}{})
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

	if trip.Status != domain.PendingDelivery {
		return nil, nil, domain.ErrNotPendingTrip
		// return trip, nil, nil // Пока возвращаю trip
	}

	_, err = ts.petService.GrantXP(ctx, userID, int(*trip.RewardXP))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to grant xp: %w", err)
	}

	err = ts.userClient.AdjustCoins(ctx, userID, -int(*trip.RewardCoins))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to withdraw coins: %w", err)
	}

	var reward *domain.Reward
	if trip.RewardCode != nil {
		reward, err = ts.userClient.ClaimReward(ctx, trip.UserID, *trip.RewardCode)
		if err != nil {
			log.Printf("[LAST TRIP] lost reward trip `%d`", trip.ID)
		}
	}

	if err = ts.tripRepository.MarkTripDelivered(ctx, trip.ID); err != nil {
		return nil, nil, fmt.Errorf("failed to mark trip delivered: %w", err)
	}

	return trip, reward, nil
}
