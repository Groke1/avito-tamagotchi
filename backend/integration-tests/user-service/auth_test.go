package user_service

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const testPassword = "StrongPass123"

func TestRegisterLoginRefreshAndProfile(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	username := fmt.Sprintf("integration-user-%s", suffix)
	email := fmt.Sprintf("integration-%s@example.com", suffix)

	registered := registerUser(t, cfg, username, email, testPassword)
	require.NotEmpty(t, registered.AccessToken)
	require.NotEmpty(t, registered.RefreshToken)

	profile := getProfile(t, cfg, registered.AccessToken)
	require.NotEmpty(t, profile.UserID)
	require.Equal(t, username, profile.Username)
	require.Equal(t, email, profile.Email)

	streak := getStreak(t, cfg, profile.UserID)
	require.Equal(t, 1, streak.CurrentStreak)
	require.Equal(t, time.Now().UTC().Format(time.DateOnly), streak.LastActiveDate.Format(time.DateOnly))

	initialRewardsResp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/rewards", nil, registered.AccessToken)
	require.Equal(t, http.StatusOK, initialRewardsResp.StatusCode)
	initialRewards := decodeBody[userRewardsResponse](t, initialRewardsResp)
	require.Len(t, initialRewards.Items, 1)
	require.Equal(t, "active", initialRewards.Items[0].Status)
	require.NotEmpty(t, initialRewards.Items[0].PromoCode)
	require.NotEmpty(t, initialRewards.Items[0].Name)

	loginResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/login", map[string]any{
		"email":    email,
		"password": testPassword,
	}, "")
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	loggedIn := decodeBody[authResponse](t, loginResp)
	require.NotEmpty(t, loggedIn.AccessToken)
	require.NotEmpty(t, loggedIn.RefreshToken)

	refreshResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/refresh", map[string]any{
		"refresh_token": loggedIn.RefreshToken,
	}, "")
	require.Equal(t, http.StatusOK, refreshResp.StatusCode)

	refreshed := decodeBody[authResponse](t, refreshResp)
	require.NotEmpty(t, refreshed.AccessToken)
	require.NotEmpty(t, refreshed.RefreshToken)
	require.NotEqual(t, loggedIn.RefreshToken, refreshed.RefreshToken)

	refreshedProfile := getProfile(t, cfg, refreshed.AccessToken)
	require.Equal(t, profile.UserID, refreshedProfile.UserID)

	reusedResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/refresh", map[string]any{
		"refresh_token": loggedIn.RefreshToken,
	}, "")
	require.Equal(t, http.StatusUnauthorized, reusedResp.StatusCode)

	reusedErr := decodeBody[apiError](t, reusedResp)
	require.Equal(t, "INVALID_REFRESH_TOKEN", reusedErr.Code)
}

func TestRegisterDuplicateUser(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	username := fmt.Sprintf("duplicate-user-%s", suffix)
	email := fmt.Sprintf("duplicate-%s@example.com", suffix)

	registerUser(t, cfg, username, email, testPassword)

	resp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/register", map[string]any{
		"username": username,
		"email":    email,
		"password": testPassword,
	}, "")
	require.Equal(t, http.StatusConflict, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "USER_ALREADY_EXISTS", apiErr.Code)
}

func TestLoginWithWrongPassword(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	email := fmt.Sprintf("wrong-password-%s@example.com", suffix)
	registerUser(t, cfg, fmt.Sprintf("wrong-password-%s", suffix), email, testPassword)

	resp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/login", map[string]any{
		"email":    email,
		"password": "WrongPassword123",
	}, "")
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "INVALID_CREDENTIALS", apiErr.Code)
}

func TestProfileRequiresAccessToken(t *testing.T) {
	cfg := setup(t)

	resp := jsonReq(t, http.MethodGet, cfg.Users.APIURL+"/profile", nil, "")
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "UNAUTHORIZED", apiErr.Code)
}

func TestRegisterNormalizesUsernameAndEmail(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	rawUsername := fmt.Sprintf("  normalized-user-%s  ", suffix)
	rawEmail := fmt.Sprintf("  NORMALIZED-%s@EXAMPLE.COM  ", suffix)

	auth := registerUser(t, cfg, rawUsername, rawEmail, testPassword)
	profile := getProfile(t, cfg, auth.AccessToken)

	require.Equal(t, fmt.Sprintf("normalized-user-%s", suffix), profile.Username)
	require.Equal(t, fmt.Sprintf("normalized-%s@example.com", suffix), profile.Email)

	loginResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/login", map[string]any{
		"email":    rawEmail,
		"password": testPassword,
	}, "")
	require.Equal(t, http.StatusOK, loginResp.StatusCode)
	_ = decodeBody[authResponse](t, loginResp)
}

func TestRegisterValidationCases(t *testing.T) {
	cfg := setup(t)

	testCases := []struct {
		name string
		body string
	}{
		{
			name: "too short username",
			body: `{"username":"a","email":"short-` + uniqueSuffix(t) + `@example.com","password":"` + testPassword + `"}`,
		},
		{
			name: "invalid email",
			body: `{"username":"invalid-email-` + uniqueSuffix(t) + `","email":"not-an-email","password":"` + testPassword + `"}`,
		},
		{
			name: "short password",
			body: `{"username":"short-password-` + uniqueSuffix(t) + `","email":"short-password-` + uniqueSuffix(t) + `@example.com","password":"short"}`,
		},
		{
			name: "unknown field",
			body: `{"username":"unknown-field-` + uniqueSuffix(t) + `","email":"unknown-field-` + uniqueSuffix(t) + `@example.com","password":"` + testPassword + `","extra":1}`,
		},
		{
			name: "multiple json objects",
			body: `{"username":"multi-json-` + uniqueSuffix(t) + `","email":"multi-json-` + uniqueSuffix(t) + `@example.com","password":"` + testPassword + `"}` +
				`{"username":"second"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := rawReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/register", tc.body, "")
			require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

			apiErr := decodeBody[apiError](t, resp)
			require.Equal(t, "VALIDATION_ERROR", apiErr.Code)
		})
	}
}

func TestLoginValidationCases(t *testing.T) {
	cfg := setup(t)

	testCases := []struct {
		name string
		body string
	}{
		{
			name: "invalid email",
			body: `{"email":"bad-email","password":"` + testPassword + `"}`,
		},
		{
			name: "empty password",
			body: `{"email":"` + fmt.Sprintf("empty-password-%s@example.com", uniqueSuffix(t)) + `","password":""}`,
		},
		{
			name: "unknown field",
			body: `{"email":"` + fmt.Sprintf("unknown-login-%s@example.com", uniqueSuffix(t)) + `","password":"` + testPassword + `","extra":true}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := rawReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/login", tc.body, "")
			require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

			apiErr := decodeBody[apiError](t, resp)
			require.Equal(t, "VALIDATION_ERROR", apiErr.Code)
		})
	}
}

func TestRefreshValidationCases(t *testing.T) {
	cfg := setup(t)

	testCases := []struct {
		name string
		body string
	}{
		{name: "empty body", body: `{}`},
		{name: "blank refresh token", body: `{"refresh_token":"   "}`},
		{name: "unknown field", body: `{"refresh_token":"token","extra":1}`},
		{name: "multiple json objects", body: `{"refresh_token":"token"}{"refresh_token":"another"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := rawReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/refresh", tc.body, "")
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			apiErr := decodeBody[apiError](t, resp)
			require.Equal(t, "INVALID_REFRESH_TOKEN", apiErr.Code)
		})
	}
}

func TestProfileRejectsMalformedAuthorization(t *testing.T) {
	cfg := setup(t)

	testCases := []struct {
		name          string
		authorization string
	}{
		{name: "basic scheme", authorization: "Basic token"},
		{name: "bearer without token", authorization: "Bearer"},
		{name: "empty bearer token", authorization: "Bearer   "},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, cfg.Users.APIURL+"/profile", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", tc.authorization)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

			apiErr := decodeBody[apiError](t, resp)
			require.Equal(t, "UNAUTHORIZED", apiErr.Code)
		})
	}
}

func TestLogoutInvalidatesRefreshToken(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("logout-user-%s", suffix),
		fmt.Sprintf("logout-%s@example.com", suffix),
		testPassword,
	)

	logoutResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/logout", map[string]any{
		"refresh_token": auth.RefreshToken,
	}, auth.AccessToken)
	require.Equal(t, http.StatusNoContent, logoutResp.StatusCode)
	requireEmptyBody(t, logoutResp)

	refreshResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/refresh", map[string]any{
		"refresh_token": auth.RefreshToken,
	}, "")
	require.Equal(t, http.StatusUnauthorized, refreshResp.StatusCode)

	apiErr := decodeBody[apiError](t, refreshResp)
	require.Equal(t, "INVALID_REFRESH_TOKEN", apiErr.Code)
}

func TestLogoutIsIdempotent(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("logout-idempotent-%s", suffix),
		fmt.Sprintf("logout-idempotent-%s@example.com", suffix),
		testPassword,
	)

	for attempt := 0; attempt < 2; attempt++ {
		resp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/logout", map[string]any{
			"refresh_token": auth.RefreshToken,
		}, auth.AccessToken)
		require.Equalf(t, http.StatusNoContent, resp.StatusCode, "attempt=%d", attempt+1)
		requireEmptyBody(t, resp)
	}
}

func TestLogoutDoesNotRevokeAnotherUsersSession(t *testing.T) {
	cfg := setup(t)

	firstSuffix := uniqueSuffix(t)
	firstAuth := registerUser(
		t,
		cfg,
		fmt.Sprintf("logout-owner-%s", firstSuffix),
		fmt.Sprintf("logout-owner-%s@example.com", firstSuffix),
		testPassword,
	)

	secondSuffix := uniqueSuffix(t)
	secondAuth := registerUser(
		t,
		cfg,
		fmt.Sprintf("logout-attacker-%s", secondSuffix),
		fmt.Sprintf("logout-attacker-%s@example.com", secondSuffix),
		testPassword,
	)

	logoutResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/logout", map[string]any{
		"refresh_token": firstAuth.RefreshToken,
	}, secondAuth.AccessToken)
	require.Equal(t, http.StatusNoContent, logoutResp.StatusCode)
	requireEmptyBody(t, logoutResp)

	refreshResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/refresh", map[string]any{
		"refresh_token": firstAuth.RefreshToken,
	}, "")
	require.Equal(t, http.StatusOK, refreshResp.StatusCode)
	_ = decodeBody[authResponse](t, refreshResp)
}

func TestLogoutRequiresAccessToken(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("logout-no-auth-%s", suffix),
		fmt.Sprintf("logout-no-auth-%s@example.com", suffix),
		testPassword,
	)

	resp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/logout", map[string]any{
		"refresh_token": auth.RefreshToken,
	}, "")
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "UNAUTHORIZED", apiErr.Code)
}

func TestLogoutValidationCases(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("logout-validation-%s", suffix),
		fmt.Sprintf("logout-validation-%s@example.com", suffix),
		testPassword,
	)

	testCases := []struct {
		name string
		body string
	}{
		{name: "empty body", body: `{}`},
		{name: "blank refresh token", body: `{"refresh_token":"   "}`},
		{name: "unknown field", body: `{"refresh_token":"token","extra":1}`},
		{name: "multiple json objects", body: `{"refresh_token":"token"}{"refresh_token":"another"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := rawReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/logout", tc.body, auth.AccessToken)
			require.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

			apiErr := decodeBody[apiError](t, resp)
			require.Equal(t, "VALIDATION_ERROR", apiErr.Code)
		})
	}
}
