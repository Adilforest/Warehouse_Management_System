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
var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // Убедитесь, что значение берётся из .env

// GenerateToken создаёт новый JWT-токен для пользователя.
// Принимается userID типа string или primitive.ObjectID.
func GenerateToken(userID interface{}, email string, role string) (string, error) {
	var uid string
	switch v := userID.(type) {
	case string:
		uid = v
	case primitive.ObjectID:
		uid = v.Hex()
	default:
		return "", fmt.Errorf("unsupported userID type")
	}

	// Создаём токен с заданными claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": uid,                                // идентификатор пользователя
		"email":   email,                              // email пользователя
		"role":    role,                               // роль пользователя
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // срок действия 24 часа
	})
	// Подписываем токен секретным ключом.
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

// VerifyToken проверяет и валидирует JWT-токен.
func VerifyToken(tokenString string) (*jwt.Token, error) {
	// Парсим токен с проверкой метода подписи.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token, nil
}

// AuthMiddleware проверяет JWT-токен, извлекает данные и устанавливает их в контекст Gin.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлекаем заголовок Authorization и проверяем префикс "Bearer ".
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Валидируем токен.
		token, err := VerifyToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Извлекаем claims.
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Извлекаем user_id.
		uid, ok := claims["user_id"].(string)
		if !ok || uid == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id in token"})
			c.Abort()
			return
		}
		// Если строка соответствует формату ObjectID, сохраняем и преобразованное значение.
		if isValidObjectID(uid) {
			objectID, err := primitive.ObjectIDFromHex(uid)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse user_id"})
				c.Abort()
				return
			}
			c.Set("userID", objectID)
		}
		// Сохраняем оригинальный user_id как строку.
		c.Set("user_id", uid)

		// Извлекаем email.
		email, ok := claims["email"].(string)
		if !ok || email == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email in token"})
			c.Abort()
			return
		}
		c.Set("email", email)

		// Извлекаем роль пользователя.
		role, ok := claims["role"].(string)
		if !ok || role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role in token"})
			c.Abort()
			return
		}
		c.Set("role", role)

		c.Next()
	}
}

// isValidObjectID проверяет, является ли строка корректным MongoDB ObjectID.
func isValidObjectID(id string) bool {
	if len(id) != 24 {
		return false
	}
	for _, char := range id {
		if !((char >= '0' && char <= '9') ||
			(char >= 'a' && char <= 'f') ||
			(char >= 'A' && char <= 'F')) {
			return false
		}
	}
	return true
}

// AdminOnlyMiddleware разрешает доступ только администраторам.
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
