package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres"
)

type (
	coinsClient interface {
		UpdateCoins(ctx context.Context, req entity.UpdateCoinsRequest) (*entity.UpdateCoinsResponse, error)
	}
	transactor interface {
		WithTx(ctx context.Context, f func(ctx context.Context) error) error
	}
)

type (
	CompleteTaskQuery struct {
		TaskID string
		UserID string
	}
)
type AwardedDTO struct {
	Coins int   `json:"coins"`
	XP    int64 `json:"xp"`
}
type BalanceDTO struct {
	Coins int64 `json:"coins"`
}

type CompleteTaskResponse struct {
	Task    TaskDTO     `json:"task"`
	Awarded AwardedDTO  `json:"awarded"`
	Balance *BalanceDTO `json:"balance"`
	// TODO: here must be also a PetDTO, but they are not implemented yet
}

type CompleteTaskHandler struct {
	taskRepo    *postgres.TaskRepository
	coinsClient coinsClient
	transactor  transactor
}

func NewCompleteTaskHandler(repo *postgres.TaskRepository, coinsClient coinsClient, transactor transactor) *CompleteTaskHandler {
	return &CompleteTaskHandler{
		taskRepo:    repo,
		coinsClient: coinsClient,
		transactor:  transactor,
	}
}

func (h *CompleteTaskHandler) Handle(ctx context.Context, query CompleteTaskQuery) (*CompleteTaskResponse, error) {
	var resp *CompleteTaskResponse

	err := h.transactor.WithTx(ctx, func(ctx context.Context) error {
		completedTask, err := h.taskRepo.CompleteTask(ctx, query.UserID, query.TaskID)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		coinsResp, err := h.coinsClient.UpdateCoins(ctx, entity.UpdateCoinsRequest{
			UserID: query.UserID,
			Delta:  completedTask.Task.RewardCoins,
		})
		if err != nil {
			fmt.Println(err.Error())
			return fmt.Errorf("update coins: %w", err)
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
			},
			Awarded: AwardedDTO{
				Coins: completedTask.Task.RewardCoins,
				XP:    completedTask.Task.RewardXP,
			},
			Balance: &BalanceDTO{Coins: coinsResp.Coins},
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
