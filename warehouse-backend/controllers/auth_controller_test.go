package controllers

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"warehouse-backend/database"
	"warehouse-backend/models"
)

// Mock MongoDB collection for testing
type MockCollection struct {
	InsertOneFunc func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOneFunc   func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.InsertOneFunc(ctx, document, opts...)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.FindOneFunc(ctx, filter, opts...)
}

// Test CreateUser function
func TestCreateUser(t *testing.T) {
	mockCollection := &MockCollection{
		InsertOneFunc: func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
			return &mongo.InsertOneResult{InsertedID: "mockID"}, nil
		},
		FindOneFunc: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
			return &mongo.SingleResult{Err: mongo.ErrNoDocuments} // Simulate no existing user
		},
	}

	// Override GetCollection to return the mock collection
	database.GetCollection = func(databaseName, collectionName string) interface{} {
		return mockCollection
	}

	// Call CreateUser
	user, err := CreateUser("John Doe", "john@example.com", "password123")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// Assertions
	if user.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", user.Name)
	}
	if user.Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", user.Email)
	}
	if _, err := primitive.ObjectIDFromHex(user.ID.Hex()); err != nil {
		t.Errorf("Invalid ObjectID generated: %v", err)
	}
}

// Test LoginUser function
func TestLoginUser(t *testing.T) {
	mockCollection := &MockCollection{
		FindOneFunc: func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
			user := models.User{
				ID:       primitive.NewObjectID(),
				Email:    "john@example.com",
				Password: hashPassword("password123"),
				Verified: true,
				Role:     "user",
			}
			return &mongo.SingleResult{Err: nil}
		},
	}

	// Override GetCollection to return the mock collection
	database.GetCollection = func(databaseName, collectionName string) interface{} {
		return mockCollection
	}

	// Call LoginUser
	token, err := LoginUser("john@example.com", "password123")
	if err != nil {
		t.Fatalf("LoginUser failed: %v", err)
	}

	// Assertions
	if token == "" {
		t.Errorf("Expected non-empty token, got empty")
	}
}

// Helper function to hash password
func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
