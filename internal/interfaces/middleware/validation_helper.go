package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"minigo/internal/interfaces/response"
)

// ValidateAndBindJSON 验证并绑定JSON请求体
func ValidateAndBindJSON(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Warnf("绑定JSON参数失败: %v", err)
		response.ErrorCode(c, "BIND_JSON_ERROR", "参数错误")
		return false
	}
	return true
}

// ValidateAndBindQuery 验证并绑定查询参数
func ValidateAndBindQuery(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindQuery(req); err != nil {
		logrus.Warnf("绑定查询参数失败: %v", err)
		response.ErrorCode(c, "BIND_QUERY_ERROR", "参数错误")
		return false
	}
	return true
}

// ValidateAndBindURI 验证并绑定URI参数
func ValidateAndBindURI(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindUri(req); err != nil {
		logrus.Warnf("绑定URI参数失败: %v", err)
		response.ErrorCode(c, "BIND_URI_ERROR", "参数错误")
		return false
	}
	return true
}

// ValidateIDParam 验证ID参数
func ValidateIDParam(c *gin.Context, paramName string) (int64, bool) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		logrus.Warnf("无效的ID参数: %s", idStr)
		response.ErrorCode(c, "INVALID_ID_ERROR", "无效的ID参数")
		return 0, false
	}
	return id, true
}

// ValidatePagination 验证并设置分页参数默认值
func ValidatePagination(page, size *int) {
	if *page <= 0 {
		*page = 1
	}
	if *size <= 0 {
		*size = 10
	}
}

// ValidateRequiredString 验证必需的字符串参数
func ValidateRequiredString(c *gin.Context, value, fieldName string) bool {
	if value == "" {
		response.ErrorCode(c, "MISSING_REQUIRED_ERROR", fieldName+"不能为空")
		return false
	}
	return true
}

// ValidatePositiveInt64 验证正整数参数
func ValidatePositiveInt64(c *gin.Context, value int64, fieldName string) bool {
	if value <= 0 {
		response.ErrorCode(c, "INVALID_POSITIVE_ERROR", fieldName+"必须大于0")
		return false
	}
	return true
}
