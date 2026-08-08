package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

func (s *userService) UpdateStreak(ctx context.Context, userID, occurredAt string) error {
	eventTime, err := time.Parse(time.RFC3339, occurredAt)
	if err != nil {
		return fmt.Errorf("parse occurredAt: %w", err)
	}

	businessDate := dateOnly(eventTime)

	var isStreakChanged bool
	var streak *entity.Streak
	err = s.transactor.WithTx(ctx, func(ctx context.Context) error {
		streak, err = s.streakRepository.GetStreakByUserIDForUpdate(ctx, userID)

		if err != nil {
			if !errors.Is(err, entity.ErrUserNotFound) {
				return err
			}
			streak = &entity.Streak{
				UserID:         userID,
				CurrentStreak:  1,
				LastActiveDate: businessDate,
			}
			isStreakChanged = true
		} else {
			isStreakChanged = updateStreak(streak, businessDate)
		}

		if isStreakChanged {
			err = s.streakRepository.UpdateStreak(ctx, streak)
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
		if err = s.petClient.SendDailyBonus(ctx, streak.UserID, streak.CurrentStreak); err != nil {
			return fmt.Errorf("send daily bonus: %w", err)
		}
	}
	return nil
}

func dateOnly(t time.Time) time.Time {
	t = t.UTC()

	return time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0, time.UTC,
	)
}

func updateStreak(streak *entity.Streak, currentDate time.Time) bool {
	if sameDate(streak.LastActiveDate, currentDate) {
		return false
	}

	yesterday := currentDate.AddDate(0, 0, -1)

	if sameDate(streak.LastActiveDate, yesterday) {
		streak.CurrentStreak++
	} else {
		streak.CurrentStreak = 1
	}

	streak.LastActiveDate = currentDate

	return true
}

func sameDate(a, b time.Time) bool {
	return a.Year() == b.Year() &&
		a.Month() == b.Month() &&
		a.Day() == b.Day()
}
