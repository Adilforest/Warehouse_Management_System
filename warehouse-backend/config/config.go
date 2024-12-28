package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     string
}

// GetConfig загружает конфигурацию из переменных окружения
func GetConfig() Config {
	return Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}

// GetMongoDBURI формирует MongoDB connection string
func (c Config) GetMongoDBURI() string {
	if c.DBUser != "" && c.DBPassword != "" {
		// Строка подключения с учётными данными
		return fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority",
			c.DBUser, c.DBPassword, c.DBHost, c.DBName)
	}
	// Строка подключения без учётных данных
	return fmt.Sprintf("mongodb+srv://%s/%s?retryWrites=true&w=majority", c.DBHost, c.DBName)
}
