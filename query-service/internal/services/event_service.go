package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/JustBrowsing/query-service/config"
	"github.com/JustBrowsing/query-service/internal/models"
	"github.com/JustBrowsing/query-service/pkg/kafka"
	"go.uber.org/zap"
)

// EventService handles incoming events from Kafka
type EventService struct {
	logger        *zap.Logger
	productService *ProductService
	orderService   *OrderService
	topics        config.KafkaTopicConfig
}

// NewEventService creates a new event service
func NewEventService(
	logger *zap.Logger,
	productService *ProductService,
	orderService *OrderService,
	topics config.KafkaTopicConfig,
) *EventService {
	return &EventService{
		logger:        logger,
		productService: productService,
		orderService:   orderService,
		topics:        topics,
	}
}

// HandleMessage handles an incoming Kafka message
func (s *EventService) HandleMessage(topic string, key []byte, value []byte, timestamp time.Time) error {
	s.logger.Debug("handling message", 
		zap.String("topic", topic), 
		zap.String("key", string(key)),
		zap.Time("timestamp", timestamp))

	// Handle different topics
	switch topic {
	case s.topics.Product:
		return s.handleProductEvent(key, value, timestamp)
	case s.topics.Inventory:
		return s.handleInventoryEvent(key, value, timestamp)
	case s.topics.Order:
		return s.handleOrderEvent(key, value, timestamp)
	default:
		s.logger.Warn("unknown topic", zap.String("topic", topic))
	}

	return nil
}

// handleProductEvent handles product events
func (s *EventService) handleProductEvent(key []byte, value []byte, timestamp time.Time) error {
	// Parse event type from key
	eventType := string(key)
	ctx := context.Background()

	switch eventType {
	case "ProductCreated":
		var event models.ProductEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return err
		}
		event.Timestamp = timestamp
		return s.productService.HandleProductCreated(ctx, &event)

	case "ProductUpdated":
		var event models.ProductEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return err
		}
		event.Timestamp = timestamp
		return s.productService.HandleProductUpdated(ctx, &event)

	case "ProductTagAdded":
		var event models.ProductTagEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return err
		}
		event.Timestamp = timestamp
		return s.productService.HandleProductTagAdded(ctx, &event)

	case "ProductTagRemoved":
		var event models.ProductTagEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return err
		}
		event.Timestamp = timestamp
		return s.productService.HandleProductTagRemoved(ctx, &event)

	default:
		s.logger.Warn("unknown product event type", zap.String("type", eventType))
	}

	return nil
}

// handleInventoryEvent handles inventory events
func (s *EventService) handleInventoryEvent(key []byte, value []byte, timestamp time.Time) error {
	// Parse event type from key
	eventType := string(key)
	ctx := context.Background()

	switch eventType {
	case "InventoryUpdated":
		var event models.InventoryEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return err
		}
		event.Timestamp = timestamp
		return s.productService.HandleInventoryUpdated(ctx, &event)

	default:
		s.logger.Warn("unknown inventory event type", zap.String("type", eventType))
	}

	return nil
}

// handleOrderEvent handles order events
func (s *EventService) handleOrderEvent(key []byte, value []byte, timestamp time.Time) error {
	// Parse event type from key
	eventType := string(key)
	ctx := context.Background()

	switch eventType {
	case "OrderCreated":
		var event models.OrderEvent
		if err := json.Unmarshal(value, &event); err != nil {
			return err
		}
		event.Timestamp = timestamp
		return s.orderService.HandleOrderCreated(ctx, &event)

	default:
		s.logger.Warn("unknown order event type", zap.String("type", eventType))
	}

	return nil
}

// StartConsumer starts the Kafka consumer
func (s *EventService) StartConsumer(ctx context.Context, cfg config.KafkaConfig) error {
	// Create a new consumer
	consumer, err := kafka.NewConsumer(cfg, s.logger, s.HandleMessage)
	if err != nil {
		return err
	}

	// Define topics to consume
	topics := []string{
		s.topics.Product,
		s.topics.Inventory,
		s.topics.Order,
	}

	// Start consuming
	return consumer.Consume(ctx, topics)
}