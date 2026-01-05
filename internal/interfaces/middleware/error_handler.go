package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	apperrors "minigo/internal/domain/errors"
	"minigo/internal/infrastructure/logging"
	resp "minigo/internal/interfaces/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware catches panics and returns unified error response.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// 获取请求ID用于日志追踪
		requestID := c.GetString("request_id")

		// 记录panic详情和堆栈
		logging.L().WithFields(map[string]interface{}{
			"request_id": requestID,
			"panic":      recovered,
			"stack":      string(debug.Stack()),
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
		}).Error("panic_recovered")

		// 返回统一错误响应
		resp.Error(c, http.StatusInternalServerError, "系统内部错误")
		c.Abort()
	})
}

// HandleError processes errors and returns appropriate response
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	requestID := c.GetString("request_id")

	// Check if it's an AppError
	if appErr, ok := apperrors.AsAppError(err); ok {
		// 根据错误类型决定日志级别
		switch appErr.Type {
		case apperrors.SystemError:
			logging.L().WithFields(map[string]interface{}{
				"request_id": requestID,
				"error_code": appErr.Code,
				"error_type": appErr.Type,
			}).WithError(err).Error("system_error")
		case apperrors.BusinessError, apperrors.ValidationError:
			logging.L().WithFields(map[string]interface{}{
				"request_id": requestID,
				"error_code": appErr.Code,
				"error_type": appErr.Type,
			}).Warn("business_error")
		default:
			logging.L().WithFields(map[string]interface{}{
				"request_id": requestID,
				"error_code": appErr.Code,
				"error_type": appErr.Type,
			}).Info("app_error")
		}

		// 使用AppError的HTTP状态码和错误信息
		resp.ErrorFromAppError(c, appErr)
		return
	}

	// 处理未知错误
	logging.L().WithFields(map[string]interface{}{
		"request_id": requestID,
		"error":      err.Error(),
	}).Error("unexpected_error")

	// 返回通用内部错误
	resp.Error(c, http.StatusInternalServerError, "系统内部错误")
}

// AbortWithError 中止请求并返回错误
func AbortWithError(c *gin.Context, err error) {
	HandleError(c, err)
	c.Abort()
}

// AbortWithAppError 中止请求并返回AppError
func AbortWithAppError(c *gin.Context, appErr *apperrors.AppError) {
	HandleError(c, appErr)
	c.Abort()
}

// AbortWithBusinessError 中止请求并返回业务错误
func AbortWithBusinessError(c *gin.Context, code, message string) {
	err := apperrors.NewBusinessError(code, message)
	AbortWithAppError(c, err)
}

// MustNotError 如果有错误则中止请求
func MustNotError(c *gin.Context, err error) bool {
	if err != nil {
		AbortWithError(c, err)
		return false
	}
	return true
}

// RecoverToError 将panic转换为error
func RecoverToError() error {
	if r := recover(); r != nil {
		return fmt.Errorf("panic recovered: %v", r)
	}
	return nil
}
