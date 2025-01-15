package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

// RateLimiterMiddleware создает middleware для ограничения частоты запросов
func RateLimiterMiddleware(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b) // r - количество запросов в секунду, b - размер "буфера"

	return func(c *gin.Context) {
		if !limiter.Allow() {
			// Если лимит превышен, возвращаем ошибку 429
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
			})
			c.Abort() // Прерываем выполнение следующих обработчиков
			return
		}
		c.Next() // Передаем управление следующему middleware или обработчику
	}
}
