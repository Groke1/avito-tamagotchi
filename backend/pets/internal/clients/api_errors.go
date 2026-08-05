package clients

import "net/http"

type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

var (
	ErrValidationError = APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Code:       "VALIDATION_ERROR",
		Message:    "Проверьте переданные данные",
	}

	ErrUserNotFound = APIError{
		StatusCode: http.StatusNotFound,
		Code:       "USER_NOT_FOUND",
		Message:    "Пользователь не найден",
	}

	ErrNotEnoughCoins = APIError{
		StatusCode: http.StatusConflict,
		Code:       "INSUFFICIENT_COINS",
		Message:    "Недостаточно монет",
	}

	ErrInternalError = APIError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    "Внутренняя ошибка сервер",
	}
)
