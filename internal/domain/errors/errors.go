package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorType 定义错误类型
type ErrorType string

const (
	// SystemError 系统错误（数据库连接失败、网络错误等）
	SystemError ErrorType = "SYSTEM_ERROR"
	// BusinessError 业务错误（用户不存在、余额不足等）
	BusinessError ErrorType = "BUSINESS_ERROR"
	// ValidationError 验证错误（参数格式错误等）
	ValidationError ErrorType = "VALIDATION_ERROR"
	// AuthError 认证错误（未登录、权限不足等）
	AuthError ErrorType = "AUTH_ERROR"
	// NotFoundError 资源不存在错误
	NotFoundError ErrorType = "NOT_FOUND_ERROR"
)

// AppError 应用错误结构
type AppError struct {
	Type    ErrorType `json:"type"`
	Code    string    `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Cause   error     `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 支持errors.Unwrap
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Is 支持errors.Is
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return false
}

// GetHTTPStatus 根据错误类型返回HTTP状态码
func (e *AppError) GetHTTPStatus() int {
	switch e.Type {
	case ValidationError:
		return http.StatusBadRequest
	case AuthError:
		return http.StatusUnauthorized
	case NotFoundError:
		return http.StatusNotFound
	case BusinessError:
		return http.StatusBadRequest
	case SystemError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// 构造函数

// NewSystemError 创建系统错误
func NewSystemError(code, message string, cause error) *AppError {
	return &AppError{
		Type:    SystemError,
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewBusinessError 创建业务错误
func NewBusinessError(code, message string) *AppError {
	return &AppError{
		Type:    BusinessError,
		Code:    code,
		Message: message,
	}
}

// NewValidationError 创建验证错误
func NewValidationError(code, message string) *AppError {
	return &AppError{
		Type:    ValidationError,
		Code:    code,
		Message: message,
	}
}

// NewAuthError 创建认证错误
func NewAuthError(code, message string) *AppError {
	return &AppError{
		Type:    AuthError,
		Code:    code,
		Message: message,
	}
}

// NewNotFoundError 创建资源不存在错误
func NewNotFoundError(code, message string) *AppError {
	return &AppError{
		Type:    NotFoundError,
		Code:    code,
		Message: message,
	}
}

// 预定义的通用错误

var (
	/* ---系统错误--- */

	ErrDatabase = NewSystemError("SYS_001", "数据库操作失败", nil)
	ErrNetwork  = NewSystemError("SYS_002", "网络连接失败", nil)
	ErrInternal = NewSystemError("SYS_003", "内部服务错误", nil)

	/* ---通用业务错误--- */

	ErrResourceNotFound  = NewNotFoundError("BIZ_001", "资源不存在")
	ErrInvalidOperation  = NewBusinessError("BIZ_002", "操作无效")
	ErrDuplicateResource = NewBusinessError("BIZ_003", "资源已存在")

	/* ---验证错误--- */

	ErrInvalidParams = NewValidationError("VAL_001", "参数验证失败")
	ErrInvalidFormat = NewValidationError("VAL_002", "格式错误")

	/* ---认证错误--- */
	ErrUnauthorized = NewAuthError("AUTH_001", "未授权访问")
	ErrForbidden    = NewAuthError("AUTH_002", "权限不足")
	ErrTokenExpired = NewAuthError("AUTH_003", "令牌已过期")
)

// 工具函数

// IsAppError 检查是否为AppError
func IsAppError(err error) bool {
	var appError *AppError
	ok := errors.As(err, &appError)
	return ok
}

// AsAppError 转换为AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// WrapSystemError 包装系统错误
func WrapSystemError(err error, code, message string) *AppError {
	return &AppError{
		Type:    SystemError,
		Code:    code,
		Message: message,
		Cause:   err,
	}
}

// FromStandardError 从标准错误转换
func FromStandardError(err error) *AppError {
	if err == nil {
		return nil
	}

	// 如果已经是AppError，直接返回
	if appErr, ok := AsAppError(err); ok {
		return appErr
	}

	// 否则包装为系统错误
	return WrapSystemError(err, "SYS_999", "未知系统错误")
}
