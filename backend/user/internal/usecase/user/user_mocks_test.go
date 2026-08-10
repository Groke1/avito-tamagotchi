package user

import (
	"context"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

type fakeUserRepository struct {
	getUserByIDFn   func(ctx context.Context, id string) (*entity.User, error)
	getUsersByIDsFn func(ctx context.Context, ids []string) ([]entity.User, error)
	updateCoinsFn   func(ctx context.Context, userID string, coins int64) (*entity.User, error)
}

func (f *fakeUserRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	return f.getUserByIDFn(ctx, id)
}
func (f *fakeUserRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]entity.User, error) {
	return f.getUsersByIDsFn(ctx, ids)
}
func (f *fakeUserRepository) UpdateCoins(ctx context.Context, userID string, coins int64) (*entity.User, error) {
	return f.updateCoinsFn(ctx, userID, coins)
}

type fakeStreakRepository struct {
	getStreakByUserIDFn          func(ctx context.Context, userID string) (*entity.Streak, error)
	getStreakByUserIDForUpdateFn func(ctx context.Context, userID string) (*entity.Streak, error)
	updateStreakFn               func(ctx context.Context, streak *entity.Streak) error
	updateStreakCalled           bool
}

func (f *fakeStreakRepository) GetStreakByUserID(ctx context.Context, userID string) (*entity.Streak, error) {
	return f.getStreakByUserIDFn(ctx, userID)
}
func (f *fakeStreakRepository) GetStreakByUserIDForUpdate(ctx context.Context, userID string) (*entity.Streak, error) {
	return f.getStreakByUserIDForUpdateFn(ctx, userID)
}
func (f *fakeStreakRepository) UpdateStreak(ctx context.Context, streak *entity.Streak) error {
	f.updateStreakCalled = true
	return f.updateStreakFn(ctx, streak)
}

type fakeRewardRepository struct {
	getRewardsByUserIDAndPeriodFn func(ctx context.Context, userID string, from, to time.Time) ([]entity.UserReward, error)
}

func (f *fakeRewardRepository) GetRewardsByUserIDAndPeriod(ctx context.Context, userID string, from, to time.Time) ([]entity.UserReward, error) {
	return f.getRewardsByUserIDAndPeriodFn(ctx, userID, from, to)
}

type fakeTransactor struct {
	beginErr error
}

func (t *fakeTransactor) WithTx(ctx context.Context, f func(ctx context.Context) error) error {
	if t.beginErr != nil {
		return t.beginErr
	}
	return f(ctx)
}

type fakePetClient struct {
	sendDailyBonusFn   func(ctx context.Context, userID string, streak int32) error
	getPetDailyStatFn  func(ctx context.Context, userID string) (*entity.PetStat, error)
	sendDailyBonusCalled bool
	sendDailyBonusStreak int32
}

func (f *fakePetClient) SendDailyBonus(ctx context.Context, userID string, streak int32) error {
	f.sendDailyBonusCalled = true
	f.sendDailyBonusStreak = streak
	return f.sendDailyBonusFn(ctx, userID, streak)
}
func (f *fakePetClient) GetPetDailyStat(ctx context.Context, userID string) (*entity.PetStat, error) {
	return f.getPetDailyStatFn(ctx, userID)
}

type fakeTasksClient struct {
	getCompletedTasksFn func(ctx context.Context, userID string) ([]entity.TasksStat, error)
}

func (f *fakeTasksClient) GetCompletedTasks(ctx context.Context, userID string) ([]entity.TasksStat, error) {
	return f.getCompletedTasksFn(ctx, userID)
}
