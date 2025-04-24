package services

import (
	"context"
	"encoding/json"

	"github.com/JustBrowsing/query-service/internal/models"
	"github.com/JustBrowsing/query-service/internal/repositories"
	"go.uber.org/zap"
)

// OrderService defines the operations for order management
type OrderService struct {
	logger     *zap.Logger
	orderRepo  repositories.OrderRepository
}

// NewOrderService creates a new order service
func NewOrderService(
	logger *zap.Logger,
	orderRepo repositories.OrderRepository,
) *OrderService {
	return &OrderService{
		logger:     logger,
		orderRepo:  orderRepo,
	}
}

// GetOrderByID retrieves an order by its ID
func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*models.Order, error) {
	return s.orderRepo.GetByID(ctx, id)
}

// HandleOrderCreated handles an order created event
func (s *OrderService) HandleOrderCreated(ctx context.Context, event *models.OrderEvent) error {
	s.logger.Info("handling order created event", 
		zap.String("orderId", event.OrderID),
		zap.String("orderNumber", event.OrderNumber))

	// Create a new order
	order := models.NewOrder(
		event.OrderID,
		event.OrderNumber,
		event.TotalAmount,
		event.Items,
		event.Timestamp,
	)

	// Save to MongoDB
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return err
	}

	s.logger.Debug("order saved successfully", 
		zap.String("orderId", event.OrderID),
		zap.String("orderNumber", event.OrderNumber))

	return nil
}