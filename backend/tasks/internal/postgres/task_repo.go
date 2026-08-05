package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
)

type TaskModel struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	RewardCoins int    `gorm:"not null"`
	RewardXP    int64  `gorm:"not null"`
}

func (TaskModel) TableName() string {
	return "tasks"
}

type UserTaskModel struct {
	UserID      string `gorm:"primaryKey;type:uuid"`
	TaskID      string `gorm:"primaryKey;type:uuid;not null"`
	Status      string `gorm:"not null;default:active"`
	CompletedAt *time.Time
	UpdatedAt   time.Time
}

func (UserTaskModel) TableName() string {
	return "user_tasks"
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*entity.Task, error) {
	db := GetDB(ctx, r.db)

	var model TaskModel
	err := db.Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, entity.ErrTaskNotFound
		}
		return nil, err
	}
	return r.toDomain(&model), nil
}

func (r *TaskRepository) CompleteTask(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
	db := GetDB(ctx, r.db)

	var taskModel TaskModel
	if err := db.Where("id = ?", taskID).First(&taskModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrTaskNotFound
		}
		return nil, fmt.Errorf("complete task: get task: %w", err)
	}

	var current struct {
		Status      string
		CompletedAt *time.Time
	}
	result := db.Raw(
		`SELECT status, completed_at FROM user_tasks WHERE user_id = ? AND task_id = ? FOR UPDATE`,
		userID,
		taskID,
	).Scan(&current)
	if result.Error != nil {
		return nil, fmt.Errorf("complete task: get user task: %w", result.Error)
	}

	now := time.Now().UTC()
	if result.RowsAffected == 0 {
		if err := db.Exec(
			`INSERT INTO user_tasks (user_id, task_id, status, completed_at, updated_at) VALUES (?, ?, 'completed', ?, NOW())`,
			userID,
			taskID,
			now,
		).Error; err != nil {
			return nil, fmt.Errorf("complete task: insert user task: %w", err)
		}
	} else {
		if current.Status == string(entity.StatusCompleted) {
			return nil, entity.ErrTaskAlreadyCompleted
		}
		if err := db.Exec(
			`UPDATE user_tasks SET status = 'completed', completed_at = ?, updated_at = NOW() WHERE user_id = ? AND task_id = ?`,
			now,
			userID,
			taskID,
		).Error; err != nil {
			return nil, fmt.Errorf("complete task: update user task: %w", err)
		}
	}

	return &entity.UserTask{
		Task: entity.Task{
			ID:          taskModel.ID,
			Title:       taskModel.Title,
			Description: taskModel.Description,
			RewardCoins: taskModel.RewardCoins,
			RewardXP:    taskModel.RewardXP,
		},
		Status:      entity.StatusCompleted,
		CompletedAt: &now,
	}, nil
}

type taskWithUserStatusRow struct {
	ID          string
	Title       string
	Description string
	RewardCoins int
	RewardXP    int64
	Status      string
	CompletedAt *time.Time
}

func (r *TaskRepository) FindByIDForUser(ctx context.Context, userID, taskID string) (*entity.UserTask, error) {
	rows, err := r.listRows(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, entity.ErrTaskNotFound
	}

	return r.rowToUserTask(rows[0]), nil
}

func (r *TaskRepository) ListForUser(ctx context.Context, userID string) ([]entity.UserTask, error) {
	rows, err := r.listRows(ctx, userID, "")
	if err != nil {
		return nil, err
	}

	tasks := make([]entity.UserTask, 0, len(rows))
	for _, row := range rows {
		tasks = append(tasks, *r.rowToUserTask(row))
	}

	return tasks, nil
}

func (r *TaskRepository) listRows(ctx context.Context, userID, taskID string) ([]taskWithUserStatusRow, error) {
	db := GetDB(ctx, r.db)

	query := `
SELECT
    t.id,
    t.title,
    t.description,
    t.reward_coins,
    t.reward_xp,
    COALESCE(ut.status, 'active') AS status,
    ut.completed_at
FROM tasks t
LEFT JOIN user_tasks ut ON ut.task_id = t.id AND ut.user_id = ?
`
	args := []any{userID}
	if taskID != "" {
		query += " WHERE t.id = ?"
		args = append(args, taskID)
	}
	query += " ORDER BY t.title"

	var rows []taskWithUserStatusRow
	if err := db.Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *TaskRepository) rowToUserTask(row taskWithUserStatusRow) *entity.UserTask {
	return &entity.UserTask{
		Task: entity.Task{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			RewardCoins: row.RewardCoins,
			RewardXP:    row.RewardXP,
		},
		Status:      entity.Status(row.Status),
		CompletedAt: row.CompletedAt,
	}
}

func (r *TaskRepository) toDomain(model *TaskModel) *entity.Task {
	return entity.ReconstructTask(model.ID, model.Title, model.Description, model.RewardCoins, model.RewardXP)
}
