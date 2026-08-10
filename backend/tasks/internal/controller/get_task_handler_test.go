package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
)

func TestGetTaskHandler_Handle_Success(t *testing.T) {
	task := entity.ReconstructUserTask(
		entity.Task{ID: "task-1", Title: "Закажи с доставкой", Description: "desc", RewardCoins: 100, RewardXP: 500, Type: "Доставка"},
		entity.StatusInProgress,
		nil,
	)

	repo := &fakeTaskRepository{
		findByIDForUserFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			if userID != "user-1" || taskID != "task-1" {
				t.Fatalf("unexpected args: userID=%s taskID=%s", userID, taskID)
			}
			return &task, nil
		},
	}

	handler := NewGetTaskHandler(repo)
	dto, err := handler.Handle(context.Background(), GetTaskQuery{TaskID: "task-1", UserID: "user-1"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dto.ID != "task-1" || dto.Status != entity.StatusInProgress || dto.RewardCoins != 100 {
		t.Fatalf("unexpected dto: %+v", dto)
	}
}

func TestGetTaskHandler_Handle_NotFound(t *testing.T) {
	repo := &fakeTaskRepository{
		findByIDForUserFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			return nil, entity.ErrTaskNotFound
		},
	}

	handler := NewGetTaskHandler(repo)
	dto, err := handler.Handle(context.Background(), GetTaskQuery{TaskID: "missing", UserID: "user-1"})

	if !errors.Is(err, entity.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
	if dto != nil {
		t.Fatalf("expected nil dto, got %+v", dto)
	}
}

func TestGetTaskHandler_Handle_RepositoryError(t *testing.T) {
	repoErr := errors.New("db unavailable")
	repo := &fakeTaskRepository{
		findByIDForUserFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			return nil, repoErr
		},
	}

	handler := NewGetTaskHandler(repo)
	_, err := handler.Handle(context.Background(), GetTaskQuery{TaskID: "task-1", UserID: "user-1"})

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error to propagate, got %v", err)
	}
}
