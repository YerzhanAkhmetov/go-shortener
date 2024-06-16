// internal/config/config.go

package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Debug    bool   `env:"DEBUG" envDefault:"false"`
	HttpPort string `env:"HTTP_PORT" envDefault:"8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}

	// Load from environment variables
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse env vars: %w", err)
	}

	// Override with command line arguments
	flag.StringVar(&cfg.HttpPort, "a", cfg.HttpPort, "HTTP server address (e.g., localhost:8888)")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL for shortened links (e.g., http://localhost:8000)")

	flag.Parse()

	return &cfg, nil
}
