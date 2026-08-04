package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Pet struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Level       int       `json:"level" db:"level"`
	XP          int       `json:"xp" db:"xp"`
	NextLevelXP int       `json:"next_level_xp" db:"next_level_xp"`
	Satiety     int       `json:"satiety" db:"satiety"`
	Happiness   int       `json:"happiness" db:"happiness"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type PetRepository struct {
	db *sqlx.DB
}

func NewPetRepository(db *sqlx.DB) *PetRepository {
	return &PetRepository{db: db}
}

func (pr *PetRepository) GetPet(ctx context.Context, userID int64) (*Pet, error) {
	query := `
				SELECT *
				FROM pets p
				WHERE p.user_id = $1
	`

	var pet Pet
	err := pr.db.GetContext(ctx, &pet, query, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrPetNotFound
	}

	return &pet, err
}

func (pr *PetRepository) CreatePet(ctx context.Context, petName string, userID int64) (*Pet, error) {
	query := `
				INSERT INTO pets (name, user_id)
				VALUES ($1, $2)
				RETURNING *
	`

	var pet Pet
	err := pr.db.GetContext(ctx, &pet, query, petName, userID)
	if err != nil {
		return nil, domain.ErrPetAlreadyExists
	}

	return &pet, err
}
