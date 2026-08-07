package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HttpPort       string `env:"HTTP_PORT,required"`
	DatabaseURL    string `env:"DATABASE_URL,required"`
	JwtSecret      string `env:"JWT_SECRET,required"`
	UserServiceURL string `env:"USER_SERVICE_URL"`
}

func NewConfig() *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	return &cfg
}
