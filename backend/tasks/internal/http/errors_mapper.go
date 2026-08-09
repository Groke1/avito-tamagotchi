package http

import (
	"errors"
	"net/http"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type ErrorDTO struct {
	Status  int    `json:"-"`
	ErrName string `json:"code"`
	Message string `json:"message"`
}

func SendError(c *gin.Context, err error) {
	response := mapError(err)
	c.JSON(response.Status, ErrorResponse{
		Code:    response.ErrName,
		Message: response.Message,
	})
}

func mapError(err error) ErrorDTO {
	var (
		status  int
		errName string
		message string
	)

	switch {
	case errors.Is(err, controller.ErrExpiredToken) ||
		errors.Is(err, controller.ErrInvalidToken) ||
		errors.Is(err, controller.ErrUnauthorized):
		status = http.StatusUnauthorized
		errName = controller.ErrUnauthorized.Error()
		message = "Требуется авторизация"

	case errors.Is(err, entity.ErrTaskAlreadyCompleted):
		status = http.StatusConflict
		errName = entity.ErrTaskAlreadyCompleted.Error()
		message = "Награда за эту задачу уже получена"

	case errors.Is(err, entity.ErrTaskNotFound) ||
		errors.Is(err, entity.ErrUserNotFound):
		status = http.StatusNotFound
		errName = entity.ErrTaskNotFound.Error()
		message = "Задача не найдена"

	case errors.Is(err, entity.ErrInvalidID) ||
		errors.Is(err, entity.ErrInvalidTitle) ||
		errors.Is(err, entity.ErrInvalidDescription) ||
		errors.Is(err, entity.ErrInvalidRewardCoins) ||
		errors.Is(err, entity.ErrInvalidRewardXP) ||
		errors.Is(err, entity.ErrInvalidTaskType):

		status = http.StatusBadRequest
		errName = controller.ErrInvalidRequest.Error()
		message = "Некорректный запрос"

	default:
		status = http.StatusInternalServerError
		errName = controller.ErrInternal.Error()
		message = "Внутренняя ошибка сервера"
	}
	return ErrorDTO{
		Status:  status,
		ErrName: errName,
		Message: message,
	}
}
