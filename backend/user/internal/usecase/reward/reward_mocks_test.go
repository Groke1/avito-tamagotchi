package reward

import (
	"context"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

type fakeRewardRepository struct {
	getUserRewardsByUserIDFn      func(ctx context.Context, userID string) ([]entity.UserReward, error)
	getActiveRewardsByUserIDFn    func(ctx context.Context, userID string) ([]entity.UserReward, error)
	getRewardByUserIDAndRewardIDFn func(ctx context.Context, userID, rewardID string) (*entity.UserReward, error)
	getRewardDefinitionByCodeFn   func(ctx context.Context, code string) (*entity.RewardDefinition, error)
	redeemUserRewardFn            func(ctx context.Context, userID, promoCode string) error
	addUserRewardFn                func(ctx context.Context, userID string, promoCode string, rewardID int32, expiresAt *time.Time) (*entity.UserReward, error)

	addUserRewardCalls []string
}

func (f *fakeRewardRepository) GetUserRewardsByUserID(ctx context.Context, userID string) ([]entity.UserReward, error) {
	return f.getUserRewardsByUserIDFn(ctx, userID)
}

func (f *fakeRewardRepository) GetActiveRewardsByUserID(ctx context.Context, userID string) ([]entity.UserReward, error) {
	return f.getActiveRewardsByUserIDFn(ctx, userID)
}

func (f *fakeRewardRepository) GetRewardByUserIDAndRewardID(ctx context.Context, userID, rewardID string) (*entity.UserReward, error) {
	return f.getRewardByUserIDAndRewardIDFn(ctx, userID, rewardID)
}

func (f *fakeRewardRepository) GetRewardDefinitionByCode(ctx context.Context, code string) (*entity.RewardDefinition, error) {
	return f.getRewardDefinitionByCodeFn(ctx, code)
}

func (f *fakeRewardRepository) RedeemUserReward(ctx context.Context, userID, promoCode string) error {
	return f.redeemUserRewardFn(ctx, userID, promoCode)
}

func (f *fakeRewardRepository) AddUserReward(ctx context.Context, userID string, promoCode string, rewardID int32, expiresAt *time.Time) (*entity.UserReward, error) {
	f.addUserRewardCalls = append(f.addUserRewardCalls, promoCode)
	return f.addUserRewardFn(ctx, userID, promoCode, rewardID, expiresAt)
}
