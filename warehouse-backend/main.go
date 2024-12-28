package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	// Элегантное завершение сервера с обработкой сигналов
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
	// Освобождение ресурсов (например, закрытие подключений к БД)
}

func setupDatabase() {
	database.ConnectPostgres()
}

func setupRoutes() *gin.Engine {
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:63342"}, // Измените на URL фронтенда в production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Главная страница
	router.GET("/", handleHome)

	// Тестовые маршруты
	router.GET("/get", handleGetRequest)
	router.POST("/post", handlePostRequest)

	// Группа маршрутов для работы с продуктами
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

// Обработчик GET-запроса
func handleGetRequest(c *gin.Context) {
	response := createResponse("success", "GET request processed successfully", nil)
	c.JSON(http.StatusOK, response)
}

// Обработчик POST-запроса
func handlePostRequest(c *gin.Context) {
	var rawData map[string]interface{}

	if err := c.ShouldBindJSON(&rawData); err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid JSON payload", nil))
		return
	}

	messageValue, exists := rawData["message"]
	if !exists {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Missing 'message' field", nil))
		return
	}

	message, ok := messageValue.(string)
	if !ok || message == "" {
		c.JSON(http.StatusOK, createResponse("success", "Empty or invalid message received", nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "Data received successfully with message: "+message, nil))
}

// createResponse создает стандартный ответ API
func createResponse(status, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// CRUD Handlers for Product

// Создание продукта
func createProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid JSON payload", nil))
		return
	}

	if err := database.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, createResponse("fail", "Failed to create product", nil))
		return
	}

	c.JSON(http.StatusCreated, createResponse("success", "Product created successfully", product))
}

// Получение продукта по ID
func getProductHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid product ID", nil))
		return
	}

	product, err := database.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, createResponse("fail", "Product not found", nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "Product retrieved successfully", product))
}

// Пагинация при получении всех продуктов
func getAllProductsHandler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default page size = 10
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid product ID", nil))
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid JSON payload", nil))
		return
	}

	err = database.UpdateProduct(uint(id), &updatedProduct)
	if err != nil {
		c.JSON(http.StatusNotFound, createResponse("fail", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "Product updated successfully", nil))
}

// Удаление продукта по ID
func deleteProductHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid product ID", nil))
		return
	}

	err = database.DeleteProduct(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, createResponse("fail", err.Error(), nil))
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
