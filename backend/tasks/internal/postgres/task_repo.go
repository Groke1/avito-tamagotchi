package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	sqlctask "github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres/sqlc"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) FindByIDForUser(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
	userId, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	taskId, err := uuid.Parse(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %w", err)
	}
	q := sqlctask.New(r.pool)
	row, err := q.FindByIDForUser(ctx, sqlctask.FindByIDForUserParams{
		UserID: pgtype.UUID{Bytes: userId, Valid: true},
		ID:     pgtype.UUID{Bytes: taskId, Valid: true},
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrTaskNotFound
		}
		return nil, err
	}
	var completedAt *time.Time
	if row.CompletedAt.Valid {
		completedAt = &row.CompletedAt.Time
	}

	return &entity.UserTask{
		Task: entity.Task{
			ID:          row.ID.String(),
			Title:       row.Title,
			Description: row.Description,
			RewardCoins: int(row.RewardCoins),
			RewardXP:    row.RewardXp,
		},
		Status:      entity.Status(row.Status),
		CompletedAt: completedAt,
	}, nil
}
func (r *TaskRepository) ListForUser(ctx context.Context, userIDStr string) ([]entity.UserTask, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	q := sqlctask.New(r.pool)
	rows, err := q.ListRandomTasksForUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}

	tasks := make([]entity.UserTask, 0, len(rows))
	for _, row := range rows {
		var completedAt *time.Time
		if row.CompletedAt.Valid {
			completedAt = &row.CompletedAt.Time
		}

		tasks = append(tasks, entity.UserTask{
			Task: entity.Task{
				ID:          row.ID.String(),
				Title:       row.Title,
				Description: row.Description,
				RewardCoins: int(row.RewardCoins),
				RewardXP:    row.RewardXp,
			},
			Status:      entity.Status(row.Status),
			CompletedAt: completedAt,
		})
	}

	return tasks, nil
}

func (r *TaskRepository) CompleteTask(ctx context.Context, userIDStr, taskIDStr string) (*entity.UserTask, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	q := sqlctask.New(tx)

	taskModel, err := q.GetTaskByID(ctx, pgtype.UUID{Bytes: taskID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrTaskNotFound
		}
		return nil, fmt.Errorf("get task: %w", err)
	}

	userTask, err := q.GetUserTaskForUpdate(ctx, sqlctask.GetUserTaskForUpdateParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		TaskID: pgtype.UUID{Bytes: taskID, Valid: true},
	})

	now := time.Now().UTC()
	pgNow := pgtype.Timestamptz{Time: now, Valid: true}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = q.InsertUserTaskCompleted(ctx, sqlctask.InsertUserTaskCompletedParams{
				UserID:      pgtype.UUID{Bytes: userID, Valid: true},
				TaskID:      pgtype.UUID{Bytes: taskID, Valid: true},
				CompletedAt: pgNow,
			})
			if err != nil {
				return nil, fmt.Errorf("insert user task: %w", err)
			}
		} else {
			return nil, fmt.Errorf("get user task: %w", err)
		}
	} else {
		if userTask.Status == string(entity.StatusCompleted) {
			return nil, entity.ErrTaskAlreadyCompleted
		}

		err = q.UpdateUserTaskCompleted(ctx, sqlctask.UpdateUserTaskCompletedParams{
			CompletedAt: pgNow,
			UserID:      pgtype.UUID{Bytes: userID, Valid: true},
			TaskID:      pgtype.UUID{Bytes: taskID, Valid: true},
		})
		if err != nil {
			return nil, fmt.Errorf("update user task: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return &entity.UserTask{
		Task: entity.Task{
			ID:          taskModel.ID.String(),
			Title:       taskModel.Title,
			Description: taskModel.Description,
			RewardCoins: int(taskModel.RewardCoins),
			RewardXP:    taskModel.RewardXp,
		},
		Status:      entity.StatusCompleted,
		CompletedAt: &now,
	}, nil
}
