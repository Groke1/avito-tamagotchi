package controller

import (
	"context"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	sqlctask "github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres/sqlc"
)

type fakeTaskRepository struct {
	findByIDForUserFn    func(ctx context.Context, userID, taskID string) (*entity.UserTask, error)
	listForUserFn        func(ctx context.Context, userID string) ([]entity.UserTask, error)
	completeTaskFn       func(ctx context.Context, userID, taskID string) (*entity.UserTask, error)
	listCompletedTodayFn func(ctx context.Context, userID string) ([]sqlctask.GetTodayCompletedTasksForUserRow, error)
}

func (f *fakeTaskRepository) FindByIDForUser(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
	return f.findByIDForUserFn(ctx, userID, taskID)
}

func (f *fakeTaskRepository) ListForUser(ctx context.Context, userID string) ([]entity.UserTask, error) {
	return f.listForUserFn(ctx, userID)
}

func (f *fakeTaskRepository) CompleteTask(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
	return f.completeTaskFn(ctx, userID, taskID)
}

func (f *fakeTaskRepository) ListCompletedToday(ctx context.Context, userID string) ([]sqlctask.GetTodayCompletedTasksForUserRow, error) {
	return f.listCompletedTodayFn(ctx, userID)
}

type fakeCoinsClient struct {
	updateCoinsFn func(ctx context.Context, req UpdateCoinsRequest) (*UpdateCoinsResponse, error)
	calls         []UpdateCoinsRequest
}

func (f *fakeCoinsClient) UpdateCoins(ctx context.Context, req UpdateCoinsRequest) (*UpdateCoinsResponse, error) {
	f.calls = append(f.calls, req)
	return f.updateCoinsFn(ctx, req)
}

type fakeXPClient struct {
	updateXPFn func(ctx context.Context, req UpdateXPRequest) error
	calls      []UpdateXPRequest
}

func (f *fakeXPClient) UpdateXP(ctx context.Context, req UpdateXPRequest) error {
	f.calls = append(f.calls, req)
	return f.updateXPFn(ctx, req)
}

type fakeTransactor struct {
	beginErr error
}

func (t *fakeTransactor) WithTx(ctx context.Context, f func(ctx context.Context) error) error {
	if t.beginErr != nil {
		return t.beginErr
	}
	return f(ctx)
}
