package controller

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type ErrorCode string

const (
	ErrUnauthorized           ErrorCode = "UNAUTHORIZED"
	ErrTaskNotFound           ErrorCode = "TASK_NOT_FOUND"
	ErrTaskAlreadyCompleted   ErrorCode = "TASK_ALREADY_COMPLETED"
	ErrUserNotFound           ErrorCode = "USER_NOT_FOUND"
	ErrInvalidRewardRequest   ErrorCode = "INVALID_REWARD_REQUEST"
	ErrUserServiceUnavailable ErrorCode = "USER_SERVICE_UNAVAILABLE"
	ErrInternal               ErrorCode = "INTERNAL_ERROR"
)

func (e ErrorCode) Message() string {
	switch e {
	case ErrUnauthorized:
		return "Требуется повторная авторизация"

	case ErrTaskNotFound:
		return "Задание не найдено"

	case ErrTaskAlreadyCompleted:
		return "Награда за это задание уже получена"

	case ErrUserNotFound:
		return "Пользователь не найден"

	case ErrInvalidRewardRequest:
		return "Не удалось начислить награду"

	case ErrUserServiceUnavailable:
		return "Сервис пользователей временно недоступен"

	case ErrInternal:
		return "Внутренняя ошибка сервера"

	default:
		return "Unknown error"
	}
}
func WriteError(w http.ResponseWriter, status int, code ErrorCode) {
	WriteJSON(w, status, ErrorResponse{
		Code:    string(code),
		Message: code.Message(),
	})
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
