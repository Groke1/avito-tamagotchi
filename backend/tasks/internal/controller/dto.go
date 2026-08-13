package controller

import "time"

type (
	UpdateCoinsRequest struct {
		UserID string `json:"user_id"`
		Delta  int    `json:"delta"`
	}
	NotifyActionRequest struct {
		UserID     string    `json:"user_id"`
		OccurredAt time.Time `json:"occurred_at"`
	}
	UpdateCoinsResponse struct {
		UserID string `json:"user_id"`
		Coins  int64  `json:"coins"`
	}
	UpdateXPRequest struct {
		UserID string `json:"user_id"`
		XP     int    `json:"xp"`
	}
	CompleteTaskQuery struct {
		TaskID string `json:"task_id"`
		UserID string `json:"user_id"`
	}
	AwardedDTO struct {
		Coins int   `json:"coins"`
		XP    int64 `json:"xp"`
	}
	BalanceDTO struct {
		Coins int64 `json:"coins"`
	}
	CompleteTaskResponse struct {
		Task    TaskDTO     `json:"task"`
		Awarded AwardedDTO  `json:"awarded"`
		Balance *BalanceDTO `json:"balance"`
		// TODO: here must be also a PetDTO, but they are not implemented yet
	}
)
