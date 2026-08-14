package reward

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

const maxPromoCodeGenerationAttempts = 5

//go:generate mockgen -source=reward.go -destination=mocks/reward_mocks.go -package=mocks
type (
	rewardRepository interface {
		GetUserRewardsByUserID(ctx context.Context, userID string) ([]entity.UserReward, error)
		GetActiveRewardsByUserID(ctx context.Context, userID string) ([]entity.UserReward, error)
		GetRewardByUserIDAndRewardID(ctx context.Context, userID, rewardID string) (*entity.UserReward, error)
		GetRewardDefinitionByCode(ctx context.Context, code string) (*entity.RewardDefinition, error)
		RedeemUserReward(ctx context.Context, userID, promoCode string) error
		AddUserReward(ctx context.Context, userID string, promoCode string, earnedReason string,
			rewardID int32, expiresAt *time.Time) (*entity.UserReward, error)
	}

	eventRepository interface {
		AddUserEvent(ctx context.Context, event entity.UserEvent) error
	}

	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) error
	}
)

type rewardService struct {
	rewardRepository rewardRepository
	eventRepository  eventRepository
	transactor       transactor
}

func NewRewardService(
	rewardRepository rewardRepository,
	eventRepository eventRepository,
	transactor transactor,
) *rewardService {
	return &rewardService{
		rewardRepository: rewardRepository,
		eventRepository:  eventRepository,
		transactor:       transactor,
	}
}

func (s *rewardService) GetAllRewards(ctx context.Context, userID string) ([]entity.UserReward, error) {
	rewards, err := s.rewardRepository.GetUserRewardsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get all rewards: %w", err)
	}

	now := time.Now()

	for i := range rewards {
		rewards[i].Status = updateStatus(&rewards[i], now)
	}

	return rewards, nil
}

func (s *rewardService) GetActiveRewards(ctx context.Context, userID string) ([]entity.UserReward, error) {
	rewards, err := s.rewardRepository.GetActiveRewardsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get active rewards: %w", err)
	}

	return rewards, nil
}

func (s *rewardService) GetReward(ctx context.Context, userID string, rewardID string) (*entity.UserReward, error) {
	reward, err := s.rewardRepository.GetRewardByUserIDAndRewardID(ctx, userID, rewardID)
	if err != nil {
		return nil, fmt.Errorf("get reward: %w", err)
	}

	reward.Status = updateStatus(reward, time.Now())

	return reward, nil
}

func (s *rewardService) RedeemReward(ctx context.Context, userID string, promoCode string) error {
	if err := s.rewardRepository.RedeemUserReward(ctx, userID, promoCode); err != nil {
		return fmt.Errorf("redeem reward: %w", err)
	}

	return nil
}

func (s *rewardService) GetDefinition(ctx context.Context, code string) (*entity.RewardDefinition, error) {
	definition, err := s.rewardRepository.GetRewardDefinitionByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get reward definition: %w", err)
	}

	return definition, nil
}

func (s *rewardService) GrantReward(ctx context.Context, userID, code, earnedReason string) (*entity.UserReward, error) {
	definition, err := s.rewardRepository.GetRewardDefinitionByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("grant reward: get definition: %w", err)
	}

	for attempt := 0; attempt < maxPromoCodeGenerationAttempts; attempt++ {
		var promoCode string
		promoCode, err = generatePromoCode(definition.Code)
		if err != nil {
			return nil, fmt.Errorf("grant reward: generate promo code: %w", err)
		}

		var reward *entity.UserReward

		err = s.transactor.WithTx(ctx, func(ctx context.Context) error {
			expiresAt := time.Now().UTC().Add(2 * time.Minute)
			reward, err = s.rewardRepository.AddUserReward(ctx, userID, promoCode, earnedReason, definition.ID, &expiresAt)
			if err != nil {
				return fmt.Errorf("grant reward: %w", err)
			}
			err = s.eventRepository.AddUserEvent(ctx, entity.UserEvent{
				UserID:       userID,
				Type:         entity.NewReward,
				UserRewardID: &reward.ID,
			})

			if err != nil {
				return fmt.Errorf("add reward event: %w", err)
			}

			return nil
		})

		if err == nil {
			reward.Definition = *definition
			return reward, nil
		}

		if errors.Is(err, entity.ErrPromoCodeAlreadyExists) {
			continue
		}

		return nil, fmt.Errorf("grant reward: add user reward: %w", err)
	}

	return nil, fmt.Errorf("can not generate promo code %w", err)
}

func updateStatus(reward *entity.UserReward, now time.Time) entity.Status {
	if reward.Status == entity.StatusRedeemed {
		return entity.StatusRedeemed
	}

	if reward.ExpiresAt != nil && !reward.ExpiresAt.After(now) {
		return entity.StatusExpired
	}

	return entity.StatusActive
}
