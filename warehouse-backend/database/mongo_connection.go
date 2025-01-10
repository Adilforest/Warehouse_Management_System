package database

import (
	"context"
	"time"
	"warehouse-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateProduct создает новый продукт в коллекции "products"
func CreateProduct(product *models.Product) error {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product.ID = primitive.NewObjectID()
	_, err := collection.InsertOne(ctx, product)
	return err
}

// GetProductByID возвращает продукт по его ID
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

// GetProductsPaginated возвращает список продуктов с фильтрацией, сортировкой и пагинацией
func GetProductsPaginated(limit, offset int, productType string, minPrice, maxPrice float64, brand, ram, storage, processor, color, sortBy, sortOrder string) ([]models.Product, error) {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем фильтр для MongoDB
	filter := bson.M{}

	if productType != "" {
		filter["type"] = productType
	}
	if minPrice > 0 || maxPrice > 0 {
		priceFilter := bson.M{}
		if minPrice > 0 {
			priceFilter["$gte"] = minPrice
		}
		if maxPrice > 0 {
			priceFilter["$lte"] = maxPrice
		}
		filter["price"] = priceFilter
	}
	if brand != "" {
		filter["brand"] = brand
	}
	if ram != "" {
		filter["ram"] = ram
	}
	if storage != "" {
		filter["storage"] = storage
	}
	if processor != "" {
		filter["processor"] = processor
	}
	if color != "" {
		filter["color"] = color
	}

	// Создаем опции для сортировки
	options := options.Find()
	if sortBy != "" {
		order := 1 // По умолчанию сортировка по возрастанию
		if sortOrder == "desc" {
			order = -1
		}
		options.SetSort(bson.D{{sortBy, order}})
	}

	// Пагинация
	options.SetLimit(int64(limit))
	options.SetSkip(int64(offset))

	// Выполняем запрос к MongoDB
	cursor, err := collection.Find(ctx, filter, options)
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

// UpdateProduct обновляет продукт по его ID
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

// DeleteProduct удаляет продукт по его ID
func DeleteProduct(id primitive.ObjectID) error {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// DeleteAllProducts удаляет все продукты из коллекции
func DeleteAllProducts() error {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, bson.M{})
	return err
}
