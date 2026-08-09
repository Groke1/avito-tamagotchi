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
		status   int
		err_name string
		message  string
	)

	switch {
	case errors.Is(err, controller.ErrExpiredToken) ||
		errors.Is(err, controller.ErrInvalidToken) ||
		errors.Is(err, controller.ErrUnauthorized):
		status = http.StatusUnauthorized
		err_name = controller.ErrUnauthorized.Error()
		message = "Требуется авторизация"

	case errors.Is(err, entity.ErrTaskAlreadyCompleted):
		status = http.StatusConflict
		err_name = entity.ErrTaskAlreadyCompleted.Error()
		message = "Награда за эту задачу уже получена"

	case errors.Is(err, entity.ErrTaskNotFound) ||
		errors.Is(err, entity.ErrUserNotFound):
		status = http.StatusNotFound
		err_name = entity.ErrTaskNotFound.Error()
		message = "Задача не найдена"

	case errors.Is(err, entity.ErrInvalidID) ||
		errors.Is(err, entity.ErrInvalidTitle) ||
		errors.Is(err, entity.ErrInvalidDescription) ||
		errors.Is(err, entity.ErrInvalidRewardCoins) ||
		errors.Is(err, entity.ErrInvalidRewardXP) ||
		errors.Is(err, entity.ErrInvalidTaskType):

		status = http.StatusBadRequest
		err_name = controller.ErrInvalidRequest.Error()
		message = "Некорректный запрос"

	case errors.Is(err, controller.ErrUserServiceUnavailable) ||
		errors.Is(err, controller.ErrPetServiceUnavailable):
		status = http.StatusInternalServerError
		err_name = controller.ErrInternal.Error()
		message = "Внутренняя ошибка сервера"

	default:
		status = http.StatusInternalServerError
		err_name = controller.ErrInternal.Error()
		message = "Внутренняя ошибка сервера"
	}

	return ErrorDTO{
		Status:  status,
		ErrName: err_name,
		Message: message,
	}
}
