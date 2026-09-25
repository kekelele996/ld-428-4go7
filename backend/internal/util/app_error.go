package util

import (
	"errors"
	"fmt"
)

// AppError 业务错误。
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d message=%s cause=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

func Wrap(code int, message string, err error) error {
	return &AppError{Code: code, Message: message, Err: err}
}

// IsAppError 判断错误链中是否包含指定错误码。
func IsAppError(err error, code int) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}
