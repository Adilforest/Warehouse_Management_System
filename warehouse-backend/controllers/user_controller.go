package controllers

import (
	"context"
	"errors"
	"time"
	"warehouse-backend/database"
	"warehouse-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(name, email, password string) (models.User, error) {
	user := models.User{
		ID:       primitive.NewObjectID(),
		Name:     name,
		Email:    email,
		Password: password,
	}

	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, errors.New("failed to create user: " + err.Error())
	}

	return user, nil
}

func GetUserByID(id string) (models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("invalid user ID")
	}

	var user models.User
	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func UpdateUserByID(id string, name, email, password string) (models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, errors.New("invalid user ID")
	}

	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"name":     name,
			"email":    email,
			"password": password,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return models.User{}, errors.New("failed to update user: " + err.Error())
	}

	return GetUserByID(id)
}

func DeleteUserByID(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid user ID")
	}

	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return errors.New("failed to delete user: " + err.Error())
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User

	collection := database.GetCollection("warehouse", "users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errors.New("failed to retrieve users: " + err.Error())
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		return nil, errors.New("failed to parse users: " + err.Error())
	}

	return users, nil
}
