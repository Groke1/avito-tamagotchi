package http

import "github.com/cayman444/avito-gamification-hackathon.user/internal/entity"

type getCompletedTasksResponse struct {
	UserID string             `json:"user_id"`
	Items  []entity.TasksStat `json:"items"`
}
