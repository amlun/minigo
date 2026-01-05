package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"minigo/internal/infrastructure/logging"
)

// RequestLoggerMiddleware logs basic request info using logrus.
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		lat := time.Since(start)
		status := c.Writer.Status()
		logging.L().WithFields(map[string]interface{}{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"status": status,
			"lat_ms": float64(lat.Milliseconds()),
			"client": c.ClientIP(),
		}).Info("http_request")
	}
}
