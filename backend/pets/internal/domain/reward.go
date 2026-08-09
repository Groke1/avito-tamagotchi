package domain

import "time"

type Reward struct {
	ID          string
	PromoCode   string
	Name        string
	Description string
	Status      string
	ExpiresAt   string
	RedeemedAt  *time.Time
}
