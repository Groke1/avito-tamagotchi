package controller

import (
	"context"
	"errors"
	"testing"

	sqlctask "github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres/sqlc"
)

func TestGetCompletedTasksHandler_Handle_Success(t *testing.T) {
	rows := make([]sqlctask.GetTodayCompletedTasksForUserRow, 2)

	repo := &fakeTaskRepository{
		listCompletedTodayFn: func(ctx context.Context, userID string) ([]sqlctask.GetTodayCompletedTasksForUserRow, error) {
			if userID != "user-1" {
				t.Fatalf("unexpected userID: %s", userID)
			}
			return rows, nil
		},
	}

	handler := NewGetCompletedTasksHandler(repo)
	resp, err := handler.Handle(context.Background(), GetCompletedTasksQuery{UserID: "user-1"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.UserID != "user-1" {
		t.Fatalf("expected user_id to be echoed back, got %q", resp.UserID)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(resp.Items))
	}
}

func TestGetCompletedTasksHandler_Handle_RepositoryError(t *testing.T) {
	repoErr := errors.New("db unavailable")
	repo := &fakeTaskRepository{
		listCompletedTodayFn: func(ctx context.Context, userID string) ([]sqlctask.GetTodayCompletedTasksForUserRow, error) {
			return nil, repoErr
		},
	}

	handler := NewGetCompletedTasksHandler(repo)
	_, err := handler.Handle(context.Background(), GetCompletedTasksQuery{UserID: "user-1"})

	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repository error to propagate, got %v", err)
	}
}
