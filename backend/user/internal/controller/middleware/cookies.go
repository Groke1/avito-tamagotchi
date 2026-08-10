package middleware

import (
	"net/http"
	"time"

	"github.com/cayman444/avito-gamification-hackathon.user/internal/entity"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

type CookieConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func SetAuthCookies(w http.ResponseWriter, r *http.Request, tokens *entity.JWT, cfg CookieConfig) {
	if tokens == nil {
		return
	}

	setCookie(w, r, AccessTokenCookieName, tokens.AccessToken, cfg.AccessTokenTTL)
	setCookie(w, r, RefreshTokenCookieName, tokens.RefreshToken, cfg.RefreshTokenTTL)
}

func ClearAuthCookies(w http.ResponseWriter, r *http.Request) {
	clearCookie(w, r, AccessTokenCookieName)
	clearCookie(w, r, RefreshTokenCookieName)
}

func setCookie(w http.ResponseWriter, r *http.Request, name, value string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().UTC().Add(ttl),
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
}

func clearCookie(w http.ResponseWriter, r *http.Request, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})
}
