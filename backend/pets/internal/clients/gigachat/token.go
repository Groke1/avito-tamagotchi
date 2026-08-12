package gigachat

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// tokenManager кэширует access_token GigaChat и обновляет его по мере
// истечения срока действия. Использовать напрямую refresh token тут не
// нужно — GigaChat выдаёт короткоживущий access_token сразу по client
// credentials, каждый раз заново.
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

// getToken возвращает валидный access_token, обновляя его при необходимости.
// forceRefresh используется, когда предыдущий токен был отклонён (401).
func (tm *tokenManager) getToken(ctx context.Context, forceRefresh bool) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	//nolint:mnd // обновляем токен заранее, за минуту до истечения
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tm.cfg.OAuthURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", time.Time{}, err
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
		return "", time.Time{}, fmt.Errorf("gigachat oauth failed: status %d", resp.StatusCode)
	}

	var tokenResp oauthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", time.Time{}, fmt.Errorf("failed to decode gigachat oauth response: %w", err)
	}

	return tokenResp.AccessToken, time.UnixMilli(tokenResp.ExpiresAt), nil
}

// newRqUID генерирует уникальный идентификатор запроса в формате UUID v4,
// который GigaChat требует на каждый вызов /oauth.
func newRqUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)

	// проставляем версию (4) и вариант (RFC 4122), как того требует формат UUID
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
