package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

func newService(
	userRepo *fakeUserRepository,
	tokenRepo *fakeTokenRepository,
	rewardRepo *fakeAuthRewardRepository,
	streakRepo *fakeAuthStreakRepository,
	rewardSvc *fakeAuthRewardService,
	tx *fakeAuthTransactor,
) *authService {
	return NewAuthService(userRepo, tokenRepo, tx, rewardRepo, streakRepo, rewardSvc, Config{
		JWTSecret:              []byte("test-secret"),
		AccessTokenTTL:         time.Hour,
		RefreshTokenTTL:        24 * time.Hour,
		RegistrationBonusCoins: 100,
	})
}

func TestRegister_Success(t *testing.T) {
	userRepo := &fakeUserRepository{
		addUserFn: func(ctx context.Context, arg entity.User) (string, error) { return "user-1", nil },
	}
	tokenRepo := &fakeTokenRepository{
		addTokenFn: func(ctx context.Context, userID string, token entity.RefreshToken) error { return nil },
	}
	rewardRepo := &fakeAuthRewardRepository{
		getRewardDefinitionsFn: func(ctx context.Context) ([]entity.RewardDefinition, error) {
			return []entity.RewardDefinition{{Code: "WELCOME"}}, nil
		},
	}
	streakRepo := &fakeAuthStreakRepository{
		updateStreakFn: func(ctx context.Context, streak *entity.Streak) error { return nil },
	}
	rewardSvc := &fakeAuthRewardService{
		grantRewardFn: func(ctx context.Context, userID string, code string) (*entity.UserReward, error) {
			return &entity.UserReward{}, nil
		},
	}
	s := newService(userRepo, tokenRepo, rewardRepo, streakRepo, rewardSvc, &fakeAuthTransactor{})

	tokens, err := s.Register(context.Background(), entity.User{Username: "ivan", Email: "ivan@example.com", Password: "strong-password"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens == nil || tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Fatalf("expected non-empty tokens, got %+v", tokens)
	}
	if len(userRepo.addUserCalls) != 1 {
		t.Fatalf("expected AddUser to be called once, got %d", len(userRepo.addUserCalls))
	}
	created := userRepo.addUserCalls[0]
	if created.PasswordHash == "" {
		t.Errorf("expected PasswordHash to be set before persisting")
	}
	if created.Coins != 100 {
		t.Errorf("expected registration bonus coins 100, got %d", created.Coins)
	}
	if len(rewardSvc.grantRewardCalls) != 1 || rewardSvc.grantRewardCalls[0] != "WELCOME" {
		t.Errorf("expected GrantReward to be called with WELCOME, got %+v", rewardSvc.grantRewardCalls)
	}
	if streakRepo.updateStreakArg == nil || streakRepo.updateStreakArg.CurrentStreak != 1 {
		t.Errorf("expected initial streak of 1, got %+v", streakRepo.updateStreakArg)
	}
}

func TestRegister_NoRewardDefinitions_ReturnsErrorBeforeTx(t *testing.T) {
	userRepo := &fakeUserRepository{
		addUserFn: func(ctx context.Context, arg entity.User) (string, error) {
			t.Fatalf("AddUser should not be called if reward definitions can't be loaded")
			return "", nil
		},
	}
	rewardRepo := &fakeAuthRewardRepository{
		getRewardDefinitionsFn: func(ctx context.Context) ([]entity.RewardDefinition, error) {
			return nil, errors.New("db down")
		},
	}
	s := newService(userRepo, &fakeTokenRepository{}, rewardRepo, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	_, err := s.Register(context.Background(), entity.User{Username: "ivan", Email: "ivan@example.com", Password: "pw"})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRegister_AddUserFails_PropagatesError(t *testing.T) {
	addErr := entity.ErrEmailAlreadyExists
	userRepo := &fakeUserRepository{
		addUserFn: func(ctx context.Context, arg entity.User) (string, error) { return "", addErr },
	}
	rewardRepo := &fakeAuthRewardRepository{
		getRewardDefinitionsFn: func(ctx context.Context) ([]entity.RewardDefinition, error) {
			return []entity.RewardDefinition{{Code: "WELCOME"}}, nil
		},
	}
	s := newService(userRepo, &fakeTokenRepository{}, rewardRepo, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	_, err := s.Register(context.Background(), entity.User{Username: "ivan", Email: "ivan@example.com", Password: "pw"})
	if !errors.Is(err, addErr) {
		t.Fatalf("expected wrapped %v, got %v", addErr, err)
	}
}

func TestLogin_Success(t *testing.T) {
	hash, err := hashPassword("strong-password")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	userRepo := &fakeUserRepository{
		getUserByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: email, PasswordHash: hash}, nil
		},
	}
	tokenRepo := &fakeTokenRepository{
		addTokenFn: func(ctx context.Context, userID string, token entity.RefreshToken) error { return nil },
	}
	s := newService(userRepo, tokenRepo, &fakeAuthRewardRepository{}, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	tokens, err := s.Login(context.Background(), "ivan@example.com", "strong-password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens == nil || tokens.AccessToken == "" {
		t.Fatalf("expected tokens, got %+v", tokens)
	}
}

func TestLogin_UserNotFound_ReturnsInvalidCredentials(t *testing.T) {
	userRepo := &fakeUserRepository{
		getUserByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
			return nil, entity.ErrUserNotFound
		},
	}
	s := newService(userRepo, &fakeTokenRepository{}, &fakeAuthRewardRepository{}, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	_, err := s.Login(context.Background(), "ghost@example.com", "whatever")
	if !errors.Is(err, entity.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_WrongPassword_ReturnsInvalidCredentials(t *testing.T) {
	hash, err := hashPassword("correct-password")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}
	userRepo := &fakeUserRepository{
		getUserByEmailFn: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: email, PasswordHash: hash}, nil
		},
	}
	s := newService(userRepo, &fakeTokenRepository{}, &fakeAuthRewardRepository{}, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	_, err = s.Login(context.Background(), "ivan@example.com", "wrong-password")
	if !errors.Is(err, entity.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestRefresh_Success(t *testing.T) {
	tokenRepo := &fakeTokenRepository{
		consumeRefreshTokenFn: func(ctx context.Context, hash string) (*entity.RefreshToken, error) {
			return &entity.RefreshToken{UserID: "user-1", ExpiresAt: time.Now().UTC().Add(time.Hour)}, nil
		},
		addTokenFn: func(ctx context.Context, userID string, token entity.RefreshToken) error { return nil },
	}
	s := newService(&fakeUserRepository{}, tokenRepo, &fakeAuthRewardRepository{}, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	tokens, err := s.Refresh(context.Background(), "some-refresh-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokens == nil || tokens.AccessToken == "" {
		t.Fatalf("expected tokens, got %+v", tokens)
	}
}

func TestRefresh_TokenNotFound(t *testing.T) {
	tokenRepo := &fakeTokenRepository{
		consumeRefreshTokenFn: func(ctx context.Context, hash string) (*entity.RefreshToken, error) {
			return nil, entity.ErrRefreshTokenNotFound
		},
	}
	s := newService(&fakeUserRepository{}, tokenRepo, &fakeAuthRewardRepository{}, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	_, err := s.Refresh(context.Background(), "unknown-token")
	if !errors.Is(err, entity.ErrInvalidRefreshToken) {
		t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
	}
}

func TestRefresh_TokenExpired(t *testing.T) {
	tokenRepo := &fakeTokenRepository{
		consumeRefreshTokenFn: func(ctx context.Context, hash string) (*entity.RefreshToken, error) {
			return &entity.RefreshToken{UserID: "user-1", ExpiresAt: time.Now().UTC().Add(-time.Minute)}, nil
		},
	}
	s := newService(&fakeUserRepository{}, tokenRepo, &fakeAuthRewardRepository{}, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	_, err := s.Refresh(context.Background(), "expired-token")
	if !errors.Is(err, entity.ErrInvalidRefreshToken) {
		t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
	}
}

func TestLogout_Success(t *testing.T) {
	tokenRepo := &fakeTokenRepository{
		deleteSessionFn: func(ctx context.Context, userID, tokenHash string) error { return nil },
	}
	s := newService(&fakeUserRepository{}, tokenRepo, &fakeAuthRewardRepository{}, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	err := s.Logout(context.Background(), "user-1", "some-refresh-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenRepo.deleteSessionUserID != "user-1" {
		t.Errorf("expected DeleteSession called with user-1, got %q", tokenRepo.deleteSessionUserID)
	}
	if tokenRepo.deleteSessionTokenHash == "" {
		t.Errorf("expected a non-empty token hash to be passed")
	}
}

func TestLogout_RepositoryError_Propagates(t *testing.T) {
	repoErr := errors.New("redis unavailable")
	tokenRepo := &fakeTokenRepository{
		deleteSessionFn: func(ctx context.Context, userID, tokenHash string) error { return repoErr },
	}
	s := newService(&fakeUserRepository{}, tokenRepo, &fakeAuthRewardRepository{}, &fakeAuthStreakRepository{}, &fakeAuthRewardService{}, &fakeAuthTransactor{})

	err := s.Logout(context.Background(), "user-1", "token")
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}
