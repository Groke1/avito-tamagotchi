package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/clients"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/repository"
)

const (
	XPperStreak = 7
)

type EventNotifier interface {
	SendToClient(userID string, eventType domain.WsEvent, v any) bool
	BroadcastLeaderboard()
}

type PetService struct {
	petRepository  *repository.PetRepository
	tripRepository *repository.TripRepository
	client         *clients.UserClient
	eventNotifier  EventNotifier
	levelPolicy    *domain.LevelPolicy
}

func NewPetService(
	petRepository *repository.PetRepository,
	tripRepository *repository.TripRepository,
	userServiceURL string,
	eventNotifier EventNotifier) *PetService {
	return &PetService{
		petRepository:  petRepository,
		tripRepository: tripRepository,
		client:         clients.NewUserClient(userServiceURL + "/internal"),
		eventNotifier:  eventNotifier,
		levelPolicy:    domain.NewLevelPolicy(),
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

func (ps *PetService) CreatePet(ctx context.Context, userID string, petName string) (*domain.Pet, error) {
	pet, err := ps.petRepository.CreatePet(ctx, userID, petName)
	if err != nil {
		return nil, err
	}
	ps.eventNotifier.BroadcastLeaderboard()

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

	now := time.Now().UTC()
	trip, err := ps.tripRepository.GetLastTripByPetID(ctx, pet.ID)
	if err != nil && !errors.Is(err, domain.ErrTripNotFound) {
		return nil, err
	} else if !errors.Is(err, domain.ErrTripNotFound) && trip.Status == domain.InProgress && now.Before(trip.EndedAt) {
		return nil, &domain.ActionUnavailableError{
			RetryAfter: trip.EndedAt.Sub(time.Now().UTC()),
		}
	}

	pet.RecalculateState(time.Now())
	levelUp, cost, err := pet.Feed()
	if err != nil {
		return nil, err
	}

	if err = ps.petRepository.UpdatePet(ctx, tx, pet); err != nil {
		return nil, domain.ErrUnavailableAction
	}

	err = ps.client.AdjustCoins(ctx, userID, cost)
	if err != nil {
		return nil, fmt.Errorf("failed to withdraw coins: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	for _, level := range levelUp {
		if err = ps.levelUp(ctx, userID, level); err != nil {
			return nil, err
		}
	}

	if err = ps.client.NotifyActionDone(ctx, userID, time.Now().UTC()); err != nil {
		log.Printf("Action notification failed of user '%s' with %v", userID, err)
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

	now := time.Now().UTC()
	trip, err := ps.tripRepository.GetLastTripByPetID(ctx, pet.ID)
	if err != nil && !errors.Is(err, domain.ErrTripNotFound) {
		return nil, err
	} else if !errors.Is(err, domain.ErrTripNotFound) && trip.Status == domain.InProgress && now.Before(trip.EndedAt) {
		return nil, &domain.ActionUnavailableError{
			RetryAfter: trip.EndedAt.Sub(time.Now().UTC()),
		}
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

	for _, level := range levelUp {
		if err = ps.levelUp(ctx, userID, level); err != nil {
			return nil, err
		}
	}

	if err = ps.client.NotifyActionDone(ctx, userID, time.Now().UTC()); err != nil {
		log.Printf("Action notification failed of user '%s' with %v", userID, err)
	}
	ps.eventNotifier.BroadcastLeaderboard()

	return pet, nil
}

func (ps *PetService) GetLeaderboard(ctx context.Context, userID string, limit int) ([]domain.LeaderboardItem, *domain.LeaderboardItem, error) {
	records, err := ps.petRepository.GetLeaderboardWithUser(ctx, userID, limit)
	if err != nil {
		return nil, nil, err
	}

	if len(records) == 0 {
		return []domain.LeaderboardItem{}, nil, nil
	}

	var userIDs = make([]string, len(records))
	for i := range records {
		userIDs[i] = records[i].UserID
	}

	userMap, err := ps.client.GetUsernamesByIDs(ctx, userIDs)
	if err != nil {
		return nil, nil, err
	}

	var currentUserItem *domain.LeaderboardItem
	for i := range records {
		name, ok := userMap[records[i].UserID]
		if !ok {
			continue
		}

		records[i].UserName = name

		if records[i].UserID == userID {
			currentUserItem = &records[i]
		}
	}

	topRecords := append([]domain.LeaderboardItem{}, records[:min(limit, len(records))]...)
	currentUserItem.UserName = userMap[currentUserItem.UserID]

	return topRecords, currentUserItem, nil
}

func (ps *PetService) GrantXP(ctx context.Context, userID string, amount int) (*domain.Pet, error) {
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

	_ = ps.eventNotifier.SendToClient(userID, domain.EventPetUpdated, pet)

	for _, level := range levelUp {
		if err = ps.levelUp(ctx, userID, level); err != nil {
			return nil, err
		}
	}

	ps.eventNotifier.BroadcastLeaderboard()

	return pet, nil
}

func (ps *PetService) ClaimDailyBonusForStreak(ctx context.Context, userID string, streak int, coins int) error {
	amount := XPperStreak * streak
	_, err := ps.GrantXP(ctx, userID, amount)
	if err != nil {
		return err
	}

	payload := map[string]int{
		"coins": coins,
		"xp":    amount,
	}

	_ = ps.eventNotifier.SendToClient(userID, domain.EventStreakReward, payload)

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

func (ps *PetService) levelUp(ctx context.Context, userID string, level int) error {
	code := ps.levelPolicy.GetCode(level)
	reward, err := ps.client.ClaimReward(ctx, userID, code)
	if err != nil {
		return err
	}

	wsReward := struct {
		ID           string     `json:"id"`
		PromoCode    string     `json:"promo_code"`
		Name         string     `json:"name"`
		Description  string     `json:"description"`
		Status       string     `json:"status"`
		ExpiresAt    string     `json:"expires_at"`
		EarnedReason string     `json:"earned_reason"`
		RedeemedAt   *time.Time `json:"redeemed_at"`
	}{
		ID:           reward.ID,
		PromoCode:    reward.PromoCode,
		Name:         reward.Name,
		Description:  reward.Description,
		Status:       reward.Status,
		ExpiresAt:    reward.ExpiresAt,
		EarnedReason: reward.EarnedReason,
		RedeemedAt:   reward.RedeemedAt,
	}

	_ = ps.eventNotifier.SendToClient(userID, domain.EventLevelUp, wsReward)

	return nil
}

func (ps *PetService) GetNextRewardDescription(ctx context.Context, userID string) (*domain.Reward, error) {
	pet, err := ps.petRepository.GetPet(ctx, userID)
	if err != nil {
		return nil, err
	}

	code := ps.levelPolicy.GetCode(pet.Level + 1)
	reward, err := ps.client.GetRewardDescription(ctx, code)
	if err != nil {
		return nil, err
	}

	return reward, nil
}

func (ps *PetService) GetWeeklyLeaderboard(ctx context.Context, userID string, limit int) ([]domain.LeaderboardItem, *domain.LeaderboardItem, error) {
	records, err := ps.petRepository.GetWeeklyLeaderboardWithUser(ctx, userID, limit)
	if err != nil {
		return nil, nil, err
	}

	if len(records) == 0 {
		return []domain.LeaderboardItem{}, nil, nil
	}

	var userIDs = make([]string, len(records))
	for i := range records {
		userIDs[i] = records[i].UserID
	}

	userMap, err := ps.client.GetUsernamesByIDs(ctx, userIDs)
	if err != nil {
		return nil, nil, err
	}

	var currentUserItem *domain.LeaderboardItem
	for i := range records {
		name, ok := userMap[records[i].UserID]
		if !ok {
			continue
		}

		records[i].UserName = name

		if records[i].UserID == userID {
			currentUserItem = &records[i]
		}
	}

	topRecords := append([]domain.LeaderboardItem{}, records[:limit]...)
	currentUserItem.UserName = userMap[currentUserItem.UserID]

	return topRecords, currentUserItem, nil
}
