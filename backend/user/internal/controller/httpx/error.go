package httpx

import "net/http"

type ErrorCode string

const (
	ErrUserAlreadyExists   ErrorCode = "USER_ALREADY_EXISTS"
	ErrValidation          ErrorCode = "VALIDATION_ERROR"
	ErrInternal            ErrorCode = "INTERNAL_ERROR"
	ErrInvalidCredentials  ErrorCode = "INVALID_CREDENTIALS"
	ErrInvalidRefreshToken ErrorCode = "INVALID_REFRESH_TOKEN"
	ErrUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrUserNotFound        ErrorCode = "USER_NOT_FOUND"
	ErrInsufficientCoins   ErrorCode = "INSUFFICIENT_COINS"
	ErrRewardNotFound      ErrorCode = "REWARD_NOT_FOUND"
	ErrRewardUnavailable   ErrorCode = "REWARD_UNAVAILABLE"
	ErrNotFoundDefinition  ErrorCode = "NOT_FOUND_DEFINITION"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ErrorCode) Message() string {
	switch e {
	case ErrUserAlreadyExists:
		return "Пользователь с такими данными уже существует"
	case ErrValidation:
		return "Проверьте переданные данные"
	case ErrInternal:
		return "Внутренняя ошибка сервера"
	case ErrInvalidCredentials:
		return "Неверный email или пароль"
	case ErrInvalidRefreshToken:
		return "Сессия истекла. Выполните вход снова"
	case ErrUnauthorized:
		return "Требуется повторная авторизация"
	case ErrUserNotFound:
		return "Пользователь не найден"
	case ErrInsufficientCoins:
		return "Недостаточно монет"
	case ErrRewardNotFound:
		return "Награда не найдена"
	case ErrRewardUnavailable:
		return "Награда недоступна"
	case ErrNotFoundDefinition:
		return "Награда не найдена"

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
