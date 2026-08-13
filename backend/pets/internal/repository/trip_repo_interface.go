package repository

import (
	"context"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TripRepositoryInterface interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	GetTripEvents(ctx context.Context) ([]domain.TripEvent, error)
	CreateTrip(ctx context.Context, trip domain.PetTrip) error
	GetLastDeliveredTripsByPetID(ctx context.Context, petID int64, limit int) ([]domain.PetTrip, error)
	GetLastTripByPetID(ctx context.Context, petID int64) (*domain.PetTrip, error)
	GetFinishedTrips(ctx context.Context) ([]domain.PetTrip, error)
	MarkTripPendingDelivery(ctx context.Context, tripID int64) error
	MarkTripDelivered(ctx context.Context, tripID int64) error
}
