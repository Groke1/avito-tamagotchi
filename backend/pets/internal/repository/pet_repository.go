package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Pet struct {
	ID          int64     `db:"id"`
	UserID      string    `db:"user_id"`
	Name        string    `db:"name"`
	Level       int       `db:"level"`
	XP          int       `db:"xp"`
	NextLevelXP int       `db:"next_level_xp"`
	Satiety     int       `db:"satiety"`
	Happiness   int       `db:"happiness"`
	CreatedAt   time.Time `db:"created_at"`
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
		ID:          dbPet.ID,
		UserID:      dbPet.UserID,
		Name:        dbPet.Name,
		Level:       dbPet.Level,
		XP:          dbPet.XP,
		NextLevelXP: dbPet.NextLevelXP,
		Satiety:     dbPet.Satiety,
		Happiness:   dbPet.Happiness,
		CreatedAt:   dbPet.CreatedAt,
	}

	return &pet, nil
}

func (pr *PetRepository) CreatePet(ctx context.Context, petName string, userID string) (*domain.Pet, error) {
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
		ID:          dbPet.ID,
		UserID:      dbPet.UserID,
		Name:        dbPet.Name,
		Level:       dbPet.Level,
		XP:          dbPet.XP,
		NextLevelXP: dbPet.NextLevelXP,
		Satiety:     dbPet.Satiety,
		Happiness:   dbPet.Happiness,
		CreatedAt:   dbPet.CreatedAt,
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
	}

	pet := domain.Pet{
		ID:          dbPet.ID,
		UserID:      dbPet.UserID,
		Name:        dbPet.Name,
		Level:       dbPet.Level,
		XP:          dbPet.XP,
		NextLevelXP: dbPet.NextLevelXP,
		Satiety:     dbPet.Satiety,
		Happiness:   dbPet.Happiness,
		CreatedAt:   dbPet.CreatedAt,
	}

	return &pet, nil
}

func (pr *PetRepository) UpdatePet(ctx context.Context, tx *sqlx.Tx, pet *domain.Pet) error {
	query := `
				UPDATE pets
				SET satiety = $1, happiness = $2, xp = $3, next_level_xp = $4, level = $5
				WHERE id = $6
	`

	_, err := tx.ExecContext(ctx, query, pet.Satiety, pet.Happiness, pet.XP, pet.NextLevelXP, pet.Level, pet.ID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return fmt.Errorf("failed to update pet: %v", err)
	}

	return nil
}
