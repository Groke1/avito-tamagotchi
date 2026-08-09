package postgres

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	sqlctask "github.com/cayman444/avito-gamification-hackathon.tasks/internal/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	TaskLimitRange = 3
	MinTaskLimit   = 3
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

func (r *TaskRepository) FindByIDForUser(ctx context.Context, userIDStr, taskIDStr string) (*entity.UserTask, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", entity.ErrInvalidID)
	}
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid task id: %w", entity.ErrInvalidID)
	}
	row, err := r.queries.FindByIDForUser(ctx, sqlctask.FindByIDForUserParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		ID:     pgtype.UUID{Bytes: taskID, Valid: true},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("task_id %q for user %q: %w", taskIDStr, userIDStr, entity.ErrTaskNotFound)
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

	limit := int32(rand.Intn(TaskLimitRange) + MinTaskLimit)
	rows, err := r.queries.CreateUserTasksBatch(ctx, sqlctask.CreateUserTasksBatchParams{
		UserID:      pgUserID,
		RandomLimit: limit,
	})
	if err != nil {
		return nil, fmt.Errorf("create user tasks batch %q: %w", userIDStr, entity.ErrInvalidID)
	}

	tasks := make([]entity.UserTask, 0, len(rows))
	var newTasksIDs []pgtype.UUID
	for _, row := range rows {
		compAt := &row.CompletedAt.Time
		if entity.Status(row.Status) != entity.StatusCompleted {
			compAt = nil
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
			CompletedAt: compAt,
		})
		if row.IsNew {
			newTasksIDs = append(newTasksIDs, row.ID)
		}
	}
	if len(newTasksIDs) > 0 {
		err = r.queries.InsertUserTasksBatch(ctx, sqlctask.InsertUserTasksBatchParams{
			UserID:  pgUserID,
			TaskIds: newTasksIDs,
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

	userTask, err := r.queries.GetUserTaskForUpdate(
		ctx,
		sqlctask.GetUserTaskForUpdateParams{
			UserID: pgtype.UUID{Bytes: userID, Valid: true},
			TaskID: pgtype.UUID{Bytes: taskID, Valid: true},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get user task: %w", err)
	}

	if userTask.Status == string(entity.StatusCompleted) {
		return nil, entity.ErrTaskAlreadyCompleted
	}

	err = r.queries.UpdateUserTaskCompleted(
		ctx,
		sqlctask.UpdateUserTaskCompletedParams{
			UserID: pgtype.UUID{Bytes: userID, Valid: true},
			TaskID: pgtype.UUID{Bytes: taskID, Valid: true},
		})
	if err != nil {
		return nil, fmt.Errorf("update user task: %w", err)
	}

	now := time.Now().UTC()

	return &entity.UserTask{
		Task: entity.Task{
			ID:          userTask.ID.String(),
			Title:       userTask.Title,
			Description: userTask.Description,
			RewardCoins: int(userTask.RewardCoins),
			RewardXP:    userTask.RewardXp,
			Type:        userTask.TaskType.String,
		},
		Status:      entity.StatusCompleted,
		CompletedAt: &now,
	}, nil
}

func (r *TaskRepository) ListCompletedToday(ctx context.Context, userIDStr string) ([]sqlctask.GetTodayCompletedTasksForUserRow, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", entity.ErrInvalidID)
	}
	rows, err := r.queries.GetTodayCompletedTasksForUser(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("list completed today: %w", err)
	}
	return rows, nil
}
