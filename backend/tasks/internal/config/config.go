package config

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v10"
)

const (
	ServerReadTimeout = 5 * time.Second
	HTTPClientTimeout = 3 * time.Second
)

type (
	Config struct {
		HTTP struct {
			Addr string `env:"HTTP_ADDR_T" envDefault:":8081"`
		}
		JWT struct {
			Secret string `env:"JWT_SECRET,required"`
		}

		PG struct {
			Host     string `env:"POSTGRES_HOST_T" envDefault:"tasks-postgres"`
			Port     string `env:"POSTGRES_PORT_T" envDefault:"5432"`
			DB       string `env:"POSTGRES_DB_T" envDefault:"tasks"`
			User     string `env:"POSTGRES_USER_T" envDefault:"tasks_user"`
			Password string `env:"POSTGRES_PASSWORD_T" envDefault:"12345"`
		}

		Client struct {
			CoinsServiceURL string `env:"USER_SERVICE_URL_T" envDefault:"http://users:8080"`
		}
		Pet struct {
			PetServiceURL string `env:"PET_SERVICE_URL_T" envDefault:"http://pets:8080"`
		}
	}
)

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
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
