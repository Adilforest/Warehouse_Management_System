package database

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "transaction-service/models"
)

// CreateTransaction создает новую транзакцию в базе данных
func CreateTransaction(transaction *models.Transaction) error {
    collection := GetCollection("warehouse", "transactions")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    transaction.ID = primitive.NewObjectID()
    transaction.CreatedAt = time.Now()
    transaction.UpdatedAt = time.Now()

    _, err := collection.InsertOne(ctx, transaction)
    if err != nil {
        log.Printf("Failed to create transaction: %v", err)
        return err
    }
    return nil
}

// GetTransactionByID возвращает транзакцию по её ID
func GetTransactionByID(id primitive.ObjectID) (*models.Transaction, error) {
    collection := GetCollection("warehouse", "transactions")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var transaction models.Transaction
    err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&transaction)
    if err != nil {
        log.Printf("Failed to get transaction by ID: %v", err)
        return nil, err
    }
    return &transaction, nil
}

// UpdateTransactionStatus обновляет статус транзакции
func UpdateTransactionStatus(id primitive.ObjectID, status string) error {
    collection := GetCollection("warehouse", "transactions")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err := collection.UpdateOne(
        ctx,
        bson.M{"_id": id},
        bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
    )
    if err != nil {
        log.Printf("Failed to update transaction status: %v", err)
        return err
    }
    return nil
}