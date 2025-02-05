package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name              string             `bson:"name" json:"name"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bson:"password" json:"-"`
	Verified          bool               `bson:"verified" json:"verified"`
	VerificationToken string             `bson:"verification_token" json:"-"`
	Role              string             `bson:"role" json:"role"`
}
