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

	petClient interface {
		SendDailyBonus(ctx context.Context, userID string, streak int32) error
		GetPetDailyStat(ctx context.Context, userID string) (*entity.PetStat, error)
	}

	tasksClient interface {
		GetCompletedTasks(ctx context.Context, userID string) ([]entity.TasksStat, error)
	}
)

type userService struct {
	userRepository   userRepository
	streakRepository streakRepository
	rewardRepository rewardRepository
	transactor       transactor
	petClient        petClient
	tasksClient      tasksClient
}

func NewUserService(
	userRepository userRepository,
	streakRepository streakRepository,
	rewardRepository rewardRepository,
	transactor transactor,
	petClient petClient,
	tasksClient tasksClient,
) *userService {
	return &userService{
		userRepository:   userRepository,
		streakRepository: streakRepository,
		rewardRepository: rewardRepository,
		transactor:       transactor,
		petClient:        petClient,
		tasksClient:      tasksClient,
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
	streak, err := s.streakRepository.GetStreakByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get daily stat: %w", err)
	}

	from, to := dayBounds(time.Now())

	rewards, err := s.rewardRepository.GetRewardsByUserIDAndPeriod(ctx, userID, from, to)

	if err != nil {
		return nil, fmt.Errorf("get daily stat: %w", err)
	}

	rewardStats := make([]entity.RewardsStat, 0, len(rewards))

	for _, reward := range rewards {
		if reward.CreatedAt.Compare(from) >= 0 &&
			reward.CreatedAt.Compare(to) < 0 {
			rewardStats = append(rewardStats, entity.RewardsStat{
				PromoCode:    reward.PromoCode,
				Name:         reward.Definition.Name,
				Description:  reward.Definition.Description,
				FinishedDesc: reward.Definition.EarnedDescription,
				CreatedTime:  reward.CreatedAt,
			})
		}

		if reward.RedeemedAt != nil &&
			reward.RedeemedAt.Compare(from) >= 0 &&
			reward.RedeemedAt.Compare(to) < 0 {
			rewardStats = append(rewardStats, entity.RewardsStat{
				PromoCode:    reward.PromoCode,
				Name:         reward.Definition.Name,
				Description:  reward.Definition.Description,
				FinishedDesc: reward.Definition.EarnedDescription,
				CreatedTime:  *reward.RedeemedAt,
			})
		}
	}

	tasks, err := s.tasksClient.GetCompletedTasks(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get daily stat: %w", err)
	}

	pet, err := s.petClient.GetPetDailyStat(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get daily stat: %w", err)
	}

	return &entity.DailyStat{
		UserID:  userID,
		Streak:  streak.CurrentStreak,
		Rewards: rewardStats,
		Tasks:   tasks,
		Pet:     *pet,
	}, nil
}

func dayBounds(t time.Time) (time.Time, time.Time) {
	t = t.UTC()

	from := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	to := from.AddDate(0, 0, 1)

	return from, to
}
