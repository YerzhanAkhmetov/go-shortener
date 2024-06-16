// internal/app/app.go

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

// App contains application components
type App struct {
	Config  *config.Config
	Handler *handler.Handler
	Router  *mux.Router
	Server  *server.Server
}

// NewApp initializes a new App instance
func NewApp(cfg *config.Config) *App {
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	handler := handler.NewHandler(urlUsecase)

	router := mux.NewRouter()
	server := server.NewServer(handler)

	return &App{
		Config:  cfg,
		Handler: handler,
		Router:  router,
		Server:  server,
	}
}

// Run starts the application server
func (app *App) Run() {
	port := fmt.Sprintf(":%s", app.Config.HttpPort)
	fmt.Println("Starting server on " + port)
	log.Fatal(http.ListenAndServe(port, app.Server.Router))
}
