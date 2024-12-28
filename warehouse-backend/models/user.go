package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User представляет пользователя в MongoDB
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`  // MongoDB ObjectID
	Name     string             `bson:"name" json:"name"`         // Имя пользователя
	Email    string             `bson:"email" json:"email"`       // Электронная почта
	Password string             `bson:"password" json:"password"` // Хэш пароля
}
