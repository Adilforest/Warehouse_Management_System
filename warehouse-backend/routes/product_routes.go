package routes

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"warehouse-backend/database"
	"warehouse-backend/models"
)

const (
	productCreatedMessage     = "Product created successfully"
	allProductsDeletedMessage = "All products deleted successfully"
)

// HandleCreateProduct создает новый продукт
func HandleCreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}

	product.ID = primitive.NewObjectID()

	collection := database.GetCollection("warehouse", "products")
	if collection == nil {
		handleError(c, http.StatusInternalServerError, "Failed to get collection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create product: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     productCreatedMessage,
		"product":     product,
		"inserted_id": result.InsertedID,
		"status":      "success",
	})
}

// HandleGetProduct возвращает продукт по ID
func HandleGetProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid product ID: "+err.Error())
		return
	}

	collection := database.GetCollection("warehouse", "products")
	if collection == nil {
		handleError(c, http.StatusInternalServerError, "Failed to get collection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product models.Product
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		handleError(c, http.StatusNotFound, "Product not found: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

// HandleGetAllProducts возвращает список продуктов с фильтрацией, сортировкой и пагинацией
func HandleGetAllProducts(c *gin.Context) {
	// Параметры пагинации
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Параметры фильтрации
	productType := c.Query("type")
	minPrice, _ := strconv.ParseFloat(c.Query("minPrice"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("maxPrice"), 64)
	brand := c.Query("brand")
	ram := c.Query("ram")
	storage := c.Query("storage")
	processor := c.Query("processor")
	color := c.Query("color")

	// Параметры сортировки
	sortBy := c.Query("sortBy")   // Например, "price", "brand", "model"
	sortOrder := c.Query("order") // "asc" или "desc"

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

	collection := database.GetCollection("warehouse", "products")
	if collection == nil {
		handleError(c, http.StatusInternalServerError, "Failed to get collection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to fetch products: "+err.Error())
		return
	}

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to decode products: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Products retrieved successfully",
		"data":    products,
	})
}

// HandleUpdateProduct обновляет продукт по ID
func HandleUpdateProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid product ID: "+err.Error())
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}

	collection := database.GetCollection("warehouse", "products")
	if collection == nil {
		handleError(c, http.StatusInternalServerError, "Failed to get collection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": updatedProduct}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update product: "+err.Error())
		return
	}

	if result.MatchedCount == 0 {
		handleError(c, http.StatusNotFound, "Product not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product updated successfully",
		"data":    updatedProduct,
	})
}

// HandleDeleteProduct удаляет продукт по ID
func HandleDeleteProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid product ID: "+err.Error())
		return
	}

	collection := database.GetCollection("warehouse", "products")
	if collection == nil {
		handleError(c, http.StatusInternalServerError, "Failed to get collection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete product: "+err.Error())
		return
	}

	if result.DeletedCount == 0 {
		handleError(c, http.StatusNotFound, "Product not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}

// HandleDeleteAllProducts удаляет все продукты
func HandleDeleteAllProducts(c *gin.Context) {
	confirmation := c.Query("confirm")
	if confirmation != "yes" {
		handleError(c, http.StatusBadRequest,
			"Action not confirmed. Add ?confirm=yes to delete all products.")
		return
	}

	collection := database.GetCollection("warehouse", "products")
	if collection == nil {
		handleError(c, http.StatusInternalServerError, "Failed to get collection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete all products: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       allProductsDeletedMessage,
		"deleted_count": result.DeletedCount,
		"status":        "success",
	})
}

// SetupProductRoutes настраивает маршруты для работы с продуктами
func SetupProductRoutes(router *gin.Engine) {
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("/", HandleCreateProduct)
		productRoutes.GET("/:id", HandleGetProduct)
		productRoutes.GET("/", HandleGetAllProducts)
		productRoutes.PUT("/:id", HandleUpdateProduct)
		productRoutes.DELETE("/:id", HandleDeleteProduct)
		productRoutes.DELETE("/deleteAll", HandleDeleteAllProducts)
	}
}

// handleError обрабатывает ошибки и возвращает JSON-ответ
func handleError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": errorMessage,
	})
}
