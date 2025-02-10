package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"warehouse-backend/controllers"
	"warehouse-backend/database"
	"warehouse-backend/logger"
	"warehouse-backend/middleware"
	"warehouse-backend/models"
	"warehouse-backend/routes"
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
	// Load environment variables from .env file
	err := godotenv.Load("/app/.env") // Explicitly specify the path
	if err != nil {
		logger.Log.Error("setup_database", "Error loading .env file", map[string]interface{}{
			"error": err.Error(),
		}, err)
		logger.Log.Fatal("Error loading .env file: ", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		logger.Log.Error("setup_database", "MONGO_URI is not set in .env", map[string]interface{}{}, nil)
		logger.Log.Fatal("MONGO_URI is not set in .env")
	}

	logger.Log.Info("Connecting to MongoDB", map[string]interface{}{
		"mongoURI": mongoURI,
	})

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
			"ip":     c.ClientIP(),
			"agent":  c.Request.UserAgent(),
		}
		logger.LogInfo("http_request", "Request received", fields)
	})

	// Сервировка статичных файлов и маршрута для index.html
	router.Static("/static", "./static") // Папка для CSS, JS и прочих статичных файлов

	// Главная страница
	router.GET("/", func(c *gin.Context) {
		// Логирование посещения главной страницы
		logger.LogInfo("page_visit", "Index page visited", map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"ip":     c.ClientIP(),
			"agent":  c.Request.UserAgent(),
		})
		c.File("./public/index.html") // Отсылка index.html (путь до файла)
	})

	// Маршруты для работы с продуктами
	productRoutes := router.Group("/products")
	{
		productRoutes.POST("/create", createProductHandler)
		productRoutes.GET("/:id", getProductHandler)
		productRoutes.GET("/", getAllProductsHandler)
		productRoutes.PUT("/:id", updateProductHandler)
		productRoutes.DELETE("/deleteAll", deleteAllProductsHandler)
		productRoutes.DELETE("/:id", deleteProductHandler)
	}

	// Настройка rate limiter: 1 запрос в 15 секунд
	limiter := middleware.RateLimiterMiddleware(rate.Every(15*time.Second), 1)

	// Маршрут для обработки запросов от формы "Contact Us" с rate limiting
	router.POST("/api/contact", limiter, controllers.ContactController)

	// Настройка маршрутов для аутентификации
	routes.SetupAuthRoutes(router)

	// Настройка защищенных маршрутов (требуется авторизация)
	protectedRoutes := router.Group("/protected")
	protectedRoutes.Use(middleware.AuthMiddleware()) // Middleware для проверки JWT
	{
		// Маршруты для профиля пользователя
		protectedRoutes.GET("/profile", controllers.GetProfile)    // Получение данных профиля
		protectedRoutes.PUT("/profile", controllers.UpdateProfile) // Обновление данных профиля
	}

	// Middleware для проверки роли администратора
	adminOnlyMiddleware := func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to perform this action"})
			c.Abort()
			return
		}
		c.Next()
	}

	// Маршруты для работы с пользователями
	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware()) // Middleware для проверки JWT
	{
		userRoutes.GET("/", adminOnlyMiddleware, controllers.GetAllUsers)                                                      // Только администраторы могут просматривать всех пользователей
		userRoutes.GET("/:id", adminOnlyMiddleware, controllers.GetUserByID)                                                   // Только администраторы могут просматривать пользователя по ID
		userRoutes.PUT("/:id", adminOnlyMiddleware, controllers.UpdateUser)                                                    // Только администраторы могут обновлять пользователя
		userRoutes.DELETE("/:id", adminOnlyMiddleware, controllers.DeleteUser)                                                 // Только администраторы могут удалять пользователя
		userRoutes.DELETE("/deleteAll", adminOnlyMiddleware, controllers.DeleteAllUsers)                                       // Только администраторы могут удалять всех пользователей
		userRoutes.PUT("/:id/role", middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware(), controllers.UpdateUserRole) // Только администраторы могут изменять роль пользователя
	}

	return router
}

func createProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		logger.Log.Error("create_product", "Invalid payload", map[string]interface{}{
			"error": err.Error(),
		}, err) // Добавлен четвертый аргумент
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
		logger.Log.Error("get_product", "Invalid product ID", map[string]interface{}{
			"product_id": id,
		}, err) // Добавлен четвертый аргумент
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
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "30"))
	if err != nil || limit <= 0 || limit > 100 {
		logger.LogWarning("get_all_products", "Invalid limit value, using default value", map[string]interface{}{
			"limit": c.Query("limit"),
		})
		limit = 30 // Значение по умолчанию
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		logger.LogWarning("get_all_products", "Invalid offset value, using default value", map[string]interface{}{
			"offset": c.Query("offset"),
		})
		offset = 0 // Значение по умолчанию
	}

	// Параметры фильтрации
	productType := c.DefaultQuery("type", "")
	priceRanges := c.QueryArray("price") // Получаем массив выбранных диапазонов цен
	brand := c.DefaultQuery("brand", "")
	ram := c.DefaultQuery("ram", "")
	storage := c.DefaultQuery("storage", "")
	processor := c.DefaultQuery("processor", "")
	color := c.DefaultQuery("color", "")

	// Параметры сортировки
	sortBy := c.DefaultQuery("sortBy", "")   // Например, "price", "brand", "model"
	sortOrder := c.DefaultQuery("order", "") // "asc" или "desc"

	// Логирование параметров запроса
	logger.LogInfo("get_all_products", "Fetching products with filters", map[string]interface{}{
		"limit":       limit,
		"offset":      offset,
		"type":        productType,
		"priceRanges": priceRanges,
		"brand":       brand,
		"ram":         ram,
		"storage":     storage,
		"processor":   processor,
		"color":       color,
		"sort_by":     sortBy,
		"order":       sortOrder,
	})

	// Вызов функции для получения продуктов (работа с базой данных)
	products, total, err := database.GetProductsPaginated(
		limit,
		offset,
		productType,
		priceRanges,
		brand,
		ram,
		storage,
		processor,
		color,
		sortBy,
		sortOrder,
	)

	if err != nil {
		// Обработка ошибки базы данных
		logger.Log.Error("get_all_products", "Failed to fetch products from database", map[string]interface{}{
			"error": err.Error(),
		}, err)
		c.JSON(http.StatusInternalServerError, createResponse("fail", "Failed to fetch products", nil))
		return
	}

	// Проверка: пустой результат
	if len(products) == 0 {
		logger.LogWarning("get_all_products", "No products match the filter", map[string]interface{}{
			"filters": map[string]interface{}{
				"limit":       limit,
				"offset":      offset,
				"type":        productType,
				"priceRanges": priceRanges,
				"brand":       brand,
				"ram":         ram,
				"storage":     storage,
				"processor":   processor,
				"color":       color,
				"sort_by":     sortBy,
				"order":       sortOrder,
			},
		})
		c.JSON(http.StatusNotFound, createResponse("fail", "No products match the filter", nil))
		return
	}

	// Логирование успешного получения данных
	logger.LogInfo("get_all_products", "Products fetched successfully", map[string]interface{}{
		"total_products": total,
		"filters": map[string]interface{}{
			"limit":       limit,
			"offset":      offset,
			"type":        productType,
			"priceRanges": priceRanges,
			"brand":       brand,
			"ram":         ram,
			"storage":     storage,
			"processor":   processor,
			"color":       color,
			"sort_by":     sortBy,
			"order":       sortOrder,
		},
	})

	// Отправка успешного ответа клиенту
	c.JSON(http.StatusOK, createResponse("success", "Products retrieved successfully", gin.H{
		"data":  products,
		"total": total,
	}))
}

func updateProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error("update_product", "Invalid product ID", map[string]interface{}{
			"product_id": id,
		}, err) // Добавлен четвертый аргумент
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid product ID", nil))
		return
	}

	var updatedProduct models.Product
	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		logger.Log.Error("update_product", "Invalid JSON payload", map[string]interface{}{
			"error": err.Error(),
		}, err) // Добавлен четвертый аргумент
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid JSON payload", nil))
		return
	}

	err = database.UpdateProduct(objectID, &updatedProduct)
	if err != nil {
		logger.LogDBError("update_product", "Failed to update product", err)
		c.JSON(http.StatusNotFound, createResponse("fail", "Failed to update product: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "Product updated successfully", updatedProduct))
}

func deleteProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Log.Error("delete_product", "Invalid product ID", map[string]interface{}{
			"product_id": id,
		}, err) // Добавлен четвертый аргумент
		c.JSON(http.StatusBadRequest, createResponse("fail", "Invalid product ID", nil))
		return
	}

	err = database.DeleteProduct(objectID)
	if err != nil {
		logger.LogDBError("delete_product", "Failed to delete product", err)
		c.JSON(http.StatusNotFound, createResponse("fail", "Product not found: "+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "Product deleted successfully", nil))
}

func deleteAllProductsHandler(c *gin.Context) {
	err := database.DeleteAllProducts()
	if err != nil {
		logger.LogDBError("delete_all_products", "Failed to delete all products", err)
		c.JSON(http.StatusInternalServerError, createResponse("fail", "Failed to delete all products", nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success", "All products deleted successfully", nil))
}

func createResponse(status, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
