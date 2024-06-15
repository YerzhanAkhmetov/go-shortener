package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlers "github.com/YerzhanAkhmetov/go-shortener/internal/handler"
	"github.com/YerzhanAkhmetov/go-shortener/internal/repository"
	"github.com/YerzhanAkhmetov/go-shortener/internal/storage"
	"github.com/YerzhanAkhmetov/go-shortener/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURLHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	h := handlers.NewHandler(urlUsecase)

	r := mux.NewRouter()
	r.HandleFunc("/", h.CreateShortURL).Methods("POST")

	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "valid URL",
			body: "https://practicum.yandex.ru/",
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
			},
		},
		{
			name: "empty body",
			body: "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.statusCode == http.StatusCreated {
				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				shortURL := string(body)
				assert.True(t, strings.HasPrefix(shortURL, "http://localhost:8080/"))
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	h := handlers.NewHandler(urlUsecase)

	store.SaveURL("test1", "https://practicum.yandex.ru/")

	r := mux.NewRouter()
	r.HandleFunc("/{id}", h.Redirect).Methods("GET")

	type want struct {
		statusCode int
		location   string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "URL found",
			request: "/test1",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				location:   "https://practicum.yandex.ru/",
			},
		},
		{
			name:    "URL not found",
			request: "/nonexistent",
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			}
		})
	}
}
