package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TripRepository struct {
	db *sqlx.DB
}

func NewTripRepository(db *sqlx.DB) *TripRepository {
	return &TripRepository{db: db}
}

func (r *TripRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *TripRepository) GetTripEvents(ctx context.Context) ([]domain.TripEvent, error) {
	query := `
		SELECT id, description, is_negative
		FROM trip_events;
	`

	var dbTripEvents []tripEvent
	if err := r.db.SelectContext(ctx, &dbTripEvents, query); err != nil {
		return nil, fmt.Errorf("failed to get trip events: %w", err)
	}
	tripEvents := make([]domain.TripEvent, len(dbTripEvents))
	for i, event := range dbTripEvents {
		tripEvents[i] = domain.TripEvent{
			ID:          event.ID,
			Description: event.Description,
			IsNegative:  event.IsNegative,
		}
	}
	return tripEvents, nil
}

// Создать новое путешествие для питомца

func (r *TripRepository) CreateTrip(ctx context.Context, trip domain.PetTrip) error {
	const query = `
		INSERT INTO pet_trips (pet_id, user_id, location,
			reward_xp, reward_coins, reward_code, 
			story, ended_at 
		    
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	if _, err := r.db.ExecContext(ctx, query, trip.PetID, trip.UserID, trip.Location,
		trip.RewardXP, trip.RewardCoins, trip.RewardCode, trip.Story, trip.EndedAt,
	); err != nil {
		return fmt.Errorf("failed to insert pet trip: %w", err)
	}

	return nil
}

// Получить `limit` последних доставленных путешествий питомца, чтобы передать их истории в нейронку
func (r *TripRepository) GetLastDeliveredTripsByPetID(ctx context.Context, petID int64, limit int) ([]domain.PetTrip, error) {
	const query = `
		SELECT * FROM pet_trips
		WHERE pet_id = $1 AND status = 'delivered'
		ORDER BY ended_at DESC
		LIMIT $2
	`

	var dbPetTrips []petTrip
	if err := r.db.SelectContext(ctx, &dbPetTrips, query, petID, limit); err != nil {
		return nil, fmt.Errorf("failed to get last delivered trips: %w", err)
	}

	petTrips := make([]domain.PetTrip, len(dbPetTrips))
	for i, trip := range dbPetTrips {
		petTrips[i] = *convertToPetTrip(&trip)
	}

	return petTrips, nil
}

// Получить последнее путешествие для питомца. Например, чтобы проверить 'in-progres' ли оно, чтобы запретить действия (кормление, поглаживание) над питомцем
// Если питомец еще ни разу не отправлялся в путешествия - верну ошибку ErrTripNotFound
func (r *TripRepository) GetLastTripByPetID(ctx context.Context, petID int64) (*domain.PetTrip, error) {
	const query = `
		SELECT * FROM pet_trips
		WHERE pet_id = $1 ORDER BY started_at DESC
		LIMIT 1
	`

	var trip petTrip
	if err := r.db.GetContext(ctx, &trip, query, petID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrTripNotFound
		}

		return nil, fmt.Errorf("failed to get trip: %w", err)
	}

	return convertToPetTrip(&trip), nil
}

// Получить список завершенных трипов (для воркера-2 - затем он выставляет всем статусы либо delivered, либо pending_delivery)
func (r *TripRepository) GetFinishedTrips(ctx context.Context) ([]domain.PetTrip, error) {
	const query = `
		SELECT * FROM pet_trips
		WHERE ended_at <= CURRENT_TIMESTAMP
		  AND status IN 'in_progress'
		ORDER BY ended_at
	`

	var dbTrips []petTrip
	if err := r.db.SelectContext(ctx, &dbTrips, query); err != nil {
		return nil, err
	}

	trips := make([]domain.PetTrip, len(dbTrips))
	for i, trip := range dbTrips {
		trips[i] = *convertToPetTrip(&trip)
	}

	return trips, nil
}

// Пометить, что путешествие не получилось доставить на фронт по вебсокету
func (r *TripRepository) MarkTripPendingDelivery(ctx context.Context, tripID int64) error {
	const query = `UPDATE pet_trips
		SET status = 'pending_delivery' WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, tripID)
	if err != nil {
		return fmt.Errorf("failed to mark trip %d as pending_delivery: %w", tripID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrTripNotFound
	}

	return nil
}

// Пометить, что путешествие доставлено на фронт
func (r *TripRepository) MarkTripDelivered(ctx context.Context, tripID int64) error {
	const query = `
		UPDATE pet_trips
		SET status = 'delivered'
		WHERE id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, tripID)
	if err != nil {
		return fmt.Errorf("failed to mark trip %d as delivered: %w", tripID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrTripNotFound
	}

	return nil
}

func convertToPetTrip(dbPetTrip *petTrip) *domain.PetTrip {
	return &domain.PetTrip{
		ID:          dbPetTrip.ID,
		PetID:       dbPetTrip.PetID,
		UserID:      dbPetTrip.UserID,
		Location:    dbPetTrip.Location,
		RewardXP:    dbPetTrip.RewardXP,
		RewardCoins: dbPetTrip.RewardCoins,
		RewardCode:  dbPetTrip.RewardCode,
		Story:       dbPetTrip.Story,
		Status:      domain.TripStatus(dbPetTrip.Status),
		StartedAt:   dbPetTrip.StartedAt,
		EndedAt:     dbPetTrip.EndedAt,
	}
}
