package service

import (
	"context"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/clients"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	repository "github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
)

type PetService struct {
	petRepository *repository.PetRepository
	client        *clients.UserClient
}

func NewPetService(petRepository *repository.PetRepository) *PetService {
	return &PetService{
		petRepository: petRepository,
		client:        clients.NewUserClient("http://localhost:8080/internal"),
	}
}

func (ps *PetService) GetPet(ctx context.Context, userID string) (*domain.Pet, error) {
	pet, err := ps.petRepository.GetPet(ctx, userID)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (ps *PetService) CreatePet(ctx context.Context, petName string, userID string) (*domain.Pet, error) {
	pet, err := ps.petRepository.CreatePet(ctx, petName, userID)
	if err != nil {
		return nil, err
	}

	return pet, nil
}

func (ps *PetService) FeedPet(ctx context.Context, userID string) (*domain.Pet, error) {
	err := ps.client.WithdrawCoins(ctx, userID, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw coins: %v", err)
	}

	tx, err := ps.petRepository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	defer tx.Rollback()

	pet, err := ps.petRepository.GetPetForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, domain.ErrPetNotFound
	}

	levelUp, err := pet.Feed()
	if err != nil {
		return nil, domain.ErrUnavailableAction
	}

	if err = ps.petRepository.UpdatePet(ctx, tx, pet); err != nil {
		return nil, domain.ErrUnavailableAction
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	if levelUp {
		// TODO сообщить о награде
	}

	return pet, nil
}

func (ps *PetService) StrokePet(ctx context.Context, userID string) (*domain.Pet, error) {
	err := ps.client.WithdrawCoins(ctx, userID, 7)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw coins: %v", err)
	}

	tx, err := ps.petRepository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	pet, err := ps.petRepository.GetPetForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, domain.ErrPetNotFound
	}

	levelUp, err := pet.Stroke()
	if err != nil {
		return nil, domain.ErrUnavailableAction
	}

	if err = ps.petRepository.UpdatePet(ctx, tx, pet); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	if levelUp {
		// TODO сообщить о награде
	}

	return pet, nil
}
