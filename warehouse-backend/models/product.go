package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Type        string             `bson:"type" json:"type"`
	Brand       string             `bson:"brand" json:"brand"`
	Model       string             `bson:"model" json:"model"`
	Description string             `bson:"description" json:"description"`
	Processor   string             `bson:"processor" json:"processor"`
	RAM         string             `bson:"ram" json:"ram"`
	Storage     string             `bson:"storage" json:"storage"`
	Color       string             `bson:"color" json:"color"`
	Price       float64            `bson:"price" json:"price"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Warranty    int                `bson:"warranty" json:"warranty"`
	Link        string             `bson:"link" json:"link"`
}

func (p *Product) Debug() string {
	return fmt.Sprintf(
		"Product[ID=%s, Name='%s', Type='%s', Brand='%s', Model='%s', Description='%s', Processor='%s', RAM='%s', Storage='%s', Price=%.2f, Quantity=%d]",
		p.ID.Hex(), p.Name, p.Type, p.Brand, p.Model, p.Description, p.Processor, p.RAM, p.Storage, p.Price, p.Quantity,
	)
}
