package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidTitle         = errors.New("title is required")
	ErrInvalidDescription   = errors.New("description is required")
	ErrInvalidRewardCoins   = errors.New("coins must be nonnegative")
	ErrInvalidRewardXP      = errors.New("xp must be nonnegative")
	ErrTaskNotFound         = errors.New("task not found")
	ErrTaskAlreadyCompleted = errors.New("task already completed")
)

type Status string

const (
	StatusInProgress Status = "active"
	StatusCompleted  Status = "completed"
)

type Task struct {
	ID          string
	Title       string
	Description string
	RewardCoins int
	RewardXP    int64
}

type UserTask struct {
	Task        Task
	Status      Status
	CompletedAt *time.Time
}

func NewTask(id, title, description string, rewardCoins int, rewardXP int64) (*Task, error) {
	if title == "" {
		return nil, ErrInvalidTitle
	}
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if rewardCoins < 0 {
		return nil, ErrInvalidRewardCoins
	}
	if rewardXP < 0 {
		return nil, ErrInvalidRewardXP
	}

	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		RewardCoins: rewardCoins,
		RewardXP:    rewardXP,
	}, nil
}

func ReconstructTask(id, title, description string, rewardCoins int, rewardXP int64) *Task {
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		RewardCoins: rewardCoins,
		RewardXP:    rewardXP,
	}
}

func ReconstructUserTask(task Task, status Status, completedAt *time.Time) UserTask {
	return UserTask{
		Task:        task,
		Status:      status,
		CompletedAt: completedAt,
	}
}
