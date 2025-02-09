package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
	"warehouse-backend/controllers"
	"warehouse-backend/database"
	"warehouse-backend/logger"
	"warehouse-backend/middleware"
	"warehouse-backend/models"
)

// SetupAuthRoutes настраивает маршруты для аутентификации
func SetupAuthRoutes(router *gin.Engine) {
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", SignupHandler)
		authRoutes.GET("/verify", VerifyHandler)
		authRoutes.POST("/login", LoginHandler)
	}
}

// Helper function to send standardized JSON responses
func sendResponse(c *gin.Context, status int, message string, data interface{}, err error) {
	if err != nil {
		logger.Log.Error("response_handler", message, map[string]interface{}{"error": err.Error()}, err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	logger.LogInfo("response_handler", message, map[string]interface{}{"data": data})
	c.JSON(status, gin.H{"message": message, "data": data})
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
		sendResponse(c, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}
	// Проверяем совпадение паролей
	if req.Password != req.ConfirmPassword {
		sendResponse(c, http.StatusBadRequest, "Passwords do not match", nil, nil)
		return
	}
	// Вызываем функцию для создания пользователя
	user, err := controllers.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, "Failed to create user", nil, err)
		return
	}
	// Генерируем JWT-токен с ролью пользователя
	token, err := middleware.GenerateToken(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, "Failed to generate token", nil, err)
		return
	}
	sendResponse(c, http.StatusCreated, "User registered and logged in successfully", gin.H{"token": token}, nil)
}

// VerifyHandler обрабатывает запрос на верификацию email
func VerifyHandler(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		sendResponse(c, http.StatusBadRequest, "Verification token is required", nil, nil)
		return
	}
	logger.LogInfo("verify_handler", "Verifying token", map[string]interface{}{"token": token})
	// Находим пользователя по токену
	var user models.User
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"verification_token": token}).Decode(&user)
	if err != nil {
		sendResponse(c, http.StatusNotFound, "Invalid or expired verification token", nil, err)
		return
	}
	// Обновляем статус верификации
	update := bson.M{"$set": bson.M{"verified": true, "verification_token": ""}}
	_, err = collection.UpdateOne(ctx, bson.M{"verification_token": token}, update)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, "Failed to verify user", nil, err)
		return
	}
	sendResponse(c, http.StatusOK, "Email verified successfully", nil, nil)
}

// LoginHandler обрабатывает запрос на вход
func LoginHandler(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendResponse(c, http.StatusBadRequest, "Invalid request payload", nil, err)
		return
	}
	// Вызываем функцию для входа
	token, err := controllers.LoginUser(req.Email, req.Password)
	if err != nil {
		sendResponse(c, http.StatusUnauthorized, "Failed to login user", nil, err)
		return
	}
	// Получаем роль пользователя из базы данных
	var user models.User
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, "Failed to fetch user role", nil, err)
		return
	}
	// Возвращаем успешный ответ с токеном и ролью
	sendResponse(c, http.StatusOK, "Login successful", gin.H{
		"token": token,
		"role":  user.Role,
	}, nil)
}
