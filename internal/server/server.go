package server

import (
	handlers "github.com/YerzhanAkhmetov/go-shortener/internal/handler"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
}

func NewServer(h *handlers.Handler) *Server {
	router := mux.NewRouter()
	router.HandleFunc("/", h.CreateShortURL).Methods("POST")
	router.HandleFunc("/{id}", h.Redirect).Methods("GET")
	return &Server{Router: router}
}
