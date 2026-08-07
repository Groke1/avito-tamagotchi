package entity

import "time"

type Streak struct {
	UserID         string
	CurrentStreak  int32
	LastActiveDate time.Time
}
