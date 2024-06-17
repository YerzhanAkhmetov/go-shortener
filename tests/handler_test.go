package tests

import (
	"encoding/json"
	"io"
	"net"
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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortURLHandler(t *testing.T) {
	// Создаем конфигурацию с использованием случайного порта и переопределенными адресом сервера и базовым URL
	// Получаем доступный порт
	port, err := getAvailablePort()
	require.NoError(t, err)

	// Обновляем конфигурацию с новым портом
	cfg := &config.Config{
		Debug:         false,
		HTTPPort:      ":" + port,
		ServerAddress: "localhost:" + port,
		BaseURL:       "http://localhost:" + port,
	}

	// Создаем хранилище, репозиторий, usecase и хендлер
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	h := handler.NewHandler(urlUsecase, cfg)

	// Создаем маршрутизатор и добавляем хендлер для тестирования
	r := mux.NewRouter()
	r.HandleFunc("/", h.CreateShortURL).Methods("POST")

	// Определяем ожидаемые результаты тестов
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

	// Проходим по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			// Проверяем ожидаемые результаты
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.body != nil {
				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				expectedBody, err := json.Marshal(tt.want.body)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedBody), string(body))
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	// Получаем доступный порт
	port, err := getAvailablePort()
	require.NoError(t, err)

	// Обновляем конфигурацию с новым портом
	cfg := &config.Config{
		Debug:         false,
		HTTPPort:      ":" + port,
		ServerAddress: "localhost:" + port,
		BaseURL:       "http://localhost:" + port,
	}

	// Создаем хранилище, репозиторий, usecase и хендлер
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	h := handler.NewHandler(urlUsecase, cfg)

	// Сохраняем тестовый URL в хранилище
	store.SaveURL("test1", "https://practicum.yandex.ru/")

	// Создаем маршрутизатор и добавляем хендлер для тестирования
	r := mux.NewRouter()
	r.HandleFunc("/{id}", h.Redirect).Methods("GET")

	// Определяем ожидаемые результаты тестов
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

	// Проходим по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()

			// Проверяем ожидаемые результаты
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			} else if tt.want.body != nil {
				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				expectedBody, err := json.Marshal(tt.want.body)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedBody), string(body))
			}
		})
	}
}

func getAvailablePort() (string, error) {
	// Создаем прослушиватель на случайном порту
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	defer listener.Close()

	// Получаем адрес прослушивателя и извлекаем порт
	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return "", err
	}

	return portStr, nil
}
