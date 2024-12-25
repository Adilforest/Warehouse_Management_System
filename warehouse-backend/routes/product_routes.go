package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"warehouse-backend/database"
	"warehouse-backend/models"
)

const productCreatedMessage = "Product created successfully"

// HandleCreateProduct creates a new product
func HandleCreateProduct(c *gin.Context) {
	var product models.Product

	// Bind the JSON payload to the product model
	if err := c.ShouldBindJSON(&product); err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Insert the product into the database
	if result := database.DB.Create(&product); result.Error != nil {
		handleError(c, http.StatusInternalServerError, result.Error.Error())
		return
	}

	// Return success response with created product
	c.JSON(http.StatusOK, gin.H{
		"message": productCreatedMessage,
		"product": product,
	})
}

// DeleteAllProducts deletes all products from the database
func HandleDeleteAllProducts(c *gin.Context) {
	err := database.DeleteAllProducts() // Метод, который удаляет все продукты в БД
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": "fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All products deleted successfully"})
}

// SetupProductRoutes sets up routes for product-related actions
func SetupProductRoutes(router *gin.Engine) {
	// Route for creating a product
	router.POST("/products", HandleCreateProduct)

	// Route for deleting all products
	router.DELETE("/products/deleteAll", HandleDeleteAllProducts)
}

// handleError is a helper function to return error responses
func handleError(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{"error": errorMessage})
}
