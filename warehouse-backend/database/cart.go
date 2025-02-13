package database

import (
	"context"
	"time"
	"warehouse-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// getCartCollection возвращает коллекцию "carts" из базы "warehouse".
func getCartCollection() *mongo.Collection {
	return GetCollection("warehouse", "carts")
}

// GetCartByUserID возвращает корзину для указанного пользователя.
func GetCartByUserID(userID primitive.ObjectID) (*models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart models.Cart
	err := getCartCollection().FindOne(ctx, bson.M{"userId": userID}).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

// CreateCart создаёт новую корзину.
func CreateCart(cart *models.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := getCartCollection().InsertOne(ctx, cart)
	return err
}

// UpdateCart обновляет корзину (например, изменяя список Items).
func UpdateCart(cart *models.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": cart.ID}
	update := bson.M{"$set": bson.M{"items": cart.Items}}
	_, err := getCartCollection().UpdateOne(ctx, filter, update)
	return err
}

// DeleteCart удаляет корзину по её ID.
func DeleteCart(cartID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := getCartCollection().DeleteOne(ctx, bson.M{"_id": cartID})
	return err
}

// AddItemToCart добавляет товар в корзину или увеличивает его количество, если он уже там есть.
func AddItemToCart(userID, productID primitive.ObjectID, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := getCartCollection()

	// Пытаемся обновить количество, если элемент уже есть.
	filter := bson.M{"userId": userID, "items.productId": productID}
	update := bson.M{"$inc": bson.M{"items.$.quantity": quantity}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// Если элемент не найден, добавляем новый.
	if result.ModifiedCount == 0 {
		newItem := models.CartItem{
			ID:        primitive.NewObjectID(),
			ProductID: productID,
			Quantity:  quantity,
		}
		filter = bson.M{"userId": userID}
		update = bson.M{"$push": bson.M{"items": newItem}}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveItemFromCart удаляет товар из корзины.
func RemoveItemFromCart(userID, productID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"userId": userID}
	update := bson.M{"$pull": bson.M{"items": bson.M{"productId": productID}}}
	_, err := getCartCollection().UpdateOne(ctx, filter, update)
	return err
}
