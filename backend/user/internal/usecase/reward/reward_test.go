package reward

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/usecase/reward/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRewardService_GetAllRewards_UpdatesStatus(t *testing.T) {
	t.Parallel()

	expiredAt := time.Now().Add(-time.Minute)

	deps := newTestDeps(t)
	deps.rewardRepository.EXPECT().
		GetUserRewardsByUserID(gomock.Any(), "user-id").
		Return([]entity.UserReward{
			{
				Status:    entity.StatusActive,
				ExpiresAt: &expiredAt,
			},
			{
				Status: entity.StatusRedeemed,
			},
		}, nil)

	rewards, err := deps.service.GetAllRewards(context.Background(), "user-id")

	require.NoError(t, err)
	require.Len(t, rewards, 2)
	require.Equal(t, entity.StatusExpired, rewards[0].Status)
	require.Equal(t, entity.StatusRedeemed, rewards[1].Status)
}

func TestRewardService_GrantReward(t *testing.T) {
	t.Parallel()

	definition := &entity.RewardDefinition{
		ID:   42,
		Code: "AVITO",
		Name: "Reward",
	}

	deps := newTestDeps(t)

	deps.rewardRepository.EXPECT().
		GetRewardDefinitionByCode(gomock.Any(), "AVITO").
		Return(definition, nil)

	deps.transactor.EXPECT().
		WithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	deps.rewardRepository.EXPECT().
		AddUserReward(
			gomock.Any(),
			"user-id",
			gomock.Any(),
			"achievement",
			int32(42),
			gomock.Any(),
		).
		DoAndReturn(func(
			_ context.Context,
			userID, promoCode, earnedReason string,
			rewardID int32,
			expiresAt *time.Time,
		) (*entity.UserReward, error) {
			require.Equal(t, "user-id", userID)
			require.Equal(t, "achievement", earnedReason)
			require.Equal(t, int32(42), rewardID)
			require.True(t, strings.HasPrefix(promoCode, "AVITO-"))
			require.Len(t, promoCode, len("AVITO-")+promoCodeRandomLength)
			require.NotNil(t, expiresAt)

			return &entity.UserReward{
				ID:        "reward-id",
				PromoCode: promoCode,
			}, nil
		})

	deps.eventRepository.EXPECT().
		AddUserEvent(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, event entity.UserEvent) error {
			require.Equal(t, "user-id", event.UserID)
			require.Equal(t, entity.NewReward, event.Type)
			require.NotNil(t, event.UserRewardID)
			require.Equal(t, "reward-id", *event.UserRewardID)
			return nil
		})

	reward, err := deps.service.GrantReward(
		context.Background(),
		"user-id",
		"AVITO",
		"achievement",
	)

	require.NoError(t, err)
	require.Equal(t, definition.Code, reward.Definition.Code)
	require.Equal(t, definition.Name, reward.Definition.Name)
}

func TestRewardService_GrantReward_RetryOnPromoCodeCollision(t *testing.T) {
	t.Parallel()

	deps := newTestDeps(t)

	deps.rewardRepository.EXPECT().
		GetRewardDefinitionByCode(gomock.Any(), "AVITO").
		Return(&entity.RewardDefinition{
			ID:   42,
			Code: "AVITO",
		}, nil)

	deps.transactor.EXPECT().
		WithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Times(2)

	deps.rewardRepository.EXPECT().
		AddUserReward(
			gomock.Any(),
			"user-id",
			gomock.Any(),
			"achievement",
			int32(42),
			gomock.Any(),
		).
		Return(nil, entity.ErrPromoCodeAlreadyExists)

	deps.rewardRepository.EXPECT().
		AddUserReward(
			gomock.Any(),
			"user-id",
			gomock.Any(),
			"achievement",
			int32(42),
			gomock.Any(),
		).
		Return(&entity.UserReward{ID: "reward-id"}, nil)

	deps.eventRepository.EXPECT().
		AddUserEvent(gomock.Any(), gomock.Any()).
		Return(nil)

	reward, err := deps.service.GrantReward(
		context.Background(),
		"user-id",
		"AVITO",
		"achievement",
	)

	require.NoError(t, err)
	require.Equal(t, "reward-id", reward.ID)
}

func Test_UpdateStatus(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 14, 10, 0, 0, 0, time.UTC)
	expiredAt := now.Add(-time.Second)
	activeUntil := now.Add(time.Hour)

	tests := []struct {
		name     string
		reward   entity.UserReward
		expected entity.Status
	}{
		{
			name: "redeemed stays redeemed",
			reward: entity.UserReward{
				Status:    entity.StatusRedeemed,
				ExpiresAt: &expiredAt,
			},
			expected: entity.StatusRedeemed,
		},
		{
			name: "expired",
			reward: entity.UserReward{
				Status:    entity.StatusActive,
				ExpiresAt: &expiredAt,
			},
			expected: entity.StatusExpired,
		},
		{
			name: "active",
			reward: entity.UserReward{
				Status:    entity.StatusActive,
				ExpiresAt: &activeUntil,
			},
			expected: entity.StatusActive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expected, updateStatus(&tt.reward, now))
		})
	}
}

func TestGeneratePromoCode(t *testing.T) {
	t.Parallel()

	code, err := generatePromoCode("AVITO")

	require.NoError(t, err)
	require.True(t, strings.HasPrefix(code, "AVITO-"))
	require.Len(t, code, len("AVITO-")+promoCodeRandomLength)

	randomPart := strings.TrimPrefix(code, "AVITO-")
	for _, ch := range randomPart {
		require.Contains(t, promoCodeAlphabet, string(ch))
	}
}

type testDeps struct {
	service          *rewardService
	rewardRepository *mocks.MockrewardRepository
	eventRepository  *mocks.MockeventRepository
	transactor       *mocks.Mocktransactor
}

func newTestDeps(t *testing.T) testDeps {
	t.Helper()

	ctrl := gomock.NewController(t)

	rewardRepository := mocks.NewMockrewardRepository(ctrl)
	eventRepository := mocks.NewMockeventRepository(ctrl)
	transactor := mocks.NewMocktransactor(ctrl)

	return testDeps{
		service: NewRewardService(
			rewardRepository,
			eventRepository,
			transactor,
		),
		rewardRepository: rewardRepository,
		eventRepository:  eventRepository,
		transactor:       transactor,
	}
}
