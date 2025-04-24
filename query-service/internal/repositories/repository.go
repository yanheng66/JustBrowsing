package repositories

import (
	"context"

	"github.com/JustBrowsing/query-service/internal/models"
)

// ProductRepository defines the operations for product data access
type ProductRepository interface {
	// GetByID retrieves a product by its ID
	GetByID(ctx context.Context, id string) (*models.Product, error)
	
	// SearchByTags searches for products that match all the specified tags
	SearchByTags(ctx context.Context, tags map[string]string) ([]*models.Product, int, error)
	
	// Save saves a product
	Save(ctx context.Context, product *models.Product) error
	
	// Update updates a product
	Update(ctx context.Context, product *models.Product) error
	
	// UpdateInventory updates the product's inventory
	UpdateInventory(ctx context.Context, productID string, quantity int) error
	
	// AddTag adds a tag to a product
	AddTag(ctx context.Context, productID string, tag models.ProductTag) error
	
	// RemoveTag removes a tag from a product
	RemoveTag(ctx context.Context, productID string, tagID string) error
}

// OrderRepository defines the operations for order data access
type OrderRepository interface {
	// GetByID retrieves an order by its ID
	GetByID(ctx context.Context, id string) (*models.Order, error)
	
	// Save saves an order
	Save(ctx context.Context, order *models.Order) error
}