package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
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

	// Привязка JSON-контента из тела запроса к модели Product
	if err := c.ShouldBindJSON(&product); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid JSON format: "+err.Error())
		return
	}

	// Генерируем новый ObjectID для продукта
	product.ID = primitive.NewObjectID()

	// Получаем коллекцию "products"
	collection := database.GetCollection("warehouse", "products")
	if collection == nil {
		handleError(c, http.StatusInternalServerError, "Failed to get collection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Вставляем новый продукт в коллекцию
	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create product: "+err.Error())
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message":     productCreatedMessage,
		"product":     product,
		"inserted_id": result.InsertedID, // Добавляем ID нового объекта в ответ
		"status":      "success",
	})
}

// HandleDeleteAllProducts удаляет все продукты из базы данных
func HandleDeleteAllProducts(c *gin.Context) {
	confirmation := c.Query("confirm") // Проверяем подтверждение через query-параметр
	if confirmation != "yes" {
		handleError(c, http.StatusBadRequest, "Action not confirmed. Add ?confirm=yes to delete all products.")
		return
	}

	// Получаем коллекцию "products"
	collection := database.GetCollection("warehouse", "products")
	if collection == nil {
		handleError(c, http.StatusInternalServerError, "Failed to get collection")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выполняем удаление всех записей в коллекции
	result, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete all products: "+err.Error())
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"message":       allProductsDeletedMessage,
		"deleted_count": result.DeletedCount,
		"status":        "success",
	})
}

// SetupProductRoutes настраивает маршруты для управления продуктами
func SetupProductRoutes(router *gin.Engine) {
	// Группа маршрутов для продуктов
	productRoutes := router.Group("/products")
	{
		// Маршрут для создания продукта
		productRoutes.POST("/", HandleCreateProduct)

		// Маршрут для удаления всех продуктов
		productRoutes.DELETE("/deleteAll", HandleDeleteAllProducts)
	}
}

// handleError — вспомогательная функция для возврата ошибок в одном формате
func handleError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": errorMessage,
	})
}
