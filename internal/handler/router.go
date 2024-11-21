package handlers

import (
	"context"
	"github.com/YerzhanAkhmetov/go-shortener/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Rest отвечает за запуск HTTP сервера
type Rest struct {
	server *http.Server
	logger zap.Logger
}

// NewRest создает новый экземпляр Rest API
func NewRest() *Rest {
	return &Rest{}
}

// Start запускает сервер
func (r *Rest) Start(lAddr string, h *Handler) {
	logger.Log.Info("Running server", zap.String("address", lAddr))

	// Устанавливаем режим Gin
	gin.SetMode(gin.ReleaseMode)

	// Создаем маршрутизатор
	engine := gin.Default()
	engine.Use(logger.RequestLogger(logger.Log), gin.Recovery())

	// Устанавливаем маршруты
	SetRoutes(engine, h)

	// Настройка и запуск HTTP-сервера
	r.server = &http.Server{
		Addr:    lAddr,
		Handler: engine,
	}

	if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error("Failed to start server", zap.Error(err))
	}
}

// Stop завершает работу сервера
func (r *Rest) Stop(ctx context.Context) error {
	if err := r.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// SetRoutes регистрирует маршруты
func SetRoutes(r *gin.Engine, h *Handler) {
	r.POST("/", h.CreateShortURL)         // Создание короткой ссылки
	r.GET("/:id", h.Redirect)             // Перенаправление по короткой ссылке
	r.POST("/api/shorten", h.ShortenJSON) // Пример дополнительного маршрута
}
