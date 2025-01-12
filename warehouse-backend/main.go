package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"warehouse-backend/database"
	"warehouse-backend/logger"
	"warehouse-backend/models"
)

const serverPort = ":8080"

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Инициализация логгера
	logger.InitLogger()
	logger.Log.Info("Initializing logger...")

	// Настройка базы данных
	setupDatabase()

	// Настройка маршрутов
	router := setupRoutes()

	// Запуск сервера в горутине
	go func() {
		logger.LogInfo("server_start", "Server is starting...", map[string]interface{}{
			"port": serverPort,
		})
		if err := router.Run(serverPort); err != nil {
			logger.Log.Fatal(err)
		}
	}()

	// Закрытие сервера по сигналам
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Log.Warn("Shutting down server...")
	database.DisconnectMongoDB()
}

func setupDatabase() {
	// Загрузка переменных из .env файла
	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Error loading .env file: ", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		logger.Log.Fatal("MONGO_URI is not set in .env")
	}

	err = database.InitMongoDB(mongoURI)
	if err != nil {
		logger.LogDBError("connect_to_mongo", "Failed to connect to MongoDB", err)
		logger.Log.Fatal(err)
	}
	logger.Log.Info("Connected to MongoDB successfully")
}

func setupRoutes() *gin.Engine {
	router := gin.Default()

	// Настройка CORS
	allowOrigins := os.Getenv("CORS_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3000,http://localhost:8080"
	}
	origins := strings.Split(allowOrigins, ",")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Логирование запросов
	router.Use(func(c *gin.Context) {
		c.Next()
		fields := map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": c.Writer.Status(),
		}
		logger.LogInfo("http_request", "Request received", fields)
	})

	// Маршруты
	router.GET("/", handleHome)

	productRoutes := router.Group("/products")
	{
		productRoutes.POST("/create", createProductHandler)
		productRoutes.GET("/:id", getProductHandler)
		productRoutes.GET("/", getAllProductsHandler)
		productRoutes.PUT("/:id", updateProductHandler)
		productRoutes.DELETE("/deleteAll", deleteAllProductsHandler)
		productRoutes.DELETE("/:id", deleteProductHandler)
	}

	return router
}

func handleHome(c *gin.Context) {
	logger.LogRequest("GET", "/", "200")
	c.String(http.StatusOK, "Welcome to the Warehouse Backend!")
}

func createProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		logger.LogError("create_product", "Invalid payload", map[string]interface{}{
			"error": err.Error(),
		})
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid JSON payload", nil))
		return
	}

	product.ID = primitive.NewObjectID()

	err := database.CreateProduct(&product)
	if err != nil {
		logger.LogDBError("create_product", "Failed to save product", err)
		c.JSON(http.StatusInternalServerError, createResponse("fail", "Failed to create product", nil))
		return
	}

	logger.LogInfo("create_product", "Product created successfully", map[string]interface{}{
		"product_id": product.ID,
		"product":    product,
	})
	c.JSON(http.StatusCreated, createResponse("success", "Product created successfully", product))
}

func getProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.LogError("get_product", "Invalid product ID", map[string]interface{}{
			"product_id": id,
		})
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid product ID", nil))
		return
	}

	product, err := database.GetProductByID(objectID)
	if err != nil {
		logger.LogDBError("get_product", "Product not found", err)
		c.JSON(http.StatusNotFound, createResponse("fail", "Product not found", nil))
		return
	}

	logger.LogInfo("get_product", "Product fetched successfully", map[string]interface{}{
		"product": product,
	})
	c.JSON(http.StatusOK, createResponse("success", "Product retrieved successfully", product))
}

func getAllProductsHandler(c *gin.Context) {
	// Параметры пагинации
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
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

	// Вызов GetProductsPaginated с новыми параметрами
	products, err := database.GetProductsPaginated(
		limit,       // limit
		offset,      // offset
		productType, // productType
		minPrice,    // minPrice
		maxPrice,    // maxPrice
		brand,       // brand
		ram,         // ram
		storage,     // storage
		processor,   // processor
		color,       // color
		sortBy,      // sortBy
		sortOrder,   // sortOrder
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createResponse("fail",
			"Failed to fetch products", nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success",
		"Products retrieved successfully", products))
}

func updateProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail",
			"Invalid product ID", nil))
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail",
			"Invalid JSON payload", nil))
		return
	}

	err = database.UpdateProduct(objectID, &updatedProduct)
	if err != nil {
		c.JSON(http.StatusNotFound, createResponse("fail",
			"Failed to update product: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success",
		"Product updated successfully", updatedProduct))
}

func deleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail",
			"Invalid product ID", nil))
		return
	}

	err = database.DeleteProduct(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, createResponse("fail",
			"Product not found: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success",
		"Product deleted successfully", nil))
}

func deleteAllProductsHandler(c *gin.Context) {
	err := database.DeleteAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, createResponse("fail",
			"Failed to delete all products", nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success",
		"All products deleted successfully", nil))
}

func createResponse(status, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
