package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
	"time"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // Убедитесь, что это значение берётся из .env

func GenerateToken(userID string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Токен действителен 24 часа
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken проверяет и валидирует JWT-токен
func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Убедитесь, что метод подписи совпадает
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// AuthMiddleware проверяет JWT-токен и устанавливает данные пользователя в контекст
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		userID := claims["user_id"].(string)
		email := claims["email"].(string)

		// Сохраняем данные пользователя в контексте
		c.Set("user_id", userID)
		c.Set("email", email)

		c.Next()
	}
}
