package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client // Глобальная переменная для хранения клиента MongoDB

// InitMongoDB устанавливает соединение с MongoDB
func InitMongoDB(uri string) error {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	MongoClient = client
	log.Println("Successfully connected to MongoDB!")
	return nil
}

// GetCollection возвращает коллекцию из MongoDB
func GetCollection(databaseName, collectionName string) *mongo.Collection {
	if MongoClient == nil {
		log.Fatalf("MongoClient is not initialized. Did you call InitMongoDB?")
	}
	return MongoClient.Database(databaseName).Collection(collectionName)
}

// DisconnectMongoDB закрывает соединение с MongoDB
func DisconnectMongoDB() {
	if MongoClient == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := MongoClient.Disconnect(ctx); err != nil {
		log.Printf("Error while disconnecting from MongoDB: %v", err)
	} else {
		log.Println("Disconnected from MongoDB.")
	}
}
