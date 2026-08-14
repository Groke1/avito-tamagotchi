package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/auth/mocks"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var errService = errors.New("service error")

func TestController_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		serviceErr     error
		expectedStatus int
		expectedCode   httpx.ErrorCode
	}{
		{
			name:           "success",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "user already exists",
			serviceErr:     entity.ErrUsernameAlreadyExists,
			expectedStatus: http.StatusConflict,
			expectedCode:   httpx.ErrUserAlreadyExists,
		},
		{
			name:           "service error",
			serviceErr:     errService,
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   httpx.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tokens := &entity.JWT{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			}

			deps.service.EXPECT().
				Register(gomock.Any(), entity.User{
					Username: "test-user",
					Email:    "user@example.com",
					Password: "password123",
				}).
				Return(tokens, tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := newJSONRequest(http.MethodPost, "/auth/register", `{
				"username": "test-user",
				"email": "user@example.com",
				"password": "password123"
			}`)

			deps.controller.Register(recorder, request)

			if tt.serviceErr == nil {
				checkTokensResponse(t, recorder, tt.expectedStatus, tokens)
				return
			}

			checkErrorResponse(t, recorder, tt.expectedStatus, tt.expectedCode)
		})
	}
}

func TestController_Login(t *testing.T) {
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
			name:           "invalid credentials",
			serviceErr:     entity.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   httpx.ErrInvalidCredentials,
		},
		{
			name:           "service error",
			serviceErr:     errService,
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   httpx.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tokens := &entity.JWT{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
			}

			deps.service.EXPECT().
				Login(gomock.Any(), "user@example.com", "password123").
				Return(tokens, tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := newJSONRequest(http.MethodPost, "/auth/login", `{
				"email": "USER@EXAMPLE.COM",
				"password": "password123"
			}`)

			deps.controller.Login(recorder, request)

			if tt.serviceErr == nil {
				checkTokensResponse(t, recorder, tt.expectedStatus, tokens)
				return
			}

			checkErrorResponse(t, recorder, tt.expectedStatus, tt.expectedCode)
		})
	}
}

func TestController_Refresh(t *testing.T) {
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
			name:           "invalid refresh token",
			serviceErr:     entity.ErrInvalidRefreshToken,
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   httpx.ErrInvalidRefreshToken,
		},
		{
			name:           "service error",
			serviceErr:     errService,
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   httpx.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tokens := &entity.JWT{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
			}

			deps.service.EXPECT().
				Refresh(gomock.Any(), "refresh-token").
				Return(tokens, tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := newJSONRequest(http.MethodPost, "/auth/refresh", `{
				"refresh_token": "refresh-token"
			}`)

			deps.controller.Refresh(recorder, request)

			if tt.serviceErr == nil {
				checkTokensResponse(t, recorder, tt.expectedStatus, tokens)
				return
			}

			checkErrorResponse(t, recorder, tt.expectedStatus, tt.expectedCode)
		})
	}
}

func TestController_Logout(t *testing.T) {
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
			name:           "service error",
			serviceErr:     errService,
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   httpx.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)

			deps.service.EXPECT().
				Logout(gomock.Any(), "user-id", "refresh-token").
				Return(tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := newJSONRequest(http.MethodPost, "/auth/logout", `{
				"refresh_token": "refresh-token"
			}`)

			handler := middleware.RequireAccessToken(staticAccessTokenValidator{
				userID: "user-id",
			})(http.HandlerFunc(deps.controller.Logout))

			request.Header.Set("Authorization", "Bearer access-token")
			handler.ServeHTTP(recorder, request)

			if tt.serviceErr == nil {
				require.Equal(t, tt.expectedStatus, recorder.Code)
				require.Empty(t, recorder.Body.String())
				return
			}

			checkErrorResponse(t, recorder, tt.expectedStatus, tt.expectedCode)
		})
	}
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
		service: service,
		controller: NewController(
			zap.NewNop(),
			service,
			time.Hour,
			24*time.Hour,
		),
	}
}

func newJSONRequest(method, target, body string) *http.Request {
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func checkTokensResponse(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	expectedStatus int,
	expectedTokens *entity.JWT,
) {
	t.Helper()

	require.Equal(t, expectedStatus, recorder.Code)
	require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var response tokensResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	require.Equal(t, expectedTokens.AccessToken, response.AccessToken)
	require.Equal(t, expectedTokens.RefreshToken, response.RefreshToken)
}

func checkErrorResponse(
	t *testing.T,
	recorder *httptest.ResponseRecorder,
	expectedStatus int,
	expectedCode httpx.ErrorCode,
) {
	t.Helper()

	require.Equal(t, expectedStatus, recorder.Code)
	require.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

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
