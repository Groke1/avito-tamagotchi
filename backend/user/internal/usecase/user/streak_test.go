package user

import (
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

func TestUpdateStreak_SameDay_NoChange(t *testing.T) {
	today := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	streak := &entity.Streak{CurrentStreak: 5, LastActiveDate: today}

	changed := updateStreak(streak, today)

	if changed {
		t.Fatalf("expected no change for the same business date")
	}
	if streak.CurrentStreak != 5 {
		t.Errorf("expected streak to stay 5, got %d", streak.CurrentStreak)
	}
}

func TestUpdateStreak_ConsecutiveDay_Increments(t *testing.T) {
	yesterday := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	today := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	streak := &entity.Streak{CurrentStreak: 5, LastActiveDate: yesterday}

	changed := updateStreak(streak, today)

	if !changed {
		t.Fatalf("expected a change")
	}
	if streak.CurrentStreak != 6 {
		t.Errorf("expected streak to become 6, got %d", streak.CurrentStreak)
	}
	if !streak.LastActiveDate.Equal(today) {
		t.Errorf("expected LastActiveDate to be updated to today")
	}
}

func TestUpdateStreak_GapInDays_ResetsToOne(t *testing.T) {
	longAgo := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	today := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	streak := &entity.Streak{CurrentStreak: 5, LastActiveDate: longAgo}

	changed := updateStreak(streak, today)

	if !changed {
		t.Fatalf("expected a change")
	}
	if streak.CurrentStreak != 1 {
		t.Errorf("expected streak to reset to 1, got %d", streak.CurrentStreak)
	}
}

func TestUpdateStreak_FirstEverActivity_ZeroValueDate(t *testing.T) {
	var zero time.Time
	today := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	streak := &entity.Streak{CurrentStreak: 0, LastActiveDate: zero}

	changed := updateStreak(streak, today)

	if !changed {
		t.Fatalf("expected a change")
	}
	if streak.CurrentStreak != 1 {
		t.Errorf("expected streak to start at 1, got %d", streak.CurrentStreak)
	}
}

func TestSameDate(t *testing.T) {
	a := time.Date(2026, 8, 10, 5, 0, 0, 0, time.UTC)
	b := time.Date(2026, 8, 10, 23, 59, 0, 0, time.UTC)
	c := time.Date(2026, 8, 11, 0, 0, 0, 0, time.UTC)

	if !sameDate(a, b) {
		t.Errorf("expected same calendar day regardless of time-of-day")
	}
	if sameDate(a, c) {
		t.Errorf("expected different calendar days to not match")
	}
}

func TestDayBounds(t *testing.T) {
	t0 := time.Date(2026, 8, 10, 15, 45, 0, 0, time.UTC)

	from, to := dayBounds(t0)

	wantFrom := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	wantTo := time.Date(2026, 8, 11, 0, 0, 0, 0, time.UTC)
	if !from.Equal(wantFrom) {
		t.Errorf("expected from=%v, got %v", wantFrom, from)
	}
	if !to.Equal(wantTo) {
		t.Errorf("expected to=%v, got %v", wantTo, to)
	}
}
