package repository

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	apperrors "minigo/internal/domain/errors"
)

// CheckUpdateResult checks the result of an update operation and returns appropriate error
func CheckUpdateResult(result sql.Result, err error) error {
	if err != nil {
		return ConvertError(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ConvertError(err)
	}
	if rowsAffected == 0 {
		return apperrors.ErrResourceNotFound
	}
	return nil
}

// CheckDeleteResult checks the result of a delete operation and returns appropriate error
func CheckDeleteResult(result sql.Result, err error) error {
	if err != nil {
		return ConvertError(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ConvertError(err)
	}
	if rowsAffected == 0 {
		return apperrors.ErrResourceNotFound
	}
	return nil
}

// Now returns current time for consistent timestamp management
func Now() time.Time {
	return time.Now()
}

// ConvertError 转换数据库错误为应用错误
func ConvertError(err error) error {
	if err == nil {
		return nil
	}

	// 将sql.ErrNoRows转换为NotFound错误
	if errors.Is(err, sql.ErrNoRows) {
		return apperrors.ErrResourceNotFound
	}

	slog.Error("数据库错误", "error", err)

	// 其他数据库错误转换为系统错误
	return apperrors.WrapSystemError(err, "DB_001", "数据库操作失败")
}

// ConvertQueryError 转换查询错误
func ConvertQueryError(err error) error {
	return ConvertError(err)
}

// ConvertExecError 转换执行错误
func ConvertExecError(err error) error {
	return ConvertError(err)
}
