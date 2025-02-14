package controllers

import (
	"net/http"
	"strconv"

	"warehouse-backend/database"
	"warehouse-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddToCartHandler добавляет товар в корзину
func AddToCartHandler(c *gin.Context) {
    userID, exists := c.Get("userID") // Получаем ID пользователя из middleware
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Убедитесь, что userID имеет тип primitive.ObjectID
    userObjectID, ok := userID.(primitive.ObjectID)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID"})
        return
    }

    productIDStr := c.Param("product_id")
    productID, err := primitive.ObjectIDFromHex(productIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    quantityStr := c.Query("quantity")
    quantity, err := strconv.Atoi(quantityStr)
    if err != nil || quantity <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
        return
    }

    // Получаем корзину пользователя
    cart, err := database.GetCartByUserID(userObjectID)
    if err != nil {
        // Если корзины нет, создаем новую
        cart, err = database.CreateCart(userObjectID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
            return
        }
    }

    // Получаем продукт
    product, err := database.GetProductByID(productID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    // Создаем элемент корзины
    item := models.CartItem{
        Product:  *product,
        Quantity: quantity,
    }

    // Добавляем товар в корзину
    err = database.AddItemToCart(cart.ID, item)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

// GetCartHandler возвращает корзину пользователя
func GetCartHandler(c *gin.Context) {
    userID, exists := c.Get("userID") // Получаем ID пользователя из middleware
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Убедитесь, что userID имеет тип primitive.ObjectID
    userObjectID, ok := userID.(primitive.ObjectID)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID"})
        return
    }

    // Получаем корзину пользователя
    cart, err := database.GetCartByUserID(userObjectID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "data":   cart,
    })
}

// RemoveFromCartHandler удаляет товар из корзины
func RemoveFromCartHandler(c *gin.Context) {
    userID, exists := c.Get("userID") // Получаем ID пользователя из middleware
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Убедитесь, что userID имеет тип primitive.ObjectID
    userObjectID, ok := userID.(primitive.ObjectID)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID"})
        return
    }

    productIDStr := c.Param("product_id")
    productID, err := primitive.ObjectIDFromHex(productIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    // Удаляем товар из корзины
    err = database.RemoveItemFromCart(userObjectID, productID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

// ClearCartHandler очищает корзину пользователя
func ClearCartHandler(c *gin.Context) {
    userID, exists := c.Get("userID") // Получаем ID пользователя из middleware
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Убедитесь, что userID имеет тип primitive.ObjectID
    userObjectID, ok := userID.(primitive.ObjectID)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID"})
        return
    }

    // Очищаем корзину
    err := database.ClearCart(userObjectID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}

func PayCartHandler(c *gin.Context) {
	// Получаем ID пользователя из middleware
	userID, exists := c.Get("userID")
	if !exists {
		logrus.Warn("User not authenticated, cannot process payment")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Убедимся, что userID имеет тип primitive.ObjectID
	userObjectID, ok := userID.(primitive.ObjectID)
	if !ok {
		logrus.Error("Failed to extract user ID in payment handler")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract user ID"})
		return
	}

	logrus.Infof("Processing payment for user %s", userObjectID.Hex())
	// Вызываем функцию для оплаты корзины
	success, err := database.PayShoppingCart(userObjectID)
	if err != nil {
		logrus.Errorf("Payment processing error for user %s: %v", userObjectID.Hex(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if success {
		// Если оплата прошла успешно, очищаем корзину
		if err := database.ClearCart(userObjectID); err != nil {
			logrus.Errorf("Payment succeeded but failed to clear cart for user %s: %v", userObjectID.Hex(), err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment succeeded but failed to clear cart"})
			return
		}
		logrus.Infof("Payment successful and cart cleared for user %s", userObjectID.Hex())
		c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
	} else {
		logrus.Warnf("Payment failed for user %s", userObjectID.Hex())
		c.JSON(http.StatusBadRequest, gin.H{"message": "Payment failed"})
	}
}