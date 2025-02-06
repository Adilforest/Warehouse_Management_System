package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"warehouse-backend/database"
	"warehouse-backend/models"
)

// GetProfile возвращает данные текущего пользователя
func GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	email := c.GetString("email")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Проверка роли администратора
	isAdmin := false
	if email == "231441@astanait.edu.kz" { // Замените на вашу логику проверки администратора
		isAdmin = true
	}

	// Возвращаем данные пользователя
	c.JSON(http.StatusOK, gin.H{
		"name":    user.Name,
		"email":   user.Email,
		"isAdmin": isAdmin, // Возвращаем флаг администратора
	})
}

// UpdateProfile обновляет данные текущего пользователя
func UpdateProfile(c *gin.Context) {
	// Структура для парсинга запроса
	type UpdateProfileRequest struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	// Парсим запрос
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request payload: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Получаем user_id из контекста
	userID := c.GetString("user_id")

	// Преобразуем строковый ID в ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получаем коллекцию пользователей
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ищем пользователя в базе данных
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		log.Printf("User not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Проверяем текущий пароль, если он предоставлен
	if req.CurrentPassword != "" {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
		if err != nil {
			log.Printf("Invalid current password: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid current password"})
			return
		}
	}

	// Подготавливаем данные для обновления
	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Email != "" {
		update["email"] = req.Email
	}
	if req.NewPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash new password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}
		update["password"] = string(hashedPassword)
	}

	// Если нет полей для обновления, возвращаем ошибку
	if len(update) == 0 {
		log.Printf("No fields to update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	// Обновляем данные пользователя в базе данных
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		log.Printf("Failed to update profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Логирование успешного обновления
	log.Printf("Profile updated successfully for user: %s", user.Email)

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
