package service

import (
	"context"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
)

type TripService struct {
	tripRepository *repository.TripRepository
	eventNotifier  EventNotifier
}

func NewTripService(tripRepository *repository.TripRepository, eventNotifier EventNotifier) *TripService {
	return &TripService{tripRepository: tripRepository, eventNotifier: eventNotifier}
}

func (ts *TripService) CompleteTrip(ctx context.Context, trip *domain.Trip) error {
	return nil
}
