package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/YerzhanAkhmetov/go-shortener/internal/errs"
	"github.com/YerzhanAkhmetov/go-shortener/internal/usecase"
	"github.com/gin-gonic/gin"
)

// Handler обрабатывает HTTP запросы
type Handler struct {
	usecase usecase.URLUsecase
	BaseURL string
}

// NewHandler создает новый экземпляр Handler
func NewHandler(usc usecase.URLUsecase, baseURL string) *Handler {
	return &Handler{
		usecase: usc,
		BaseURL: baseURL,
	}
}

type Req struct {
	URL string `json:"url"`
}

// CreateShortURL обрабатывает запрос на создание короткой ссылки
func (h *Handler) CreateShortURL(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil || len(body) == 0 {
		c.JSON(http.StatusBadRequest, errs.NewError("Invalid request body", http.StatusBadRequest, "Bad Request"))
		return
	}

	originalURL := string(body)
	url, err := h.usecase.Create(originalURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError("Error generating URL ID", http.StatusInternalServerError, "Internal Server Error"))
		return
	}

	shortURL := h.BaseURL + "/" + url.ID
	c.String(http.StatusCreated, shortURL)
}

// Redirect обрабатывает запрос на перенаправление по короткой ссылке
func (h *Handler) Redirect(c *gin.Context) {
	id := c.Param("id")
	url, exists := h.usecase.GetByID(id)
	if !exists {
		c.JSON(http.StatusNotFound, errs.NewError("URL not found", http.StatusNotFound, "Not Found"))
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url.OriginalURL)
}

// ShortenJSON обрабатывает JSON-запрос на сокращение ссылки
func (h *Handler) ShortenJSON(c *gin.Context) {
	// Чтение тела запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		c.JSON(http.StatusBadRequest, errs.NewError("Invalid JSON body", http.StatusBadRequest, "Bad Request"))
		return
	}

	// Парсинг JSON
	var req Req
	err = json.Unmarshal(body, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Unmarshall body",
			"error":   err.Error(),
		})
		return
	}

	// Проверяем, передан ли URL
	if req.URL == "" {
		c.JSON(http.StatusBadRequest, errs.NewError("Invalid JSON body", http.StatusBadRequest, "Bad Request"))
		return
	}

	// Генерация короткого URL
	url, err := h.usecase.Create(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errs.NewError("Error generating URL ID", http.StatusInternalServerError, "Internal Server Error"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"result": h.BaseURL + "/" + url.ID})
}
