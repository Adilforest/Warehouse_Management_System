package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"warehouse-backend/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionRequest struct {
	CartID primitive.ObjectID  `json:"cart_id"`
	UserID primitive.ObjectID  `json:"user_id"`
	Items  []models.CartItem   `json:"items"`
	Total  float64             `json:"total"`
}

// CreateCart создает новую корзину для пользователя
func CreateCart(userID primitive.ObjectID) (*models.Cart, error) {
	collection := GetCollection("warehouse", "carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cart := models.Cart{
		UserID:  userID,
		Items:   []models.CartItem{},
		Total:   0,
		Created: time.Now(),
		Updated: time.Now(),
	}

	result, err := collection.InsertOne(ctx, cart)
	if err != nil {
		return nil, err
	}

	cart.ID = result.InsertedID.(primitive.ObjectID)
	return &cart, nil
}

// GetCartByUserID возвращает корзину пользователя по его ID
func GetCartByUserID(userID primitive.ObjectID) (*models.Cart, error) {
	collection := GetCollection("warehouse", "carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart models.Cart
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

// AddItemToCart добавляет товар в корзину
func AddItemToCart(cartID primitive.ObjectID, item models.CartItem) error {
	collection := GetCollection("warehouse", "carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Проверяем, существует ли уже такой товар в корзине
	// Обратите внимание, что если идентификатор продукта хранится в поле "_id" внутри вложенного документа "product",
	// то фильтр должен использовать "items.product._id"
	filter := bson.M{"_id": cartID, "items.product._id": item.Product.ID}
	update := bson.M{
		"$inc": bson.M{"items.$.quantity": item.Quantity},
		"$set": bson.M{"updated_at": time.Now()},
	}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		// Если товар не найден, добавляем новый элемент
		filter = bson.M{"_id": cartID}
		update = bson.M{
			"$push": bson.M{"items": item},
			"$set":  bson.M{"updated_at": time.Now()},
		}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}

	// Обновляем общую стоимость корзины
	return UpdateCartTotal(cartID)
}

// RemoveItemFromCart удаляет товар из корзины
func RemoveItemFromCart(userID, productID primitive.ObjectID) error {
	collection := GetCollection("warehouse", "carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Используем правильное имя поля "product._id"
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$pull": bson.M{"items": bson.M{"product._id": productID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// Получаем корзину по userID, чтобы обновить общую сумму
	cart, err := GetCartByUserID(userID)
	if err != nil {
		return err
	}

	return UpdateCartTotal(cart.ID)
}

// ClearCart очищает корзину пользователя
func ClearCart(userID primitive.ObjectID) error {
	collection := GetCollection("warehouse", "carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"items":   []models.CartItem{},
			"total":   0,
			"updated": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func UpdateCartTotal(cartID primitive.ObjectID) error {
	collection := GetCollection("warehouse", "carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{"$match": bson.M{"_id": cartID}},
		{"$unwind": "$items"},
		{"$group": bson.M{
			"_id":   "$_id",
			"total": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$items.product.price", "$items.quantity"}}},
		}},
	}

	var result struct {
		ID    primitive.ObjectID `bson:"_id"`
		Total float64            `bson:"total"`
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return err
		}
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": cartID}, bson.M{"$set": bson.M{"total": result.Total}})
	return err
}

// PayShoppingCart отправляет данные корзины и клиента в микросервис транзакций.
// Возвращает true, если транзакция прошла успешно, и false в противном случае.
func PayShoppingCart(userID primitive.ObjectID) (bool, error) {
	// Получаем корзину пользователя
	cart, err := GetCartByUserID(userID)
	if err != nil {
		logrus.Errorf("Failed to retrieve cart for user %s: %v", userID.Hex(), err)
		return false, fmt.Errorf("failed to retrieve cart: %v", err)
	}
	logrus.Infof("Retrieved cart for user %s", userID.Hex())

	// Генерируем новый идентификатор транзакции
	transactionID := primitive.NewObjectID()
	logrus.Infof("Generated transaction ID: %s", transactionID.Hex())

	// Формируем пейлоад для микросервиса транзакций
	payload := map[string]interface{}{
		"transaction_id": transactionID.Hex(),
		"user_id":        userID.Hex(),
		"cart":           cart, // Можно передать только необходимые поля
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorf("Failed to marshal payload for user %s: %v", userID.Hex(), err)
		return false, fmt.Errorf("failed to marshal payload: %v", err)
	}

	// URL микросервиса транзакций (при необходимости измените)
	transactionURL := "http://localhost:8081/transactions"
	logrus.Infof("Sending payment request to %s for user %s", transactionURL, userID.Hex())

	// Отправляем POST-запрос в микросервис транзакций
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", transactionURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorf("Failed to create request for user %s: %v", userID.Hex(), err)
		return false, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Payment request failed for user %s: %v", userID.Hex(), err)
		return false, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Transaction microservice returned status %d for user %s", resp.StatusCode, userID.Hex())
		return false, fmt.Errorf("transaction microservice returned status %d", resp.StatusCode)
	}

	// Ожидаем JSON-ответ, например: {"success": true}
	var res struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		logrus.Errorf("Failed to decode response for user %s: %v", userID.Hex(), err)
		return false, fmt.Errorf("failed to decode response: %v", err)
	}

	logrus.Infof("Payment microservice response for user %s: success=%v", userID.Hex(), res.Success)
	return res.Success, nil
}
