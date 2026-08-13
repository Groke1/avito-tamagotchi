package repository

import "time"

type tripEvent struct {
	ID          int32  `db:"id"`
	Description string `db:"description"`
	IsNegative  bool   `db:"is_negative"`
}

type petTrip struct {
	ID          int64     `db:"id"`
	PetID       int64     `db:"pet_id"`
	UserID      string    `db:"user_id"`
	Location    string    `db:"location"`
	RewardXP    *int32    `db:"reward_xp"`
	RewardCoins *int32    `db:"reward_coins"`
	RewardCode  *string   `db:"reward_code"`
	Story       string    `db:"story"`
	Status      string    `db:"status"`
	StartedAt   time.Time `db:"started_at"`
	EndedAt     time.Time `db:"ended_at"`
}
