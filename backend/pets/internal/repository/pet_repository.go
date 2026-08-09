package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Pet struct {
	ID               int64     `db:"id"`
	UserID           string    `db:"user_id"`
	Name             string    `db:"name"`
	Level            int       `db:"level"`
	XP               int       `db:"xp"`
	NextLevelXP      int       `db:"next_level_xp"`
	Satiety          int       `db:"satiety"`
	Happiness        int       `db:"happiness"`
	CreatedAt        time.Time `db:"created_at"`
	LastCalculatedAt time.Time `db:"last_calculated_at"`
	LastFeedAt       time.Time `db:"last_feed_at"`
	LastStrokeAt     time.Time `db:"last_stroke_at"`
}

type PetForLeaderboard struct {
	UserID string `db:"user_id"`
	Name   string `db:"name"`
	Level  int    `db:"level"`
	Rank   int    `db:"rank"`
	XP     int    `db:"xp"`
}

type PetRepository struct {
	db *sqlx.DB
}

func NewPetRepository(db *sqlx.DB) *PetRepository {
	return &PetRepository{db: db}
}

func (pr *PetRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return pr.db.BeginTxx(ctx, nil)
}

func (pr *PetRepository) GetPet(ctx context.Context, userID string) (*domain.Pet, error) {
	query := `
				SELECT *
				FROM pets p
				WHERE p.user_id = $1
	`

	var dbPet Pet
	err := pr.db.GetContext(ctx, &dbPet, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPetNotFound
	}

	pet := domain.Pet{
		ID:               dbPet.ID,
		UserID:           dbPet.UserID,
		Name:             dbPet.Name,
		Level:            dbPet.Level,
		XP:               dbPet.XP,
		NextLevelXP:      dbPet.NextLevelXP,
		Satiety:          dbPet.Satiety,
		Happiness:        dbPet.Happiness,
		CreatedAt:        dbPet.CreatedAt,
		LastCalculatedAt: dbPet.LastCalculatedAt,
		LastFeedAt:       dbPet.LastFeedAt,
		LastStrokeAt:     dbPet.LastStrokeAt,
	}

	return &pet, nil
}

func (pr *PetRepository) CreatePet(ctx context.Context, userID string, petName string) (*domain.Pet, error) {
	query := `
				INSERT INTO pets (name, user_id)
				VALUES ($1, $2)
				RETURNING *
	`

	var dbPet Pet
	err := pr.db.GetContext(ctx, &dbPet, query, petName, userID)
	if err != nil {
		return nil, domain.ErrPetAlreadyExists
	}

	pet := domain.Pet{
		ID:               dbPet.ID,
		UserID:           dbPet.UserID,
		Name:             dbPet.Name,
		Level:            dbPet.Level,
		XP:               dbPet.XP,
		NextLevelXP:      dbPet.NextLevelXP,
		Satiety:          dbPet.Satiety,
		Happiness:        dbPet.Happiness,
		CreatedAt:        dbPet.CreatedAt,
		LastCalculatedAt: dbPet.LastCalculatedAt,
		LastFeedAt:       dbPet.LastFeedAt,
		LastStrokeAt:     dbPet.LastStrokeAt,
	}

	return &pet, nil
}

func (pr *PetRepository) GetPetForUpdate(ctx context.Context, tx *sqlx.Tx, userID string) (*domain.Pet, error) {
	query := `
				SELECT *
				FROM pets p
				WHERE p.user_id = $1
				FOR UPDATE
	`

	var dbPet Pet
	err := tx.GetContext(ctx, &dbPet, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPetNotFound
	} else if err != nil {
		print(fmt.Errorf("failed to get pet for update: %w", err))
		return nil, fmt.Errorf("failed to get pet for update: %w", err)
	}

	pet := domain.Pet{
		ID:               dbPet.ID,
		UserID:           dbPet.UserID,
		Name:             dbPet.Name,
		Level:            dbPet.Level,
		XP:               dbPet.XP,
		NextLevelXP:      dbPet.NextLevelXP,
		Satiety:          dbPet.Satiety,
		Happiness:        dbPet.Happiness,
		CreatedAt:        dbPet.CreatedAt,
		LastCalculatedAt: dbPet.LastCalculatedAt,
		LastFeedAt:       dbPet.LastFeedAt,
		LastStrokeAt:     dbPet.LastStrokeAt,
	}

	return &pet, nil
}

func (pr *PetRepository) UpdatePet(ctx context.Context, tx *sqlx.Tx, pet *domain.Pet) error {
	query := `
				UPDATE pets
				SET 
					satiety = $1, 
					happiness = $2, 
					xp = $3, 
					next_level_xp = $4, 
					level = $5,
					last_calculated_at = $6,
					last_feed_at = $7,
					last_stroke_at = $8
				WHERE id = $9
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		pet.Satiety,
		pet.Happiness,
		pet.XP,
		pet.NextLevelXP,
		pet.Level,
		pet.LastCalculatedAt,
		pet.LastFeedAt,
		pet.LastStrokeAt,
		pet.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update pet: %w", err)
	}

	if pet.LastGainedXP > 0 {
		query = `
				INSERT INTO pets_daily_xp (pet_id, date, gained_xp)
				VALUES ($1, (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date, $2)
				ON CONFLICT (pet_id, date)
				DO UPDATE SET gained_xp = pets_daily_xp.gained_xp + EXCLUDED.gained_xp
		`
		_, err := tx.ExecContext(ctx, query, pet.ID, pet.LastGainedXP)
		if err != nil {
			return fmt.Errorf("failed to upsert last gained xp: %w", err)
		}
	}

	return nil
}

func (pr *PetRepository) GetLeaderboardWithUser(ctx context.Context, userID string, limit int) ([]domain.LeaderboardItem, error) {
	query := `
				WITH ranked_pets AS (
					SELECT  
						xp
						name, 
						user_id, 
						level,
						DENSE_RANK() OVER (
							ORDER BY 
								level DESC, 
								xp DESC, 
								id ASC
						) as rank
					FROM pets
				)
				SELECT *
				FROM ranked_pets
				WHERE rank <= $1
				OR user_id = $2
				ORDER BY rank
	`

	var dbPets []PetForLeaderboard
	err := pr.db.SelectContext(ctx, &dbPets, query, limit, userID)
	if err != nil {
		return nil, err
	}

	pets := make([]domain.LeaderboardItem, len(dbPets))
	for i, dbPet := range dbPets {
		pets[i] = domain.LeaderboardItem{
			Rank:    dbPet.Rank,
			Level:   dbPet.Level,
			UserID:  dbPet.UserID,
			PetName: dbPet.Name,
			XP:      dbPet.XP,
		}
	}

	return pets, nil
}

func (pr *PetRepository) GetDailyGainedXPByUserID(ctx context.Context, userID string) (int, error) {
	query := `
				SELECT COALESCE(pdx.gained_xp, 0)
				FROM pets p
				LEFT JOIN pets_daily_xp pdx
				ON p.id = pdx.pet_id AND pdx.date = (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date
				WHERE p.user_id = $1
	`
	var gainedXP int
	err := pr.db.GetContext(ctx, &gainedXP, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.ErrPetNotFound
		}
		return 0, err
	}

	return gainedXP, nil
}

func (pr *PetRepository) GetWeeklyLeaderboardWithUser(ctx context.Context, userID string, limit int) ([]domain.LeaderboardItem, error) {
	query := `
				WITH weekly_xp AS (
					SELECT
						pdx.pet_id,
						SUM(pdx.gained_xp) AS weekly_xp
					FROM pets_daily_xp pdx
					WHERE pdx.date >= (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date - INTERVAL '6 days'
					GROUP BY pdx.pet_id
				),
				ranked_pets AS (
					SELECT
						p.name,
						p.user_id,
						p.level,
						COALESCE(wxp.weekly_xp, 0) AS xp,
						DENSE_RANK() OVER (
							ORDER BY
								p.level DESC,
								COALESCE(wxp.weekly_xp, 0) DESC,
								p.id ASC
						) AS rank
					FROM pets p
					LEFT JOIN weekly_xp wxp ON p.id = wxp.pet_id
				)
				SELECT *
				FROM ranked_pets
				WHERE rank <= $1
				OR user_id = $2
				ORDER BY rank;
	`

	var dbPets []PetForLeaderboard
	err := pr.db.SelectContext(ctx, &dbPets, query, limit, userID)
	if err != nil {
		return nil, err
	}

	pets := make([]domain.LeaderboardItem, len(dbPets))
	for i, dbPet := range dbPets {
		pets[i] = domain.LeaderboardItem{
			Rank:    dbPet.Rank,
			Level:   dbPet.Level,
			UserID:  dbPet.UserID,
			PetName: dbPet.Name,
			XP:      dbPet.XP,
		}
	}

	return pets, nil
}
