package user

import (
	"context"
	"fmt"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

//go:generate mockgen -source=user.go -destination=mocks/user_mocks.go -package=mocks

type (
	userRepository interface {
		GetUserByID(ctx context.Context, id string) (*entity.User, error)
		GetUsersByIDs(ctx context.Context, ids []string) ([]entity.User, error)
		UpdateCoins(ctx context.Context, userID string, coins int64) (*entity.User, error)
	}

	streakRepository interface {
		GetStreakByUserIDForUpdate(ctx context.Context, userID string) (*entity.Streak, error)
		UpdateStreak(ctx context.Context, streak *entity.Streak) error
	}

	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) error
	}
)

type userService struct {
	userRepository   userRepository
	streakRepository streakRepository
	transactor       transactor
}

func NewUserService(
	userRepository userRepository,
	streakRepository streakRepository,
	transactor transactor,
) *userService {
	return &userService{
		userRepository:   userRepository,
		streakRepository: streakRepository,
		transactor:       transactor,
	}
}

func (s *userService) Profile(ctx context.Context, userID string) (*entity.User, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("profile: get user by id: %w", err)
	}

	return user, nil
}

func (s *userService) GetUsers(ctx context.Context, userIDs []string) ([]entity.User, error) {
	users, err := s.userRepository.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	return users, nil
}

func (s *userService) UpdateCoins(ctx context.Context, userID string, deltaCoins int64) (*entity.User, error) {
	user, err := s.userRepository.UpdateCoins(ctx, userID, deltaCoins)
	if err != nil {
		return nil, fmt.Errorf("update coins: %w", err)
	}

	return user, nil
}
