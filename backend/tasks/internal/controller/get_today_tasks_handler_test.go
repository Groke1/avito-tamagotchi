package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
)

func TestGetTodayTasksHandler_Handle_Success(t *testing.T) {
	tasks := []entity.UserTask{
		entity.ReconstructUserTask(entity.Task{ID: "t1", Title: "A", RewardCoins: 10, RewardXP: 20, Type: "Продажи"}, entity.StatusInProgress, nil),
		entity.ReconstructUserTask(entity.Task{ID: "t2", Title: "B", RewardCoins: 30, RewardXP: 40, Type: "Поиск"}, entity.StatusInProgress, nil),
	}

	repo := &fakeTaskRepository{
		listForUserFn: func(ctx context.Context, userID string) ([]entity.UserTask, error) {
			if userID != "user-1" {
				t.Fatalf("unexpected userID: %s", userID)
			}
			return tasks, nil
		},
	}

	handler := NewGetTodayTasksHandler(repo)
	resp, err := handler.Handle(context.Background(), GetTodayTasksQuery{UserID: "user-1"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}
	if resp.Items[0].ID != "t1" || resp.Items[1].ID != "t2" {
		t.Fatalf("unexpected items order/content: %+v", resp.Items)
	}
	if _, err := time.Parse("2006-01-02", resp.Date); err != nil {
		t.Fatalf("expected valid date format, got %q: %v", resp.Date, err)
	}
}

func TestGetTodayTasksHandler_Handle_EmptyList(t *testing.T) {
	repo := &fakeTaskRepository{
		listForUserFn: func(ctx context.Context, userID string) ([]entity.UserTask, error) {
			return []entity.UserTask{}, nil
		},
	}

	handler := NewGetTodayTasksHandler(repo)
	resp, err := handler.Handle(context.Background(), GetTodayTasksQuery{UserID: "user-1"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Items) != 0 {
		t.Fatalf("expected empty items, got %+v", resp.Items)
	}
}

func TestGetTodayTasksHandler_Handle_RepositoryError(t *testing.T) {
	repoErr := errors.New("db unavailable")
	repo := &fakeTaskRepository{
		listForUserFn: func(ctx context.Context, userID string) ([]entity.UserTask, error) {
			return nil, repoErr
		},
	}

	handler := NewGetTodayTasksHandler(repo)
	_, err := handler.Handle(context.Background(), GetTodayTasksQuery{UserID: "user-1"})

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error to propagate, got %v", err)
	}
}
