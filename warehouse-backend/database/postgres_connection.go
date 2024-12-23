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

func ConnectPostgres() {
	db := initializeDatabaseConnection()
	Database = db
	fmt.Println(dbConnectionSuccessMsg)

	if err := Database.AutoMigrate(&models.Product{}); err != nil {
		log.Fatalf("Error migrating the database: %v", err)
	}
	fmt.Println("Database migrated successfully!")
}

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

func CreateProduct(product *models.Product) error {
	result := Database.Create(product)
	if result.Error != nil {
		return fmt.Errorf("error creating product: %v", result.Error)
	}
	return nil
}

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

func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	result := Database.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("error fetching products: %v", result.Error)
	}
	return products, nil
}
