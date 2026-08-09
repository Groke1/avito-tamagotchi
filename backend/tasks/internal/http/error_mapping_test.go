package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/controller"
	"github.com/cayman444/avito-gamification-hackathon.tasks/internal/entity"
	"github.com/gin-gonic/gin"
)

func TestMapError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{"expired token", controller.ErrExpiredToken, http.StatusUnauthorized, controller.ErrUnauthorized.Error()},
		{"invalid token", controller.ErrInvalidToken, http.StatusUnauthorized, controller.ErrUnauthorized.Error()},
		{"unauthorized", controller.ErrUnauthorized, http.StatusUnauthorized, controller.ErrUnauthorized.Error()},
		{"task already completed", entity.ErrTaskAlreadyCompleted, http.StatusConflict, entity.ErrTaskAlreadyCompleted.Error()},
		{"task not found", entity.ErrTaskNotFound, http.StatusNotFound, entity.ErrTaskNotFound.Error()},
		{"user not found maps to task not found response", entity.ErrUserNotFound, http.StatusNotFound, entity.ErrTaskNotFound.Error()},
		{"invalid title", entity.ErrInvalidTitle, http.StatusBadRequest, controller.ErrInvalidRequest.Error()},
		{"invalid reward coins", entity.ErrInvalidRewardCoins, http.StatusBadRequest, controller.ErrInvalidRequest.Error()},
		{"unknown error falls back to internal", errors.New("boom"), http.StatusInternalServerError, controller.ErrInternal.Error()},
		{"wrapped known error is still recognized", errors.Join(errors.New("ctx"), entity.ErrTaskNotFound), http.StatusNotFound, entity.ErrTaskNotFound.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapError(tt.err)
			if got.Status != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, got.Status)
			}
			if got.ErrName != tt.wantCode {
				t.Errorf("expected code %q, got %q", tt.wantCode, got.ErrName)
			}
			if got.Message == "" {
				t.Errorf("expected non-empty message")
			}
		})
	}
}

func TestSendError_WritesExpectedJSONAndStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)

	SendError(c, entity.ErrTaskNotFound)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}

	var body ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}
	if body.Code != entity.ErrTaskNotFound.Error() {
		t.Errorf("expected code %q, got %q", entity.ErrTaskNotFound.Error(), body.Code)
	}
	if body.Message != "Задача не найдена" {
		t.Errorf("unexpected message: %q", body.Message)
	}
}
