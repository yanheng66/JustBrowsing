package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger returns a middleware that logs request details
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		
		// Process request
		c.Next()
		
		// Post-processing
		end := time.Now()
		latency := end.Sub(start)
		
		if len(c.Errors) > 0 {
			// Log errors
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		} else {
			// Log request details
			logger.Info("request completed",
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
			)
		}
	}
}