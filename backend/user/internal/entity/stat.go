package entity

import (
	"time"
)

type DailyStat struct {
	UserID  string        `json:"user_id"`
	Streak  int32         `json:"streak"`
	Tasks   []TasksStat   `json:"tasks"`
	Pet     PetStat       `json:"pet"`
	Rewards []RewardsStat `json:"rewards"`
}

type PetStat struct {
	DailyGainedXP int64 `json:"daily_gained_xp"`
}

type TasksStat struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	RewardCoins  int32     `json:"reward_coins"`
	RewardXp     int64     `json:"reward_xp"`
	FinishedDesc string    `json:"finished_desc"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RewardsStat struct {
	PromoCode    string    `json:"promo_code"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	Description  string    `json:"description"`
	FinishedDesc string    `json:"finished_desc"`
	CreatedTime  time.Time `json:"created_time"`
}
