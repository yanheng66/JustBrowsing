package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/JustBrowsing/query-service/api"
	"github.com/JustBrowsing/query-service/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery returns a middleware that recovers from panics
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log stack trace
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
				)

				// Respond with 500
				response := api.Response{
					Status:  http.StatusInternalServerError,
					Error:   "INTERNAL_ERROR",
					Message: "An unexpected error occurred",
					Path:    c.Request.URL.Path,
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			}
		}()

		c.Next()
	}
}