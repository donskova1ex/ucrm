package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	Auth AuthConfig
}

type DBConfig struct {
	DSN string
}

type AuthConfig struct {
	Secret string
}

func NewConfig(logger *slog.Logger) (*Config, error) {
	err := godotenv.Load(".env.local")
	if err != nil {
	}
	return &Config{}, nil
}
