package entity

import "time"

type JWT struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}
