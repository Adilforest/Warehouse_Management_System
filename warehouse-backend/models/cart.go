package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartItem struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID primitive.ObjectID `bson:"productId,omitempty" json:"productId"`
	Quantity  int                `bson:"quantity" json:"quantity"`
}

type Cart struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"userId,omitempty" json:"userId"`
	Items  []CartItem         `bson:"items" json:"items"`
}
