package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/JustBrowsing/query-service/internal/models"
	"github.com/JustBrowsing/query-service/internal/repositories"
	"github.com/JustBrowsing/query-service/pkg/errors"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

// ProductService defines the operations for product management
type ProductService struct {
	logger         *zap.Logger
	productRepo    repositories.ProductRepository
	redisClient    *redis.Client
	elasticClient  *elastic.Client
	cacheTTL       time.Duration
	elasticIndex   string
}

// NewProductService creates a new product service
func NewProductService(
	logger *zap.Logger,
	productRepo repositories.ProductRepository,
	redisClient *redis.Client,
	elasticClient *elastic.Client,
	cacheTTL time.Duration,
	elasticIndex string,
) *ProductService {
	return &ProductService{
		logger:         logger,
		productRepo:    productRepo,
		redisClient:    redisClient,
		elasticClient:  elasticClient,
		cacheTTL:       cacheTTL,
		elasticIndex:   elasticIndex,
	}
}

// GetProductByID retrieves a product by its ID
func (s *ProductService) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("product:%s", id)
	cachedData, err := s.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		// Cache hit
		var product models.Product
		if err := json.Unmarshal(cachedData, &product); err == nil {
			s.logger.Debug("cache hit for product", zap.String("id", id))
			return &product, nil
		}
		// If we can't unmarshal, just log and continue to fetch from DB
		s.logger.Warn("failed to unmarshal cached product", zap.String("id", id), zap.Error(err))
	} else if err != redis.Nil {
		// If error is not "key not found", log it
		s.logger.Warn("cache error", zap.Error(err))
	}

	// Cache miss or error, fetch from database
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache for next time
	if productJSON, err := json.Marshal(product); err == nil {
		if err := s.redisClient.Set(ctx, cacheKey, productJSON, s.cacheTTL).Err(); err != nil {
			s.logger.Warn("failed to cache product", zap.String("id", id), zap.Error(err))
		}
	}

	return product, nil
}

// SearchProductsByTags searches for products by tags
func (s *ProductService) SearchProductsByTags(ctx context.Context, tagParams string) ([]*models.Product, int, error) {
	// Parse tag parameters
	tags, err := parseTags(tagParams)
	if err != nil {
		return nil, 0, err
	}

	if len(tags) == 0 {
		return []*models.Product{}, 0, nil
	}

	// Try to get from cache first
	cacheKey := fmt.Sprintf("products:tags:%s", tagParams)
	cachedData, err := s.redisClient.Get(ctx, cacheKey).Bytes()
	if err == nil {
		// Cache hit
		var cachedResult struct {
			Products []*models.Product `json:"products"`
			Total    int               `json:"total"`
		}
		if err := json.Unmarshal(cachedData, &cachedResult); err == nil {
			s.logger.Debug("cache hit for products search", zap.String("tags", tagParams))
			return cachedResult.Products, cachedResult.Total, nil
		}
		// If we can't unmarshal, just log and continue to search
		s.logger.Warn("failed to unmarshal cached search result", zap.String("tags", tagParams), zap.Error(err))
	} else if err != redis.Nil {
		// If error is not "key not found", log it
		s.logger.Warn("cache error", zap.Error(err))
	}

	// Cache miss or error, search in Elasticsearch first
	var products []*models.Product
	var total int

	// Build Elasticsearch query
	var musts []elastic.Query
	for name, value := range tags {
		musts = append(musts, elastic.NewNestedQuery("tags", 
			elastic.NewBoolQuery().Must(
				elastic.NewTermQuery("tags.name", name),
				elastic.NewTermQuery("tags.value", value),
			),
		))
	}

	boolQuery := elastic.NewBoolQuery().Must(musts...)
	
	// Search in Elasticsearch
	searchResult, err := s.elasticClient.Search().
		Index(s.elasticIndex).
		Query(boolQuery).
		From(0).
		Size(100). // Limit results
		Do(ctx)
	
	if err != nil {
		s.logger.Warn("elasticsearch search failed, falling back to MongoDB", zap.Error(err))
		// Fall back to MongoDB
		products, total, err = s.productRepo.SearchByTags(ctx, tags)
		if err != nil {
			return nil, 0, err
		}
	} else {
		// Process Elasticsearch results
		total = int(searchResult.TotalHits())
		for _, hit := range searchResult.Hits.Hits {
			var product models.Product
			if err := json.Unmarshal(hit.Source, &product); err != nil {
				s.logger.Warn("failed to unmarshal elasticsearch hit", zap.Error(err))
				continue
			}
			products = append(products, &product)
		}
	}

	// Store in cache for next time
	cacheResult := struct {
		Products []*models.Product `json:"products"`
		Total    int               `json:"total"`
	}{
		Products: products,
		Total:    total,
	}

	if resultJSON, err := json.Marshal(cacheResult); err == nil {
		if err := s.redisClient.Set(ctx, cacheKey, resultJSON, s.cacheTTL).Err(); err != nil {
			s.logger.Warn("failed to cache search result", zap.String("tags", tagParams), zap.Error(err))
		}
	}

	return products, total, nil
}

// HandleProductCreated handles a product created event
func (s *ProductService) HandleProductCreated(ctx context.Context, event *models.ProductEvent) error {
	s.logger.Info("handling product created event", zap.String("productId", event.ProductID))

	// Create a new product
	product := models.NewProduct(
		event.ProductID,
		event.SKU,
		event.Name,
		event.Description,
		event.Price,
		event.Timestamp,
	)

	// Add tags if any
	for _, tag := range event.Tags {
		product.Tags = append(product.Tags, tag)
	}

	// Save to MongoDB
	if err := s.productRepo.Save(ctx, product); err != nil {
		return err
	}

	// Index in Elasticsearch
	if _, err := s.elasticClient.Index().
		Index(s.elasticIndex).
		Id(event.ProductID).
		BodyJson(product).
		Do(ctx); err != nil {
		s.logger.Error("failed to index product in elasticsearch", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%s", event.ProductID)
	if err := s.redisClient.Del(ctx, cacheKey).Err(); err != nil {
		s.logger.Warn("failed to invalidate product cache", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	}

	return nil
}

// HandleProductUpdated handles a product updated event
func (s *ProductService) HandleProductUpdated(ctx context.Context, event *models.ProductEvent) error {
	s.logger.Info("handling product updated event", zap.String("productId", event.ProductID))

	// Get existing product
	product, err := s.productRepo.GetByID(ctx, event.ProductID)
	if err != nil {
		return err
	}

	// Update product fields
	product.Name = event.Name
	product.Description = event.Description
	product.Price = event.Price
	product.UpdatedAt = event.Timestamp

	// Update in MongoDB
	if err := s.productRepo.Update(ctx, product); err != nil {
		return err
	}

	// Update in Elasticsearch
	if _, err := s.elasticClient.Index().
		Index(s.elasticIndex).
		Id(event.ProductID).
		BodyJson(product).
		Do(ctx); err != nil {
		s.logger.Error("failed to update product in elasticsearch", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("product:%s", event.ProductID)
	if err := s.redisClient.Del(ctx, cacheKey).Err(); err != nil {
		s.logger.Warn("failed to invalidate product cache", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	}

	return nil
}

// HandleProductTagAdded handles a product tag added event
func (s *ProductService) HandleProductTagAdded(ctx context.Context, event *models.ProductTagEvent) error {
	s.logger.Info("handling product tag added event", 
		zap.String("productId", event.ProductID),
		zap.String("tagId", event.TagID))

	// Create tag
	tag := models.ProductTag{
		ID:    event.TagID,
		Name:  event.TagName,
		Value: event.TagValue,
	}

	// Add tag to product
	if err := s.productRepo.AddTag(ctx, event.ProductID, tag); err != nil {
		return err
	}

	// Get updated product for Elasticsearch
	product, err := s.productRepo.GetByID(ctx, event.ProductID)
	if err != nil {
		s.logger.Error("failed to get product for elasticsearch update", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	} else {
		// Update in Elasticsearch
		if _, err := s.elasticClient.Index().
			Index(s.elasticIndex).
			Id(event.ProductID).
			BodyJson(product).
			Do(ctx); err != nil {
			s.logger.Error("failed to update product in elasticsearch", 
				zap.String("productId", event.ProductID), 
				zap.Error(err))
		}
	}

	// Invalidate product cache
	cacheKey := fmt.Sprintf("product:%s", event.ProductID)
	if err := s.redisClient.Del(ctx, cacheKey).Err(); err != nil {
		s.logger.Warn("failed to invalidate product cache", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	}

	// Invalidate tag search caches
	pattern := "products:tags:*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		s.logger.Warn("failed to get tag search cache keys", zap.Error(err))
	} else if len(keys) > 0 {
		if err := s.redisClient.Del(ctx, keys...).Err(); err != nil {
			s.logger.Warn("failed to invalidate tag search caches", zap.Error(err))
		}
	}

	return nil
}

// HandleProductTagRemoved handles a product tag removed event
func (s *ProductService) HandleProductTagRemoved(ctx context.Context, event *models.ProductTagEvent) error {
	s.logger.Info("handling product tag removed event", 
		zap.String("productId", event.ProductID),
		zap.String("tagId", event.TagID))

	// Remove tag from product
	if err := s.productRepo.RemoveTag(ctx, event.ProductID, event.TagID); err != nil {
		return err
	}

	// Get updated product for Elasticsearch
	product, err := s.productRepo.GetByID(ctx, event.ProductID)
	if err != nil {
		s.logger.Error("failed to get product for elasticsearch update", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	} else {
		// Update in Elasticsearch
		if _, err := s.elasticClient.Index().
			Index(s.elasticIndex).
			Id(event.ProductID).
			BodyJson(product).
			Do(ctx); err != nil {
			s.logger.Error("failed to update product in elasticsearch", 
				zap.String("productId", event.ProductID), 
				zap.Error(err))
		}
	}

	// Invalidate product cache
	cacheKey := fmt.Sprintf("product:%s", event.ProductID)
	if err := s.redisClient.Del(ctx, cacheKey).Err(); err != nil {
		s.logger.Warn("failed to invalidate product cache", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	}

	// Invalidate tag search caches
	pattern := "products:tags:*"
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		s.logger.Warn("failed to get tag search cache keys", zap.Error(err))
	} else if len(keys) > 0 {
		if err := s.redisClient.Del(ctx, keys...).Err(); err != nil {
			s.logger.Warn("failed to invalidate tag search caches", zap.Error(err))
		}
	}

	return nil
}

// HandleInventoryUpdated handles an inventory updated event
func (s *ProductService) HandleInventoryUpdated(ctx context.Context, event *models.InventoryEvent) error {
	s.logger.Info("handling inventory updated event", 
		zap.String("productId", event.ProductID),
		zap.Int("quantity", event.Quantity))

	// Update inventory in MongoDB
	if err := s.productRepo.UpdateInventory(ctx, event.ProductID, event.Quantity); err != nil {
		return err
	}

	// Get updated product for Elasticsearch
	product, err := s.productRepo.GetByID(ctx, event.ProductID)
	if err != nil {
		s.logger.Error("failed to get product for elasticsearch update", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	} else {
		// Update in Elasticsearch
		if _, err := s.elasticClient.Index().
			Index(s.elasticIndex).
			Id(event.ProductID).
			BodyJson(product).
			Do(ctx); err != nil {
			s.logger.Error("failed to update product in elasticsearch", 
				zap.String("productId", event.ProductID), 
				zap.Error(err))
		}
	}

	// Invalidate product cache
	cacheKey := fmt.Sprintf("product:%s", event.ProductID)
	if err := s.redisClient.Del(ctx, cacheKey).Err(); err != nil {
		s.logger.Warn("failed to invalidate product cache", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	}

	// Set inventory cache with a shorter TTL (5 minutes)
	inventoryCacheKey := fmt.Sprintf("inventory:%s", event.ProductID)
	if err := s.redisClient.Set(ctx, inventoryCacheKey, event.Quantity, 5*time.Minute).Err(); err != nil {
		s.logger.Warn("failed to cache inventory", 
			zap.String("productId", event.ProductID), 
			zap.Error(err))
	}

	return nil
}

// Helper function to parse tag parameters from the query string
// Format is "tagName1:tagValue1,tagName2:tagValue2"
func parseTags(tagParams string) (map[string]string, error) {
	if tagParams == "" {
		return map[string]string{}, nil
	}

	result := make(map[string]string)
	tagPairs := strings.Split(tagParams, ",")
	
	for _, pair := range tagPairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			return nil, errors.BadRequest("Invalid tag format. Expected format: tagName:tagValue")
		}
		
		name := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		
		if name == "" || value == "" {
			return nil, errors.BadRequest("Tag name and value cannot be empty")
		}
		
		result[name] = value
	}
	
	return result, nil
}