package api

import (
	"net/http"
	"time"

	"github.com/JustBrowsing/query-service/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Response represents a general API response
type Response struct {
	Timestamp time.Time   `json:"timestamp"`
	Status    int         `json:"status"`
	Error     string      `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Details   []string    `json:"details,omitempty"`
	Path      string      `json:"path,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(data interface{}) Response {
	return Response{
		Timestamp: time.Now().UTC(),
		Status:    http.StatusOK,
		Data:      data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error, path string) Response {
	var code string
	var details []string
	status := errors.Status(err)
	code = errors.Code(err)
	details = errors.Details(err)

	return Response{
		Timestamp: time.Now().UTC(),
		Status:    status,
		Error:     code,
		Message:   err.Error(),
		Details:   details,
		Path:      path,
	}
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, NewSuccessResponse(data))
}

// SendCreated sends a created response
func SendCreated(c *gin.Context, data interface{}) {
	response := NewSuccessResponse(data)
	c.JSON(http.StatusCreated, response)
}

// SendError sends an error response
func SendError(c *gin.Context, err error, logger *zap.Logger) {
	status := errors.Status(err)
	response := NewErrorResponse(err, c.Request.URL.Path)

	// Log error except for 404s (which are not really errors)
	if status != http.StatusNotFound {
		logger.Error("api error",
			zap.Int("status", status),
			zap.String("path", c.Request.URL.Path),
			zap.Error(err),
		)
	}

	c.JSON(status, response)
	c.Abort()
}