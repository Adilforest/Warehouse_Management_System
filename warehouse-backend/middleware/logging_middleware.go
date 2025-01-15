package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware логирует информацию о входящих запросах.
func LoggingMiddleware(next func(c *gin.Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Засекаем время начала обработки запроса
		start := time.Now()

		// Логируем информацию о запросе
		log.Printf(
			"Начало обработки запроса: Метод=%s, Путь=%s, Время=%s",
			r.Method,
			r.URL.Path,
			start.Format("2006-01-02 15:04:05"),
		)

		// Логируем время выполнения запроса
		duration := time.Since(start)
		log.Printf(
			"Завершение обработки запроса: Метод=%s, Путь=%s, Длительность=%s",
			r.Method,
			r.URL.Path,
			duration,
		)
	}
}
