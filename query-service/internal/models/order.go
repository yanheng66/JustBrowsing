package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID   string  `bson:"productId" json:"productId"`
	ProductName string  `bson:"productName" json:"productName"`
	SKU         string  `bson:"sku" json:"sku"`
	Quantity    int     `bson:"quantity" json:"quantity"`
	UnitPrice   float64 `bson:"unitPrice" json:"unitPrice"`
	TotalPrice  float64 `bson:"totalPrice" json:"totalPrice"`
}

// Order represents an order in the system
type Order struct {
	Base         `bson:",inline"`
	OrderID      string      `bson:"orderId" json:"orderId"`
	OrderNumber  string      `bson:"orderNumber" json:"orderNumber"`
	TotalAmount  float64     `bson:"totalAmount" json:"totalAmount"`
	Items        []OrderItem `bson:"items" json:"items"`
}

// NewOrder creates a new order from an order created event
func NewOrder(orderID, orderNumber string, totalAmount float64, items []OrderItem, createdAt time.Time) *Order {
	return &Order{
		Base: Base{
			ID:        primitive.NewObjectID(),
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
		OrderID:     orderID,
		OrderNumber: orderNumber,
		TotalAmount: totalAmount,
		Items:       items,
	}
}