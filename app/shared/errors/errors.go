package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"
)

// AppError represents an application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"`
	Err        error     `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Wrap wraps an existing error with AppError
func Wrap(err error, code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Err:        err,
	}
}

// Predefined error constructors
var (
	ErrNotFound      = New(ErrCodeNotFound, "Resource not found", http.StatusNotFound)
	ErrUnauthorized  = New(ErrCodeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden     = New(ErrCodeForbidden, "Forbidden", http.StatusForbidden)
	ErrInternal      = New(ErrCodeInternal, "Internal server error", http.StatusInternalServerError)
	ErrBadRequest    = New(ErrCodeBadRequest, "Bad request", http.StatusBadRequest)
	ErrConflict      = New(ErrCodeConflict, "Resource conflict", http.StatusConflict)
)

// NewNotFound creates a not found error
func NewNotFound(message string) *AppError {
	return New(ErrCodeNotFound, message, http.StatusNotFound)
}

// NewValidation creates a validation error
func NewValidation(message string) *AppError {
	return New(ErrCodeValidation, message, http.StatusBadRequest)
}

// NewUnauthorized creates an unauthorized error
func NewUnauthorized(message string) *AppError {
	if message == "" {
		message = "Unauthorized"
	}
	return New(ErrCodeUnauthorized, message, http.StatusUnauthorized)
}

// NewForbidden creates a forbidden error
func NewForbidden(message string) *AppError {
	if message == "" {
		message = "Forbidden"
	}
	return New(ErrCodeForbidden, message, http.StatusForbidden)
}

// NewInternal creates an internal error
func NewInternal(message string, err error) *AppError {
	return Wrap(err, ErrCodeInternal, message, http.StatusInternalServerError)
}

// NewConflict creates a conflict error
func NewConflict(message string) *AppError {
	return New(ErrCodeConflict, message, http.StatusConflict)
}

// NewBadRequest creates a bad request error
func NewBadRequest(message string) *AppError {
	return New(ErrCodeBadRequest, message, http.StatusBadRequest)
}

// IsAppError checks if error is AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError converts error to AppError
func AsAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return NewInternal("Unexpected error", err)
}
