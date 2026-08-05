package service

import (
	"context"
	"time"

	repository "github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/repository"
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

type PetService struct {
	petRepository *repository.PetRepository
}

func NewPetService(petRepository *repository.PetRepository) *PetService {
	return &PetService{
		petRepository: petRepository,
	}
}

func (ps *PetService) GetPet(ctx context.Context, userID int64) (*Pet, error) {
	dbPet, err := ps.petRepository.GetPet(ctx, userID)
	if err != nil {
		return nil, err
	}

	pet := Pet{
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

	return &pet, err
}

func (ps *PetService) CreatePet(ctx context.Context, petName string, userID int64) (*Pet, error) {
	dbPet, err := ps.petRepository.CreatePet(ctx, petName, userID)
	if err != nil {
		return nil, err
	}

	pet := Pet{
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

	return &pet, err
}
