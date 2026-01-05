package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"minigo/internal/infrastructure/auth"
	resp "minigo/internal/interfaces/response"
)

// Context keys for auth
const (
	ContextUserIDKey   = "user_id"
	ContextUserRoleKey = "user_type"
)

// AuthMiddleware parses JWT and injects user info.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp.Error(c, http.StatusUnauthorized, "未提供认证令牌")
			c.Abort()
			return
		}
		if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
			authHeader = after
		}
		claims, err := auth.ParseToken(authHeader)
		if err != nil {
			resp.Error(c, http.StatusUnauthorized, "无效的认证令牌")
			c.Abort()
			return
		}
		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserRoleKey, claims.UserRole)
		c.Next()
	}
}

func GetUserIDFromContext(c *gin.Context) int64 {
	userIDVal, ok := c.Get(ContextUserIDKey)
	if !ok {
		return 0
	}
	userID, _ := userIDVal.(int64)
	return userID
}
