package clients

import "time"

type UserNamesResponse struct {
	Users []User `json:"users"`
}

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
}

type RewardResponse struct {
	RewardID     string     `json:"reward_id"`
	PromoCode    string     `json:"promo_code"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	ExpiresAt    string     `json:"expires_at"`
	EarnedReason string     `json:"earned_reason"`
	RedeemedAt   *time.Time `Json:"redeemed_at"`
}

type RewardDescriptionResponse struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
