package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorCode_Message_KnownCodes(t *testing.T) {
	tests := []struct {
		code ErrorCode
		want string
	}{
		{ErrUserAlreadyExists, "Пользователь с такими данными уже существует"},
		{ErrValidation, "Проверьте переданные данные"},
		{ErrInternal, "Внутренняя ошибка сервера"},
		{ErrInvalidCredentials, "Неверный email или пароль"},
		{ErrInvalidRefreshToken, "Сессия истекла. Выполните вход снова"},
		{ErrUnauthorized, "Требуется повторная авторизация"},
		{ErrUserNotFound, "Пользователь не найден"},
		{ErrInsufficientCoins, "Недостаточно монет"},
		{ErrRewardNotFound, "Награда не найдена"},
		{ErrNotFoundDefinition, "Награда не найдена"},
		{ErrRewardUnavailable, "Награда недоступна"},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			if got := tt.code.Message(); got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestErrorCode_Message_UnknownCode(t *testing.T) {
	unknown := ErrorCode("SOMETHING_NEW")
	if got := unknown.Message(); got != "Unknown error" {
		t.Errorf("expected fallback message, got %q", got)
	}
}

func TestWriteError_WritesStatusAndBody(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, http.StatusConflict, ErrUserAlreadyExists)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", rec.Code)
	}

	var body ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if body.Code != string(ErrUserAlreadyExists) {
		t.Errorf("expected code %q, got %q", ErrUserAlreadyExists, body.Code)
	}
	if body.Message != ErrUserAlreadyExists.Message() {
		t.Errorf("expected message %q, got %q", ErrUserAlreadyExists.Message(), body.Message)
	}
}
