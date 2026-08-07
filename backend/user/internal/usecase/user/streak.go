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
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return fmt.Errorf("time.LoadLocation: %w", err)
	}
	businessDate := dateOnly(eventTime, loc)

	var isStreakChanged bool
	err = s.transactor.WithTx(ctx, func(ctx context.Context) error {
		var streak *entity.Streak
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
		// TODO: Send to pet
	}
	return nil
}

func dateOnly(t time.Time, loc *time.Location) time.Time {
	localTime := t.In(loc)

	return time.Date(
		localTime.Year(), localTime.Month(), localTime.Day(),
		0, 0, 0, 0, loc,
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
