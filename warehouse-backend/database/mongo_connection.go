package database

import (
	"context"
	"log"
	"strings"
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
func GetProductsPaginated(limit, offset int, productType string, minPrice, maxPrice float64, brand, ram, storage, processor, color, sortBy, sortOrder string) ([]models.Product, int64, error) {
	collection := GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем фильтр для MongoDB
	filter := bson.M{}

	// Обработка множественных значений для type
	if productType != "" {
		types := strings.Split(productType, ",")
		filter["type"] = bson.M{"$in": types}
	}

	// Обработка множественных значений для brand
	if brand != "" {
		brands := strings.Split(brand, ",")
		filter["brand"] = bson.M{"$in": brands}
	}

	// Обработка множественных значений для ram
	if ram != "" {
		rams := strings.Split(ram, ",")
		filter["ram"] = bson.M{"$in": rams}
	}

	// Обработка множественных значений для storage
	if storage != "" {
		storages := strings.Split(storage, ",")
		filter["storage"] = bson.M{"$in": storages}
	}

	// Обработка множественных значений для processor
	if processor != "" {
		processors := strings.Split(processor, ",")
		filter["processor"] = bson.M{"$in": processors}
	}

	// Обработка множественных значений для color
	if color != "" {
		colors := strings.Split(color, ",")
		filter["color"] = bson.M{"$in": colors}
	}

	// Фильтр по цене
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

	// Логирование фильтра
	log.Printf("Filter: %+v", filter)

	// Получаем общее количество товаров
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
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

	// Логирование опций
	log.Printf("Options: %+v", options)

	// Выполняем запрос к MongoDB
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, err
	}

	return products, total, nil
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
