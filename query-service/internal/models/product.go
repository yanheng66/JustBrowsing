package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductTag represents a product tag
type ProductTag struct {
	ID    string `bson:"id,omitempty" json:"id,omitempty"`
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}

// Product represents a product in the system
type Product struct {
	Base            `bson:",inline"`
	ProductID       string       `bson:"productId" json:"productId"`
	SKU             string       `bson:"sku" json:"sku"`
	Name            string       `bson:"name" json:"name"`
	Description     string       `bson:"description" json:"description"`
	Price           float64      `bson:"price" json:"price"`
	Tags            []ProductTag `bson:"tags" json:"tags"`
	CurrentInventory int          `bson:"currentInventory" json:"currentInventory"`
	Images          []string     `bson:"images" json:"images"`
}

// NewProduct creates a new product from a product created event
func NewProduct(productID, sku, name, description string, price float64, createdAt time.Time) *Product {
	return &Product{
		Base: Base{
			ID:        primitive.NewObjectID(),
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
		ProductID:       productID,
		SKU:             sku,
		Name:            name,
		Description:     description,
		Price:           price,
		Tags:            []ProductTag{},
		CurrentInventory: 0,
		Images:          []string{},
	}
}