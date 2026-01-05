package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	resp "minigo/internal/interfaces/response"
)

// RequireRoleMiddleware 要求特定角色权限
func RequireRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get(ContextUserRoleKey)
		if !ok {
			resp.Error(c, http.StatusUnauthorized, "未认证的请求")
			c.Abort()
			return
		}

		userType, ok := val.(string)
		if !ok {
			resp.Error(c, http.StatusUnauthorized, "用户类型无效")
			c.Abort()
			return
		}

		// 检查用户角色是否在允许的角色列表中
		for _, role := range allowedRoles {
			if userType == role {
				c.Next()
				return
			}
		}

		resp.Error(c, http.StatusForbidden, "权限不足")
		c.Abort()
	}
}
