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

	testCases := []struct {
		name           string
		occurredAt     string
		expectedStreak int
		expectedDate   string
	}{
		{
			name:           "first action creates streak",
			occurredAt:     "2026-08-01T08:15:00Z",
			expectedStreak: 1,
			expectedDate:   "2026-08-01",
		},
		{
			name:           "same moscow day keeps streak",
			occurredAt:     "2026-08-01T20:10:00Z",
			expectedStreak: 1,
			expectedDate:   "2026-08-01",
		},
		{
			name:           "next moscow day increments streak",
			occurredAt:     "2026-08-01T22:30:00Z",
			expectedStreak: 2,
			expectedDate:   "2026-08-02",
		},
		{
			name:           "skipped day resets streak",
			occurredAt:     "2026-08-04T09:00:00Z",
			expectedStreak: 1,
			expectedDate:   "2026-08-04",
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
		resp := rawReq(t, http.MethodPost, cfg.Users.InternalURL+"/action", `{"user_id":"1","occurred_at":"2026-08-01T10:00:00Z","extra":true}`, "")
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
