package controller

import (
	"context"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres"
	sqlctask "github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres/sqlc"
)

type GetCompletedTasksQuery struct {
	UserID string
}
type GetCompletedTasksResponse struct {
	UserID string                                      `json:"user_id"`
	Items  []sqlctask.GetTodayCompletedTasksForUserRow `json:"items"`
}

type GetCompletedTasksHandler struct {
	taskRepo postgres.TaskRepoInterface
}

func NewGetCompletedTasksHandler(repo postgres.TaskRepoInterface) *GetCompletedTasksHandler {
	return &GetCompletedTasksHandler{taskRepo: repo}
}

func (h *GetCompletedTasksHandler) Handle(ctx context.Context, query GetCompletedTasksQuery) (*GetCompletedTasksResponse, error) {
	tasks, err := h.taskRepo.ListCompletedToday(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	return &GetCompletedTasksResponse{
		UserID: query.UserID,
		Items:  tasks,
	}, nil
}
