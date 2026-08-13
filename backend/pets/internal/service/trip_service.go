package service

import (
	"context"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
)

type TripService struct {
	tripRepository *repository.TripRepository
	petRepository  *repository.PetRepository
	eventNotifier  EventNotifier
}

func NewTripService(tripRepository *repository.TripRepository, petRepository *repository.PetRepository, eventNotifier EventNotifier) *TripService {
	return &TripService{tripRepository: tripRepository, petRepository: petRepository, eventNotifier: eventNotifier}
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

func (ts *TripService) GetLastTrip(ctx context.Context, userID string) (*domain.PetTrip, error) {
	pet, err := ts.petRepository.GetPet(ctx, userID)
	if err != nil {
		return nil, err
	}

	trip, err := ts.tripRepository.GetLastTripByPetID(ctx, pet.ID)
	if err != nil {
		return nil, err
	}

	if trip.Status == domain.PendingDelivery {
		return trip, nil
	}

	return nil, domain.ErrNotPendingTrip
}
