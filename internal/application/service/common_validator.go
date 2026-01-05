package service

import (
	"errors"
	"fmt"
	"regexp"
)

// CommonValidator 通用验证器
type CommonValidator struct{}

// NewCommonValidator 创建通用验证器实例
func NewCommonValidator() *CommonValidator {
	return &CommonValidator{}
}

// ValidateShopID 验证店铺ID
func (v *CommonValidator) ValidateShopID(shopID int64) error {
	if shopID <= 0 {
		return errors.New("无效的店铺ID")
	}
	return nil
}

// ValidateUserID 验证用户ID
func (v *CommonValidator) ValidateUserID(userID int64) error {
	if userID <= 0 {
		return errors.New("无效的用户ID")
	}
	return nil
}

// ValidatePhone 验证手机号格式
func (v *CommonValidator) ValidatePhone(phone string) error {
	if phone == "" {
		return errors.New("手机号不能为空")
	}

	// 简单的手机号格式验证
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("手机号格式不正确")
	}

	return nil
}

// ValidatePassword 验证密码强度
func (v *CommonValidator) ValidatePassword(password string) error {
	if password == "" {
		return errors.New("密码不能为空")
	}

	if len(password) < 6 {
		return errors.New("密码长度不能少于6位")
	}

	if len(password) > 20 {
		return errors.New("密码长度不能超过20位")
	}

	return nil
}

// ValidateUsername 验证用户名
func (v *CommonValidator) ValidateUsername(username string) error {
	if username == "" {
		return errors.New("用户名不能为空")
	}

	if len(username) < 2 {
		return errors.New("用户名长度不能少于2位")
	}

	if len(username) > 20 {
		return errors.New("用户名长度不能超过20位")
	}

	return nil
}

// ValidatePagination 验证分页参数
func (v *CommonValidator) ValidatePagination(page, size int) (int, int, error) {
	if page <= 0 {
		page = 1
	}

	if size <= 0 {
		size = 20
	}

	if size > 100 {
		return 0, 0, errors.New("每页数量不能超过100")
	}

	return page, size, nil
}

// ValidateRequiredString 验证必需的字符串字段
func (v *CommonValidator) ValidateRequiredString(value, fieldName string) error {
	if value == "" {
		return fmt.Errorf("%s不能为空", fieldName)
	}
	return nil
}

// ValidateStringLength 验证字符串长度
func (v *CommonValidator) ValidateStringLength(value, fieldName string, minLen, maxLen int) error {
	if len(value) < minLen {
		return fmt.Errorf("%s长度不能少于%d位", fieldName, minLen)
	}

	if len(value) > maxLen {
		return fmt.Errorf("%s长度不能超过%d位", fieldName, maxLen)
	}

	return nil
}

// ValidateEnum 验证枚举值
func (v *CommonValidator) ValidateEnum(value interface{}, validValues []interface{}, fieldName string) error {
	for _, validValue := range validValues {
		if value == validValue {
			return nil
		}
	}
	return fmt.Errorf("%s的值无效", fieldName)
}
