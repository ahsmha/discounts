package errors

import (
	"errors"
	"fmt"
)

// Error types for better error handling
var (
	ErrValidation = errors.New("validation error")
	ErrNotFound   = errors.New("not found")
	ErrInternal   = errors.New("internal error")
)

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) error {
	return ValidationError{Message: message}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var validationErr ValidationError
	return errors.As(err, &validationErr)
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) error {
	return NotFoundError{Message: message}
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	var notFoundErr NotFoundError
	return errors.As(err, &notFoundErr)
}

// InternalError represents an internal server error
type InternalError struct {
	Message string
	Cause   error
}

func (e InternalError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// NewInternalError creates a new internal error
func NewInternalError(message string, cause error) error {
	return InternalError{Message: message, Cause: cause}
}

// IsInternalError checks if an error is an internal error
func IsInternalError(err error) bool {
	var internalErr InternalError
	return errors.As(err, &internalErr)
}