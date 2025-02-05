package controllers

import (
	"context"
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
	if email == "231441@astanait.edu.kz" {
		isAdmin = true
	}

	c.JSON(http.StatusOK, gin.H{
		"name":    user.Name,
		"email":   user.Email,
		"isAdmin": isAdmin,
	})
}

// UpdateProfile обновляет данные текущего пользователя
func UpdateProfile(c *gin.Context) {
	type UpdateProfileRequest struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userID := c.GetString("user_id")

	// Преобразуем строковый ID в ObjectID
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

	// Проверяем текущий пароль
	if req.CurrentPassword != "" {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid current password"})
			return
		}
	}

	// Обновляем данные пользователя
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}
		update["password"] = string(hashedPassword)
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
