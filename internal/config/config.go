package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

// Config структура для хранения конфигурационных параметров
type Config struct {
	Debug    bool   `env:"DEBUG" envDefault:"false"`
	HTTPPort string `env:"HTTP_PORT" envDefault:":8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

// LoadConfig загружает конфигурацию из переменных окружения и аргументов командной строки
func LoadConfig() (*Config, error) {
	cfg := Config{}

	// Загрузка параметров из переменных окружения
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}

	// Переопределение параметров с помощью аргументов командной строки
	flag.StringVar(&cfg.HTTPPort, "a", cfg.HTTPPort, "HTTP server address (e.g., :8888)")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for shortened links (e.g., http://localhost:8000)")

	flag.Parse()

	return &cfg, nil
}
