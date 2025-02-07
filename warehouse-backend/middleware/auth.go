package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
	"time"
)

// jwtSecret хранит секретный ключ для подписи JWT-токенов
var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // Убедитесь, что это значение берётся из .env

// GenerateToken создает новый JWT-токен для пользователя
func GenerateToken(userID string, email string) (string, error) {
	// Создаем токен с указанными claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,                                // ID пользователя
		"email":   email,                                 // Email пользователя
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	})

	// Подписываем токен с использованием секретного ключа
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken проверяет и валидирует JWT-токен
func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	// Обрабатываем ошибки парсинга
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Проверяем валидность токена
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

// AuthMiddleware проверяет JWT-токен и устанавливает данные пользователя в контекст
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Убираем префикс "Bearer " из токена
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Проверяем токен
		token, err := VerifyToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Извлекаем данные пользователя из токена
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Получаем user_id и email из claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid user_id in token"})
			c.Abort()
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid email in token"})
			c.Abort()
			return
		}

		// Сохраняем данные пользователя в контексте
		c.Set("user_id", userID)
		c.Set("email", email)

		// Продолжаем выполнение запроса
		c.Next()
	}
}
