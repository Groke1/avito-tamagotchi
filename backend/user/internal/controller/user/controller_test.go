package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/user/mocks"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestController_Profile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		serviceErr     error
		expectedStatus int
		expectedCode   httpx.ErrorCode
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user not found",
			serviceErr:     entity.ErrUserNotFound,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   httpx.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			deps.service.EXPECT().
				Profile(gomock.Any(), "user-id").
				Return(&entity.User{
					ID:            "user-id",
					Username:      "test-user",
					Email:         "user@example.com",
					Coins:         100,
					CurrentStreak: 5,
				}, tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/profile", nil)

			serveAuthenticated(recorder, request, "user-id", deps.controller.Profile)

			if tt.serviceErr == nil {
				require.Equal(t, tt.expectedStatus, recorder.Code)

				var response profileResponse
				require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
				require.Equal(t, "user-id", response.UserID)
				require.Equal(t, uint64(100), response.Coins)
				return
			}

			checkErrorResponse(t, recorder, tt.expectedStatus, tt.expectedCode)
		})
	}
}

func TestController_UpdateCoins(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		serviceErr     error
		expectedStatus int
		expectedCode   httpx.ErrorCode
	}{
		{
			name:           "success",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "insufficient coins",
			serviceErr:     entity.ErrInsufficientCoins,
			expectedStatus: http.StatusConflict,
			expectedCode:   httpx.ErrInsufficientCoins,
		},
		{
			name:           "user not found",
			serviceErr:     entity.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedCode:   httpx.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			deps.service.EXPECT().
				UpdateCoins(gomock.Any(), "user-id", int64(-50)).
				Return(&entity.User{
					ID:    "user-id",
					Coins: 50,
				}, tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := newJSONRequest(http.MethodPut, "/update-coins", `{
				"user_id": "user-id",
				"delta": -50
			}`)

			deps.controller.UpdateCoins(recorder, request)

			if tt.serviceErr == nil {
				require.Equal(t, tt.expectedStatus, recorder.Code)

				var response updateCoinsResponse
				require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
				require.Equal(t, "user-id", response.UserID)
				require.Equal(t, uint64(50), response.Coins)
				return
			}

			checkErrorResponse(t, recorder, tt.expectedStatus, tt.expectedCode)
		})
	}
}

func TestController_Usernames(t *testing.T) {
	t.Parallel()

	deps := newTestDeps(t)
	deps.service.EXPECT().
		GetUsers(gomock.Any(), []string{"user-1", "user-2"}).
		Return([]entity.User{
			{ID: "user-1", Username: "alice"},
			{ID: "user-2", Username: "bob"},
		}, nil)

	recorder := httptest.NewRecorder()
	request := newJSONRequest(http.MethodPost, "/usernames", `{
		"user_ids": ["user-1", "user-2"]
	}`)

	deps.controller.Usernames(recorder, request)

	require.Equal(t, http.StatusOK, recorder.Code)

	var response usernamesResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	require.Equal(t, []usernameResponse{
		{ID: "user-1", Username: "alice"},
		{ID: "user-2", Username: "bob"},
	}, response.Users)
}

func TestController_DailyStat(t *testing.T) {
	t.Parallel()

	deps := newTestDeps(t)
	expected := &entity.DailyStat{}

	deps.service.EXPECT().
		GetDailyStat(gomock.Any(), "user-id").
		Return(expected, nil)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/daily-stat", nil)

	serveAuthenticated(recorder, request, "user-id", deps.controller.DailyStat)

	require.Equal(t, http.StatusOK, recorder.Code)

	var response entity.DailyStat
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	require.Equal(t, *expected, response)
}

func TestController_Action(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		serviceErr     error
		expectedStatus int
		expectedCode   httpx.ErrorCode
	}{
		{
			name:           "success",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "user not found",
			serviceErr:     entity.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedCode:   httpx.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			deps.service.EXPECT().
				UpdateStreak(gomock.Any(), "user-id", "2026-08-14T10:00:00Z").
				Return(tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := newJSONRequest(http.MethodPost, "/action", `{
				"user_id": "user-id",
				"occurred_at": "2026-08-14T10:00:00Z"
			}`)

			deps.controller.Action(recorder, request)

			if tt.serviceErr == nil {
				require.Equal(t, tt.expectedStatus, recorder.Code)
				require.Empty(t, recorder.Body.String())
				return
			}

			checkErrorResponse(t, recorder, tt.expectedStatus, tt.expectedCode)
		})
	}
}

func TestController_ProtectedAction(t *testing.T) {
	t.Parallel()

	deps := newTestDeps(t)
	deps.service.EXPECT().
		UpdateStreak(gomock.Any(), "user-id", gomock.Any()).
		Return(nil)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/action", nil)

	serveAuthenticated(recorder, request, "user-id", deps.controller.ProtectedAction)

	require.Equal(t, http.StatusNoContent, recorder.Code)
	require.Empty(t, recorder.Body.String())
}

type testDeps struct {
	service    *mocks.MockService
	controller *controller
}

func newTestDeps(t *testing.T) testDeps {
	t.Helper()

	ctrl := gomock.NewController(t)
	service := mocks.NewMockService(ctrl)

	return testDeps{
		service:    service,
		controller: NewController(zap.NewNop(), service),
	}
}

func serveAuthenticated(
	recorder *httptest.ResponseRecorder,
	request *http.Request,
	userID string,
	handler http.HandlerFunc,
) {
	authHandler := middleware.RequireAccessToken(staticAccessTokenValidator{
		userID: userID,
	})(handler)

	request.Header.Set("Authorization", "Bearer access-token")
	authHandler.ServeHTTP(recorder, request)
}

func newJSONRequest(method, target, body string) *http.Request {
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func checkErrorResponse(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	expectedStatus int,
	expectedCode httpx.ErrorCode,
) {
	t.Helper()

	require.Equal(t, expectedStatus, recorder.Code)

	var response httpx.ErrorResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	require.Equal(t, string(expectedCode), response.Code)
}

type staticAccessTokenValidator struct {
	userID string
}

func (v staticAccessTokenValidator) ValidateAccessToken(context.Context, string) (string, error) {
	return v.userID, nil
}
