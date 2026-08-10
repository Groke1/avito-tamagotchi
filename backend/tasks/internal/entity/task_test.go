package entity

import (
	"errors"
	"strings"
	"testing"
)

func TestNewTask_Valid(t *testing.T) {
	task, err := NewTask("id-1", "Заголовок", "Описание", 100, 500, "Продажи")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.ID != "id-1" || task.Title != "Заголовок" || task.Description != "Описание" {
		t.Fatalf("unexpected task fields: %+v", task)
	}
	if task.RewardCoins != 100 || task.RewardXP != 500 || task.Type != "Продажи" {
		t.Fatalf("unexpected reward/type fields: %+v", task)
	}
}

func TestNewTask_ValidationRules(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		rewardCoins int
		rewardXP    int64
		taskType    string
		wantErr     error
	}{
		{
			name:        "empty title",
			title:       "",
			description: "desc",
			rewardCoins: 10,
			rewardXP:    10,
			taskType:    "Продажи",
			wantErr:     ErrInvalidTitle,
		},
		{
			name:        "empty description",
			title:       "title",
			description: "",
			rewardCoins: 10,
			rewardXP:    10,
			taskType:    "Продажи",
			wantErr:     ErrInvalidDescription,
		},
		{
			name:        "negative reward coins",
			title:       "title",
			description: "desc",
			rewardCoins: -1,
			rewardXP:    10,
			taskType:    "Продажи",
			wantErr:     ErrInvalidRewardCoins,
		},
		{
			name:        "negative reward xp",
			title:       "title",
			description: "desc",
			rewardCoins: 10,
			rewardXP:    -1,
			taskType:    "Продажи",
			wantErr:     ErrInvalidRewardXP,
		},
		{
			name:        "empty task type",
			title:       "title",
			description: "desc",
			rewardCoins: 10,
			rewardXP:    10,
			taskType:    "",
			wantErr:     ErrInvalidTaskType,
		},
		{
			name:        "task type too long (>=25 chars)",
			title:       "title",
			description: "desc",
			rewardCoins: 10,
			rewardXP:    10,
			taskType:    strings.Repeat("a", 25),
			wantErr:     ErrInvalidTaskType,
		},
		{
			name:        "task type at boundary (24 chars) is valid",
			title:       "title",
			description: "desc",
			rewardCoins: 10,
			rewardXP:    10,
			taskType:    strings.Repeat("a", 24),
			wantErr:     nil,
		},
		{
			name:        "zero reward is valid (minimum boundary)",
			title:       "title",
			description: "desc",
			rewardCoins: 0,
			rewardXP:    0,
			taskType:    "Продажи",
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask("id", tt.title, tt.description, tt.rewardCoins, tt.rewardXP, tt.taskType)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if task == nil {
					t.Fatalf("expected task to be created")
				}
				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if task != nil {
				t.Fatalf("expected nil task on validation error, got %+v", task)
			}
		})
	}
}

func TestReconstructUserTask(t *testing.T) {
	task := Task{ID: "id-1", Title: "t", Description: "d", RewardCoins: 10, RewardXP: 20, Type: "Продажи"}

	userTask := ReconstructUserTask(task, StatusCompleted, nil)

	if userTask.Task != task {
		t.Fatalf("expected embedded task to match, got %+v", userTask.Task)
	}
	if userTask.Status != StatusCompleted {
		t.Fatalf("expected status completed, got %v", userTask.Status)
	}
	if userTask.CompletedAt != nil {
		t.Fatalf("expected nil completedAt, got %v", userTask.CompletedAt)
	}
}
