package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

// Config структура для хранения конфигурационных параметров
type Config struct {
	Debug         bool   `env:"DEBUG" envDefault:"false"`
	HTTPPort      string `env:"HTTP_PORT" envDefault:":8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8888"`
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
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "адрес запуска HTTP-сервера (например, localhost:8888)")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Базовый URL для сокращенных ссылок (например, http://localhost:8000)")

	flag.Parse()

	return &cfg, nil
}
