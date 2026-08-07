package entity

import "errors"

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrInvalidCredentials   = errors.New("invalid credentials")

	ErrUserNotFound          = errors.New("user not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")

	ErrInsufficientCoins = errors.New("insufficient coins")

	ErrRewardNotFound           = errors.New("reward not found")
	ErrRewardUnavailable        = errors.New("reward unavailable")
	ErrPromoCodeAlreadyExists   = errors.New("promo code already exists")
	ErrRewardDefinitionNotFound = errors.New("reward definition not found")
)
