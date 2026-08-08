package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/clients"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
)

var (
	XPperStreak = 7
)

type EventNotifier interface {
	SendToClient(string, string, any)
	BroadcastLeaderboard()
}

type PetService struct {
	petRepository *repository.PetRepository
	client        *clients.UserClient
	eventNotifier EventNotifier
}

func NewPetService(petRepository *repository.PetRepository, userServiceURL string, eventNotifier EventNotifier) *PetService {
	return &PetService{
		petRepository: petRepository,
		client:        clients.NewUserClient(userServiceURL + "/internal"),
		eventNotifier: eventNotifier,
	}
}

func (ps *PetService) GetPet(ctx context.Context, userID string) (*domain.Pet, error) {
	pet, err := ps.petRepository.GetPet(ctx, userID)
	if err != nil {
		return nil, domain.ErrPetNotFound
	}

	changed := pet.RecalculateState(time.Now())
	if changed {
		updatedPet, _, err := ps.RecalculateState(ctx, userID)
		if err != nil {
			return nil, err
		}

		return updatedPet, nil
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
	tx, err := ps.petRepository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	pet, err := ps.petRepository.GetPetForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, domain.ErrPetNotFound
	}

	pet.RecalculateState(time.Now())
	levelUp, cost, err := pet.Feed()
	if err != nil {
		return nil, err
	}

	if err = ps.petRepository.UpdatePet(ctx, tx, pet); err != nil {
		return nil, domain.ErrUnavailableAction
	}

	err = ps.client.WithdrawCoins(ctx, userID, cost)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw coins: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if levelUp {
		// TODO сообщить о награде
	}

	ps.eventNotifier.BroadcastLeaderboard()

	return pet, nil
}

func (ps *PetService) StrokePet(ctx context.Context, userID string) (*domain.Pet, error) {
	tx, err := ps.petRepository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	pet, err := ps.petRepository.GetPetForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, domain.ErrPetNotFound
	}

	pet.RecalculateState(time.Now())
	levelUp, err := pet.Stroke()
	if err != nil {
		return nil, err
	}

	if err = ps.petRepository.UpdatePet(ctx, tx, pet); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	if levelUp {
		// TODO сообщить о награде
	}

	ps.eventNotifier.BroadcastLeaderboard()

	return pet, nil
}

func (ps *PetService) GetLeaderboard(ctx context.Context, limit int, userID string) ([]domain.LeaderboardItem, *domain.LeaderboardItem, error) {
	records, err := ps.petRepository.GetLeaderboardWithUser(ctx, limit, userID)
	if err != nil {
		return nil, nil, err
	}

	var currentUserItem domain.LeaderboardItem
	var userIDs = make([]string, len(records))
	for i := range records {
		userIDs[i] = records[i].UserID
		if records[i].UserID == userID {
			currentUserItem = records[i]
		}
	}

	userMap, err := ps.client.GetUsernamesByIDs(ctx, userIDs)
	if err != nil {
		return nil, nil, err
	}

	for i := range records {
		if name, ok := userMap[records[i].UserID]; ok {
			records[i].UserName = name
		} else {
			return nil, nil, err
		}
	}
	currentUserItem.UserName = userMap[currentUserItem.UserID]

	return records, &currentUserItem, nil
}

func (ps *PetService) GrantXP(ctx context.Context, amount int, userID string) (*domain.Pet, error) {
	tx, err := ps.petRepository.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	pet, err := ps.petRepository.GetPetForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, domain.ErrPetNotFound
	}

	pet.RecalculateState(time.Now())
	levelUp := pet.AddXP(amount)

	if err = ps.petRepository.UpdatePet(ctx, tx, pet); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	ps.eventNotifier.SendToClient(userID, "pet.updated", pet)

	if levelUp {
		// TODO сообщить о награде
	}

	ps.eventNotifier.BroadcastLeaderboard()

	return pet, nil
}

func (ps *PetService) ClaimDailyBonus(ctx context.Context, streak int, userID string) error {
	amount := XPperStreak * streak
	_, err := ps.GrantXP(ctx, amount, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PetService) RecalculateState(ctx context.Context, userID string) (*domain.Pet, bool, error) {
	tx, err := ps.petRepository.BeginTx(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	pet, err := ps.petRepository.GetPetForUpdate(ctx, tx, userID)
	if err != nil {
		return nil, false, domain.ErrPetNotFound
	}

	changed := pet.RecalculateState(time.Now())
	if changed {
		err = ps.petRepository.UpdatePet(ctx, tx, pet)
		if err != nil {
			return nil, false, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("[PET SERVICE] Pet retrieved. Changed=%t for userID: '%s'", changed, userID)

	return pet, changed, nil
}

func (ps *PetService) ClaimDailyGainedXP(ctx context.Context, userID string) (int, error) {
	gainedXP, err := ps.petRepository.GetDailyGainedXPByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return gainedXP, nil
}
