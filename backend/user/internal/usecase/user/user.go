package user

import (
	"context"
	"fmt"
	"time"

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
		GetStreakByUserID(ctx context.Context, userID string) (*entity.Streak, error)
		GetStreakByUserIDForUpdate(ctx context.Context, userID string) (*entity.Streak, error)
		UpdateStreak(ctx context.Context, streak *entity.Streak) error
	}

	rewardRepository interface {
		GetRewardsByUserIDAndPeriod(ctx context.Context, userID string, from, to time.Time) ([]entity.UserReward, error)
	}

	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) error
	}

	petService interface {
		SendDailyBonus(ctx context.Context, userID string, streak int32) error
	}
)

type userService struct {
	userRepository   userRepository
	streakRepository streakRepository
	rewardRepository rewardRepository
	transactor       transactor
	petService       petService
}

func NewUserService(
	userRepository userRepository,
	streakRepository streakRepository,
	rewardRepository rewardRepository,
	transactor transactor,
	petService petService,
) *userService {
	return &userService{
		userRepository:   userRepository,
		streakRepository: streakRepository,
		rewardRepository: rewardRepository,
		transactor:       transactor,
		petService:       petService,
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

func (s *userService) GetDailyStat(ctx context.Context, userID string) (*entity.DailyStat, error) {
	_, err := s.streakRepository.GetStreakByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get daily stat: %w", err)
	}
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return nil, fmt.Errorf("get daily stat: %w", err)
	}
	from, to := dayBounds(time.Now(), loc)

	_, err = s.rewardRepository.GetRewardsByUserIDAndPeriod(ctx, userID, from, to)

	if err != nil {
		return nil, fmt.Errorf("get daily stat: %w", err)
	}

	// Todo: запросы в другие сервисы
	panic("implement me")
}

func dayBounds(now time.Time, loc *time.Location) (time.Time, time.Time) {
	localNow := now.In(loc)

	from := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)

	to := from.AddDate(0, 0, 1)

	return from, to
}
