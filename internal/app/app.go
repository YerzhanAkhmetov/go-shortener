package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/YerzhanAkhmetov/go-shortener/internal/config"
	handler "github.com/YerzhanAkhmetov/go-shortener/internal/handler"
	"github.com/YerzhanAkhmetov/go-shortener/internal/repository"
	"github.com/YerzhanAkhmetov/go-shortener/internal/server"
	"github.com/YerzhanAkhmetov/go-shortener/internal/storage"
	"github.com/YerzhanAkhmetov/go-shortener/internal/usecase"
	"github.com/gorilla/mux"
)

// App содержит компоненты приложения
type App struct {
	Config  *config.Config
	Handler *handler.Handler
	Router  *mux.Router
	Server  *server.Server
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
	handler := handler.NewHandler(urlUsecase, cfg)

	// Создание маршрутизатора
	router := mux.NewRouter()

	// Создание сервера для обработки HTTP запросов
	server := server.NewServer(handler, cfg.ServerAddress, cfg.BaseURL)

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
	addr := app.Config.ServerAddress
	fmt.Println("Starting server on " + addr)

	// Запуск сервера на указанном адресе с маршрутизатором приложения
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
