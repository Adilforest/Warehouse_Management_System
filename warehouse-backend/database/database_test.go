package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitMongoDB(t *testing.T) {
	// Используем тестовый URI
	err := InitMongoDB("mongodb://localhost:27017")
	assert.NoError(t, err)
	assert.NotNil(t, MongoClient)
}

func TestGetCollection(t *testing.T) {
	// Инициализация клиента
	InitMongoDB("mongodb://localhost:27017")

	// Получение коллекции
	collection := GetCollection("testdb", "testcollection")
	assert.NotNil(t, collection)
}
