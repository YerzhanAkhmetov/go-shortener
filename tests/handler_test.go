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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateShortURLHandler проверяет обработчик CreateShortURL
func TestCreateShortURLHandler(t *testing.T) {
	// Создание временного хранилища, репозитория, usecase и обработчика
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	h := handler.NewHandler(urlUsecase, cfg)

	// Создание маршрутизатора и установка обработчика для POST запросов
	r := mux.NewRouter()
	r.HandleFunc("/", h.CreateShortURL).Methods("POST")

	// Определение ожидаемых результатов для различных тестовых случаев
	type want struct {
		contentType string      // Ожидаемый тип содержимого ответа
		statusCode  int         // Ожидаемый HTTP статус код
		body        interface{} // Ожидаемое тело ответа
	}
	tests := []struct {
		name string // Название тестового случая
		body string // Тело HTTP запроса
		want want   // Ожидаемые результаты
	}{
		{
			name: "valid URL",                    // Название теста: валидный URL
			body: "https://practicum.yandex.ru/", // Тело запроса: валидный URL
			want: want{
				contentType: "text/plain",       // Ожидаемый тип содержимого: текстовый plain
				statusCode:  http.StatusCreated, // Ожидаемый HTTP статус код: 201 Created
				body:        nil,                // Ожидаемое тело: пустое (т.к. тут ожидается короткий URL)
			},
		},
		{
			name: "empty body", // Название теста: пустое тело запроса
			body: "",           // Тело запроса: пустое
			want: want{
				contentType: "application/json",                                                          // Ожидаемый тип содержимого: JSON
				statusCode:  http.StatusBadRequest,                                                       // Ожидаемый HTTP статус код: 400 Bad Request
				body:        errs.NewError("Invalid request body", http.StatusBadRequest, "Bad Request"), // Ожидаемое тело: ошибка "Invalid request body"
			},
		},
	}

	// Итерация по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание HTTP POST запроса с указанным телом
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			// Получение результата HTTP запроса
			result := w.Result()
			defer result.Body.Close()

			// Проверка ожидаемого HTTP статус кода
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			// Проверка ожидаемого типа содержимого ответа
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			// Если ожидается статус код 201 (Created), проверяем содержимое тела ответа
			if tt.want.statusCode == http.StatusCreated {
				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				shortURL := string(body)
				// Проверка, что возвращенный короткий URL начинается с заданного базового URL
				assert.True(t, strings.HasPrefix(shortURL, "http://localhost:8080/"))
			} else if tt.want.body != nil {
				// Если ожидается тело с ошибкой, проверяем JSON формат тела ответа
				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				expectedBody, err := json.Marshal(tt.want.body)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedBody), string(body))
			}
		})
	}
}

// TestRedirectHandler проверяет обработчик Redirect
func TestRedirectHandler(t *testing.T) {
	// Создание временного хранилища, репозитория, usecase и обработчика
	store := storage.NewMemoryStorage()
	repo := repository.NewURLRepository(store)
	urlUsecase := usecase.NewURLUsecase(repo)
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	h := handler.NewHandler(urlUsecase, cfg)

	// Сохранение тестовой ссылки в хранилище
	store.SaveURL("test1", "https://practicum.yandex.ru/")

	// Создание маршрутизатора и установка обработчика для GET запросов
	r := mux.NewRouter()
	r.HandleFunc("/{id}", h.Redirect).Methods("GET")

	// Определение ожидаемых результатов для различных тестовых случаев
	type want struct {
		statusCode int         // Ожидаемый HTTP статус код
		location   string      // Ожидаемое значение заголовка Location для перенаправления
		body       interface{} // Ожидаемое тело ответа
	}
	tests := []struct {
		name    string // Название тестового случая
		request string // Путь для HTTP GET запроса
		want    want   // Ожидаемые результаты
	}{
		{
			name:    "URL found", // Название теста: найден URL
			request: "/test1",    // Путь запроса: /test1
			want: want{
				statusCode: http.StatusTemporaryRedirect,   // Ожидаемый HTTP статус код: 307 Temporary Redirect
				location:   "https://practicum.yandex.ru/", // Ожидаемое значение заголовка Location
				body:       nil,                            // Ожидаемое тело ответа: пустое (т.к. ожидается перенаправление)
			},
		},
		{
			name:    "URL not found", // Название теста: URL не найден
			request: "/nonexistent",  // Путь запроса: /nonexistent
			want: want{
				statusCode: http.StatusNotFound,                                              // Ожидаемый HTTP статус код: 404 Not Found
				body:       errs.NewError("URL not found", http.StatusNotFound, "Not Found"), // Ожидаемое тело ответа: ошибка "URL not found"
			},
		},
	}

	// Итерация по всем тестовым случаям
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создание HTTP GET запроса с указанным путем
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			// Получение результата HTTP запроса
			result := w.Result()
			defer result.Body.Close()

			// Проверка ожидаемого HTTP статус кода
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if tt.want.location != "" {
				// Если ожидается перенаправление, проверяем заголовок Location
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			} else if tt.want.body != nil {
				// Если ожидается тело с ошибкой, проверяем JSON формат тела ответа
				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				expectedBody, err := json.Marshal(tt.want.body)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedBody), string(body))
			}
		})
	}
}
