package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// 常用的正则表达式
var (
	// EmailRegex 邮箱验证正则
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// PhoneRegex 中国手机号验证正则
	PhoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)
	// IDCardRegex 身份证号验证正则
	IDCardRegex = regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)
	// URLRegex URL验证正则
	URLRegex = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	// IPRegex IP地址验证正则
	IPRegex = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
)

// IsEmail 验证是否为有效邮箱
func IsEmail(email string) bool {
	return EmailRegex.MatchString(email)
}

// IsPhone 验证是否为有效中国手机号
func IsPhone(phone string) bool {
	return PhoneRegex.MatchString(phone)
}

// IsIDCard 验证是否为有效身份证号
func IsIDCard(idCard string) bool {
	return IDCardRegex.MatchString(idCard)
}

// IsURL 验证是否为有效URL
func IsURL(url string) bool {
	return URLRegex.MatchString(url)
}

// IsIP 验证是否为有效IP地址
func IsIP(ip string) bool {
	return IPRegex.MatchString(ip)
}

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// IsNotEmpty 检查字符串是否不为空
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}

// InRange 检查数字是否在指定范围内
func InRange(value, min, max int64) bool {
	return value >= min && value <= max
}

// InSlice 检查元素是否在切片中
func InSlice[T comparable](item T, slice []T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// LengthBetween 检查字符串长度是否在指定范围内
func LengthBetween(s string, min, max int) bool {
	length := utf8.RuneCountInString(s)
	return length >= min && length <= max
}

// MinLength 检查字符串最小长度
func MinLength(s string, min int) bool {
	return utf8.RuneCountInString(s) >= min
}

// MaxLength 检查字符串最大长度
func MaxLength(s string, max int) bool {
	return utf8.RuneCountInString(s) <= max
}

// IsAlpha 检查是否只包含字母
func IsAlpha(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
			return false
		}
	}
	return true
}

// IsAlphanumeric 检查是否只包含字母和数字
func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

// IsNumeric 检查是否只包含数字
func IsNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// ContainsAny 检查字符串是否包含任意一个子串
func ContainsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// ContainsAll 检查字符串是否包含所有子串
func ContainsAll(s string, substrs []string) bool {
	for _, substr := range substrs {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}

// IsStrongPassword 检查是否为强密码（至少8位，包含大小写字母、数字）
func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// ValidatePassword 验证密码强度（可自定义规则）
func ValidatePassword(password string, minLength int, requireUpper, requireLower, requireDigit, requireSpecial bool) bool {
	if len(password) < minLength {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if requireUpper && !hasUpper {
		return false
	}
	if requireLower && !hasLower {
		return false
	}
	if requireDigit && !hasDigit {
		return false
	}
	if requireSpecial && !hasSpecial {
		return false
	}

	return true
}
