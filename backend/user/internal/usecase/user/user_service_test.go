package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

func newService(
	userRepo *fakeUserRepository,
	streakRepo *fakeStreakRepository,
	rewardRepo *fakeRewardRepository,
	tx *fakeTransactor,
	pet *fakePetClient,
	tasks *fakeTasksClient,
) *userService {
	return NewUserService(userRepo, streakRepo, rewardRepo, tx, pet, tasks)
}

func TestProfile_Success(t *testing.T) {
	userRepo := &fakeUserRepository{
		getUserByIDFn: func(ctx context.Context, id string) (*entity.User, error) {
			return &entity.User{ID: id, Username: "ivan"}, nil
		},
	}
	s := newService(userRepo, &fakeStreakRepository{}, &fakeRewardRepository{}, &fakeTransactor{}, &fakePetClient{}, &fakeTasksClient{})

	got, err := s.Profile(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Username != "ivan" {
		t.Fatalf("unexpected profile: %+v", got)
	}
}

func TestProfile_NotFound(t *testing.T) {
	userRepo := &fakeUserRepository{
		getUserByIDFn: func(ctx context.Context, id string) (*entity.User, error) { return nil, entity.ErrUserNotFound },
	}
	s := newService(userRepo, &fakeStreakRepository{}, &fakeRewardRepository{}, &fakeTransactor{}, &fakePetClient{}, &fakeTasksClient{})

	_, err := s.Profile(context.Background(), "missing")
	if !errors.Is(err, entity.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestGetUsers_Success(t *testing.T) {
	userRepo := &fakeUserRepository{
		getUsersByIDsFn: func(ctx context.Context, ids []string) ([]entity.User, error) {
			return []entity.User{{ID: ids[0], Username: "ivan"}}, nil
		},
	}
	s := newService(userRepo, &fakeStreakRepository{}, &fakeRewardRepository{}, &fakeTransactor{}, &fakePetClient{}, &fakeTasksClient{})

	got, err := s.GetUsers(context.Background(), []string{"user-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Username != "ivan" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestGetUsers_RepositoryError(t *testing.T) {
	repoErr := errors.New("db down")
	userRepo := &fakeUserRepository{
		getUsersByIDsFn: func(ctx context.Context, ids []string) ([]entity.User, error) { return nil, repoErr },
	}
	s := newService(userRepo, &fakeStreakRepository{}, &fakeRewardRepository{}, &fakeTransactor{}, &fakePetClient{}, &fakeTasksClient{})

	_, err := s.GetUsers(context.Background(), []string{"user-1"})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}

func TestUpdateCoins_Success(t *testing.T) {
	userRepo := &fakeUserRepository{
		updateCoinsFn: func(ctx context.Context, userID string, coins int64) (*entity.User, error) {
			return &entity.User{ID: userID, Coins: 175}, nil
		},
	}
	s := newService(userRepo, &fakeStreakRepository{}, &fakeRewardRepository{}, &fakeTransactor{}, &fakePetClient{}, &fakeTasksClient{})

	got, err := s.UpdateCoins(context.Background(), "user-1", 25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Coins != 175 {
		t.Fatalf("unexpected coins: %d", got.Coins)
	}
}

func TestUpdateCoins_InsufficientCoins(t *testing.T) {
	userRepo := &fakeUserRepository{
		updateCoinsFn: func(ctx context.Context, userID string, coins int64) (*entity.User, error) {
			return nil, entity.ErrInsufficientCoins
		},
	}
	s := newService(userRepo, &fakeStreakRepository{}, &fakeRewardRepository{}, &fakeTransactor{}, &fakePetClient{}, &fakeTasksClient{})

	_, err := s.UpdateCoins(context.Background(), "user-1", -1000)
	if !errors.Is(err, entity.ErrInsufficientCoins) {
		t.Fatalf("expected ErrInsufficientCoins, got %v", err)
	}
}

func TestUpdateStreak_InvalidOccurredAt_ReturnsErrorWithoutTouchingRepos(t *testing.T) {
	streakRepo := &fakeStreakRepository{
		getStreakByUserIDForUpdateFn: func(ctx context.Context, userID string) (*entity.Streak, error) {
			t.Fatalf("repository should not be called for an unparseable occurredAt")
			return nil, nil
		},
	}
	s := newService(&fakeUserRepository{}, streakRepo, &fakeRewardRepository{}, &fakeTransactor{}, &fakePetClient{}, &fakeTasksClient{})

	err := s.UpdateStreak(context.Background(), "user-1", "not-a-timestamp")
	if err == nil {
		t.Fatalf("expected an error for malformed occurredAt")
	}
}

func TestUpdateStreak_SameDay_DoesNotCallPetClient(t *testing.T) {
	today := time.Now().UTC()
	streakRepo := &fakeStreakRepository{
		getStreakByUserIDForUpdateFn: func(ctx context.Context, userID string) (*entity.Streak, error) {
			return &entity.Streak{UserID: userID, CurrentStreak: 3, LastActiveDate: time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)}, nil
		},
		updateStreakFn: func(ctx context.Context, streak *entity.Streak) error {
			t.Fatalf("UpdateStreak should not be persisted when the business date hasn't changed")
			return nil
		},
	}
	pet := &fakePetClient{
		sendDailyBonusFn: func(ctx context.Context, userID string, streak int32) error {
			t.Fatalf("pet client should not be called when the streak didn't change")
			return nil
		},
	}
	s := newService(&fakeUserRepository{}, streakRepo, &fakeRewardRepository{}, &fakeTransactor{}, pet, &fakeTasksClient{})

	err := s.UpdateStreak(context.Background(), "user-1", today.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pet.sendDailyBonusCalled {
		t.Errorf("expected pet client not to be called")
	}
}

func TestUpdateStreak_NewDay_PersistsAndNotifiesPet(t *testing.T) {
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	today := time.Now().UTC()
	streakRepo := &fakeStreakRepository{
		getStreakByUserIDForUpdateFn: func(ctx context.Context, userID string) (*entity.Streak, error) {
			return &entity.Streak{UserID: userID, CurrentStreak: 3, LastActiveDate: time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)}, nil
		},
		updateStreakFn: func(ctx context.Context, streak *entity.Streak) error { return nil },
	}
	pet := &fakePetClient{
		sendDailyBonusFn: func(ctx context.Context, userID string, streak int32) error { return nil },
	}
	s := newService(&fakeUserRepository{}, streakRepo, &fakeRewardRepository{}, &fakeTransactor{}, pet, &fakeTasksClient{})

	err := s.UpdateStreak(context.Background(), "user-1", today.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !streakRepo.updateStreakCalled {
		t.Errorf("expected streak to be persisted")
	}
	if !pet.sendDailyBonusCalled {
		t.Errorf("expected pet client to be notified")
	}
	if pet.sendDailyBonusStreak != 4 {
		t.Errorf("expected incremented streak 4 to be sent, got %d", pet.sendDailyBonusStreak)
	}
}

func TestUpdateStreak_GetStreakError_PropagatesAndSkipsPetClient(t *testing.T) {
	repoErr := errors.New("db down")
	streakRepo := &fakeStreakRepository{
		getStreakByUserIDForUpdateFn: func(ctx context.Context, userID string) (*entity.Streak, error) { return nil, repoErr },
	}
	pet := &fakePetClient{
		sendDailyBonusFn: func(ctx context.Context, userID string, streak int32) error {
			t.Fatalf("pet client should not be called if fetching the streak failed")
			return nil
		},
	}
	s := newService(&fakeUserRepository{}, streakRepo, &fakeRewardRepository{}, &fakeTransactor{}, pet, &fakeTasksClient{})

	err := s.UpdateStreak(context.Background(), "user-1", time.Now().Format(time.RFC3339))
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}

func TestUpdateStreak_PetClientError_Propagates(t *testing.T) {
	yesterday := time.Now().UTC().AddDate(0, 0, -1)
	petErr := errors.New("pet service unavailable")
	streakRepo := &fakeStreakRepository{
		getStreakByUserIDForUpdateFn: func(ctx context.Context, userID string) (*entity.Streak, error) {
			return &entity.Streak{UserID: userID, CurrentStreak: 3, LastActiveDate: time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)}, nil
		},
		updateStreakFn: func(ctx context.Context, streak *entity.Streak) error { return nil },
	}
	pet := &fakePetClient{
		sendDailyBonusFn: func(ctx context.Context, userID string, streak int32) error { return petErr },
	}
	s := newService(&fakeUserRepository{}, streakRepo, &fakeRewardRepository{}, &fakeTransactor{}, pet, &fakeTasksClient{})

	err := s.UpdateStreak(context.Background(), "user-1", time.Now().Format(time.RFC3339))
	if !errors.Is(err, petErr) {
		t.Fatalf("expected wrapped pet client error, got %v", err)
	}
}

func TestGetDailyStat_Success_FiltersRewardsByDayBounds(t *testing.T) {
	now := time.Now().UTC()
	todayNoon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, time.UTC)
	yesterdayNoon := todayNoon.AddDate(0, 0, -1)

	streakRepo := &fakeStreakRepository{
		getStreakByUserIDFn: func(ctx context.Context, userID string) (*entity.Streak, error) {
			return &entity.Streak{UserID: userID, CurrentStreak: 7}, nil
		},
	}
	rewardRepo := &fakeRewardRepository{
		getRewardsByUserIDAndPeriodFn: func(ctx context.Context, userID string, from, to time.Time) ([]entity.UserReward, error) {
			return []entity.UserReward{
				{PromoCode: "TODAY", CreatedAt: todayNoon, Definition: entity.RewardDefinition{Name: "A"}},
				{PromoCode: "YESTERDAY", CreatedAt: yesterdayNoon, Definition: entity.RewardDefinition{Name: "B"}},
				{PromoCode: "REDEEMED_TODAY", CreatedAt: yesterdayNoon, RedeemedAt: &todayNoon, Definition: entity.RewardDefinition{Name: "C"}},
			}, nil
		},
	}
	tasks := &fakeTasksClient{
		getCompletedTasksFn: func(ctx context.Context, userID string) ([]entity.TasksStat, error) { return []entity.TasksStat{{}}, nil },
	}
	pet := &fakePetClient{
		getPetDailyStatFn: func(ctx context.Context, userID string) (*entity.PetStat, error) { return &entity.PetStat{}, nil },
	}
	s := newService(&fakeUserRepository{}, streakRepo, rewardRepo, &fakeTransactor{}, pet, tasks)

	stat, err := s.GetDailyStat(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stat.Streak != 7 {
		t.Errorf("expected streak 7, got %d", stat.Streak)
	}
	if len(stat.Tasks) != 1 {
		t.Errorf("expected 1 task stat, got %d", len(stat.Tasks))
	}
	if len(stat.Rewards) != 2 {
		t.Fatalf("expected 2 reward stat entries, got %d: %+v", len(stat.Rewards), stat.Rewards)
	}
	codes := map[string]bool{}
	for _, r := range stat.Rewards {
		codes[r.PromoCode] = true
	}
	if !codes["TODAY"] || !codes["REDEEMED_TODAY"] {
		t.Errorf("unexpected reward codes in stat: %+v", stat.Rewards)
	}
	if codes["YESTERDAY"] {
		t.Errorf("reward created yesterday and never redeemed should not appear")
	}
}

func TestGetDailyStat_StreakRepositoryError_StopsEarly(t *testing.T) {
	repoErr := errors.New("db down")
	streakRepo := &fakeStreakRepository{
		getStreakByUserIDFn: func(ctx context.Context, userID string) (*entity.Streak, error) { return nil, repoErr },
	}
	rewardRepo := &fakeRewardRepository{
		getRewardsByUserIDAndPeriodFn: func(ctx context.Context, userID string, from, to time.Time) ([]entity.UserReward, error) {
			t.Fatalf("reward repository should not be called if the streak lookup already failed")
			return nil, nil
		},
	}
	s := newService(&fakeUserRepository{}, streakRepo, rewardRepo, &fakeTransactor{}, &fakePetClient{}, &fakeTasksClient{})

	_, err := s.GetDailyStat(context.Background(), "user-1")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}

func TestGetDailyStat_TasksClientError_Propagates(t *testing.T) {
	tasksErr := errors.New("tasks service unavailable")
	streakRepo := &fakeStreakRepository{
		getStreakByUserIDFn: func(ctx context.Context, userID string) (*entity.Streak, error) { return &entity.Streak{}, nil },
	}
	rewardRepo := &fakeRewardRepository{
		getRewardsByUserIDAndPeriodFn: func(ctx context.Context, userID string, from, to time.Time) ([]entity.UserReward, error) {
			return nil, nil
		},
	}
	tasks := &fakeTasksClient{
		getCompletedTasksFn: func(ctx context.Context, userID string) ([]entity.TasksStat, error) { return nil, tasksErr },
	}
	pet := &fakePetClient{
		getPetDailyStatFn: func(ctx context.Context, userID string) (*entity.PetStat, error) {
			t.Fatalf("pet client should not be called if tasks client already failed")
			return nil, nil
		},
	}
	s := newService(&fakeUserRepository{}, streakRepo, rewardRepo, &fakeTransactor{}, pet, tasks)

	_, err := s.GetDailyStat(context.Background(), "user-1")
	if !errors.Is(err, tasksErr) {
		t.Fatalf("expected wrapped tasks client error, got %v", err)
	}
}

func TestGetDailyStat_PetClientError_Propagates(t *testing.T) {
	petErr := errors.New("pet service unavailable")
	streakRepo := &fakeStreakRepository{
		getStreakByUserIDFn: func(ctx context.Context, userID string) (*entity.Streak, error) { return &entity.Streak{}, nil },
	}
	rewardRepo := &fakeRewardRepository{
		getRewardsByUserIDAndPeriodFn: func(ctx context.Context, userID string, from, to time.Time) ([]entity.UserReward, error) {
			return nil, nil
		},
	}
	tasks := &fakeTasksClient{
		getCompletedTasksFn: func(ctx context.Context, userID string) ([]entity.TasksStat, error) { return nil, nil },
	}
	pet := &fakePetClient{
		getPetDailyStatFn: func(ctx context.Context, userID string) (*entity.PetStat, error) { return nil, petErr },
	}
	s := newService(&fakeUserRepository{}, streakRepo, rewardRepo, &fakeTransactor{}, pet, tasks)

	_, err := s.GetDailyStat(context.Background(), "user-1")
	if !errors.Is(err, petErr) {
		t.Fatalf("expected wrapped pet client error, got %v", err)
	}
}
