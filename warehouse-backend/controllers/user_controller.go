package controllers

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"warehouse-backend/logger"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"warehouse-backend/database"
	"warehouse-backend/models"
)

// GetAllUsers возвращает список всех пользователей
func GetAllUsers(c *gin.Context) {
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUserByID возвращает пользователя по ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogError("get_user_by_id", "Invalid user ID", map[string]interface{}{
			"user_id": id,
		}, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.LogError("get_user_by_id", "User not found", map[string]interface{}{
				"user_id": id,
			}, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logger.LogDBError("get_user_by_id", "Failed to fetch user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	logger.LogInfo("get_user_by_id", "User fetched successfully", map[string]interface{}{
		"user_id": id,
	})
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser обновляет данные пользователя
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogError("update_user", "Invalid user ID", map[string]interface{}{
			"user_id": id,
		}, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var updateData bson.M
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logger.LogError("update_user", "Invalid request payload", map[string]interface{}{
			"error": err.Error(),
		}, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Если обновляется пароль, хэшируем его
	if newPassword, ok := updateData["password"].(string); ok {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			logger.LogDBError("update_user", "Failed to hash new password", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}
		updateData["password"] = string(hashedPassword)
	}

	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updateData})
	if err != nil {
		logger.LogDBError("update_user", "Failed to update user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	if result.MatchedCount == 0 {
		logger.LogError("update_user", "User not found", map[string]interface{}{
			"user_id": id,
		}, nil)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	logger.LogInfo("update_user", "User updated successfully", map[string]interface{}{
		"user_id": id,
	})
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser удаляет пользователя по ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogError("delete_user", "Invalid user ID", map[string]interface{}{
			"user_id": id,
		}, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		logger.LogDBError("delete_user", "Failed to delete user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.DeletedCount == 0 {
		logger.LogError("delete_user", "User not found", map[string]interface{}{
			"user_id": id,
		}, nil)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	logger.LogInfo("delete_user", "User deleted successfully", map[string]interface{}{
		"user_id": id,
	})
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func DeleteAllUsers(c *gin.Context) {
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		logger.LogDBError("delete_all_users", "Failed to delete all users", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete all users"})
		return
	}

	logger.LogInfo("delete_all_users", "All users deleted successfully", map[string]interface{}{
		"deleted_count": result.DeletedCount,
	})
	c.JSON(http.StatusOK, gin.H{"message": "All users deleted successfully"})
}
