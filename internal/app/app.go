package app

import (
	"fmt"
	"github.com/YerzhanAkhmetov/go-shortener/internal/config"
	shortHandler "github.com/YerzhanAkhmetov/go-shortener/internal/handler"
	"github.com/YerzhanAkhmetov/go-shortener/internal/logger"
	"github.com/YerzhanAkhmetov/go-shortener/internal/repository"
	"github.com/YerzhanAkhmetov/go-shortener/internal/storage"
	"github.com/YerzhanAkhmetov/go-shortener/internal/usecase"
	"strings"
)

// App содержит компоненты приложения
type App struct {
	Config  *config.Config
	Handler *shortHandler.Handler
	Rest    *shortHandler.Rest
}

// NewApp инициализирует новый экземпляр приложения
func NewApp(cfg *config.Config) *App {
	//Инициализация логгера
	if err := logger.Initialize(cfg.LogLevel); err != nil {
		panic(err)
	}

	// Создание хранилища данных в памяти
	store := storage.NewMemoryStorage()
	// Создание репозитория для работы с URL
	repo := repository.NewURLRepository(store)
	// Создание usecase для работы с URL
	usc := usecase.NewURLUsecase(repo)
	// Создание обработчика запросов
	h := shortHandler.NewHandler(usc, cfg.BaseURL)

	return &App{
		Config:  cfg,
		Handler: h,
		Rest:    shortHandler.NewRest(),
	}
}

// Run запускает сервер приложения
func (app *App) Run() {
	addr := app.Config.ServerAddress
	if !strings.Contains(addr, ":") {
		addr += ":" + app.Config.HTTPPort
	}
	fmt.Println("Starting server on address " + addr)
	app.Rest.Start(addr, app.Handler)
}
