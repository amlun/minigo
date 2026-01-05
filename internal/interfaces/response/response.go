package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "minigo/internal/domain/errors"
)

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	Items interface{} `json:"items"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// Ok 返回统一成功响应，HTTP 状态码固定 200
func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    "SUCCESS",
		Message: "成功",
		Data:    data,
	})
}

// OkWithMessage 返回自定义消息的成功响应
func OkWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    "SUCCESS",
		Message: message,
		Data:    data,
	})
}

// OkWithPage 返回分页数据的成功响应
func OkWithPage(c *gin.Context, items interface{}, total int, page, size int) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    "SUCCESS",
		Message: "成功",
		Data: PageData{
			Items: items,
			Total: total,
			Page:  page,
			Size:  size,
		},
	})
}

// Error 返回统一错误响应，HTTP 状态码，code 使用状态码映射业务码
// 入参 status 兼容旧用法（传入 HTTP 状态码），内部映射为业务码
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: false,
		Code:    statusToCode(status),
		Message: message,
	})
}

// ErrorCode 返回自定义业务码的错误响应（HTTP 始终 200）
func ErrorCode(c *gin.Context, code string, message string) {
	c.JSON(http.StatusOK, Response{
		Success: false,
		Code:    code,
		Message: message,
	})
}

// ErrorFromAppError 从AppError返回错误响应
func ErrorFromAppError(c *gin.Context, err *apperrors.AppError) {
	c.JSON(err.GetHTTPStatus(), Response{
		Success: false,
		Code:    err.Code,
		Message: err.Message,
	})
}

// HandleError 智能处理错误，自动识别错误类型
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 如果是AppError，使用专门的处理方法
	if appErr, ok := apperrors.AsAppError(err); ok {
		ErrorFromAppError(c, appErr)
		return
	}

	// 其他错误统一返回内部错误
	Error(c, http.StatusInternalServerError, err.Error())
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
