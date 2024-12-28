package database

import (
	"context"
	"time"
	"warehouse-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateProduct добавляет новый продукт в коллекцию MongoDB.
func CreateProduct(product *models.Product) error {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product.ID = primitive.NewObjectID() // Автоматически создаем ObjectID для продукта.
	_, err := collection.InsertOne(ctx, product)
	return err
}

// GetProductByID возвращает продукт по его ObjectID.
func GetProductByID(id primitive.ObjectID) (*models.Product, error) {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product models.Product
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// GetProductsPaginated возвращает список продуктов с учетом пагинации.
func GetProductsPaginated(limit, offset int) ([]models.Product, error) {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Опции для пагинации
	opts := bson.M{}
	cursor, err := collection.Find(ctx, opts, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// UpdateProduct обновляет существующий продукт в MongoDB по ObjectID.
func UpdateProduct(id primitive.ObjectID, updatedProduct *models.Product) error {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updatedProduct},
	)
	return err
}

// DeleteProduct удаляет продукт в MongoDB по ObjectID.
func DeleteProduct(id primitive.ObjectID) error {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteAllProducts удаляет все продукты из MongoDB.
func DeleteAllProducts() error {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, bson.M{})
	return err
}
