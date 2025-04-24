package handlers

import (
	"github.com/JustBrowsing/query-service/api"
	"github.com/JustBrowsing/query-service/internal/services"
	"github.com/JustBrowsing/query-service/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OrderHandler handles order-related requests
type OrderHandler struct {
	logger       *zap.Logger
	orderService *services.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(logger *zap.Logger, orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		logger:       logger,
		orderService: orderService,
	}
}

// GetOrder handles get order requests
func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("orderId")
	if orderID == "" {
		api.SendError(c, errors.BadRequest("Order ID is required"), h.logger)
		return
	}

	order, err := h.orderService.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		api.SendError(c, err, h.logger)
		return
	}

	api.SendSuccess(c, order)
}