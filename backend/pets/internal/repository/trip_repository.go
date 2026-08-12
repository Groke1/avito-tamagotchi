package repository

import (
	"context"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TripRepository struct {
	db *sqlx.DB
}

func NewTripRepository(db *sqlx.DB) *TripRepository {
	return &TripRepository{db: db}
}

func (tr *TripRepository) GetExpiredTrips(ctx context.Context) ([]*domain.Trip, error) {

	return []*domain.Trip{}, nil
}

func (tr *TripRepository) GetTrip(ctx context.Context, userID string) (*domain.Trip, error) {
	return nil, nil
}
