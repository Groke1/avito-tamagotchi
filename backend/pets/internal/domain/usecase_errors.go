package domain

import (
	"errors"
	"time"
)

var (
	ErrPetNotFound       = errors.New("pet not found")
	ErrPetAlreadyExists  = errors.New("pet already exists")
	ErrUnavailableAction = errors.New("action is not allowed")
	ErrPetIsFull         = errors.New("pet is already full")
	ErrPetIsTooHappy     = errors.New("pet is too happy")
	ErrPetIsTooHungry    = errors.New("pet is too hungry")
	ErrUserNotFound      = errors.New("user not found")
	ErrNotEnoguhCoins    = errors.New("insuffient amount of coins")
	ErrTripNotFound      = errors.New("trip not found")
	ErrNotPendingTrip    = errors.New("trip is not pending")
)

type ActionUnavailableError struct {
	RetryAfter time.Duration
}

func (e *ActionUnavailableError) Error() string {
	return "Action unavailable"
}
