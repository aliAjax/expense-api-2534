package model

import "net/http"

const (
	CodeSuccess        = 0
	CodeInvalidRequest = 40000
	CodeNotFound       = 40400
	CodeConflict       = 40900
	CodeInternal       = 50000
)

type APIError struct {
	Status  int
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(status, code int, message string) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func NewValidationError(message string) *APIError {
	return NewAPIError(http.StatusBadRequest, CodeInvalidRequest, message)
}

func NewNotFoundError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, CodeInternal, message)
}

func NewConflictError(message string) *APIError {
	return NewAPIError(http.StatusConflict, CodeConflict, message)
}

func NewInternalError(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, CodeInternal, message)
}
