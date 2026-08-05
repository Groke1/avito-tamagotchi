package entity

import (
	"errors"
	"time"
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")

	ErrUserNotFound         = errors.New("user not found")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrInsufficientCoins    = errors.New("insufficient coins")
)

type User struct {
	ID           string
	Username     string
	Email        string
	Password     string
	PasswordHash string
	Coins        uint64
}

type JWT struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}
