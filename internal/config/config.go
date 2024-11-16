package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Server struct {
		Host string `env:"SERVER_HOST" envDefault:"localhost"`
		Port string `env:"SERVER_PORT" envDefault:"8080"`
	}

	DB struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     string `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER" envDefault:"postgres"`
		Password string `env:"DB_PASSWORD" envDefault:"postgres"`
		Name     string `env:"DB_NAME" envDefault:"referral_db"`
		SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"` // disable, require, verify-full, verify-ca
	}

	JWT struct {
		Secret string `env:"JWT_SECRET" envDefault:"your-secret-key"`
		TTL    int    `env:"JWT_TTL" envDefault:"24"` // часы
	}
)

// Config определяет конфигурацию сервера и базы данных
type Config struct {
	Server
	DB
	JWT
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() (*Config, error) {
	// Загружаем конфигурацию из .env файла
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}
	return cfg, nil
}
