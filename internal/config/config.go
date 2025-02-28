package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

// Config структура для хранения конфигурационных параметров
type Config struct {
	Debug         bool   `env:"DEBUG" envDefault:"false"`
	HTTPPort      string `env:"HTTP_PORT" envDefault:"8080"`
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"info"`
}

// LoadConfig загружает конфигурацию из переменных окружения и аргументов командной строки
func LoadConfig() (*Config, error) {
	cfg := Config{}

	// Проверяем наличие файла .env и загружаем переменные из него
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	// Загрузка параметров из переменных окружения
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}

	// Override with command line arguments
	flag.StringVar(&cfg.HTTPPort, "p", cfg.HTTPPort, "HTTP port (e.g., 8080)")                                   //./myapp -p 9090
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "HTTP server address (e.g., localhost)")          //./myapp -a 127.0.0.1
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for shortened links (e.g., http://localhost:8888)") //./myapp -b http://mydomain.com
	flag.StringVar(&cfg.LogLevel, "l", cfg.LogLevel, "Log level (debug, info, warn, error, fatal)")              //./myapp -l debug
	flag.Parse()

	// Override with environment variable SERVER_PORT if it exists
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.HTTPPort = port
	}

	return &cfg, nil
}
