package user_service

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestActionUpdatesStreakLifecycle(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("action-user-%s", suffix),
		fmt.Sprintf("action-%s@example.com", suffix),
		testPassword,
	)
	profile := getProfile(t, cfg, auth.AccessToken)
	createPet(t, cfg, auth.AccessToken, "Action Pet")

	initialStreak := getStreak(t, cfg, profile.UserID)
	baseDate := initialStreak.LastActiveDate

	require.Equal(t, 1, initialStreak.CurrentStreak)
	require.Equal(t, time.Now().UTC().Format(time.DateOnly), baseDate.Format(time.DateOnly))

	testCases := []struct {
		name           string
		occurredAt     string
		expectedStreak int
		expectedDate   string
	}{
		{
			name:           "same utc day keeps streak",
			occurredAt:     time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
			expectedStreak: 1,
			expectedDate:   baseDate.Format(time.DateOnly),
		},
		{
			name:           "next utc day increments streak",
			occurredAt:     time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day()+1, 9, 0, 0, 0, time.UTC).Format(time.RFC3339),
			expectedStreak: 2,
			expectedDate:   baseDate.AddDate(0, 0, 1).Format(time.DateOnly),
		},
		{
			name:           "skipped utc day resets streak",
			occurredAt:     time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day()+3, 9, 0, 0, 0, time.UTC).Format(time.RFC3339),
			expectedStreak: 1,
			expectedDate:   baseDate.AddDate(0, 0, 3).Format(time.DateOnly),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := jsonReq(t, http.MethodPost, cfg.Users.InternalURL+"/action", map[string]any{
				"user_id":     profile.UserID,
				"occurred_at": tc.occurredAt,
			}, "")
			require.Equal(t, http.StatusNoContent, resp.StatusCode)
			requireEmptyBody(t, resp)

			streak := getStreak(t, cfg, profile.UserID)
			require.Equal(t, tc.expectedStreak, streak.CurrentStreak)
			require.Equal(t, tc.expectedDate, streak.LastActiveDate.Format(time.DateOnly))
		})
	}
}

func TestActionStatuses(t *testing.T) {
	cfg := setup(t)

	t.Run("user not found", func(t *testing.T) {
		resp := jsonReq(t, http.MethodPost, cfg.Users.InternalURL+"/action", map[string]any{
			"user_id":     "00000000-0000-0000-0000-000000000000",
			"occurred_at": "2026-08-01T10:00:00Z",
		}, "")
		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		apiErr := decodeBody[apiError](t, resp)
		require.Equal(t, "USER_NOT_FOUND", apiErr.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		resp := rawReq(
			t,
			http.MethodPost,
			cfg.Users.InternalURL+"/action",
			`{"user_id":"1","occurred_at":"2026-08-01T10:00:00Z","extra":true}`,
			"",
		)
		require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		apiErr := decodeBody[apiError](t, resp)
		require.Equal(t, "VALIDATION_ERROR", apiErr.Code)
	})

	t.Run("invalid occurred at", func(t *testing.T) {
		resp := jsonReq(t, http.MethodPost, cfg.Users.InternalURL+"/action", map[string]any{
			"user_id":     "00000000-0000-0000-0000-000000000000",
			"occurred_at": "not-a-timestamp",
		}, "")
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		apiErr := decodeBody[apiError](t, resp)
		require.Equal(t, "INTERNAL_ERROR", apiErr.Code)
	})
}
