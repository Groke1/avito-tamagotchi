package user

import (
	"context"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.pkg/dates"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

func (s *userService) UpdateStreak(ctx context.Context, userID, occurredAt string) error {
	eventTime, err := time.Parse(time.RFC3339, occurredAt)
	if err != nil {
		return fmt.Errorf("parse occurredAt: %w", err)
	}

	businessDate := dates.DateOnly(eventTime)

	var isStreakChanged bool
	var bonusCoins int32
	var streak *entity.Streak
	err = s.transactor.WithTx(ctx, func(ctx context.Context) error {
		streak, err = s.streakRepository.GetStreakByUserIDForUpdate(ctx, userID)

		if err != nil {
			return err
		}

		isStreakChanged = updateStreak(streak, businessDate)
		bonusCoins = getBonusCoins(streak.CurrentStreak)

		if isStreakChanged {
			err = s.streakRepository.UpdateStreak(ctx, streak)
			if err != nil {
				return err
			}
			_, err = s.userRepository.UpdateCoins(ctx, userID, int64(bonusCoins))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("update streak: %w", err)
	}

	if isStreakChanged {
		err = s.petClient.SendDailyBonus(ctx, streak.UserID, streak.CurrentStreak)
		if err != nil {
			return fmt.Errorf("send daily bonus: %w", err)
		}

		err = s.eventRepository.AddUserEvent(ctx, entity.UserEvent{
			UserID: streak.UserID,
			Type:   entity.StreakReward,
			Coins:  bonusCoins,
			Streak: &streak.CurrentStreak,
		})
		if err != nil {
			return fmt.Errorf("add streak event: %w", err)
		}
	}
	return nil
}

func updateStreak(streak *entity.Streak, currentDate time.Time) bool {
	if dates.SameDate(streak.LastActiveDate, currentDate) {
		return false
	}

	if dates.SameDate(streak.LastActiveDate, dates.DayBefore(currentDate)) {
		streak.CurrentStreak++
	} else {
		streak.CurrentStreak = 1
	}

	streak.LastActiveDate = currentDate

	return true
}

func getBonusCoins(currentStreak int32) int32 {
	const week = 7
	const weekCoins = 20
	if currentStreak%week != 0 {
		return 0
	}
	return currentStreak / week * weekCoins
}
