package auth

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

//go:generate mockgen -source=auth.go -destination=mocks/middleware_auth.go -package=mocks

type (
	userRepository interface {
		AddUser(ctx context.Context, arg entity.User) (string, error)
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	}

	TokenRepository interface {
		AddToken(ctx context.Context, userID string, token entity.RefreshToken) error
		ConsumeRefreshToken(ctx context.Context, hash string) (*entity.RefreshToken, error)
		DeleteSession(ctx context.Context, userID, tokenHash string) error
	}

	rewardRepository interface {
		GetRewardDefinitions(ctx context.Context) ([]entity.RewardDefinition, error)
	}

	streakRepository interface {
		UpdateStreak(ctx context.Context, streak *entity.Streak) error
	}

	rewardService interface {
		GrantReward(ctx context.Context, userID, code, earnedReason string) (*entity.UserReward, error)
	}

	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) error
	}
)

type authService struct {
	userRepository   userRepository
	tokenRepository  TokenRepository
	rewardRepository rewardRepository
	streakRepository streakRepository
	transactor       transactor
	rewardService    rewardService
	cfg              *Config
}

type Config struct {
	JWTSecret              []byte
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	RegistrationBonusCoins uint64
}

func NewAuthService(
	userRepository userRepository,
	tokenRepository TokenRepository,
	transactor transactor,
	rewardRepository rewardRepository,
	streakRepository streakRepository,
	rewardService rewardService,
	cfg Config,
) *authService {
	return &authService{
		userRepository:   userRepository,
		tokenRepository:  tokenRepository,
		transactor:       transactor,
		rewardRepository: rewardRepository,
		streakRepository: streakRepository,
		rewardService:    rewardService,
		cfg:              &cfg,
	}
}

func (s *authService) Register(ctx context.Context, user entity.User) (*entity.JWT, error) {
	passwordHash, err := hashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("register: hash password: %w", err)
	}

	user.PasswordHash = passwordHash
	user.Coins = s.cfg.RegistrationBonusCoins

	eventTime := time.Now()

	businessDate := dateOnly(eventTime)

	rewardDefinitions, err := s.rewardRepository.GetRewardDefinitions(ctx)
	if err != nil {
		return nil, fmt.Errorf("get reward definitions: %w", err)
	}
	randRewardDef := rewardDefinitions[rand.Intn(len(rewardDefinitions))]

	var tokens *entity.JWT
	err = s.transactor.WithTx(ctx, func(ctx context.Context) error {
		var userID string
		userID, err = s.userRepository.AddUser(ctx, user)
		if err != nil {
			return err
		}

		tokens, err = s.generateTokens(ctx, userID)
		if err != nil {
			return err
		}

		_, err = s.rewardService.GrantReward(ctx, userID, randRewardDef.Code, "Награда за регистрацию")
		if err != nil {
			return err
		}

		err = s.streakRepository.UpdateStreak(ctx, &entity.Streak{
			UserID:         userID,
			CurrentStreak:  1,
			LastActiveDate: businessDate,
		})

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return tokens, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*entity.JWT, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return nil, entity.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("login: get user by email: %w", err)
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		return nil, entity.ErrInvalidCredentials
	}

	tokens, err := s.generateTokens(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("login: issue tokens: %w", err)
	}

	return tokens, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (*entity.JWT, error) {
	refreshTokenHash := hashToken(refreshToken)

	var tokens *entity.JWT
	err := s.transactor.WithTx(ctx, func(ctx context.Context) error {
		storedToken, err := s.tokenRepository.ConsumeRefreshToken(ctx, refreshTokenHash)
		if err != nil {
			if errors.Is(err, entity.ErrRefreshTokenNotFound) {
				return entity.ErrInvalidRefreshToken
			}
			return err
		}

		if time.Now().UTC().After(storedToken.ExpiresAt) {
			return entity.ErrInvalidRefreshToken
		}

		tokens, err = s.generateTokens(ctx, storedToken.UserID)
		return err
	})
	if err != nil {
		if errors.Is(err, entity.ErrInvalidRefreshToken) {
			return nil, entity.ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("refresh: %w", err)
	}

	return tokens, nil
}

func (s *authService) Logout(ctx context.Context, userID, refreshToken string) error {
	refreshTokenHash := hashToken(refreshToken)

	if err := s.tokenRepository.DeleteSession(ctx, userID, refreshTokenHash); err != nil {
		return fmt.Errorf("logout: revoke refresh token: %w", err)
	}
	return nil
}

func (s *authService) ValidateAccessToken(_ context.Context, token string) (string, error) {
	const jwtPartsAmount = 3
	parts := strings.Split(token, ".")
	if len(parts) != jwtPartsAmount {
		return "", entity.ErrInvalidAccessToken
	}

	signingInput := parts[0] + "." + parts[1]
	expectedSignature := signJWT(signingInput, s.cfg.JWTSecret)
	if subtle.ConstantTimeCompare([]byte(parts[2]), []byte(expectedSignature)) != 1 {
		return "", entity.ErrInvalidAccessToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", entity.ErrInvalidAccessToken
	}

	var claims accessTokenClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", entity.ErrInvalidAccessToken
	}

	if claims.Sub == "" || claims.Exp <= time.Now().UTC().Unix() {
		return "", entity.ErrInvalidAccessToken
	}

	return claims.Sub, nil
}

func dateOnly(t time.Time) time.Time {
	t = t.UTC()

	return time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0, time.UTC,
	)
}
