package websocket

import "github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"

type PetResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	XP          int    `json:"xp"`
	NextLevelXP int    `json:"next_level_xp"`
	Satiety     int    `json:"satiety"`
	Happiness   int    `json:"happiness"`
}

type Event struct {
	EventType domain.WsEvent `json:"event_type"`
	Payload   any            `json:"payload"`
}
