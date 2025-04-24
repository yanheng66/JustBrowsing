package api

import (
	"github.com/JustBrowsing/query-service/api/handlers"
	"github.com/JustBrowsing/query-service/api/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Router sets up the API routes
func Router(
	logger *zap.Logger,
	basePath string,
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
	healthHandler *handlers.HealthHandler,
) *gin.Engine {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Apply middlewares
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", healthHandler.Handle)

	// API routes
	api := router.Group(basePath)
	{
		// Product routes
		products := api.Group("/products")
		{
			products.GET("/:productId", productHandler.GetProduct)
			products.GET("/search", productHandler.SearchProducts)
		}

		// Order routes
		orders := api.Group("/orders")
		{
			orders.GET("/:orderId", orderHandler.GetOrder)
		}
	}

	return router
}