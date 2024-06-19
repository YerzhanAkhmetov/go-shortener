package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/YerzhanAkhmetov/go-shortener/internal/config"
	shortHandler "github.com/YerzhanAkhmetov/go-shortener/internal/handler"
	"github.com/YerzhanAkhmetov/go-shortener/internal/repository"
	"github.com/YerzhanAkhmetov/go-shortener/internal/server"
	shortServer "github.com/YerzhanAkhmetov/go-shortener/internal/server"
	"github.com/YerzhanAkhmetov/go-shortener/internal/storage"
	"github.com/YerzhanAkhmetov/go-shortener/internal/usecase"
	"github.com/gorilla/mux"
)

// App содержит компоненты приложения
type App struct {
	Config  *config.Config
	Handler *shortHandler.Handler
	Router  *mux.Router
	Server  *shortServer.Server
}

// NewApp инициализирует новый экземпляр приложения
func NewApp(cfg *config.Config) *App {
	// Создание хранилища данных в памяти
	store := storage.NewMemoryStorage()

	// Создание репозитория для работы с URL
	repo := repository.NewURLRepository(store)

	// Создание usecase для работы с URL
	urlUsecase := usecase.NewURLUsecase(repo)

	// Создание обработчика запросов
	handler := shortHandler.NewHandler(urlUsecase, cfg.BaseURL)

	// Создание маршрутизатора
	router := mux.NewRouter()

	// Создание сервера для обработки HTTP запросов
	// Создание сервера для обработки HTTP запросов
	server := server.NewServer(handler)

	// Настройка маршрутов для обработчика
	router.HandleFunc("/", handler.CreateShortURL).Methods("POST")
	router.HandleFunc("/{id}", handler.Redirect).Methods("GET")

	return &App{
		Config:  cfg,
		Handler: handler,
		Router:  router,
		Server:  server,
	}
}

// Run запускает сервер приложения
func (app *App) Run() {
	// Получение адреса сервера из конфигурации
	addr := ":" + app.Config.HTTPPort
	// Проверка наличия порта в ServerAddress
	// addr := app.Config.ServerAddress
	// if !strings.Contains(addr, ":") {
	// 	addr = addr + ":" + app.Config.HTTPPort
	// }

	fmt.Println("Starting server on address " + addr)

	// Запуск сервера на указанном адресе с маршрутизатором приложения
	// log.Fatal(http.ListenAndServe(addr, app.Router))
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
