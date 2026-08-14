package domain

import "time"

type TripStatus string

const (
	InProgress      TripStatus = "in_progress"
	PendingDelivery TripStatus = "pending_delivery"
	Delivered       TripStatus = "delivered"
)

type PetTrip struct {
	ID          int64
	PetID       int64
	UserID      string
	Location    string
	RewardXP    *int32
	RewardCoins *int32
	RewardCode  *string
	Story       string
	Status      TripStatus
	StartedAt   time.Time
	EndedAt     time.Time
}
