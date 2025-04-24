package models

import (
	"time"
)

// Event represents an event from Kafka
type Event struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Payload   []byte    `json:"payload"`
}

// ProductEvent represents a product event
type ProductEvent struct {
	ProductID   string       `json:"productId"`
	SKU         string       `json:"sku"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Price       float64      `json:"price"`
	Tags        []ProductTag `json:"tags"`
	Timestamp   time.Time    `json:"timestamp"`
}

// ProductTagEvent represents a product tag event
type ProductTagEvent struct {
	ProductID string    `json:"productId"`
	TagID     string    `json:"tagId"`
	TagName   string    `json:"tagName"`
	TagValue  string    `json:"tagValue"`
	Timestamp time.Time `json:"timestamp"`
}

// InventoryEvent represents an inventory event
type InventoryEvent struct {
	InventoryID string    `json:"inventoryId"`
	ProductID   string    `json:"productId"`
	Quantity    int       `json:"quantity"`
	Version     int       `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
}

// OrderEvent represents an order event
type OrderEvent struct {
	OrderID     string      `json:"orderId"`
	OrderNumber string      `json:"orderNumber"`
	TotalAmount float64     `json:"totalAmount"`
	Items       []OrderItem `json:"items"`
	Timestamp   time.Time   `json:"timestamp"`
}