package routes

import (
	"transaction-service/controllers"
    "github.com/gin-gonic/gin"
)

func SetupTransactionRoutes(router *gin.Engine) {
    api := router.Group("/api/v1/transactions")
    {
        api.POST("/", controllers.CreateTransaction)
        api.GET("/:id", controllers.GetTransactionByID)
        api.PUT("/:id/status", controllers.UpdateTransactionStatus)
    }
}