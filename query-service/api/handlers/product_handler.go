package handlers

import (
	"strings"

	"github.com/JustBrowsing/query-service/api"
	"github.com/JustBrowsing/query-service/internal/services"
	"github.com/JustBrowsing/query-service/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ProductHandler handles product-related requests
type ProductHandler struct {
	logger         *zap.Logger
	productService *services.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(logger *zap.Logger, productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		logger:         logger,
		productService: productService,
	}
}

// GetProduct handles get product requests
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("productId")
	if productID == "" {
		api.SendError(c, errors.BadRequest("Product ID is required"), h.logger)
		return
	}

	product, err := h.productService.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		api.SendError(c, err, h.logger)
		return
	}

	api.SendSuccess(c, product)
}

// SearchProducts handles product search requests
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	// Get tags parameter
	tags := c.Query("tags")
	if tags == "" {
		api.SendError(c, errors.BadRequest("Tags parameter is required"), h.logger)
		return
	}

	// Parse tags parameter
	tagMap, err := parseTagsParameter(tags)
	if err != nil {
		api.SendError(c, err, h.logger)
		return
	}

	// Search products
	products, total, err := h.productService.SearchProductsByTags(c.Request.Context(), tags)
	if err != nil {
		api.SendError(c, err, h.logger)
		return
	}

	// Build response
	response := gin.H{
		"items": products,
		"total": total,
	}

	api.SendSuccess(c, response)
}

// parseTagsParameter parses the tags parameter from the query string
// Format is "tagName1:tagValue1,tagName2:tagValue2"
func parseTagsParameter(tagsParam string) (map[string]string, error) {
	if tagsParam == "" {
		return nil, errors.BadRequest("Tags parameter cannot be empty")
	}

	result := make(map[string]string)
	tagPairs := strings.Split(tagsParam, ",")
	
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