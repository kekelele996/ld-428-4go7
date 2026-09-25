package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/artvault/artvault/internal/constants"
	"github.com/artvault/artvault/internal/util"
	"github.com/gin-gonic/gin"
)

// ErrorHandler 统一错误响应。
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		var appErr *util.AppError
		if errors.As(err, &appErr) {
			status := http.StatusBadRequest
			switch appErr.Code {
			case constants.CodeUnauthorized, constants.CodeTokenExpired, constants.CodeInvalidCredentials:
				status = http.StatusUnauthorized
			case constants.CodeForbidden:
				status = http.StatusForbidden
			case constants.CodeNotFound:
				status = http.StatusNotFound
			case constants.CodeConflict, constants.CodeUserExists:
				status = http.StatusConflict
			case constants.CodeInternalError, constants.CodeDatabaseError:
				status = http.StatusInternalServerError
			}
			util.Fail(c, status, appErr.Code, appErr.Message)
			return
		}
		logger.Error("unhandled error", "error", err)
		util.Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MsgInternalError)
	}
}
