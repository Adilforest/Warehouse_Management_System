package models

import "fmt"

const ProductTableName = "products"

// Product представляет товар
type Product struct {
	ID             int     `json:"id"`             // Уникальный ID
	Name           string  `json:"name"`           // Название товара
	Type           string  `json:"type"`           // Тип товара (ноутбук, телефон, наушники)
	Brand          string  `json:"brand"`          // Бренд (например, Apple)
	Model          string  `json:"model"`          // Модель (например, MacBook Air)
	Specifications string  `json:"specifications"` // Спецификации
	Color          string  `json:"color"`          // Цвет (например, черный)
	Price          float64 `json:"price"`          // Цена в USD
	Quantity       int     `json:"quantity"`       // Количество на складе
	Warranty       int     `json:"warranty"`       // Гарантия в месяцах
}

// TableName возвращает имя таблицы для модели Product
func (p *Product) TableName() string {
	return ProductTableName
}

// Debug формирует строку для отладки продукта
func (p *Product) Debug() string {
	return fmt.Sprintf("Product[ID=%d, Name='%s', Type='%s', Brand='%s', Model='%s', Price=%.2f, Quantity=%d]",
		p.ID, p.Name, p.Type, p.Brand, p.Model, p.Price, p.Quantity)
}
