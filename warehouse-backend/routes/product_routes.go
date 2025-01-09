package routes

import (
	"context"
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

func SetupProductRoutes(router *gin.Engine) {
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("/", HandleCreateProduct)

		productRoutes.DELETE("/deleteAll", HandleDeleteAllProducts)
	}
}

func handleError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": errorMessage,
	})
}
