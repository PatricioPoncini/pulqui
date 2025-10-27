package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	DolarApiUrl   string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		DolarApiUrl:   os.Getenv("DOLAR_API_URL"),
	}

	if cfg.TelegramToken == "" {
		return nil, errors.New("environment variable TELEGRAM_TOKEN is required")
	}
	if cfg.DolarApiUrl == "" {
		return nil, errors.New("environment variable DOLAR_API_URL is required")
	}

	return cfg, nil
}
