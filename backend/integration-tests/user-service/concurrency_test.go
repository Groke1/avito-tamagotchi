package user_service

import (
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcurrentRefreshOnlyOneSucceeds(t *testing.T) {
	cfg := setup(t)

	suffix := uniqueSuffix(t)
	auth := registerUser(
		t,
		cfg,
		fmt.Sprintf("refresh-race-%s", suffix),
		fmt.Sprintf("refresh-race-%s@example.com", suffix),
		testPassword,
	)

	const amount = 10

	start := make(chan struct{})
	statuses := make([]int, amount)

	var wg sync.WaitGroup
	for i := range amount {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			<-start

			resp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/refresh", map[string]any{
				"refresh_token": auth.RefreshToken,
			}, "")
			statuses[index] = resp.StatusCode
			_ = resp.Body.Close()
		}(i)
	}

	close(start)
	wg.Wait()

	var okCount int
	var unauthorizedCount int
	for _, status := range statuses {
		switch status {
		case http.StatusOK:
			okCount++
		case http.StatusUnauthorized:
			unauthorizedCount++
		default:
			t.Fatalf("unexpected refresh status: %d", status)
		}
	}

	require.Equal(t, 1, okCount)
	require.Equal(t, amount-1, unauthorizedCount)
}

func TestConcurrentRegisterLoginRefreshAndProfile(t *testing.T) {
	cfg := setup(t)

	const users = 50

	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < users; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			<-start

			suffix := fmt.Sprintf("%s-%d", uniqueSuffix(t), i)
			username := fmt.Sprintf("integration-user-%s", suffix)
			email := fmt.Sprintf("integration-%s@example.com", suffix)

			registered := registerUser(t, cfg, username, email, testPassword)
			require.NotEmpty(t, registered.AccessToken)
			require.NotEmpty(t, registered.RefreshToken)

			profile := getProfile(t, cfg, registered.AccessToken)
			require.Equal(t, username, profile.Username)
			require.Equal(t, email, profile.Email)

			loginResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/login", map[string]any{
				"email":    email,
				"password": testPassword,
			}, "")
			require.Equal(t, http.StatusOK, loginResp.StatusCode)

			loggedIn := decodeBody[authResponse](t, loginResp)

			refreshResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/refresh", map[string]any{
				"refresh_token": loggedIn.RefreshToken,
			}, "")
			require.Equal(t, http.StatusOK, refreshResp.StatusCode)

			refreshed := decodeBody[authResponse](t, refreshResp)
			require.NotEqual(t, loggedIn.RefreshToken, refreshed.RefreshToken)

			refreshedProfile := getProfile(t, cfg, refreshed.AccessToken)
			require.Equal(t, profile.UserID, refreshedProfile.UserID)

			reusedResp := jsonReq(t, http.MethodPost, cfg.Users.APIURL+"/auth/refresh", map[string]any{
				"refresh_token": loggedIn.RefreshToken,
			}, "")
			require.Equal(t, http.StatusUnauthorized, reusedResp.StatusCode)

		}(i)
	}

	close(start)
	wg.Wait()
}
