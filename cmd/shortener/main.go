// cmd/server/main.go

package main

import (
	"log"

	"github.com/YerzhanAkhmetov/go-shortener/internal/app"
	"github.com/YerzhanAkhmetov/go-shortener/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	application := app.NewApp(cfg)
	application.Run()
}
