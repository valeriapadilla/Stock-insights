package errors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound   ErrorType = "NOT_FOUND"
	ErrorTypeDatabase   ErrorType = "DATABASE_ERROR"
	ErrorTypeInternal   ErrorType = "INTERNAL_ERROR"
	ErrorTypeExternal   ErrorType = "EXTERNAL_ERROR"
)

type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewValidationError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Code:    http.StatusBadRequest,
		Err:     err,
	}
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Code:    http.StatusNotFound,
		Err:     err,
	}
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeDatabase,
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

func NewExternalError(message string, err error) *AppError {
	return &AppError{
		Type:    ErrorTypeExternal,
		Message: message,
		Code:    http.StatusBadGateway,
		Err:     err,
	}
}
