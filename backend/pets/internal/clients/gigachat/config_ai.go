package gigachat

import (
	"net/http"
	"os"
	"time"
)

type Config struct {
	AuthKey string
	Scope   string
	Model   string

	OAuthURL   string
	APIBaseURL string

	RequestTimeout time.Duration

	InsecureSkipVerify bool
	HTTPClient         *http.Client
}

func NewConfigFromEnv() *Config {
	return &Config{
		AuthKey:            os.Getenv("GIGACHAT_AUTH_KEY"),
		Scope:              envOrDefault("GIGACHAT_SCOPE", "GIGACHAT_API_PERS"),
		Model:              envOrDefault("GIGACHAT_MODEL", "GigaChat"),
		OAuthURL:           envOrDefault("GIGACHAT_OAUTH_URL", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"),
		APIBaseURL:         envOrDefault("GIGACHAT_API_BASE_URL", "https://api.giga.chat"),
		RequestTimeout:     15 * time.Second,
		InsecureSkipVerify: envOrDefault("GIGACHAT_INSECURE_SKIP_VERIFY", "false") == "true",
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
