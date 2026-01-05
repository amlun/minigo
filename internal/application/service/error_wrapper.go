package service

import (
	"fmt"
)

// ErrorWrapper 错误包装器
type ErrorWrapper struct{}

// NewErrorWrapper 创建错误包装器实例
func NewErrorWrapper() *ErrorWrapper {
	return &ErrorWrapper{}
}

// WrapRepositoryError 包装仓储层错误
func (e *ErrorWrapper) WrapRepositoryError(err error, operation string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s失败: %w", operation, err)
}

// WrapValidationError 包装验证错误
func (e *ErrorWrapper) WrapValidationError(err error, field string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s验证失败: %w", field, err)
}

// WrapBusinessError 包装业务逻辑错误
func (e *ErrorWrapper) WrapBusinessError(err error, operation string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", operation, err)
}

// WrapTransactionError 包装事务错误
func (e *ErrorWrapper) WrapTransactionError(err error, operation string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("事务执行失败[%s]: %w", operation, err)
}

// WrapExternalServiceError 包装外部服务错误
func (e *ErrorWrapper) WrapExternalServiceError(err error, serviceName string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("外部服务[%s]调用失败: %w", serviceName, err)
}

// CheckAndWrapError 检查错误并包装
func (e *ErrorWrapper) CheckAndWrapError(err error, operation string) error {
	if err != nil {
		return e.WrapRepositoryError(err, operation)
	}
	return nil
}

// CheckEntityNotFound 检查实体是否存在
func (e *ErrorWrapper) CheckEntityNotFound(entity interface{}, entityName string) error {
	if entity == nil {
		return fmt.Errorf("%s不存在", entityName)
	}
	return nil
}

// CheckUpdateResult 检查更新结果
func (e *ErrorWrapper) CheckUpdateResult(rowsAffected int64, entityName string) error {
	if rowsAffected == 0 {
		return fmt.Errorf("%s不存在或未发生变更", entityName)
	}
	return nil
}

// CheckDeleteResult 检查删除结果
func (e *ErrorWrapper) CheckDeleteResult(rowsAffected int64, entityName string) error {
	if rowsAffected == 0 {
		return fmt.Errorf("%s不存在", entityName)
	}
	return nil
}
