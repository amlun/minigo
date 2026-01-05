package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader 请求ID的Header名称
	RequestIDHeader = "X-Request-ID"
	// RequestIDKey 上下文中存储请求ID的key
	RequestIDKey = "request_id"
)

// RequestIDMiddleware 为每个请求生成或传递请求ID
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取请求ID
		requestID := c.GetHeader(RequestIDHeader)

		// 如果没有，则生成新的请求ID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将请求ID存入上下文
		c.Set(RequestIDKey, requestID)

		// 将请求ID写入响应头
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID 从上下文获取请求ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
