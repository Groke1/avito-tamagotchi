package clients

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon/backend/pets/internal/domain"
)

func TestUserClientWithdrawCoins(t *testing.T) {
	t.Run("zero amount skips request", func(t *testing.T) {
		called := false
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		err := client.WithdrawCoins(context.Background(), "user-1", 0)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if called {
			t.Error("expected no HTTP call for zero amount")
		}
	})

	t.Run("success", func(t *testing.T) {
		var capturedBody UpdateCoinsRequest
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("method = %s, want PUT", r.Method)
			}
			if r.URL.Path != "/update-coins" {
				t.Errorf("path = %s, want /update-coins", r.URL.Path)
			}
			_ = json.NewDecoder(r.Body).Decode(&capturedBody)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		err := client.WithdrawCoins(context.Background(), "user-1", 50)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if capturedBody.UserID != "user-1" {
			t.Errorf("userID = %q, want %q", capturedBody.UserID, "user-1")
		}
		if capturedBody.Delta != -50 {
			t.Errorf("delta = %d, want %d", capturedBody.Delta, -50)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(APIError{StatusCode: http.StatusNotFound, Code: "USER_NOT_FOUND", Message: "not found"})
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		err := client.WithdrawCoins(context.Background(), "user-1", 10)

		if !errors.Is(err, domain.ErrUserNotFound) {
			t.Errorf("err = %v, want domain.ErrUserNotFound", err)
		}
	})

	t.Run("not enough coins", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			_ = json.NewEncoder(w).Encode(APIError{StatusCode: http.StatusConflict, Code: "INSUFFICIENT_COINS", Message: "not enough"})
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		err := client.WithdrawCoins(context.Background(), "user-1", 10)

		if !errors.Is(err, domain.ErrNotEnoguhCoins) {
			t.Errorf("err = %v, want domain.ErrNotEnoguhCoins", err)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(APIError{StatusCode: http.StatusUnprocessableEntity, Code: "VALIDATION_ERROR", Message: "bad data"})
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		err := client.WithdrawCoins(context.Background(), "user-1", 10)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("internal error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(APIError{StatusCode: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "boom"})
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		err := client.WithdrawCoins(context.Background(), "user-1", 10)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("unexpected status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
			_ = json.NewEncoder(w).Encode(APIError{StatusCode: http.StatusTeapot, Code: "TEAPOT", Message: "im a teapot"})
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		err := client.WithdrawCoins(context.Background(), "user-1", 10)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("non-200 with undecodable body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte("not json"))
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		err := client.WithdrawCoins(context.Background(), "user-1", 10)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("server unreachable", func(t *testing.T) {
		client := NewUserClient("http://127.0.0.1:0")
		err := client.WithdrawCoins(context.Background(), "user-1", 10)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		client := NewUserClient(server.URL)
		err := client.WithdrawCoins(ctx, "user-1", 10)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUserClientGetUsernamesByIDs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/usernames" {
				t.Errorf("path = %s, want /usernames", r.URL.Path)
			}
			_ = json.NewEncoder(w).Encode(UserNamesResponse{Users: []User{
				{ID: "1", UserName: "alice"},
				{ID: "2", UserName: "bob"},
			}})
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		result, err := client.GetUsernamesByIDs(context.Background(), []string{"1", "2"})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["1"] != "alice" || result["2"] != "bob" {
			t.Errorf("result = %v, want map with alice and bob", result)
		}
	})

	t.Run("empty ids", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = json.NewEncoder(w).Encode(UserNamesResponse{Users: []User{}})
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		result, err := client.GetUsernamesByIDs(context.Background(), []string{})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("result = %v, want empty map", result)
		}
	})

	t.Run("server unreachable", func(t *testing.T) {
		client := NewUserClient("http://127.0.0.1:0")
		_, err := client.GetUsernamesByIDs(context.Background(), []string{"1"})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not json"))
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		_, err := client.GetUsernamesByIDs(context.Background(), []string{"1"})

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUserClientClaimReward(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		redeemedAt := time.Now()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s, want POST", r.Method)
			}
			if r.URL.Path != "/rewards" {
				t.Errorf("path = %s, want /rewards", r.URL.Path)
			}
			_ = json.NewEncoder(w).Encode(RewardResponse{
				RewardID:    "r1",
				PromoCode:   "PROMO10",
				Name:        "10% off",
				Description: "discount",
				Status:      "active",
				ExpiresAt:   "2030-01-01",
				RedeemedAt:  &redeemedAt,
			})
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		reward, err := client.ClaimReward(context.Background(), "user-1", "PROMO10")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if reward.ID != "r1" || reward.PromoCode != "PROMO10" || reward.Name != "10% off" {
			t.Errorf("reward = %+v, unexpected values", reward)
		}
	})

	t.Run("server unreachable", func(t *testing.T) {
		client := NewUserClient("http://127.0.0.1:0")
		_, err := client.ClaimReward(context.Background(), "user-1", "PROMO10")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not json"))
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		_, err := client.ClaimReward(context.Background(), "user-1", "PROMO10")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUserClientGetRewardDescription(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %s, want GET", r.Method)
			}
			if r.URL.Path != "/rewards/PROMO10" {
				t.Errorf("path = %s, want /rewards/PROMO10", r.URL.Path)
			}
			_ = json.NewEncoder(w).Encode(RewardDescriptionResponse{
				Code:        "PROMO10",
				Name:        "10% off",
				Description: "discount",
			})
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		reward, err := client.GetRewardDescription(context.Background(), "PROMO10")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if reward.PromoCode != "PROMO10" || reward.Name != "10% off" || reward.Description != "discount" {
			t.Errorf("reward = %+v, unexpected values", reward)
		}
	})

	t.Run("server unreachable", func(t *testing.T) {
		client := NewUserClient("http://127.0.0.1:0")
		_, err := client.GetRewardDescription(context.Background(), "PROMO10")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("not json"))
		}))
		defer server.Close()

		client := NewUserClient(server.URL)
		_, err := client.GetRewardDescription(context.Background(), "PROMO10")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
