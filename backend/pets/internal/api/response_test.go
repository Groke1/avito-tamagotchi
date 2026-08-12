package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	writeError(rec, ErrPetNotFound)

	if rec.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
	if ct := rec.Header().Get("Content-type"); ct != "application/json" {
		t.Errorf("content-type = %q, want %q", ct, "application/json")
	}

	var body ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body.Code != ErrPetNotFound.Code {
		t.Errorf("code = %q, want %q", body.Code, ErrPetNotFound.Code)
	}
	if body.Message != ErrPetNotFound.Message {
		t.Errorf("message = %q, want %q", body.Message, ErrPetNotFound.Message)
	}
	if body.RetryAfter != nil {
		t.Errorf("retry_after = %v, want nil", body.RetryAfter)
	}
}

func TestWriteJSONResponse(t *testing.T) {
	rec := httptest.NewRecorder()
	payload := PetResponse{ID: 1, Name: "Rex", Level: 2, XP: 10, NextLevelXP: 50, Satiety: 80, Happiness: 90}

	writeJSONResponse(rec, http.StatusCreated, payload)

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if ct := rec.Header().Get("Content-type"); ct != "application/json" {
		t.Errorf("content-type = %q, want %q", ct, "application/json")
	}

	var body PetResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body != payload {
		t.Errorf("body = %+v, want %+v", body, payload)
	}
}

func TestWriteErrorWithRetryAfter(t *testing.T) {
	rec := httptest.NewRecorder()

	writeErrorWithRetryAfter(rec, ErrUnavailableAction, 42)

	if rec.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusConflict)
	}

	var body ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body.RetryAfter == nil {
		t.Fatal("retry_after = nil, want 42")
	}
	if *body.RetryAfter != 42 {
		t.Errorf("retry_after = %d, want %d", *body.RetryAfter, 42)
	}
}

func TestHTTPErrorDefinitions(t *testing.T) {
	tests := []struct {
		name       string
		err        HTTPError
		wantStatus int
		wantCode   string
	}{
		{"unauthorized", ErrUnauthorized, http.StatusUnauthorized, "UNAUTHORIZED"},
		{"pet not found", ErrPetNotFound, http.StatusConflict, "PET_NOT_FOUND"},
		{"action unavailable", ErrUnavailableAction, http.StatusConflict, "PET_ACTION_UNAVAILABLE"},
		{"pet already exists", ErrPetAlreadyExists, http.StatusNotFound, "PET_ALREADY_EXISTS"},
		{"validation error", ErrValidationError, http.StatusUnprocessableEntity, "VALIDATION_ERROR"},
		{"user not found", ErrUserNotFound, http.StatusNotFound, "USER_NOT_FOUND"},
		{"not enough coins", ErrNotEnoughCoins, http.StatusConflict, "INSUFFICIENT_COINS"},
		{"internal error", ErrInternalError, http.StatusInternalServerError, "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", tt.err.StatusCode, tt.wantStatus)
			}
			if tt.err.Code != tt.wantCode {
				t.Errorf("code = %q, want %q", tt.err.Code, tt.wantCode)
			}
			if tt.err.Message == "" {
				t.Error("message is empty")
			}
		})
	}
}
