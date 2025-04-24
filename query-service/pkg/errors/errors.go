package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application error with context for API responses
type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
	Details []string
}

// Error returns the string representation of the error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Err.Error())
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError with the given message
func New(message string) error {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

// Wrap wraps an error with a message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}

	// If it's already an AppError, just update the message
	var appErr *AppError
	if errors.As(err, &appErr) {
		appErr.Message = fmt.Sprintf("%s: %s", message, appErr.Message)
		return appErr
	}

	// Otherwise, create a new AppError
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// NotFound creates a not found error for the given entity
func NotFound(entity, id string) error {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s with ID %s not found", entity, id),
		Status:  http.StatusNotFound,
	}
}

// BadRequest creates a bad request error with the given message
func BadRequest(message string) error {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// BadRequestWithDetails creates a bad request error with the given message and details
func BadRequestWithDetails(message string, details []string) error {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: message,
		Status:  http.StatusBadRequest,
		Details: details,
	}
}

// ValidationError creates a validation error with the given details
func ValidationError(details []string) error {
	return &AppError{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Status:  http.StatusBadRequest,
		Details: details,
	}
}

// Internal creates an internal server error
func Internal(message string) error {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Status == http.StatusNotFound
	}
	return false
}

// IsBadRequest checks if the error is a bad request error
func IsBadRequest(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Status == http.StatusBadRequest
	}
	return false
}

// Status returns the HTTP status code for the error
func Status(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Status
	}
	return http.StatusInternalServerError
}

// Code returns the error code
func Code(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return "INTERNAL_ERROR"
}

// Details returns the error details
func Details(err error) []string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Details
	}
	return nil
}