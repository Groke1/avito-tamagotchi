package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

func TestSetAuthCookies_SetsBothCookiesWithCorrectAttributes(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	cfg := CookieConfig{AccessTokenTTL: time.Hour, RefreshTokenTTL: 24 * time.Hour}

	SetAuthCookies(rec, req, &entity.JWT{AccessToken: "access-val", RefreshToken: "refresh-val"}, cfg)

	cookies := rec.Result().Cookies()
	byName := map[string]*http.Cookie{}
	for _, c := range cookies {
		byName[c.Name] = c
	}

	access, ok := byName[AccessTokenCookieName]
	if !ok {
		t.Fatalf("expected access_token cookie to be set")
	}
	if access.Value != "access-val" {
		t.Errorf("unexpected access token value: %q", access.Value)
	}
	if !access.HttpOnly {
		t.Errorf("expected access_token cookie to be HttpOnly")
	}
	if access.MaxAge != int(time.Hour.Seconds()) {
		t.Errorf("expected MaxAge %d, got %d", int(time.Hour.Seconds()), access.MaxAge)
	}

	refresh, ok := byName[RefreshTokenCookieName]
	if !ok {
		t.Fatalf("expected refresh_token cookie to be set")
	}
	if refresh.Value != "refresh-val" {
		t.Errorf("unexpected refresh token value: %q", refresh.Value)
	}
	if refresh.MaxAge != int((24 * time.Hour).Seconds()) {
		t.Errorf("expected MaxAge %d, got %d", int((24*time.Hour).Seconds()), refresh.MaxAge)
	}
}

func TestSetAuthCookies_NilTokens_NoOp(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	SetAuthCookies(rec, req, nil, CookieConfig{})

	if len(rec.Result().Cookies()) != 0 {
		t.Fatalf("expected no cookies to be set for nil tokens")
	}
}

func TestClearAuthCookies_ExpiresBothCookies(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	ClearAuthCookies(rec, req)

	cookies := rec.Result().Cookies()
	byName := map[string]*http.Cookie{}
	for _, c := range cookies {
		byName[c.Name] = c
	}

	for _, name := range []string{AccessTokenCookieName, RefreshTokenCookieName} {
		c, ok := byName[name]
		if !ok {
			t.Fatalf("expected %s cookie to be present in the clearing response", name)
		}
		if c.MaxAge != -1 {
			t.Errorf("expected MaxAge -1 to clear %s, got %d", name, c.MaxAge)
		}
		if c.Value != "" {
			t.Errorf("expected empty value for cleared %s, got %q", name, c.Value)
		}
	}
}
