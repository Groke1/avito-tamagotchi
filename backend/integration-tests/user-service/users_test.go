package user_service

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type usernamesResponse struct {
	Users []struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"users"`
}

func TestGetUsernames(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	firstUsername := fmt.Sprintf("first-user-%s", suffix)
	secondUsername := fmt.Sprintf("second-user-%s", suffix)

	firstAuth := registerUser(
		t,
		cfg,
		firstUsername,
		fmt.Sprintf("first-%s@example.com", suffix),
		testPassword,
	)
	secondAuth := registerUser(
		t,
		cfg,
		secondUsername,
		fmt.Sprintf("second-%s@example.com", suffix),
		testPassword,
	)

	firstProfile := getProfile(t, cfg, firstAuth.AccessToken)
	secondProfile := getProfile(t, cfg, secondAuth.AccessToken)

	resp := jsonReq(t, http.MethodPost, cfg.Users.BaseURL+"/usernames", map[string]any{
		"user_ids": []string{firstProfile.UserID, secondProfile.UserID},
	}, "")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := decodeBody[usernamesResponse](t, resp)
	require.Len(t, body.Users, 2)

	usersByID := make(map[string]string, len(body.Users))
	for _, user := range body.Users {
		usersByID[user.ID] = user.Username
	}

	require.Equal(t, firstUsername, usersByID[firstProfile.UserID])
	require.Equal(t, secondUsername, usersByID[secondProfile.UserID])
}

func TestGetUsernamesManyUsers(t *testing.T) {
	cfg := setup(t)

	const usersCount = 100
	suffix := uniqueSuffix(t)

	type expectedUser struct {
		ID       string
		Username string
	}

	expected := make([]expectedUser, 0, usersCount)
	userIDs := make([]string, 0, usersCount)

	for i := 0; i < usersCount; i++ {
		username := fmt.Sprintf("user-%03d-%s", i, suffix)
		email := fmt.Sprintf("user-%03d-%s@example.com", i, suffix)

		auth := registerUser(
			t,
			cfg,
			username,
			email,
			testPassword,
		)

		profile := getProfile(t, cfg, auth.AccessToken)

		expected = append(expected, expectedUser{
			ID:       profile.UserID,
			Username: username,
		})
		userIDs = append(userIDs, profile.UserID)
	}

	resp := jsonReq(t, http.MethodPost, cfg.Users.BaseURL+"/usernames", map[string]any{
		"user_ids": userIDs,
	}, "")
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body := decodeBody[usernamesResponse](t, resp)
	require.Len(t, body.Users, usersCount)

	usersByID := make(map[string]string, len(body.Users))
	for _, user := range body.Users {
		usersByID[user.ID] = user.Username
	}

	for _, exp := range expected {
		require.Equal(t, exp.Username, usersByID[exp.ID])
	}
}

func TestUpdateCoins(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("coins-user-%s", suffix),
		fmt.Sprintf("coins-%s@example.com", suffix),
		testPassword,
	)
	profileBefore := getProfile(t, cfg, auth.AccessToken)

	addResp := jsonReq(t, http.MethodPut, cfg.Users.BaseURL+"/update-coins", map[string]any{
		"user_id": profileBefore.UserID,
		"delta":   15,
	}, "")
	require.Equal(t, http.StatusNoContent, addResp.StatusCode)
	requireEmptyBody(t, addResp)

	profileAfterAdd := getProfile(t, cfg, auth.AccessToken)
	require.Equal(t, profileBefore.Coins+15, profileAfterAdd.Coins)

	subtractResp := jsonReq(t, http.MethodPut, cfg.Users.BaseURL+"/update-coins", map[string]any{
		"user_id": profileBefore.UserID,
		"delta":   -5,
	}, "")
	require.Equal(t, http.StatusNoContent, subtractResp.StatusCode)
	requireEmptyBody(t, subtractResp)

	profileAfterSubtract := getProfile(t, cfg, auth.AccessToken)
	require.Equal(t, profileBefore.Coins+10, profileAfterSubtract.Coins)
}

func TestUpdateCoinsRejectsInsufficientBalance(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("poor-user-%s", suffix),
		fmt.Sprintf("poor-%s@example.com", suffix),
		testPassword,
	)
	profile := getProfile(t, cfg, auth.AccessToken)

	resp := jsonReq(t, http.MethodPut, cfg.Users.BaseURL+"/update-coins", map[string]any{
		"user_id": profile.UserID,
		"delta":   -int64(profile.Coins) - 1,
	}, "")
	require.Equal(t, http.StatusConflict, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "INSUFFICIENT_COINS", apiErr.Code)

	profileAfter := getProfile(t, cfg, auth.AccessToken)
	require.Equal(t, profile.Coins, profileAfter.Coins)
}

func TestUpdateCoinsUserNotFound(t *testing.T) {
	cfg := setup(t)

	resp := jsonReq(t, http.MethodPut, cfg.Users.BaseURL+"/update-coins", map[string]any{
		"user_id": "00000000-0000-0000-0000-000000000000",
		"delta":   10,
	}, "")
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "USER_NOT_FOUND", apiErr.Code)
}

func TestGetUsernamesValidationCases(t *testing.T) {
	cfg := setup(t)

	tooManyIDs := make([]string, 1001)
	for i := range tooManyIDs {
		tooManyIDs[i] = fmt.Sprintf("id-%d", i)
	}

	testCases := []struct {
		name string
		body any
		raw  string
	}{
		{name: "empty ids", body: map[string]any{"user_ids": []string{}}},
		{name: "too many ids", body: map[string]any{"user_ids": tooManyIDs}},
		{name: "unknown field", raw: `{"user_ids":["1"],"extra":1}`},
		{name: "multiple json objects", raw: `{"user_ids":["1"]}{"user_ids":["2"]}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var resp *http.Response
			if tc.raw != "" {
				resp = rawReq(t, http.MethodPost, cfg.Users.BaseURL+"/usernames", tc.raw, "")
			} else {
				resp = jsonReq(t, http.MethodPost, cfg.Users.BaseURL+"/usernames", tc.body, "")
			}

			require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
			apiErr := decodeBody[apiError](t, resp)
			require.Equal(t, "VALIDATION_ERROR", apiErr.Code)
		})
	}
}

func TestUpdateCoinsRandomizedSequence(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("random-coins-user-%s", suffix),
		fmt.Sprintf("random-coins-%s@example.com", suffix),
		testPassword,
	)
	profile := getProfile(t, cfg, auth.AccessToken)
	expectedCoins := int64(profile.Coins)

	rng := rand.New(rand.NewSource(42))
	for step := 0; step < 20; step++ {
		delta := int64(rng.Intn(41) - 20)
		if delta < 0 && -delta > expectedCoins {
			delta = -expectedCoins
		}

		resp := jsonReq(t, http.MethodPut, cfg.Users.BaseURL+"/update-coins", map[string]any{
			"user_id": profile.UserID,
			"delta":   delta,
		}, "")
		require.Equalf(t, http.StatusNoContent, resp.StatusCode, "step=%d delta=%d", step, delta)
		requireEmptyBody(t, resp)

		expectedCoins += delta

		currentProfile := getProfile(t, cfg, auth.AccessToken)
		require.Equalf(t, uint64(expectedCoins), currentProfile.Coins, "step=%d delta=%d", step, delta)
	}
}

func TestUpdateCoinsValidationCases(t *testing.T) {
	cfg := setup(t)

	testCases := []struct {
		name           string
		raw            string
		expectedStatus int
		expectedCode   string
	}{
		{name: "empty body", raw: `{}`, expectedStatus: http.StatusInternalServerError, expectedCode: "INTERNAL_ERROR"},
		{name: "unknown field", raw: `{"user_id":"1","delta":1,"extra":true}`, expectedStatus: http.StatusUnprocessableEntity, expectedCode: "VALIDATION_ERROR"},
		{name: "multiple json objects", raw: `{"user_id":"1","delta":1}{"user_id":"2","delta":2}`, expectedStatus: http.StatusUnprocessableEntity, expectedCode: "VALIDATION_ERROR"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := rawReq(t, http.MethodPut, cfg.Users.BaseURL+"/update-coins", tc.raw, "")
			require.Equal(t, tc.expectedStatus, resp.StatusCode)

			apiErr := decodeBody[apiError](t, resp)
			require.Equal(t, tc.expectedCode, apiErr.Code)
		})
	}
}
