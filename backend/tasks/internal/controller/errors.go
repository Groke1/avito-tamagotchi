package controller

import (
	"errors"
)

var (
	ErrUnauthorized           = errors.New("UNAUTHORIZED")
	ErrInvalidRequest         = errors.New("INVALID_REQUEST")
	ErrUserServiceUnavailable = errors.New("USER_SERVICE_UNAVAILABLE")
	ErrPetServiceUnavailable  = errors.New("PET_SERVICE_UNAVAILABLE")
	ErrInvalidToken           = errors.New("INVALID_TOKEN")
	ErrExpiredToken           = errors.New("EXPIRED_TOKEN")

	ErrInternal = errors.New("INTERNAL_ERROR")
)
