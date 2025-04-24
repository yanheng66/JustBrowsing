package mongo

import (
	"context"
	"time"

	"github.com/JustBrowsing/query-service/internal/models"
	"github.com/JustBrowsing/query-service/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const productCollection = "products"

// ProductRepository is a MongoDB implementation of ProductRepository
type ProductRepository struct {
	db *mongo.Database
}

// NewProductRepository creates a new MongoDB product repository
func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// collection returns the products collection
func (r *ProductRepository) collection() *mongo.Collection {
	return r.db.Collection(productCollection)
}

// GetByID retrieves a product by its ID
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	var product models.Product

	filter := bson.M{"productId": id}
	err := r.collection().FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.NotFound("product", id)
		}
		return nil, errors.Wrap(err, "failed to get product")
	}

	return &product, nil
}

// SearchByTags searches for products that match all the specified tags
func (r *ProductRepository) SearchByTags(ctx context.Context, tags map[string]string) ([]*models.Product, int, error) {
	if len(tags) == 0 {
		return []*models.Product{}, 0, nil
	}

	// Build a filter for each tag
	var filters []bson.M
	for name, value := range tags {
		filters = append(filters, bson.M{
			"tags": bson.M{
				"$elemMatch": bson.M{
					"name":  name,
					"value": value,
				},
			},
		})
	}

	// Combine all filters with $and to require matching all tags
	filter := bson.M{
		"$and": filters,
	}

	// Count total matching products
	total, err := r.collection().CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to count products")
	}

	// Find matching products
	cursor, err := r.collection().Find(ctx, filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to search products")
	}
	defer cursor.Close(ctx)

	var products []*models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, 0, errors.Wrap(err, "failed to decode products")
	}

	return products, int(total), nil
}

// Save saves a product
func (r *ProductRepository) Save(ctx context.Context, product *models.Product) error {
	_, err := r.collection().InsertOne(ctx, product)
	if err != nil {
		return errors.Wrap(err, "failed to save product")
	}
	return nil
}

// Update updates a product
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	product.UpdatedAt = time.Now()

	filter := bson.M{"productId": product.ProductID}
	update := bson.M{"$set": product}

	_, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "failed to update product")
	}
	return nil
}

// UpdateInventory updates the product's inventory
func (r *ProductRepository) UpdateInventory(ctx context.Context, productID string, quantity int) error {
	filter := bson.M{"productId": productID}
	update := bson.M{
		"$set": bson.M{
			"currentInventory": quantity,
			"updated":         time.Now(),
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "failed to update inventory")
	}

	if result.MatchedCount == 0 {
		return errors.NotFound("product", productID)
	}

	return nil
}

// AddTag adds a tag to a product
func (r *ProductRepository) AddTag(ctx context.Context, productID string, tag models.ProductTag) error {
	if tag.ID == "" {
		tag.ID = primitive.NewObjectID().Hex()
	}

	filter := bson.M{"productId": productID}
	update := bson.M{
		"$push": bson.M{"tags": tag},
		"$set":  bson.M{"updated": time.Now()},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "failed to add tag")
	}

	if result.MatchedCount == 0 {
		return errors.NotFound("product", productID)
	}

	return nil
}

// RemoveTag removes a tag from a product
func (r *ProductRepository) RemoveTag(ctx context.Context, productID string, tagID string) error {
	filter := bson.M{"productId": productID}
	update := bson.M{
		"$pull": bson.M{"tags": bson.M{"id": tagID}},
		"$set":  bson.M{"updated": time.Now()},
	}

	result, err := r.collection().UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrap(err, "failed to remove tag")
	}

	if result.MatchedCount == 0 {
		return errors.NotFound("product", productID)
	}

	return nil
}