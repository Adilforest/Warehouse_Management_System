package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// jwtSecret хранит секретный ключ для подписи JWT-токенов
var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // Убедитесь, что это значение берётся из .env

// GenerateToken создает новый JWT-токен для пользователя
func GenerateToken(userID string, email string, role string) (string, error) {
	// Создаем токен с указанными claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,                                // ID пользователя
		"email":   email,                                 // Email пользователя
		"role":    role,                                  // Роль пользователя
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
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
            c.Abort()
            return
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")
        token, err := VerifyToken(tokenString)
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        userID, ok := claims["user_id"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
            c.Abort()
            return
        }

        // Преобразуем userID в primitive.ObjectID
        objectID, err := primitive.ObjectIDFromHex(userID)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id format"})
            c.Abort()
            return
        }

        // Сохраняем userID в контекст
        c.Set("userID", objectID)
        c.Next()
    }
}

func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
			c.Abort()
			return
		}
		c.Next()
	}
}
