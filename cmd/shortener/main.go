package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Хранилище для сокращенных URL
var urlStore = map[string]string{}

// Генерация случайного идентификатора
func generateID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Обработчик для создания сокращенного URL
func createShortURLHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Тело запроса не валидно", http.StatusBadRequest)
		return
	}
	originalURL := string(body)

	id, err := generateID()
	if err != nil {
		http.Error(w, "Ошибка генерации URL ID", http.StatusInternalServerError)
		return
	}
	//Помещаем в хранилище
	urlStore[id] = originalURL

	shortURL := fmt.Sprintf("http://localhost:8080/%s", id)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(shortURL))
}

// Обработчик для перенаправления по сокращенному URL
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	originalURL, exists := urlStore[id]
	if !exists {
		http.Error(w, "URL не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", createShortURLHandler).Methods("POST")
	r.HandleFunc("/{id}", redirectHandler).Methods("GET")

	port := fmt.Sprintf(":%d", 8080)
	fmt.Println("Starting server on " + port)
	log.Fatal(http.ListenAndServe(port, r))
}
