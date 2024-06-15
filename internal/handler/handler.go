package handlers

import (
	"io"
	"net/http"

	"github.com/YerzhanAkhmetov/go-shortener/internal/usecase"
	"github.com/gorilla/mux"
)

type Handler struct {
	usecase usecase.URLUsecase
}

func NewHandler(usecase usecase.URLUsecase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Тело запроса не валидно", http.StatusBadRequest)
		return
	}
	originalURL := string(body)

	url, err := h.usecase.Create(originalURL)
	if err != nil {
		http.Error(w, "Ошибка генерации URL ID", http.StatusInternalServerError)
		return
	}
	shortURL := "http://localhost:8080/" + url.ID
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	url, exists := h.usecase.GetByID(id)
	if !exists {
		http.Error(w, "URL не найден", http.StatusNotFound)
		return
	}
	w.Header().Set("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
