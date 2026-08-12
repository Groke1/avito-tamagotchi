package dates

import "time"

func DateOnly(t time.Time) time.Time {
	t = t.UTC()

	return time.Date(
		t.Year(), t.Month(), t.Day(),
		0, 0, 0, 0, time.UTC,
	)
}

func SameDate(a, b time.Time) bool {
	return a.Year() == b.Year() &&
		a.Month() == b.Month() &&
		a.Day() == b.Day()
}

func DayBounds(t time.Time) (time.Time, time.Time) {
	t = t.UTC()

	from := DateOnly(t)

	to := DayAfter(from)

	return from, to
}

func DayAfter(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

func DayBefore(t time.Time) time.Time {
	return t.AddDate(0, 0, -1)
}
