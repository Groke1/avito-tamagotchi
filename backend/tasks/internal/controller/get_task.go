package controller

import (
	"context"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres"
)

type GetTaskQuery struct {
	TaskId string
	UserId string
}

type TaskDTO struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	RewardCoins int           `json:"reward_coins"`
	RewardXP    int64         `json:"reward_xp"`
	Status      entity.Status `json:"status"`
	CompletedAt *time.Time    `json:"completed_at"`
}
type GetTaskHandler struct {
	taskRepo *postgres.TaskRepository
}

func NewGetTaskHandler(repo *postgres.TaskRepository) *GetTaskHandler {
	return &GetTaskHandler{taskRepo: repo}
}

func (h *GetTaskHandler) Handle(ctx context.Context, query GetTaskQuery) (*TaskDTO, error) {
	task, err := h.taskRepo.FindByIDForUser(ctx, query.UserId, query.TaskId)
	if err != nil {
		return nil, err
	}
	return &TaskDTO{
		ID:          task.Task.ID,
		Title:       task.Task.Title,
		Description: task.Task.Description,
		RewardCoins: task.Task.RewardCoins,
		RewardXP:    task.Task.RewardXP,
		Status:      task.Status,
		CompletedAt: task.CompletedAt,
	}, nil
}
