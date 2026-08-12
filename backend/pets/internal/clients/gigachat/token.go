package gigachat

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type tokenManager struct {
	cfg    *Config
	client *http.Client

	mu        sync.Mutex
	token     string
	expiresAt time.Time
}

func newTokenManager(cfg *Config, client *http.Client) *tokenManager {
	return &tokenManager{cfg: cfg, client: client}
}

func (tm *tokenManager) getToken(ctx context.Context, forceRefresh bool) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if !forceRefresh && tm.token != "" && time.Now().Before(tm.expiresAt.Add(-1*time.Minute)) {
		return tm.token, nil
	}

	token, expiresAt, err := tm.fetchToken(ctx)
	if err != nil {
		return "", err
	}

	tm.token = token
	tm.expiresAt = expiresAt

	return tm.token, nil
}

func (tm *tokenManager) fetchToken(ctx context.Context) (string, time.Time, error) {
	form := url.Values{}
	form.Set("scope", tm.cfg.Scope)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tm.cfg.OAuthURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("create gigachat oauth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", newRqUID())
	req.Header.Set("Authorization", "Basic "+tm.cfg.AuthKey)

	resp, err := tm.client.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("gigachat oauth unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return "", time.Time{}, fmt.Errorf(
				"gigachat oauth failed: status=%d; read response: %w",
				resp.StatusCode,
				readErr,
			)
		}

		return "", time.Time{}, fmt.Errorf(
			"gigachat oauth failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}

	var tokenResp oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", time.Time{}, fmt.Errorf(
			"decode gigachat oauth response: %w",
			err,
		)
	}

	if tokenResp.AccessToken == "" {
		return "", time.Time{}, fmt.Errorf("gigachat oauth returned empty access token")
	}

	return tokenResp.AccessToken, time.UnixMilli(tokenResp.ExpiresAt), nil
}

func newRqUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
