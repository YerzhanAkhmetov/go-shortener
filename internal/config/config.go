package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
)

// Config структура для хранения конфигурационных параметров
type Config struct {
	Debug         bool   `env:"DEBUG" envDefault:"false"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8000"`
}

// LoadConfig загружает конфигурацию из переменных окружения и аргументов командной строки
func LoadConfig() (*Config, error) {
	cfg := Config{}

	// Загрузка параметров из переменных окружения
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}

	// Переопределение параметров с помощью аргументов командной строки
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "адрес запуска HTTP-сервера (например, localhost:8080)")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Базовый URL для сокращенных ссылок (например, http://localhost:8000)")

	flag.Parse()

	return &cfg, nil
}
