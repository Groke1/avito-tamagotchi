package entity

import "time"

type EventType string

const (
	NewReward    EventType = "new_reward"
	StreakReward EventType = "streak_reward"
)

type UserEvent struct {
	ID           int64
	UserID       string
	Type         EventType
	XP           int32
	Coins        int32
	Streak       *int32
	UserRewardID *string
	CreatedAt    time.Time
}

type UserEventDetails struct {
	Type      EventType
	XP        int32
	Coins     int32
	Streak    *int32
	Reward    *UserReward
	CreatedAt time.Time
}
