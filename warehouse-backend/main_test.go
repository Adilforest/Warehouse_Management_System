package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"warehouse-backend/logger"
	"warehouse-backend/models"

	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	// Инициализация логгера
	logger.InitLogger()

	// Инициализация базы данных
	setupDatabase()

	// Настройка тестового сервера
	router := setupRoutes()

	// Создание тестового продукта
	product := models.Product{
		Name:  "Test Product",
		Price: 100.0,
		Type:  "Test Type",
	}
	productJSON, _ := json.Marshal(product)

	// Создание HTTP-запроса
	req, _ := http.NewRequest("POST", "/products/create", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	// Запись ответа
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка статуса ответа
	assert.Equal(t, http.StatusCreated, w.Code)

	// Проверка тела ответа
	var response APIResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, "Product created successfully", response.Message)
	assert.NotNil(t, response.Data)
}

func TestDeleteAllProducts(t *testing.T) {
	// Инициализация логгера
	logger.InitLogger()

	// Инициализация базы данных
	setupDatabase()

	// Настройка тестового сервера
	router := setupRoutes()

	// Удаление всех продуктов
	reqDeleteAll, _ := http.NewRequest("DELETE", "/products/deleteAll", nil)
	wDeleteAll := httptest.NewRecorder()
	router.ServeHTTP(wDeleteAll, reqDeleteAll)

	// Проверка статуса ответа
	assert.Equal(t, http.StatusOK, wDeleteAll.Code)

	// Проверка тела ответа
	var deleteAllResponse APIResponse
	json.Unmarshal(wDeleteAll.Body.Bytes(), &deleteAllResponse)
	assert.Equal(t, "success", deleteAllResponse.Status)
	assert.Equal(t, "All products deleted successfully", deleteAllResponse.Message)
}
