package domain

import "time"

type Pet struct {
	ID          int64
	UserID      int64
	Name        string
	Level       int
	XP          int
	NextLevelXP int
	Satiety     int
	Happiness   int
	CreatedAt   time.Time
}
