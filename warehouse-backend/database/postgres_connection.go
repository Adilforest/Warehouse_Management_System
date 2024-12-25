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

// ConnectPostgres establishes a connection to the PostgreSQL database and performs AutoMigrations
func ConnectPostgres() {
	db := initializeDatabaseConnection()
	Database = db
	fmt.Println(dbConnectionSuccessMsg)

	// AutoMigrate the Product model
	if err := Database.AutoMigrate(&models.Product{}); err != nil {
		log.Fatalf("Error migrating the database: %v", err)
	}
	fmt.Println("Database migrated successfully!")
}

// initializeDatabaseConnection initializes the database connection
func initializeDatabaseConnection() *gorm.DB {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg := config.GetConfig()   // Config struct is defined in your existing code
	dsn := cfg.GetPostgresDSN() // Generate DSN from loaded environment variables
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database. DSN: %v, Error: %v", dsn, err)
	}
	return db
}

// CreateProduct creates a new product in the database
func CreateProduct(product *models.Product) error {
	result := Database.Create(product)
	if result.Error != nil {
		return fmt.Errorf("error creating product: %v", result.Error)
	}
	return nil
}

// GetProductByID retrieves a product by its ID from the database
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

// UpdateProduct updates an existing product by its ID
func UpdateProduct(id uint, productData *models.Product) error {
	var product models.Product
	result := Database.First(&product, id)
	if result.Error != nil {
		if result.RowsAffected == 0 {
			return fmt.Errorf("product with ID %d not found", id)
		}
		return fmt.Errorf("error finding product: %v", result.Error)
	}

	product.Name = productData.Name
	product.Price = productData.Price
	product.Quantity = productData.Quantity

	saveResult := Database.Save(&product)
	if saveResult.Error != nil {
		return fmt.Errorf("error updating product: %v", saveResult.Error)
	}
	return nil
}

// DeleteProduct deletes a product by its ID
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

// GetAllProducts retrieves all products from the database
func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	result := Database.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching products: %v", result.Error)
	}
	return products, nil
}

// DeleteAllProducts deletes all products from the database
func DeleteAllProducts() error {
	result := Database.Exec("DELETE FROM products")
	if result.Error != nil {
		fmt.Printf("SQL Error: %v\n", result.Error) // Логируем ошибку в SQL-запросе
		return fmt.Errorf("error deleting all products: %v", result.Error)
	}
	fmt.Printf("Rows affected: %d\n", result.RowsAffected) // Логируем количество удаленных строк
	return nil
}
