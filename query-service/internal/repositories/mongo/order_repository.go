package mongo

import (
	"context"

	"github.com/JustBrowsing/query-service/internal/models"
	"github.com/JustBrowsing/query-service/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const orderCollection = "orders"

// OrderRepository is a MongoDB implementation of OrderRepository
type OrderRepository struct {
	db *mongo.Database
}

// NewOrderRepository creates a new MongoDB order repository
func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// collection returns the orders collection
func (r *OrderRepository) collection() *mongo.Collection {
	return r.db.Collection(orderCollection)
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order

	filter := bson.M{"orderId": id}
	err := r.collection().FindOne(ctx, filter).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("order", id)
		}
		return nil, errors.Wrap(err, "failed to get order")
	}

	return &order, nil
}

// Save saves an order
func (r *OrderRepository) Save(ctx context.Context, order *models.Order) error {
	_, err := r.collection().InsertOne(ctx, order)
	if err != nil {
		return errors.Wrap(err, "failed to save order")
	}
	return nil
}