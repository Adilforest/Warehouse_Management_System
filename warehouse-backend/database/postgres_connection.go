package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"warehouse-backend/config"
	"warehouse-backend/models"
)

var Database *gorm.DB

const dbConnectionSuccessMsg = "Successfully connected to PostgreSQL!"

// ConnectPostgres устанавливает соединение с базой данных PostgreSQL и выполняет AutoMigrations
func ConnectPostgres() {
	db := initializeDatabaseConnection()
	Database = db
	fmt.Println(dbConnectionSuccessMsg)

	// AutoMigrate для модели Product
	if err := Database.AutoMigrate(&models.Product{}); err != nil {
		log.Fatalf("Error migrating the database: %v", err)
	}
	fmt.Println("Database migrated successfully!")
}

// initializeDatabaseConnection инициализирует подключение к базе данных
func initializeDatabaseConnection() *gorm.DB {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := config.GetConfig()
	dsn := cfg.GetPostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database. DSN: %v, Error: %v", dsn, err)
	}
	return db
}

// CreateProduct создает новый продукт в базе данных
func CreateProduct(product *models.Product) error {
	result := Database.Create(product)
	if result.Error != nil {
		return fmt.Errorf("error creating product with fields: Name (%v), Type (%v): %v", product.Name, product.Type, result.Error)
	}
	return nil
}

// GetProductByID извлекает продукт по его ID из базы данных
func GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	result := Database.First(&product, id)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return nil, fmt.Errorf("product with ID %d not found", id)
		}
		return nil, fmt.Errorf("error finding product: %v", result.Error)
	}
	return &product, nil
}

// UpdateProduct обновляет существующий продукт по его ID
func UpdateProduct(id uint, productData *models.Product) error {
	var product models.Product
	result := Database.First(&product, id)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return fmt.Errorf("product with ID %d not found", id)
		}
		return fmt.Errorf("error finding product: %v", result.Error)
	}

	// Обновляем поля
	product.Name = productData.Name
	product.Price = productData.Price
	product.Quantity = productData.Quantity
	product.Type = productData.Type
	product.Brand = productData.Brand
	product.Model = productData.Model
	product.Specifications = productData.Specifications
	product.Color = productData.Color
	product.Warranty = productData.Warranty

	saveResult := Database.Save(&product)
	if saveResult.Error != nil {
		return fmt.Errorf("error updating product: %v", saveResult.Error)
	}
	return nil
}

// DeleteProduct удаляет продукт по его ID
func DeleteProduct(id uint) error {
	result := Database.Delete(&models.Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting product: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", id)
	}
	return nil
}

// GetAllProducts извлекает все продукты из базы данных
func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	result := Database.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching products: %v", result.Error)
	}
	return products, nil
}

// DeleteAllProducts удаляет все продукты из базы данных
func DeleteAllProducts() error {
	result := Database.Exec("DELETE FROM products")
	if result.Error != nil {
		return fmt.Errorf("error deleting all products: %v", result.Error)
	}
	return nil
}

// GetProductsPaginated возвращает список продуктов из базы данных с поддержкой пагинации
func GetProductsPaginated(limit, offset int) ([]models.Product, error) {
	var products []models.Product
	result := Database.Limit(limit).Offset(offset).Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching products with pagination: %v", result.Error)
	}
	return products, nil
}
