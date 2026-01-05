package middleware

import (
	"errors"
	"net/http"

	apperrors "minigo/internal/domain/errors"
	"minigo/internal/infrastructure/logging"
	resp "minigo/internal/interfaces/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware catches panics and returns unified error response.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logging.L().WithField("panic", recovered).Error("panic_recovered")
		resp.Error(c, http.StatusInternalServerError, "系统内部错误")
		c.Abort()
	})
}

// HandleError processes errors and returns appropriate response
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	logging.L().WithError(err).Error("error_occurred")

	// Check if it's an AppError
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		// Use the error code and message from AppError
		resp.ErrorCode(c, appErr.Code, appErr.Message)
		return
	}

	// Log unexpected errors
	logging.L().WithError(err).Error("unexpected_error")

	// Return generic internal error for unknown errors
	resp.Error(c, http.StatusInternalServerError, "系统内部错误")
}
