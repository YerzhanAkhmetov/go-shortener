package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тест для createShortURLHandler
func TestCreateShortURLHandler(t *testing.T) {
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
			h := http.HandlerFunc(createShortURLHandler)
			h(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			if tt.want.statusCode == http.StatusCreated {
				body, err := io.ReadAll(result.Body)
				require.NoError(t, err)
				require.NoError(t, result.Body.Close())

				shortURL := string(body)
				assert.True(t, strings.HasPrefix(shortURL, "http://localhost:8080/"))
			}
		})
	}
}

// Тест для redirectHandler
func TestRedirectHandler(t *testing.T) {
	// Предварительно добавляем URL в хранилище
	urlStore["test1"] = "http://example.com"

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
				location:   "http://example.com",
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

			r := mux.NewRouter()
			r.HandleFunc("/{id}", redirectHandler)
			r.ServeHTTP(w, request)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			if tt.want.location != "" {
				assert.Equal(t, tt.want.location, result.Header.Get("Location"))
			}
		})
	}
}
