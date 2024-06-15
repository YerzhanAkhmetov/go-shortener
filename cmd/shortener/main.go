package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "github.com/YerzhanAkhmetov/go-shortener/internal/handler"
	"github.com/YerzhanAkhmetov/go-shortener/internal/repository"
	"github.com/YerzhanAkhmetov/go-shortener/internal/server"
	"github.com/YerzhanAkhmetov/go-shortener/internal/storage"
	"github.com/YerzhanAkhmetov/go-shortener/internal/usecase"
)

func main() {
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	h := handlers.NewHandler(urlUsecase)

	s := server.NewServer(h)
	port := fmt.Sprintf(":%d", 8080)
	fmt.Println("Starting server on " + port)
	log.Fatal(http.ListenAndServe(port, s.Router))
}
