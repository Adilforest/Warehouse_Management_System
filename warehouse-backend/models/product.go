package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product представляет товар
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`        // Уникальный ID от MongoDB
	Name        string             `bson:"name" json:"name"`               // Название товара
	Type        string             `bson:"type" json:"type"`               // Тип товара (ноутбук, телефон, наушники)
	Brand       string             `bson:"brand" json:"brand"`             // Бренд (например, Apple)
	Model       string             `bson:"model" json:"model"`             // Модель (например, MacBook Air)
	Description string             `bson:"description" json:"description"` // Общее описание товара
	Processor   string             `bson:"processor" json:"processor"`     // Процессор (например, Apple M1) — для ноутбуков и телефонов
	RAM         string             `bson:"ram" json:"ram"`                 // ОЗУ (например, 16 GB) — для ноутбуков и телефонов
	Storage     string             `bson:"storage" json:"storage"`         // Хранилище (например, 512 GB) — для ноутбуков и телефонов
	Color       string             `bson:"color" json:"color"`             // Цвет (например, черный)
	Price       float64            `bson:"price" json:"price"`             // Цена в USD
	Quantity    int                `bson:"quantity" json:"quantity"`       // Количество на складе
	Warranty    int                `bson:"warranty" json:"warranty"`       // Гарантия в месяцах
}

// Debug формирует строку для отладки продукта
func (p *Product) Debug() string {
	return fmt.Sprintf(
		"Product[ID=%s, Name='%s', Type='%s', Brand='%s', Model='%s', Description='%s', Processor='%s', RAM='%s', Storage='%s', Price=%.2f, Quantity=%d]",
		p.ID.Hex(), p.Name, p.Type, p.Brand, p.Model, p.Description, p.Processor, p.RAM, p.Storage, p.Price, p.Quantity,
	)
}
