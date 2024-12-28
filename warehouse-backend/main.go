package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"warehouse-backend/database"
	"warehouse-backend/models"
)

const serverPort = ":8080"

// APIResponse представляет стандартную структуру ответа сервера
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Установка базы данных
	setupDatabase()

	// Настройка маршрутов
	router := setupRoutes()

	// Элегантное завершение работы сервера с обработкой сигналов
	go func() {
		log.Println("Server is running on port " + serverPort)
		if err := router.Run(serverPort); err != nil {
			log.Fatal(err)
		}
	}()

	// Обработка SIGINT и SIGTERM для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	database.DisconnectMongoDB() // Закрываем соединение с MongoDB
}

func setupDatabase() {
	err := database.InitMongoDB("mongodb://localhost:27017") // Инициализация MongoDB
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
}

func setupRoutes() *gin.Engine {
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:63342"}, // Замените URL для production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Главная страница
	router.GET("/", handleHome)

	// CRUD маршруты для продуктов
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("/create", createProductHandler)
		productRoutes.GET("/:id", getProductHandler)
		productRoutes.GET("/", getAllProductsHandler) // С поддержкой пагинации
		productRoutes.PUT("/:id", updateProductHandler)
		productRoutes.DELETE("/deleteAll", deleteAllProductsHandler)
		productRoutes.DELETE("/:id", deleteProductHandler)
	}

	return router
}

// Обработчик главной страницы
func handleHome(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to the Warehouse Backend!")
}

// createResponse создает стандартный API-ответ
func createResponse(status, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// CRUD обработчики для продуктов

// Создание продукта
func createProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid JSON payload", nil))
		return
	}

	product.ID = primitive.NewObjectID() // Генерируем ObjectID для нового продукта

	// Вставляем продукт в MongoDB
	err := database.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createResponse("fail", "Failed to create product", nil))
		return
	}

	c.JSON(http.StatusCreated, createResponse("success", "Product created successfully", product))
}

// Получение продукта по ID
func getProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id) // Конвертация ID в ObjectID
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid product ID", nil))
		return
	}

	product, err := database.GetProductByID(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, createResponse("fail", "Product not found", nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "Product retrieved successfully", product))
}

// Получение всех продуктов с пагинацией
func getAllProductsHandler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // По умолчанию 10
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := database.GetProductsPaginated(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createResponse("fail", "Failed to fetch products", nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "Products retrieved successfully", products))
}

// Обновление продукта
func updateProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id) // Конвертация ID в ObjectID
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid product ID", nil))
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid JSON payload", nil))
		return
	}

	err = database.UpdateProduct(objectID, &updatedProduct)
	if err != nil {
		c.JSON(http.StatusNotFound, createResponse("fail", "Failed to update product: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "Product updated successfully", updatedProduct))
}

// Удаление продукта по ID
func deleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id) // Конвертация ID в ObjectID
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid product ID", nil))
		return
	}

	err = database.DeleteProduct(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, createResponse("fail", "Product not found: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "Product deleted successfully", nil))
}

// Удаление всех продуктов
func deleteAllProductsHandler(c *gin.Context) {
	err := database.DeleteAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, createResponse("fail", "Failed to delete all products", nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "All products deleted successfully", nil))
}
