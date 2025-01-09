package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"warehouse-backend/database"
	"warehouse-backend/models"
)

const serverPort = ":8080"

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	setupDatabase()

	router := setupRoutes()

	go func() {
		log.Println("Server is running on port " + serverPort)
		if err := router.Run(serverPort); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	database.DisconnectMongoDB()
}

func setupDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not set in .env")
	}

	err = database.InitMongoDB(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
}

func setupRoutes() *gin.Engine {
	router := gin.Default()
	router.RedirectTrailingSlash = false

	allowOrigins := os.Getenv("CORS_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:63342,http://localhost:8080"
	}
	origins := strings.Split(allowOrigins, ",")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
	c.String(http.StatusOK, "Welcome to the Warehouse Backend!")
}

func createResponse(status, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func createProductHandler(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail",
			"Invalid JSON payload", nil))
		return
	}

	product.ID = primitive.NewObjectID()

	// Вставляем продукт в MongoDB
	err := database.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createResponse("fail",
			"Failed to create product", nil))
		return
	}

	c.JSON(http.StatusCreated, createResponse("success",
		"Product created successfully", product))
}

func getProductHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse("fail",
			"Invalid product ID", nil))
		return
	}

	product, err := database.GetProductByID(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, createResponse("fail",
			"Product not found", nil))
		return
	}

	c.JSON(http.StatusOK, createResponse("success",
		"Product retrieved successfully", product))
}

func getAllProductsHandler(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	products, err := database.GetProductsPaginated(limit, offset)
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
