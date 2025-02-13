package controllers

import (
	"net/http"
	"warehouse-backend/database"
	"warehouse-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCart возвращает корзину текущего пользователя.
func GetCart(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cart, err := database.GetCartByUserID(userID)
	if err != nil {
		// Если корзина не найдена, создаём новую.
		newCart := models.Cart{
			ID:     primitive.NewObjectID(),
			UserID: userID,
			Items:  []models.CartItem{},
		}
		if err := database.CreateCart(&newCart); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create cart"})
			return
		}
		cart = &newCart
	}

	c.JSON(http.StatusOK, cart)
}

// AddToCart добавляет товар в корзину.
func AddToCart(c *gin.Context) {
	var req struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}
	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than zero"})
		return
	}

	if err := database.AddItemToCart(userID, productID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add item to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart"})
}

// RemoveFromCart удаляет товар из корзины.
func RemoveFromCart(c *gin.Context) {
	var req struct {
		ProductID string `json:"productId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	productID, err := primitive.ObjectIDFromHex(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := database.RemoveItemFromCart(userID, productID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart"})
}

// ClearCart очищает корзину пользователя.
func ClearCart(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	cart, err := database.GetCartByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve cart"})
		return
	}
	cart.Items = []models.CartItem{}
	if err := database.UpdateCart(cart); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not clear cart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared"})
}
