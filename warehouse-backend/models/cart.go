package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

// CartItem представляет элемент в корзине
type CartItem struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"` // Уникальный ID элемента корзины
    Product  Product            `bson:"product" json:"product"`   // Продукт
    Quantity int                `bson:"quantity" json:"quantity"` // Количество
}

// Cart представляет корзину пользователя
type Cart struct {
    ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"` // Уникальный ID корзины
    UserID  primitive.ObjectID `bson:"user_id" json:"user_id"`   // ID пользователя
    Items   []CartItem         `bson:"items" json:"items"`       // Список товаров в корзине
    Total   float64            `bson:"total" json:"total"`       // Общая стоимость товаров
    Created time.Time          `bson:"created_at" json:"created_at"`
    Updated time.Time          `bson:"updated_at" json:"updated_at"`
}