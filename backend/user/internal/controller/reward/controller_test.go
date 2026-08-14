package reward

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/httpx"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/middleware"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/controller/reward/mocks"
	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

var errService = errors.New("service error")

func TestController_GetRewards(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		call       func(*controller, http.ResponseWriter, *http.Request)
		setupMocks func(*mocks.MockService)
	}{
		{
			name: "all rewards",
			call: func(c *controller, w http.ResponseWriter, r *http.Request) {
				c.GetAllRewards(w, r)
			},
			setupMocks: func(service *mocks.MockService) {
				service.EXPECT().
					GetAllRewards(gomock.Any(), "user-id").
					Return([]entity.UserReward{{}}, nil)
			},
		},
		{
			name: "active rewards",
			call: func(c *controller, w http.ResponseWriter, r *http.Request) {
				c.GetActiveRewards(w, r)
			},
			setupMocks: func(service *mocks.MockService) {
				service.EXPECT().
					GetActiveRewards(gomock.Any(), "user-id").
					Return([]entity.UserReward{{}}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			tt.setupMocks(deps.service)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/rewards", nil)

			serveAuthenticated(recorder, request, "user-id", func(w http.ResponseWriter, r *http.Request) {
				tt.call(deps.controller, w, r)
			})

			require.Equal(t, http.StatusOK, recorder.Code)

			var response userListRewardResponse
			require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
			require.Len(t, response.Items, 1)
		})
	}
}

func TestController_GetAllRewards_ServiceError(t *testing.T) {
	t.Parallel()

	deps := newTestDeps(t)
	deps.service.EXPECT().
		GetAllRewards(gomock.Any(), "user-id").
		Return(nil, errService)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/rewards", nil)

	serveAuthenticated(recorder, request, "user-id", deps.controller.GetAllRewards)

	checkErrorResponse(t, recorder, http.StatusInternalServerError, httpx.ErrInternal)
}

func TestController_GetReward(t *testing.T) {
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
			name:           "not found",
			serviceErr:     entity.ErrRewardNotFound,
			expectedStatus: http.StatusNotFound,
			expectedCode:   httpx.ErrRewardNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			deps.service.EXPECT().
				GetReward(gomock.Any(), "user-id", "reward-id").
				Return(&entity.UserReward{}, tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/rewards/reward-id", nil)
			request = mux.SetURLVars(request, map[string]string{"reward_id": "reward-id"})

			serveAuthenticated(recorder, request, "user-id", deps.controller.GetReward)

			if tt.serviceErr == nil {
				require.Equal(t, tt.expectedStatus, recorder.Code)
				return
			}

			checkErrorResponse(t, recorder, tt.expectedStatus, tt.expectedCode)
		})
	}
}

func TestController_GetDefinition(t *testing.T) {
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
			name:           "not found",
			serviceErr:     entity.ErrRewardDefinitionNotFound,
			expectedStatus: http.StatusNotFound,
			expectedCode:   httpx.ErrNotFoundDefinition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			deps.service.EXPECT().
				GetDefinition(gomock.Any(), "reward-code").
				Return(&entity.RewardDefinition{}, tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/rewards/reward-code", nil)
			request = mux.SetURLVars(request, map[string]string{"code": "reward-code"})

			deps.controller.GetDefinition(recorder, request)

			if tt.serviceErr == nil {
				require.Equal(t, tt.expectedStatus, recorder.Code)
				return
			}

			checkErrorResponse(t, recorder, tt.expectedStatus, tt.expectedCode)
		})
	}
}

func TestController_GrantReward(t *testing.T) {
	t.Parallel()

	deps := newTestDeps(t)
	deps.service.EXPECT().
		GrantReward(gomock.Any(), "user-id", "reward-code", "achievement").
		Return(&entity.UserReward{}, nil)

	recorder := httptest.NewRecorder()
	request := newJSONRequest(http.MethodPost, "/rewards", `{
		"user_id": "user-id",
		"code": "reward-code",
		"earned_reason": "achievement"
	}`)

	deps.controller.GrantReward(recorder, request)

	require.Equal(t, http.StatusCreated, recorder.Code)
}

func TestController_RedeemReward(t *testing.T) {
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
			name:           "reward unavailable",
			serviceErr:     entity.ErrRewardUnavailable,
			expectedStatus: http.StatusNotFound,
			expectedCode:   httpx.ErrRewardUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := newTestDeps(t)
			deps.service.EXPECT().
				RedeemReward(gomock.Any(), "user-id", "promo-code").
				Return(tt.serviceErr)

			recorder := httptest.NewRecorder()
			request := newJSONRequest(http.MethodPost, "/rewards/redeem", `{
				"promo_code": "promo-code"
			}`)

			serveAuthenticated(recorder, request, "user-id", deps.controller.RedeemReward)

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
