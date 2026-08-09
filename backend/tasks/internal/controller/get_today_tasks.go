package controller

import (
	"context"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres"
)

type GetTodayTasksQuery struct {
	UserID string
}

type TodayTasksResponse struct {
	Date  string    `json:"date"`
	Items []TaskDTO `json:"items"`
}

type GetTodayTasksHandler struct {
	taskRepo *postgres.TaskRepository
}

func NewGetTodayTasksHandler(repo *postgres.TaskRepository) *GetTodayTasksHandler {
	return &GetTodayTasksHandler{taskRepo: repo}
}

func (h *GetTodayTasksHandler) Handle(ctx context.Context, query GetTodayTasksQuery) (*TodayTasksResponse, error) {
	tasks, err := h.taskRepo.ListForUser(ctx, query.UserID)
	if err != nil {
		return nil, err
	}

	items := make([]TaskDTO, 0, len(tasks))
	for _, task := range tasks {
		items = append(items, TaskDTO{
			ID:          task.Task.ID,
			Title:       task.Task.Title,
			Description: task.Task.Description,
			RewardCoins: task.Task.RewardCoins,
			RewardXP:    task.Task.RewardXP,
			Status:      entity.Status(task.Status),
			CompletedAt: task.CompletedAt,
			TaskType:    task.Task.Type,
		})
	}

	return &TodayTasksResponse{
		Date:  time.Now().UTC().Format("2006-01-02"),
		Items: items,
	}, nil
}
