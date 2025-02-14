package controllers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "transaction-service/database"
    "transaction-service/models"
    "transaction-service/utils"
)

// CreateTransaction создает новую транзакцию
func CreateTransaction(c *gin.Context) {
    var transaction models.Transaction
    if err := c.ShouldBindJSON(&transaction); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // Установка начальных значений
    transaction.ID = primitive.NewObjectID()
    transaction.Status = "в ожидании оплаты"
    transaction.CreatedAt = time.Now()
    transaction.UpdatedAt = time.Now()

    // Сохранение транзакции в базе данных
    if err := database.CreateTransaction(&transaction); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully", "transaction_id": transaction.ID.Hex()})
}

// UpdateTransactionStatus обновляет статус транзакции
func UpdateTransactionStatus(c *gin.Context) {
    transactionID := c.Param("id")
    status := c.Query("status")

    // Обновление статуса в базе данных
    err := database.UpdateTransactionStatus(transactionID, status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction status"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated successfully"})
}