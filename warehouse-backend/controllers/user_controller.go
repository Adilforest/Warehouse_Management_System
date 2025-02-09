package controllers

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"warehouse-backend/database"
	"warehouse-backend/logger"
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
		logger.LogDBError("get_all_users", "Failed to fetch users", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		logger.LogDBError("get_all_users", "Failed to decode users", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
		return
	}

	logger.LogInfo("get_all_users", "All users fetched successfully", map[string]interface{}{
		"user_count": len(users),
	})
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

// DeleteAllUsers удаляет всех пользователей
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

// UpdateUserRole обновляет роль пользователя
// UpdateUserRole обновляет роль пользователя
func UpdateUserRole(c *gin.Context) {
	// Получаем email текущего пользователя из контекста
	currentUserEmail := c.GetString("email")

	// Получаем роль текущего пользователя из базы данных
	var currentUser models.User
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"email": currentUserEmail}).Decode(&currentUser)
	if err != nil {
		logger.LogError("update_user_role", "Failed to fetch current user", map[string]interface{}{
			"user_email": currentUserEmail,
		}, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to fetch current user"})
		return
	}

	// Проверяем, является ли текущий пользователь администратором
	if currentUser.Role != "admin" {
		logger.LogError("update_user_role", "Unauthorized access attempt", map[string]interface{}{
			"user_email": currentUserEmail,
		}, nil)
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
		return
	}

	// Получаем ID пользователя из параметров запроса
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogError("update_user_role", "Invalid user ID", map[string]interface{}{
			"user_id": id,
		}, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Получаем данные о новой роли из тела запроса
	type RoleUpdateRequest struct {
		Role string `json:"role"`
	}
	var req RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError("update_user_role", "Invalid request payload", map[string]interface{}{
			"error": err.Error(),
		}, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Проверяем, что роль является допустимой
	if req.Role != "admin" && req.Role != "user" {
		logger.LogError("update_user_role", "Invalid role specified", map[string]interface{}{
			"role": req.Role,
		}, nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role. Allowed roles: 'admin', 'user'"})
		return
	}

	// Находим пользователя в базе данных
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.LogError("update_user_role", "User not found", map[string]interface{}{
				"user_id": id,
			}, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		logger.LogDBError("update_user_role", "Failed to fetch user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Обновляем роль пользователя
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"role": req.Role}})
	if err != nil {
		logger.LogDBError("update_user_role", "Failed to update user role", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	logger.LogInfo("update_user_role", "User role updated successfully", map[string]interface{}{
		"user_id":  id,
		"new_role": req.Role,
	})
	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully", "new_role": req.Role})
}
