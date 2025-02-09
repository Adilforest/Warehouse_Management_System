package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// database/database.go
var MongoClient *mongo.Client

// GetCollection is a variable that can be overridden in tests
var GetCollection = func(databaseName, collectionName string) *mongo.Collection {
	if MongoClient == nil {
		log.Fatalf("MongoClient is not initialized. Did you call InitMongoDB?")
	}
	if databaseName == "" {
		log.Fatal("Database name is empty. Check your configuration.")
	}
	return MongoClient.Database(databaseName).Collection(collectionName)
}

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
