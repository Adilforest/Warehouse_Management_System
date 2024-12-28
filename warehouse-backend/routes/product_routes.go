package routes

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"warehouse-backend/database"
	"warehouse-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const productCreatedMessage = "Product created successfully"

// HandleCreateProduct создает новый продукт
func HandleCreateProduct(c *gin.Context) {
	var product models.Product

	// Привязка JSON-контента из тела запроса к модели Product
	if err := c.ShouldBindJSON(&product); err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	product.ID = primitive.NewObjectID() // Генерируем новый ObjectID для продукта

	// Получаем коллекцию "products"
	collection := database.GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Вставляем новый продукт в коллекцию
	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create product: "+err.Error())
		return
	}

	// Возвращаем успешный ответ с данными созданного объекта
	c.JSON(http.StatusOK, gin.H{
		"message": productCreatedMessage,
		"product": product,
	})
}

// HandleDeleteAllProducts удаляет все продукты из базы данных
func HandleDeleteAllProducts(c *gin.Context) {
	// Получаем коллекцию "products"
	collection := database.GetCollection("warehouse", "products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выполняем удаление всех записей в коллекции
	result, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete all products: "+err.Error())
		return
	}

	// Возвращаем успешный ответ с информацией о количестве удаленных записей
	c.JSON(http.StatusOK, gin.H{
		"message": "All products deleted successfully",
		"deleted": result.DeletedCount,
		"status":  "success",
	})
}

// SetupProductRoutes настраивает маршруты для управления продуктами
func SetupProductRoutes(router *gin.Engine) {
	// Маршрут для создания продукта
	router.POST("/products", HandleCreateProduct)

	// Маршрут для удаления всех продуктов
	router.DELETE("/products/deleteAll", HandleDeleteAllProducts)
}

// handleError — вспомогательная функция для возврата ошибок в ответах
func handleError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{"error": errorMessage})
}
