package controller

import (
	"context"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres"
)

type CompleteTaskQuery struct {
	TaskID string
	UserID string
}

type AwardedDTO struct {
	Coins int   `json:"coins"`
	XP    int64 `json:"xp"`
}

type CompleteTaskResponse struct {
	Task    TaskDTO    `json:"task"`
	Awarded AwardedDTO `json:"awarded"`
	// TODO: here must be also a PetDTO and BalanceDTO, but they are not implemented yet
}

type CompleteTaskHandler struct {
	taskRepo *postgres.TaskRepository
}

func NewCompleteTaskHandler(repo *postgres.TaskRepository) *CompleteTaskHandler {
	return &CompleteTaskHandler{taskRepo: repo}
}

func (h *CompleteTaskHandler) Handle(ctx context.Context, query CompleteTaskQuery) (*CompleteTaskResponse, error) {
	completedTask, err := h.taskRepo.CompleteTask(ctx, query.UserID, query.TaskID)
	if err != nil {
		return nil, err
	}

	completedAt := completedTask.CompletedAt
	if completedAt == nil {
		now := time.Now().UTC()
		completedAt = &now
	}

	return &CompleteTaskResponse{
		Task: TaskDTO{
			ID:          completedTask.Task.ID,
			Title:       completedTask.Task.Title,
			Description: completedTask.Task.Description,
			RewardCoins: completedTask.Task.RewardCoins,
			RewardXP:    completedTask.Task.RewardXP,
			Status:      entity.StatusCompleted,
			CompletedAt: completedAt,
		},
		Awarded: AwardedDTO{
			Coins: completedTask.Task.RewardCoins,
			XP:    completedTask.Task.RewardXP,
		},
	}, nil
}
