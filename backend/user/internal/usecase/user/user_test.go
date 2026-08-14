package user

import (
	"context"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/usecase/user/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserService_UpdateStreak(t *testing.T) {
	t.Parallel()

	currentDate := time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)
	businessDate := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		streak        *entity.Streak
		expectedValue int32
		expectedBonus int64
		changed       bool
	}{
		{
			name: "continue streak",
			streak: &entity.Streak{
				UserID:         "user-id",
				CurrentStreak:  5,
				LastActiveDate: currentDate.AddDate(0, 0, -1),
			},
			expectedValue: 6,
			expectedBonus: 0,
			changed:       true,
		},
		{
			name: "weekly bonus",
			streak: &entity.Streak{
				UserID:         "user-id",
				CurrentStreak:  6,
				LastActiveDate: currentDate.AddDate(0, 0, -1),
			},
			expectedValue: 7,
			expectedBonus: 20,
			changed:       true,
		},
		{
			name: "already active today",
			streak: &entity.Streak{
				UserID:         "user-id",
				CurrentStreak:  5,
				LastActiveDate: currentDate,
			},
			expectedValue: 5,
			changed:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)

			deps.transactor.EXPECT().
				WithTx(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
					return fn(ctx)
				})

			deps.streakRepository.EXPECT().
				GetStreakByUserIDForUpdate(gomock.Any(), "user-id").
				Return(tt.streak, nil)

			if tt.changed {
				deps.streakRepository.EXPECT().
					UpdateStreak(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, streak *entity.Streak) error {
						require.Equal(t, tt.expectedValue, streak.CurrentStreak)
						require.Equal(t, businessDate, streak.LastActiveDate)
						return nil
					})

				deps.userRepository.EXPECT().
					UpdateCoins(gomock.Any(), "user-id", tt.expectedBonus).
					Return(&entity.User{}, nil)

				deps.petClient.EXPECT().
					SendDailyBonus(gomock.Any(), "user-id", tt.expectedValue, int32(tt.expectedBonus)).
					Return(nil)

				deps.eventRepository.EXPECT().
					AddUserEvent(gomock.Any(), gomock.Any()).
					DoAndReturn(func(_ context.Context, event entity.UserEvent) error {
						require.Equal(t, "user-id", event.UserID)
						require.Equal(t, entity.StreakReward, event.Type)
						require.Equal(t, int32(tt.expectedBonus), event.Coins)
						require.NotNil(t, event.Streak)
						require.Equal(t, tt.expectedValue, *event.Streak)
						return nil
					})
			}

			err := deps.service.UpdateStreak(
				context.Background(),
				"user-id",
				currentDate.Format(time.RFC3339),
			)

			require.NoError(t, err)
			require.Equal(t, tt.expectedValue, tt.streak.CurrentStreak)
		})
	}
}

func TestUserService_UpdateStreak_Reset(t *testing.T) {
	t.Parallel()

	currentDate := time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)
	businessDate := time.Date(2026, 8, 14, 0, 0, 0, 0, time.UTC)
	streak := &entity.Streak{
		UserID:         "user-id",
		CurrentStreak:  10,
		LastActiveDate: currentDate.AddDate(0, 0, -2),
	}

	deps := newTestDeps(t)

	deps.transactor.EXPECT().
		WithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	deps.streakRepository.EXPECT().
		GetStreakByUserIDForUpdate(gomock.Any(), "user-id").
		Return(streak, nil)

	deps.streakRepository.EXPECT().
		UpdateStreak(gomock.Any(), gomock.Any()).
		Return(nil)

	deps.userRepository.EXPECT().
		UpdateCoins(gomock.Any(), "user-id", int64(0)).
		Return(&entity.User{}, nil)

	deps.petClient.EXPECT().
		SendDailyBonus(gomock.Any(), "user-id", int32(1), int32(0)).
		Return(nil)

	deps.eventRepository.EXPECT().
		AddUserEvent(gomock.Any(), gomock.Any()).
		Return(nil)

	err := deps.service.UpdateStreak(
		context.Background(),
		"user-id",
		currentDate.Format(time.RFC3339),
	)

	require.NoError(t, err)
	require.Equal(t, int32(1), streak.CurrentStreak)
	require.Equal(t, businessDate, streak.LastActiveDate)
}

func TestUserService_GetDailyStat_ExpiredStreak(t *testing.T) {
	t.Parallel()

	deps := newTestDeps(t)

	deps.streakRepository.EXPECT().
		GetStreakByUserID(gomock.Any(), "user-id").
		Return(&entity.Streak{
			UserID:         "user-id",
			CurrentStreak:  10,
			LastActiveDate: time.Now().AddDate(0, 0, -2),
		}, nil)

	deps.rewardRepository.EXPECT().
		GetRewardsByUserIDAndPeriod(gomock.Any(), "user-id", gomock.Any(), gomock.Any()).
		Return(nil, nil)

	deps.tasksClient.EXPECT().
		GetCompletedTasks(gomock.Any(), "user-id").
		Return(nil, nil)

	deps.petClient.EXPECT().
		GetPetDailyStat(gomock.Any(), "user-id").
		Return(&entity.PetStat{}, nil)

	stat, err := deps.service.GetDailyStat(context.Background(), "user-id")

	require.NoError(t, err)
	require.Equal(t, "user-id", stat.UserID)
	require.Equal(t, int32(0), stat.Streak)
	require.Empty(t, stat.Rewards)
}

func TestGetBonusCoins(t *testing.T) {
	t.Parallel()

	require.Equal(t, int32(0), getBonusCoins(6))
	require.Equal(t, int32(20), getBonusCoins(7))
	require.Equal(t, int32(40), getBonusCoins(14))
}

type testDeps struct {
	service          *userService
	userRepository   *mocks.MockuserRepository
	streakRepository *mocks.MockstreakRepository
	rewardRepository *mocks.MockrewardRepository
	eventRepository  *mocks.MockeventRepository
	transactor       *mocks.Mocktransactor
	petClient        *mocks.MockpetClient
	tasksClient      *mocks.MocktasksClient
}

func newTestDeps(t *testing.T) testDeps {
	t.Helper()

	ctrl := gomock.NewController(t)

	userRepository := mocks.NewMockuserRepository(ctrl)
	streakRepository := mocks.NewMockstreakRepository(ctrl)
	rewardRepository := mocks.NewMockrewardRepository(ctrl)
	eventRepository := mocks.NewMockeventRepository(ctrl)
	transactor := mocks.NewMocktransactor(ctrl)
	petClient := mocks.NewMockpetClient(ctrl)
	tasksClient := mocks.NewMocktasksClient(ctrl)

	return testDeps{
		service: NewUserService(
			userRepository,
			streakRepository,
			rewardRepository,
			eventRepository,
			transactor,
			petClient,
			tasksClient,
		),
		userRepository:   userRepository,
		streakRepository: streakRepository,
		rewardRepository: rewardRepository,
		eventRepository:  eventRepository,
		transactor:       transactor,
		petClient:        petClient,
		tasksClient:      tasksClient,
	}
}
