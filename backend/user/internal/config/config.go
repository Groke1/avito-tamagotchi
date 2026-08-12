package config

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v10"
)

type (
	SessionStore string

	Config struct {
		HTTP struct {
			Port string `env:"HTTP_PORT" envDefault:"8080"`
		}

		PG struct {
			Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
			Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
			DB       string `env:"POSTGRES_DB" envDefault:"db"`
			User     string `env:"POSTGRES_USER" envDefault:"db_user"`
			Password string `env:"POSTGRES_PASSWORD" envDefault:"12345"`
			MaxConns int32  `env:"POSTGRES_MAX_CONNS" envDefault:"8"`
			MinConns int32  `env:"POSTGRES_MIN_CONNS" envDefault:"1"`
		}

		Redis struct {
			Host     string `env:"REDIS_HOST" envDefault:"localhost"`
			Port     string `env:"REDIS_PORT" envDefault:"6379"`
			Password string `env:"REDIS_PASSWORD"`
			DB       int    `env:"REDIS_DB" envDefault:"0"`
		}

		Clients struct {
			PetsHTTPAddr  string `env:"PETS_HTTP_ADDR" envDefault:"pets:8080"`
			TasksHTTPAddr string `env:"TASKS_HTTP_ADDR" envDefault:"tasks:8080"`
		}

		Settings struct {
			RegistrationBonusCoins  uint64        `env:"REGISTRATION_BONUS_COINS" envDefault:"100"`
			AccessTokenTTL          time.Duration `env:"ACCESS_TOKEN_TTL" envDefault:"15m"`
			RefreshTokenTTL         time.Duration `env:"REFRESH_TOKEN_TTL" envDefault:"720h"`
			SessionStore            SessionStore  `env:"SESSION_STORE" envDefault:"postgres"`
			TokenCleanupInterval    time.Duration `env:"TOKEN_CLEANUP_INTERVAL" envDefault:"24h"`
			EventsCleanupInterval   time.Duration `env:"EVENTS_CLEANUP_INTERVAL" envDefault:"24h"`
			ShutdownTimeout         time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
			JWTSecret               string        `env:"JWT_SECRET,required"`
			PGHealthCheckPeriod     time.Duration `env:"POSTGRES_HEALTH_CHECK_PERIOD" envDefault:"30s"`
			PGMaxConnLifetime       time.Duration `env:"POSTGRES_MAX_CONN_LIFETIME" envDefault:"1h"`
			PGMaxConnIdleTime       time.Duration `env:"POSTGRES_MAX_CONN_IDLE_TIME" envDefault:"5m"`
			ServerReadHeaderTimeout time.Duration `env:"SERVER_READ_HEADER_TIMEOUT" envDefault:"5s"`
			HTTPClientTimeout       time.Duration `env:"HTTP_CLIENT_TIMEOUT" envDefault:"5s"`
		}
	}
)

const (
	SessionStorePostgres SessionStore = "postgres"
	SessionStoreRedis    SessionStore = "redis"
)

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	if !cfg.Settings.SessionStore.Valid() {
		return nil, fmt.Errorf("invalid SESSION_STORE %q", cfg.Settings.SessionStore)
	}

	return &cfg, nil
}

func (c *Config) ConstructPostgresURL() string {
	hostPort := net.JoinHostPort(c.PG.Host, c.PG.Port)
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.PG.User,
		c.PG.Password,
		hostPort,
		c.PG.DB,
	)
}

func (c *Config) ConstructRedisAddr() string {
	return net.JoinHostPort(c.Redis.Host, c.Redis.Port)
}

func (s SessionStore) Valid() bool {
	switch s {
	case SessionStorePostgres, SessionStoreRedis:
		return true
	default:
		return false
	}
}
