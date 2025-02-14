package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "transaction-service/config"
    "transaction-service/database"
    "transaction-service/routes"
)

func main() {
    // Инициализация логгера
    log.Println("Starting transaction service...")

    // Подключение к MongoDB
    if err := database.InitMongoDB(os.Getenv("MONGO_URI")); err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer database.DisconnectMongoDB()

    // Настройка Gin
    router := gin.Default()

    // Настройка CORS
    router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
            return
        }
        c.Next()
    })

    // Настройка маршрутов
    routes.SetupTransactionRoutes(router)

    // Запуск сервера
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }
    log.Printf("Transaction service is running on port %s...", port)
    log.Fatal(router.Run(":" + port))
}