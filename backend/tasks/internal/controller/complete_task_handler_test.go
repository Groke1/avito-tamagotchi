package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
)

func completedTask() *entity.UserTask {
	ut := entity.ReconstructUserTask(
		entity.Task{ID: "task-1", Title: "Закажи с доставкой", Description: "desc", RewardCoins: 100, RewardXP: 500, Type: "Доставка"},
		entity.StatusCompleted,
		nil,
	)
	return &ut
}

func TestCompleteTaskHandler_Handle_Success(t *testing.T) {
	repo := &fakeTaskRepository{
		completeTaskFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			if userID != "user-1" || taskID != "task-1" {
				t.Fatalf("unexpected args: userID=%s taskID=%s", userID, taskID)
			}
			return completedTask(), nil
		},
	}
	coins := &fakeCoinsClient{
		updateCoinsFn: func(ctx context.Context, req UpdateCoinsRequest) (*UpdateCoinsResponse, error) {
			return &UpdateCoinsResponse{UserID: req.UserID, Coins: 200}, nil
		},
	}
	xp := &fakeXPClient{
		updateXPFn: func(ctx context.Context, req UpdateXPRequest) error { return nil },
	}
	tx := &fakeTransactor{}

	handler := NewCompleteTaskHandler(repo, coins, xp, tx)
	resp, err := handler.Handle(context.Background(), CompleteTaskQuery{TaskID: "task-1", UserID: "user-1"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Task.Status != entity.StatusCompleted {
		t.Errorf("expected status completed, got %v", resp.Task.Status)
	}
	if resp.Task.CompletedAt == nil {
		t.Errorf("expected handler to fill in CompletedAt when repository didn't set it")
	}
	if resp.Awarded.Coins != 100 || resp.Awarded.XP != 500 {
		t.Errorf("unexpected awarded values: %+v", resp.Awarded)
	}
	if resp.Balance == nil || resp.Balance.Coins != 200 {
		t.Errorf("unexpected balance: %+v", resp.Balance)
	}

	if len(coins.calls) != 1 || coins.calls[0].Delta != 100 || coins.calls[0].UserID != "user-1" {
		t.Errorf("unexpected coins client calls: %+v", coins.calls)
	}
	if len(xp.calls) != 1 || xp.calls[0].XP != 500 || xp.calls[0].UserID != "user-1" {
		t.Errorf("unexpected xp client calls: %+v", xp.calls)
	}
}

func TestCompleteTaskHandler_Handle_PreservesRepositoryCompletedAt(t *testing.T) {
	fixed := time.Date(2026, 8, 9, 9, 30, 0, 0, time.UTC)
	repo := &fakeTaskRepository{
		completeTaskFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			ut := entity.ReconstructUserTask(
				entity.Task{ID: "task-1", RewardCoins: 100, RewardXP: 500},
				entity.StatusCompleted,
				&fixed,
			)
			return &ut, nil
		},
	}
	coins := &fakeCoinsClient{updateCoinsFn: func(ctx context.Context, req UpdateCoinsRequest) (*UpdateCoinsResponse, error) {
		return &UpdateCoinsResponse{Coins: 200}, nil
	}}
	xp := &fakeXPClient{updateXPFn: func(ctx context.Context, req UpdateXPRequest) error { return nil }}

	handler := NewCompleteTaskHandler(repo, coins, xp, &fakeTransactor{})
	resp, err := handler.Handle(context.Background(), CompleteTaskQuery{TaskID: "task-1", UserID: "user-1"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Task.CompletedAt == nil || !resp.Task.CompletedAt.Equal(fixed) {
		t.Errorf("expected repository-provided completedAt %v to be preserved, got %v", fixed, resp.Task.CompletedAt)
	}
}

func TestCompleteTaskHandler_Handle_TaskNotFound(t *testing.T) {
	repo := &fakeTaskRepository{
		completeTaskFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			return nil, entity.ErrTaskNotFound
		},
	}
	handler := NewCompleteTaskHandler(repo, &fakeCoinsClient{}, &fakeXPClient{}, &fakeTransactor{})

	_, err := handler.Handle(context.Background(), CompleteTaskQuery{TaskID: "missing", UserID: "user-1"})

	if !errors.Is(err, entity.ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestCompleteTaskHandler_Handle_AlreadyCompleted(t *testing.T) {
	repo := &fakeTaskRepository{
		completeTaskFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			return nil, entity.ErrTaskAlreadyCompleted
		},
	}
	handler := NewCompleteTaskHandler(repo, &fakeCoinsClient{}, &fakeXPClient{}, &fakeTransactor{})

	_, err := handler.Handle(context.Background(), CompleteTaskQuery{TaskID: "task-1", UserID: "user-1"})

	if !errors.Is(err, entity.ErrTaskAlreadyCompleted) {
		t.Fatalf("expected ErrTaskAlreadyCompleted, got %v", err)
	}
}

func TestCompleteTaskHandler_Handle_CoinsServiceUnavailable(t *testing.T) {
	repo := &fakeTaskRepository{
		completeTaskFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			return completedTask(), nil
		},
	}
	coins := &fakeCoinsClient{
		updateCoinsFn: func(ctx context.Context, req UpdateCoinsRequest) (*UpdateCoinsResponse, error) {
			return nil, errors.New("connection refused")
		},
	}
	xp := &fakeXPClient{
		updateXPFn: func(ctx context.Context, req UpdateXPRequest) error {
			t.Fatalf("xp client should not be called when coins update fails")
			return nil
		},
	}

	handler := NewCompleteTaskHandler(repo, coins, xp, &fakeTransactor{})
	_, err := handler.Handle(context.Background(), CompleteTaskQuery{TaskID: "task-1", UserID: "user-1"})

	if !errors.Is(err, ErrUserServiceUnavailable) {
		t.Fatalf("expected ErrUserServiceUnavailable, got %v", err)
	}
}

func TestCompleteTaskHandler_Handle_PetServiceUnavailable(t *testing.T) {
	repo := &fakeTaskRepository{
		completeTaskFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			return completedTask(), nil
		},
	}
	coins := &fakeCoinsClient{
		updateCoinsFn: func(ctx context.Context, req UpdateCoinsRequest) (*UpdateCoinsResponse, error) {
			return &UpdateCoinsResponse{Coins: 200}, nil
		},
	}
	xp := &fakeXPClient{
		updateXPFn: func(ctx context.Context, req UpdateXPRequest) error {
			return errors.New("pet service down")
		},
	}

	handler := NewCompleteTaskHandler(repo, coins, xp, &fakeTransactor{})
	_, err := handler.Handle(context.Background(), CompleteTaskQuery{TaskID: "task-1", UserID: "user-1"})

	if !errors.Is(err, ErrPetServiceUnavailable) {
		t.Fatalf("expected ErrPetServiceUnavailable, got %v", err)
	}
}

func TestCompleteTaskHandler_Handle_TransactionFailsToBegin(t *testing.T) {
	beginErr := errors.New("could not begin tx")
	repo := &fakeTaskRepository{
		completeTaskFn: func(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
			t.Fatalf("repository should not be called if the transaction never begins")
			return nil, nil
		},
	}

	handler := NewCompleteTaskHandler(repo, &fakeCoinsClient{}, &fakeXPClient{}, &fakeTransactor{beginErr: beginErr})
	_, err := handler.Handle(context.Background(), CompleteTaskQuery{TaskID: "task-1", UserID: "user-1"})

	if !errors.Is(err, beginErr) {
		t.Fatalf("expected begin error to propagate, got %v", err)
	}
}
