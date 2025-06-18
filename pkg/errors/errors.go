package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Domain errors
var (
	// General errors
	ErrNotFound           = NewAppError("resource not found", http.StatusNotFound)
	ErrAlreadyExists      = NewAppError("resource already exists", http.StatusConflict)
	ErrInvalidInput       = NewAppError("invalid input", http.StatusBadRequest)
	ErrUnauthorized       = NewAppError("unauthorized", http.StatusUnauthorized)
	ErrForbidden          = NewAppError("forbidden", http.StatusForbidden)
	ErrInternal           = NewAppError("internal server error", http.StatusInternalServerError)
	ErrServiceUnavailable = NewAppError("service unavailable", http.StatusServiceUnavailable)

	// Business logic errors
	ErrInsufficientStock = NewAppError("insufficient stock", http.StatusBadRequest)
	ErrInvalidPrice      = NewAppError("invalid price", http.StatusBadRequest)
	ErrOrderNotFound     = NewAppError("order not found", http.StatusNotFound)
	ErrProductNotFound   = NewAppError("product not found", http.StatusNotFound)

	// Auth errors
	ErrInvalidCredentials = NewAppError("invalid credentials", http.StatusUnauthorized)
	ErrTokenExpired       = NewAppError("token expired", http.StatusUnauthorized)
	ErrInvalidToken       = NewAppError("invalid token", http.StatusUnauthorized)
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Message    string                 `json:"message"`
	StatusCode int                    `json:"status_code"`
	Details    map[string]interface{} `json:"details,omitempty"`
	err        error
}

// NewAppError creates a new application error
func NewAppError(message string, statusCode int) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: statusCode,
	}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.err)
	}
	return e.Message
}

// WithError wraps an error
func (e *AppError) WithError(err error) *AppError {
	newErr := *e
	newErr.err = err
	return &newErr
}

// WithDetails adds additional details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	newErr := *e
	newErr.Details = details
	return &newErr
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.err
}

// Is checks if the error is of a specific type
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Message == t.Message && e.StatusCode == t.StatusCode
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError returns the AppError from an error
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
