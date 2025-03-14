package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string `env:"DB_HOST"`
	DBPort       string `env:"DB_PORT"`
	DBUser       string `env:"DB_USER"`
	DBPassword   string `env:"DB_PASSWORD"`
	DBName       string `env:"DB_NAME"`
	RedisAddr    string `env:"REDIS_ADDR"`
	DaDataAPIKey string `env:"DADATA_API_KEY"`
	DaDataURL    string `env:"DADATA_URL"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load("./.env"); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg := &Config{
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBUser:       os.Getenv("DB_USER"),
		DBPassword:   os.Getenv("DB_PASSWORD"),
		DBName:       os.Getenv("DB_NAME"),
		RedisAddr:    os.Getenv("REDIS_ADDR"),
		DaDataAPIKey: os.Getenv("DADATA_API_KEY"),
		DaDataURL:    os.Getenv("DADATA_URL"),
	}

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}
	if cfg.RedisAddr == "" {
		return nil, fmt.Errorf("missing Redis address")
	}
	if cfg.DaDataAPIKey == "" || cfg.DaDataURL == "" {
		return nil, fmt.Errorf("missing DaData API configuration")
	}

	return cfg, nil
}
