package entity

import (
	"time"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusRedeemed Status = "redeemed"
	StatusExpired  Status = "expired"
)

type RewardDefinition struct {
	ID          int32
	Code        string
	Name        string
	Description string
}

type UserReward struct {
	ID         string
	UserID     string
	PromoCode  string
	Status     Status
	Definition RewardDefinition
	ExpiresAt  *time.Time
	RedeemedAt *time.Time
}
