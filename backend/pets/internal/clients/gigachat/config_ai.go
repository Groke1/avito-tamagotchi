package gigachat

import (
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
}

func NewConfigFromEnv() *Config {
	return &Config{
		AuthKey:            os.Getenv("GIGACHAT_AUTH_KEY"),
		Scope:              envOrDefault("GIGACHAT_SCOPE", "GIGACHAT_API_PERS"),
		Model:              envOrDefault("GIGACHAT_MODEL", "GigaChat"),
		OAuthURL:           envOrDefault("GIGACHAT_OAUTH_URL", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"),
		APIBaseURL:         envOrDefault("GIGACHAT_API_BASE_URL", "https://gigachat.devices.sberbank.ru/api/v1"),
		RequestTimeout:     15 * time.Second,
		InsecureSkipVerify: os.Getenv("GIGACHAT_INSECURE_SKIP_VERIFY") == "true",
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
