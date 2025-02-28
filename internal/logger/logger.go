package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var Log = zap.NewNop()

func Initialize(level string) error {
	//парсинг уровня логирования
	levelLog, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	//содание конфигурации логера
	cfg := zap.NewProductionConfig()
	//установка уроня логирования
	cfg.Level = levelLog
	//построение логгера
	newZapLvl, err := cfg.Build()
	if err != nil {
		return err
	}
	//инициализация глобального логгера
	Log = newZapLvl
	return nil
}

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		contentLength := int64(c.Writer.Size())

		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("route", c.Request.URL.Path),
			zap.Duration("duration time", duration),
		)

		logger.Info("Response",
			zap.Int("statusCode", statusCode),
			zap.Int64("contentLength", contentLength),
		)
	}
}
