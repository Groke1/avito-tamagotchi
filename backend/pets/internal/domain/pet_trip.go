package domain

import "time"

type TripStatus string

const (
	InProgress      TripStatus = "in_progress"      // Воркер-1 создал трип
	PendingDelivery TripStatus = "pending_delivery" // Если воркер-2 не смог доставить по вебсокету
	Delivered       TripStatus = "delivered"        // Трип доставлен (либо воркером-2 по вебсокету, либо потом по ручке)
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
