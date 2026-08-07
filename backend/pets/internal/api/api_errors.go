package api

import "net/http"

type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

var (
	ErrUnauthorized = APIError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    "Требуется повторная авторизация",
	}

	ErrPetNotFound = APIError{
		StatusCode: http.StatusConflict,
		Code:       "PET_NOT_FOUND",
		Message:    "Сначала создайте питомца",
	}

	ErrUnavailableAction = APIError{
		StatusCode: http.StatusConflict,
		Code:       "PET_ACTION_UNAVAILABLE",
		Message:    "Это действие пока недоступно",
	}

	ErrPetAlreadyExists = APIError{
		StatusCode: http.StatusNotFound,
		Code:       "PET_ALREADY_EXISTS",
		Message:    "У пользователя уже есть питомец",
	}

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
