package reward

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

const maxPromoCodeGenerationAttempts = 5

type (
	rewardRepository interface {
		GetUserRewardsByUserID(ctx context.Context, userID string) ([]entity.UserReward, error)
		GetActiveRewardsByUserID(ctx context.Context, userID string) ([]entity.UserReward, error)
		GetRewardByUserIDAndRewardID(ctx context.Context, userID, rewardID string) (*entity.UserReward, error)
		GetRewardDefinitionByCode(ctx context.Context, code string) (*entity.RewardDefinition, error)
		RedeemUserReward(ctx context.Context, userID, promoCode string) error
		AddUserReward(ctx context.Context, userID string, promoCode string, rewardID int32, expiresAt *time.Time) (*entity.UserReward, error)
	}
)

type rewardService struct {
	rewardRepository rewardRepository
}

func NewRewardService(
	rewardRepository rewardRepository,
) *rewardService {
	return &rewardService{
		rewardRepository: rewardRepository,
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

func (s *rewardService) GrantReward(ctx context.Context, userID string, code string) (*entity.UserReward, error) {
	definition, err := s.rewardRepository.GetRewardDefinitionByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("grant reward: get definition: %w", err)
	}

	for attempt := 0; attempt < maxPromoCodeGenerationAttempts; attempt++ {
		promoCode, err := generatePromoCode(definition.Code)
		if err != nil {
			return nil, fmt.Errorf("grant reward: generate promo code: %w", err)
		}

		reward, err := s.rewardRepository.AddUserReward(ctx, userID, promoCode, definition.ID, nil)
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
