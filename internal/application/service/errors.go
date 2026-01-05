package service

import (
	apperrors "minigo/internal/domain/errors"
)

// 预定义一些错误
// 用户相关错误
var (
	ErrUserNotFound         = apperrors.NewNotFoundError("USER_001", "用户不存在")
	ErrUserExists           = apperrors.NewBusinessError("USER_002", "用户已存在")
	ErrInvalidCredentials   = apperrors.NewBusinessError("USER_003", "用户名或密码错误")
	ErrUserDisabled         = apperrors.NewBusinessError("USER_004", "用户已被禁用")
	ErrInvalidInviteCode    = apperrors.NewBusinessError("USER_005", "邀请码无效")
	ErrSelfReferral         = apperrors.NewBusinessError("USER_006", "不能邀请自己")
	ErrUserPaymentCheck     = apperrors.NewBusinessError("USER_007", "用户收款方式已存在")
	ErrUserPaymentNotFound  = apperrors.NewNotFoundError("USER_008", "收款方式不存在")
	ErrInvalidReferrerPhone = apperrors.NewBusinessError("USER_009", "邀请人不存在")
)

// 工具函数

// WrapRepositoryError 包装repository层返回的错误
func WrapRepositoryError(err error, defaultErr *apperrors.AppError) error {
	if err == nil {
		return nil
	}

	// 如果已经是AppError，直接返回
	if apperrors.IsAppError(err) {
		return err
	}

	// 否则返回默认的业务错误
	return defaultErr
}

// HandleRepositoryError 处理repository层错误，提供更具体的错误信息
func HandleRepositoryError(err error, resourceType string) error {
	if err == nil {
		return nil
	}

	// 如果已经是AppError，直接返回
	if apperrors.IsAppError(err) {
		return err
	}

	// 根据资源类型返回相应的NotFound错误
	switch resourceType {
	default:
		return apperrors.ErrResourceNotFound
	}
}
