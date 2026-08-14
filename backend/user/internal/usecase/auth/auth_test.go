package auth

import (
	"context"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/usecase/auth/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Register(t *testing.T) {
	t.Parallel()

	deps := newTestDeps(t)

	deps.rewardRepository.EXPECT().
		GetRewardDefinitions(gomock.Any()).
		Return([]entity.RewardDefinition{
			{ID: 1, Code: "WELCOME"},
		}, nil)

	deps.transactor.EXPECT().
		WithTx(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	deps.userRepository.EXPECT().
		AddUser(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, user entity.User) (string, error) {
			require.NotEmpty(t, user.PasswordHash)
			require.NotEqual(t, "password123", user.PasswordHash)
			require.Equal(t, uint64(100), user.Coins)
			return "user-id", nil
		})

	deps.tokenRepository.EXPECT().
		AddToken(gomock.Any(), "user-id", gomock.Any()).
		Return(nil)

	deps.rewardService.EXPECT().
		GrantReward(gomock.Any(), "user-id", "WELCOME", "Награда за регистрацию").
		Return(&entity.UserReward{}, nil)

	deps.streakRepository.EXPECT().
		UpdateStreak(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, streak *entity.Streak) error {
			require.Equal(t, "user-id", streak.UserID)
			require.Equal(t, int32(1), streak.CurrentStreak)
			return nil
		})

	tokens, err := deps.service.Register(context.Background(), entity.User{
		Username: "test-user",
		Email:    "user@example.com",
		Password: "password123",
	})

	require.NoError(t, err)
	require.NotEmpty(t, tokens.AccessToken)
	require.NotEmpty(t, tokens.RefreshToken)
}

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	passwordHash, err := hashPassword("password123")
	require.NoError(t, err)

	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:     "success",
			password: "password123",
		},
		{
			name:        "invalid password",
			password:    "wrong-password",
			expectedErr: entity.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)

			deps.userRepository.EXPECT().
				GetUserByEmail(gomock.Any(), "user@example.com").
				Return(&entity.User{
					ID:           "user-id",
					PasswordHash: passwordHash,
				}, nil)

			if tt.expectedErr == nil {
				deps.tokenRepository.EXPECT().
					AddToken(gomock.Any(), "user-id", gomock.Any()).
					Return(nil)
			}

			tokens, err := deps.service.Login(
				context.Background(),
				"user@example.com",
				tt.password,
			)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				require.Nil(t, tokens)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, tokens.AccessToken)
			require.NotEmpty(t, tokens.RefreshToken)
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		storedToken *entity.RefreshToken
		repoErr     error
		expectedErr error
	}{
		{
			name: "success",
			storedToken: &entity.RefreshToken{
				UserID:    "user-id",
				ExpiresAt: time.Now().UTC().Add(time.Hour),
			},
		},
		{
			name: "expired token",
			storedToken: &entity.RefreshToken{
				UserID:    "user-id",
				ExpiresAt: time.Now().UTC().Add(-time.Minute),
			},
			expectedErr: entity.ErrInvalidRefreshToken,
		},
		{
			name:        "token not found",
			repoErr:     entity.ErrRefreshTokenNotFound,
			expectedErr: entity.ErrInvalidRefreshToken,
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

			deps.tokenRepository.EXPECT().
				ConsumeRefreshToken(gomock.Any(), hashToken("refresh-token")).
				Return(tt.storedToken, tt.repoErr)

			if tt.expectedErr == nil {
				deps.tokenRepository.EXPECT().
					AddToken(gomock.Any(), "user-id", gomock.Any()).
					Return(nil)
			}

			tokens, err := deps.service.Refresh(context.Background(), "refresh-token")

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
				require.Nil(t, tokens)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, tokens.AccessToken)
			require.NotEmpty(t, tokens.RefreshToken)
		})
	}
}

func TestAuthService_ValidateAccessToken(t *testing.T) {
	t.Parallel()

	deps := newTestDeps(t)

	token, err := deps.service.newAccessToken("user-id")
	require.NoError(t, err)

	userID, err := deps.service.ValidateAccessToken(context.Background(), token)

	require.NoError(t, err)
	require.Equal(t, "user-id", userID)

	_, err = deps.service.ValidateAccessToken(context.Background(), token+"broken")
	require.ErrorIs(t, err, entity.ErrInvalidAccessToken)
}

func TestPasswordHash(t *testing.T) {
	t.Parallel()

	hash, err := hashPassword("password123")

	require.NoError(t, err)
	require.True(t, checkPasswordHash("password123", hash))
	require.False(t, checkPasswordHash("wrong-password", hash))
}

type testDeps struct {
	service          *authService
	userRepository   *mocks.MockuserRepository
	tokenRepository  *mocks.MockTokenRepository
	rewardRepository *mocks.MockrewardRepository
	streakRepository *mocks.MockstreakRepository
	rewardService    *mocks.MockrewardService
	transactor       *mocks.Mocktransactor
}

func newTestDeps(t *testing.T) testDeps {
	t.Helper()

	ctrl := gomock.NewController(t)

	userRepository := mocks.NewMockuserRepository(ctrl)
	tokenRepository := mocks.NewMockTokenRepository(ctrl)
	rewardRepository := mocks.NewMockrewardRepository(ctrl)
	streakRepository := mocks.NewMockstreakRepository(ctrl)
	rewardService := mocks.NewMockrewardService(ctrl)
	transactor := mocks.NewMocktransactor(ctrl)

	service := NewAuthService(
		userRepository,
		tokenRepository,
		transactor,
		rewardRepository,
		streakRepository,
		rewardService,
		Config{
			JWTSecret:              []byte("test-secret"),
			AccessTokenTTL:         time.Hour,
			RefreshTokenTTL:        24 * time.Hour,
			RegistrationBonusCoins: 100,
		},
	)

	return testDeps{
		service:          service,
		userRepository:   userRepository,
		tokenRepository:  tokenRepository,
		rewardRepository: rewardRepository,
		streakRepository: streakRepository,
		rewardService:    rewardService,
		transactor:       transactor,
	}
}
