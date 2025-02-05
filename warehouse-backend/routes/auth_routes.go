package routes

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"net/http"
	"time"
	"warehouse-backend/database"
	"warehouse-backend/logger"
	"warehouse-backend/models"

	"github.com/gin-gonic/gin"
	"warehouse-backend/controllers"
)

// SetupAuthRoutes настраивает маршруты для аутентификации
func SetupAuthRoutes(router *gin.Engine) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", SignupHandler)
		authRoutes.GET("/verify", VerifyHandler) // Новый маршрут для верификации
	}
}

// SignupHandler обрабатывает запрос на регистрацию
func SignupHandler(c *gin.Context) {
	type SignupRequest struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Проверяем совпадение паролей
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Вызываем функцию для создания пользователя
	user, err := controllers.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully. Please check your email for verification.",
		"user":    user,
	})
}

// VerifyHandler обрабатывает запрос на верификацию email
func VerifyHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		logger.LogError("verify_handler", "Verification token is missing", map[string]interface{}{}, nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	logger.LogInfo("verify_handler", "Verifying token", map[string]interface{}{
		"token": token,
	})

	// Находим пользователя по токену
	var user models.User
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"verification_token": token}).Decode(&user)
	if err != nil {
		logger.LogError("verify_handler", "Invalid or expired verification token", map[string]interface{}{
			"token": token,
		}, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid or expired verification token"})
		return
	}

	// Обновляем статус верификации
	update := bson.M{"$set": bson.M{"verified": true, "verification_token": ""}}
	_, err = collection.UpdateOne(ctx, bson.M{"verification_token": token}, update)
	if err != nil {
		logger.LogError("verify_handler", "Failed to update user verification status", map[string]interface{}{
			"token": token,
		}, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
		return
	}

	logger.LogInfo("verify_handler", "User verified successfully", map[string]interface{}{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
