package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/YerzhanAkhmetov/go-shortener/internal/config"
	"github.com/YerzhanAkhmetov/go-shortener/internal/errs"
	handler "github.com/YerzhanAkhmetov/go-shortener/internal/handler"
	"github.com/YerzhanAkhmetov/go-shortener/internal/repository"
	"github.com/YerzhanAkhmetov/go-shortener/internal/storage"
	"github.com/YerzhanAkhmetov/go-shortener/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURLHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	cfg := &config.Config{} // Используем пустую конфигурацию
	h := handler.NewHandler(urlUsecase, cfg)

	// Создаем временный HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(h.CreateShortURL))
	defer server.Close()

	type want struct {
		contentType string
		statusCode  int
		body        interface{}
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
				body:        nil,
			},
		},
		{
			name: "empty body",
			body: "",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusBadRequest,
				body:        errs.NewError("Invalid request body", http.StatusBadRequest, "Bad Request"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Post(server.URL, "text/plain", strings.NewReader(tt.body))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			if tt.want.body != nil {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				expectedBody, err := json.Marshal(tt.want.body)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedBody), string(body))
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	cfg := &config.Config{} // Используем пустую конфигурацию
	h := handler.NewHandler(urlUsecase, cfg)

	// Создаем временный HTTP-сервер
	server := httptest.NewServer(http.HandlerFunc(h.Redirect))
	defer server.Close()

	store.SaveURL("test1", "https://practicum.yandex.ru/")

	type want struct {
		statusCode int
		location   string
		body       interface{}
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
				body:       nil,
			},
		},
		{
			name:    "URL not found",
			request: "/nonexistent",
			want: want{
				statusCode: http.StatusNotFound,
				body:       errs.NewError("URL not found", http.StatusNotFound, "Not Found"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.request)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, resp.Header.Get("Location"))
			} else if tt.want.body != nil {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)
				expectedBody, err := json.Marshal(tt.want.body)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedBody), string(body))
			}
		})
	}
}
