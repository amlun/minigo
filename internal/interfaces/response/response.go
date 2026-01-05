package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ok 返回统一成功响应，HTTP 状态码固定 200
func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    "SUCCESS",
		"message": "成功",
		"data":    data,
	})
}

// Error 返回统一错误响应，HTTP 状态码，code 使用状态码映射业务码
// 入参 status 兼容旧用法（传入 HTTP 状态码），内部映射为业务码
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"code":    statusToCode(status),
		"message": message,
	})
}

// ErrorCode 返回自定义业务码的错误响应（HTTP 始终 200）
func ErrorCode(c *gin.Context, code string, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"code":    code,
		"message": message,
	})
}

// 将常见 HTTP 状态码映射为业务码（接口文档约定）
func statusToCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	default:
		return "INTERNAL_ERROR"
	}
}
