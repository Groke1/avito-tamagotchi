package gigachat

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	AuthKey      string
	ClientID     string
	ClientSecret string

	Scope string
	Model string

	OAuthURL   string
	APIBaseURL string

	RequestTimeout time.Duration
}

func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{
		AuthKey:      strings.TrimSpace(os.Getenv("GIGACHAT_AUTH_KEY")),
		ClientID:     strings.TrimSpace(os.Getenv("GIGACHAT_CLIENT_ID")),
		ClientSecret: strings.TrimSpace(os.Getenv("GIGACHAT_CLIENT_SECRET")),

		Scope:      envOrDefault("GIGACHAT_SCOPE", "GIGACHAT_API_PERS"),
		Model:      envOrDefault("GIGACHAT_MODEL", "GigaChat-2"),
		OAuthURL:   envOrDefault("GIGACHAT_OAUTH_URL", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"),
		APIBaseURL: envOrDefault("GIGACHAT_API_BASE_URL", "https://api.giga.chat/v1"),

		RequestTimeout: 15 * time.Second,
	}

	if err := cfg.prepareAuthKey(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) prepareAuthKey() error {
	if c.AuthKey != "" {
		c.AuthKey = strings.TrimPrefix(c.AuthKey, "Basic ")
		return nil
	}

	if c.ClientID == "" {
		return fmt.Errorf("GIGACHAT_CLIENT_ID is required")
	}

	if c.ClientSecret == "" {
		return fmt.Errorf("GIGACHAT_CLIENT_SECRET is required")
	}

	raw := c.ClientID + ":" + c.ClientSecret
	c.AuthKey = base64.StdEncoding.EncodeToString([]byte(raw))

	return nil
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}

	return def
}
