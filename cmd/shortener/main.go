package main

import (
	"log"

	"github.com/YerzhanAkhmetov/go-shortener/internal/app"
	"github.com/YerzhanAkhmetov/go-shortener/internal/config"
)

func main() {
	// Загрузка конфигурации из переменных окружения и аргументов командной строки
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Инициализация нового экземпляра приложения
	application := app.NewApp(cfg)

	// Запуск приложения
	application.Run()
}
