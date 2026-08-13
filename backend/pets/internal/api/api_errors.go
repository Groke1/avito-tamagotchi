package api

import "net/http"

// пока такое название, в будущем может поменяться
type HTTPError struct {
	StatusCode int
	Code       string
	Message    string
}

var (
	ErrUnauthorized = HTTPError{
		StatusCode: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    "Требуется повторная авторизация",
	}

	ErrPetNotFound = HTTPError{
		StatusCode: http.StatusConflict,
		Code:       "PET_NOT_FOUND",
		Message:    "Сначала создайте питомца",
	}

	ErrUnavailableAction = HTTPError{
		StatusCode: http.StatusConflict,
		Code:       "PET_ACTION_UNAVAILABLE",
		Message:    "Это действие пока недоступно",
	}

	ErrPetAlreadyExists = HTTPError{
		StatusCode: http.StatusNotFound,
		Code:       "PET_ALREADY_EXISTS",
		Message:    "У пользователя уже есть питомец",
	}

	ErrValidationError = HTTPError{
		StatusCode: http.StatusUnprocessableEntity,
		Code:       "VALIDATION_ERROR",
		Message:    "Проверьте переданные данные",
	}

	ErrUserNotFound = HTTPError{
		StatusCode: http.StatusNotFound,
		Code:       "USER_NOT_FOUND",
		Message:    "Пользователь не найден",
	}

	ErrNotEnoughCoins = HTTPError{
		StatusCode: http.StatusConflict,
		Code:       "INSUFFICIENT_COINS",
		Message:    "Недостаточно монет",
	}

	ErrInternalError = HTTPError{
		StatusCode: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    "Внутренняя ошибка сервере",
	}
	ErrTripGenerationError = HTTPError{
		StatusCode: http.StatusInternalServerError,
		Code:       "TRIP_GENERATION_ERROR",
		Message:    "Ошибка генерации путешествия",
	}
)
