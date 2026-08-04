package user_service

import (
	"fmt"
	"net/http"
	"testing"

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

	loginResp := jsonReq(t, http.MethodPost, cfg.Users.BaseURL+"/auth/login", map[string]any{
		"email":    email,
		"password": testPassword,
	}, "")
	require.Equal(t, http.StatusOK, loginResp.StatusCode)

	loggedIn := decodeBody[authResponse](t, loginResp)
	require.NotEmpty(t, loggedIn.AccessToken)
	require.NotEmpty(t, loggedIn.RefreshToken)

	refreshResp := jsonReq(t, http.MethodPost, cfg.Users.BaseURL+"/auth/refresh", map[string]any{
		"refresh_token": loggedIn.RefreshToken,
	}, "")
	require.Equal(t, http.StatusOK, refreshResp.StatusCode)

	refreshed := decodeBody[authResponse](t, refreshResp)
	require.NotEmpty(t, refreshed.AccessToken)
	require.NotEmpty(t, refreshed.RefreshToken)
	require.NotEqual(t, loggedIn.RefreshToken, refreshed.RefreshToken)

	refreshedProfile := getProfile(t, cfg, refreshed.AccessToken)
	require.Equal(t, profile.UserID, refreshedProfile.UserID)

	reusedResp := jsonReq(t, http.MethodPost, cfg.Users.BaseURL+"/auth/refresh", map[string]any{
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

	resp := jsonReq(t, http.MethodPost, cfg.Users.BaseURL+"/auth/register", map[string]any{
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

	resp := jsonReq(t, http.MethodPost, cfg.Users.BaseURL+"/auth/login", map[string]any{
		"email":    email,
		"password": "WrongPassword123",
	}, "")
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	apiErr := decodeBody[apiError](t, resp)
	require.Equal(t, "INVALID_CREDENTIALS", apiErr.Code)
}

func TestProfileRequiresAccessToken(t *testing.T) {
	cfg := setup(t)

	resp := jsonReq(t, http.MethodGet, cfg.Users.BaseURL+"/profile", nil, "")
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

	loginResp := jsonReq(t, http.MethodPost, cfg.Users.BaseURL+"/auth/login", map[string]any{
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
			resp := rawReq(t, http.MethodPost, cfg.Users.BaseURL+"/auth/register", tc.body, "")
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
			resp := rawReq(t, http.MethodPost, cfg.Users.BaseURL+"/auth/login", tc.body, "")
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
			resp := rawReq(t, http.MethodPost, cfg.Users.BaseURL+"/auth/refresh", tc.body, "")
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
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, cfg.Users.BaseURL+"/profile", nil)
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
