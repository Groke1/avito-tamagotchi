package domain

import "time"

type Trip struct {
	ID          int64
	PetID       int64
	UserID      string
	StartedAt   time.Time
	EndedAt     time.Time
	RewardCoins int
	RewardXp    int
	Status      string
}
