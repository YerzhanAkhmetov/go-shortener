package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
)

// Config структура для хранения конфигурационных параметров
type Config struct {
	Debug         bool   `env:"DEBUG" envDefault:"false"`
	HTTPPort      string `env:"HTTP_PORT" envDefault:":8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8888"`
}

// LoadConfig загружает конфигурацию из переменных окружения и аргументов командной строки
func LoadConfig() (*Config, error) {
	cfg := Config{}

	// Загрузка параметров из переменных окружения
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}

	// Override with command line arguments
	flag.StringVar(&cfg.HTTPPort, "p", cfg.HTTPPort, "HTTP port (e.g., :8080)")
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "HTTP server address (e.g., localhost:8000)")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for shortened links (e.g., http://localhost:8888)")

	flag.Parse()

	return &cfg, nil
}
