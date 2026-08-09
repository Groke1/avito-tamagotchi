package entity

import (
	"errors"
	"time"
)

var (
	ErrInvalidID            = errors.New("INVALID ID")
	ErrInvalidTitle         = errors.New("TITLE IS REQUIRED")
	ErrInvalidDescription   = errors.New("DESCRIPTION IS REQUIRED")
	ErrInvalidRewardCoins   = errors.New("COINS MUST BE NONNEGATIVE")
	ErrInvalidRewardXP      = errors.New("XP MUST BE NONNEGATIVE")
	ErrTaskNotFound         = errors.New("TASK NOT FOUND")
	ErrTaskAlreadyCompleted = errors.New("TASK ALREADY COMPLETED")
	ErrInvalidTaskType      = errors.New("TASK TYPE IS REQUIRED")
	ErrUserNotFound         = errors.New("USER NOT FOUND")
	ErrNoCompletedTasks     = errors.New("NO COMPLETED TASKS")
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
	Type        string
}

type UserTask struct {
	Task        Task
	Status      Status
	CompletedAt *time.Time
}

func NewTask(id, title, description string, rewardCoins int, rewardXP int64, taskType string) (*Task, error) {
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
	if taskType == "" || len(taskType) >= 25 {
		return nil, ErrInvalidTaskType
	}
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		RewardCoins: rewardCoins,
		RewardXP:    rewardXP,
		Type:        taskType,
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
