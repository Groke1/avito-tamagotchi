package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.pkg/db"
	sqlctask "github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres/sqlc"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	pool    *pgxpool.Pool
	queries *sqlctask.Queries
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		pool:    pool,
		queries: sqlctask.New(pool),
	}
}

func (r *TaskRepository) FindByIDForUser(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
	userId, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", entity.ErrInvalidID)
	}
	taskId, err := uuid.Parse(taskID)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %w", entity.ErrInvalidID)
	}
	q := sqlctask.New(r.pool)
	row, err := q.FindByIDForUser(ctx, sqlctask.FindByIDForUserParams{
		UserID: pgtype.UUID{Bytes: userId, Valid: true},
		ID:     pgtype.UUID{Bytes: taskId, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task_id %q for user %q: %w", taskID, userID, entity.ErrTaskNotFound)
		}
		return nil, fmt.Errorf("find task: %w", err)
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
			Type:        row.TaskType.String,
		},
		Status:      entity.Status(row.Status),
		CompletedAt: completedAt,
	}, nil
}

func (r *TaskRepository) ListForUser(ctx context.Context, userIDStr string) ([]entity.UserTask, error) {
	var pgUserID pgtype.UUID
	err := pgUserID.Scan(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", entity.ErrInvalidID)
	}

	q := sqlctask.New(r.pool)
	limit := int32(rand.Intn(3) + 3)
	rows, err := q.CreateUserTasksBatch(ctx, sqlctask.CreateUserTasksBatchParams{
		UserID:      pgUserID,
		RandomLimit: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("create user tasks batch %q: %w", userIDStr, entity.ErrInvalidID)
	}

	tasks := make([]entity.UserTask, 0, len(rows))
	var newTasksIds []pgtype.UUID
	for _, row := range rows {
		comp_at := &row.CompletedAt.Time
		if entity.Status(row.Status) != entity.StatusCompleted {
			comp_at = nil
		}
		tasks = append(tasks, entity.UserTask{
			Task: entity.Task{
				ID:          row.ID.String(),
				Title:       row.Title,
				Description: row.Description,
				RewardCoins: int(row.RewardCoins),
				RewardXP:    row.RewardXp,
				Type:        row.TaskType.String,
			},
			Status:      entity.Status(row.Status),
			CompletedAt: comp_at,
		})
		if row.IsNew {
			newTasksIds = append(newTasksIds, row.ID)
		}
	}
	if len(newTasksIds) > 0 {
		err = q.InsertUserTasksBatch(ctx, sqlctask.InsertUserTasksBatchParams{
			UserID:  pgUserID,
			TaskIds: newTasksIds,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save generated tasks: %w", err)
		}
	}
	return tasks, nil
}

func (r *TaskRepository) CompleteTask(ctx context.Context, userIDStr, taskIDStr string) (*entity.UserTask, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("user_id %q: %w", userIDStr, entity.ErrInvalidID)
	}
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return nil, fmt.Errorf("task_id %q: %w", taskIDStr, entity.ErrInvalidID)
	}
	q := r.querier(ctx)
	taskModel, err := q.GetTaskByID(ctx, pgtype.UUID{Bytes: taskID, Valid: true})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task_id %q: %w", taskIDStr, entity.ErrTaskNotFound)
		}
		return nil, fmt.Errorf("get task: %w", err)
	}

	userTask, err := q.GetUserTaskForUpdate(ctx, sqlctask.GetUserTaskForUpdateParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		TaskID: pgtype.UUID{Bytes: taskID, Valid: true},
	})

	now := time.Now().UTC()
	pgNow := pgtype.Timestamptz{Time: now, Valid: true}

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		if err := q.InsertUserTaskCompleted(ctx, sqlctask.InsertUserTaskCompletedParams{
			UserID:      pgtype.UUID{Bytes: userID, Valid: true},
			TaskID:      pgtype.UUID{Bytes: taskID, Valid: true},
			CompletedAt: pgNow,
		}); err != nil {
			return nil, fmt.Errorf("insert user task: %w", err)
		}
	case err != nil:
		return nil, fmt.Errorf("get user task: %w", err)
	default:
		if userTask.Status == string(entity.StatusCompleted) {
			return nil, entity.ErrTaskAlreadyCompleted
		}
		if err := q.UpdateUserTaskCompleted(ctx, sqlctask.UpdateUserTaskCompletedParams{
			UserID: pgtype.UUID{Bytes: userID, Valid: true},
			TaskID: pgtype.UUID{Bytes: taskID, Valid: true},
		}); err != nil {
			return nil, fmt.Errorf("update user task: %w", err)
		}
	}

	return &entity.UserTask{
		Task: entity.Task{
			ID:          taskModel.ID.String(),
			Title:       taskModel.Title,
			Description: taskModel.Description,
			RewardCoins: int(taskModel.RewardCoins),
			RewardXP:    taskModel.RewardXp,
			Type:        taskModel.TaskType.String,
		},
		Status:      entity.StatusCompleted,
		CompletedAt: &now,
	}, nil
}

func (r *TaskRepository) ListCompletedToday(ctx context.Context, userIDStr string) ([]sqlctask.GetTodayCompletedTasksForUserRow, error) {
	userId, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", entity.ErrInvalidID)
	}
	rows, err := r.queries.GetTodayCompletedTasksForUser(ctx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("list completed today: %w", err)
	}
	return rows, nil
}

func (r *TaskRepository) querier(ctx context.Context) *sqlctask.Queries {
	tx, err := db.ExtractTx(ctx)
	if err != nil {
		return sqlctask.New(r.pool)
	}
	return sqlctask.New(tx)
}
