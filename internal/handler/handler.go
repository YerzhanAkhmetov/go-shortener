// internal/handlers/handler.go

package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/YerzhanAkhmetov/go-shortener/internal/config"
	"github.com/YerzhanAkhmetov/go-shortener/internal/errs"
	"github.com/YerzhanAkhmetov/go-shortener/internal/usecase"
	"github.com/gorilla/mux"
)

type Handler struct {
	usecase usecase.URLUsecase
	config  *config.Config
}

func NewHandler(usecase usecase.URLUsecase, cfg *config.Config) *Handler {
	return &Handler{
		usecase: usecase,
		config:  cfg,
	}
}

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		writeError(w, errs.NewError("Invalid request body", http.StatusBadRequest, "Bad Request"))
		return
	}
	originalURL := string(body)

	url, err := h.usecase.Create(originalURL)
	if err != nil {
		writeError(w, errs.NewError("Error generating URL ID", http.StatusInternalServerError, "Internal Server Error"))
		return
	}
	shortURL := h.config.BaseURL + "/" + url.ID
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	url, exists := h.usecase.GetByID(id)
	if !exists {
		writeError(w, errs.NewError("URL not found", http.StatusNotFound, "Not Found"))
		return
	}
	w.Header().Set("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func writeError(w http.ResponseWriter, err *errs.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
