package controller

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres"
)

type (
	coinsClient interface {
		UpdateCoins(ctx context.Context, req UpdateCoinsRequest) (*UpdateCoinsResponse, error)
		NotifyActionDone(ctx context.Context, req NotifyActionRequest) error
	}
	xpClient interface {
		UpdateXP(ctx context.Context, req UpdateXPRequest) error
	}
	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) error
	}
)

type CompleteTaskHandler struct {
	taskRepo    postgres.TaskRepoInterface
	coinsClient coinsClient
	xpClient    xpClient
	transactor  transactor
}

func NewCompleteTaskHandler(repo postgres.TaskRepoInterface, coinsClient coinsClient, xpClient xpClient, transactor transactor) *CompleteTaskHandler {
	return &CompleteTaskHandler{
		taskRepo:    repo,
		coinsClient: coinsClient,
		xpClient:    xpClient,
		transactor:  transactor,
	}
}

func (h *CompleteTaskHandler) Handle(ctx context.Context, query CompleteTaskQuery) (*CompleteTaskResponse, error) {
	var resp *CompleteTaskResponse
	err := h.transactor.WithTx(ctx, func(ctx context.Context) error {
		completedTask, err := h.taskRepo.CompleteTask(ctx, query.UserID, query.TaskID)
		if err != nil {
			return err
		}

		coinsResp, err := h.coinsClient.UpdateCoins(ctx, UpdateCoinsRequest{
			UserID: query.UserID,
			Delta:  completedTask.Task.RewardCoins,
		})
		if err != nil {
			return fmt.Errorf("update coins: %w", ErrUserServiceUnavailable)
		}

		err = h.xpClient.UpdateXP(ctx, UpdateXPRequest{
			UserID: query.UserID,
			XP:     int(completedTask.Task.RewardXP),
		})
		if err != nil {
			return fmt.Errorf("update xp: %w", ErrPetServiceUnavailable)
		}

		completedAt := completedTask.CompletedAt
		if completedAt == nil {
			now := time.Now().UTC()
			completedAt = &now
		}

		resp = &CompleteTaskResponse{
			Task: TaskDTO{
				ID:          completedTask.Task.ID,
				Title:       completedTask.Task.Title,
				Description: completedTask.Task.Description,
				RewardCoins: completedTask.Task.RewardCoins,
				RewardXP:    completedTask.Task.RewardXP,
				Status:      entity.StatusCompleted,
				CompletedAt: completedAt,
				TaskType:    completedTask.Task.Type,
			},
			Awarded: AwardedDTO{
				Coins: completedTask.Task.RewardCoins,
				XP:    completedTask.Task.RewardXP,
			},
			Balance: &BalanceDTO{Coins: coinsResp.Coins},
		}

		err = h.coinsClient.NotifyActionDone(ctx, NotifyActionRequest{
			UserID:     query.UserID,
			OccurredAt: time.Now().UTC(),
		})
		if err != nil {
			log.Println("AA", err)
			return fmt.Errorf("notify action: %w", ErrUserServiceUnavailable)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
