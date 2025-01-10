package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
	"warehouse-backend/database"
	"warehouse-backend/models"
)

// CreateProduct создает новый продукт
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid JSON payload",
			"error":   err.Error(),
		})
		return
	}

	product.ID = primitive.NewObjectID()

	err := database.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to create product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Product created successfully",
		"data":    product,
	})
}

// GetProduct возвращает продукт по ID
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
		return
	}

	product, err := database.GetProductByID(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "Product not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

// GetAllProducts возвращает список продуктов с фильтрацией, сортировкой и пагинацией
func GetAllProducts(c *gin.Context) {
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

	// Получаем отфильтрованные, отсортированные и пагинированные продукты
	products, err := database.GetProductsPaginated(limit, offset, productType, minPrice, maxPrice, brand, ram, storage, processor, color, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to fetch products",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Products retrieved successfully",
		"data":    products,
	})
}

// UpdateProduct обновляет продукт по ID
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid JSON payload",
			"error":   err.Error(),
		})
		return
	}

	err = database.UpdateProduct(objectID, &updatedProduct)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "Failed to update product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product updated successfully",
		"data":    updatedProduct,
	})
}

// DeleteProduct удаляет продукт по ID
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
		return
	}

	err = database.DeleteProduct(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "Failed to delete product",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}
