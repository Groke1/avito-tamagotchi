package service

import (
	"context"

	"github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/domain"
	repository "github.com/cayman444/avito-gamification-hackathon/blob/main/backend/pets/internal/repository"
)

type PetService struct {
	petRepository *repository.PetRepository
}

func NewPetService(petRepository *repository.PetRepository) *PetService {
	return &PetService{
		petRepository: petRepository,
	}
}

func (ps *PetService) GetPet(ctx context.Context, userID int64) (*domain.Pet, error) {
	pet, err := ps.petRepository.GetPet(ctx, userID)
	if err != nil {
		return nil, err
	}

	return pet, err
}

func (ps *PetService) CreatePet(ctx context.Context, petName string, userID int64) (*domain.Pet, error) {
	pet, err := ps.petRepository.CreatePet(ctx, petName, userID)
	if err != nil {
		return nil, err
	}

	return pet, err
}

func (ps *PetService) FeedPet(ctx context.Context, userID int64) (*domain.Pet, error) {

}

func (ps *PetService) StrokePet(ctx context.Context, userID int64) (*domain.Pet, error) {

}
