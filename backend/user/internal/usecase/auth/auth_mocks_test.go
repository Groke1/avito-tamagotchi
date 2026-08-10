package auth

import (
	"context"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

type fakeUserRepository struct {
	addUserFn        func(ctx context.Context, arg entity.User) (string, error)
	getUserByEmailFn func(ctx context.Context, email string) (*entity.User, error)
	addUserCalls     []entity.User
}

func (f *fakeUserRepository) AddUser(ctx context.Context, arg entity.User) (string, error) {
	f.addUserCalls = append(f.addUserCalls, arg)
	return f.addUserFn(ctx, arg)
}

func (f *fakeUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return f.getUserByEmailFn(ctx, email)
}

type fakeTokenRepository struct {
	addTokenFn             func(ctx context.Context, userID string, token entity.RefreshToken) error
	consumeRefreshTokenFn  func(ctx context.Context, hash string) (*entity.RefreshToken, error)
	deleteSessionFn        func(ctx context.Context, userID, tokenHash string) error
	addTokenCalls          []entity.RefreshToken
	deleteSessionUserID    string
	deleteSessionTokenHash string
}

func (f *fakeTokenRepository) AddToken(ctx context.Context, userID string, token entity.RefreshToken) error {
	f.addTokenCalls = append(f.addTokenCalls, token)
	return f.addTokenFn(ctx, userID, token)
}

func (f *fakeTokenRepository) ConsumeRefreshToken(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	return f.consumeRefreshTokenFn(ctx, hash)
}

func (f *fakeTokenRepository) DeleteSession(ctx context.Context, userID, tokenHash string) error {
	f.deleteSessionUserID = userID
	f.deleteSessionTokenHash = tokenHash
	return f.deleteSessionFn(ctx, userID, tokenHash)
}

type fakeAuthRewardRepository struct {
	getRewardDefinitionsFn func(ctx context.Context) ([]entity.RewardDefinition, error)
}

func (f *fakeAuthRewardRepository) GetRewardDefinitions(ctx context.Context) ([]entity.RewardDefinition, error) {
	return f.getRewardDefinitionsFn(ctx)
}

type fakeAuthStreakRepository struct {
	updateStreakFn func(ctx context.Context, streak *entity.Streak) error
	updateStreakArg *entity.Streak
}

func (f *fakeAuthStreakRepository) UpdateStreak(ctx context.Context, streak *entity.Streak) error {
	f.updateStreakArg = streak
	return f.updateStreakFn(ctx, streak)
}

type fakeAuthRewardService struct {
	grantRewardFn func(ctx context.Context, userID string, code string) (*entity.UserReward, error)
	grantRewardCalls []string
}

func (f *fakeAuthRewardService) GrantReward(ctx context.Context, userID string, code string) (*entity.UserReward, error) {
	f.grantRewardCalls = append(f.grantRewardCalls, code)
	return f.grantRewardFn(ctx, userID, code)
}

type fakeAuthTransactor struct {
	beginErr error
}

func (t *fakeAuthTransactor) WithTx(ctx context.Context, f func(ctx context.Context) error) error {
	if t.beginErr != nil {
		return t.beginErr
	}
	return f(ctx)
}
