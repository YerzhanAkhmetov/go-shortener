package server

import (
	handlers "github.com/YerzhanAkhmetov/go-shortener/internal/handler"
	"github.com/gorilla/mux"
)

// Server представляет HTTP сервер с маршрутизатором mux
type Server struct {
	Router        *mux.Router
	ServerAddress string
	BaseURL       string
}

// NewServer создает новый экземпляр сервера с заданным обработчиком, адресом сервера и базовым URL
func NewServer(h *handlers.Handler, serverAddress, baseURL string) *Server {
	router := mux.NewRouter()
	router.HandleFunc("/", h.CreateShortURL).Methods("POST")
	router.HandleFunc("/{id}", h.Redirect).Methods("GET")

	return &Server{
		Router:        router,
		ServerAddress: serverAddress,
		BaseURL:       baseURL,
	}
}

// // Run запускает HTTP сервер на указанном адресе
// func (s *Server) Run() {
// 	// Логирование запуска сервера
// 	fmt.Printf("Starting server on %s\n", s.ServerAddress)

// 	// Запуск HTTP сервера
// 	err := http.ListenAndServe(s.ServerAddress, s.Router)
// 	if err != nil {
// 		fmt.Printf("Failed to start server: %v\n", err)
// 	}
// }
