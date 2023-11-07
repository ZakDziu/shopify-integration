package model

import "net/http"

var (
	ErrUnhealthy    = NewError(http.StatusInternalServerError, "something went wrong")
	ErrUnauthorized = NewError(http.StatusUnauthorized, "user unauthorized")
	ErrInvalidBody  = NewError(http.StatusBadRequest, "request invalid body")
)

type Error interface {
	error
	Status() int
}

// StatusError represents HTTP error.
type StatusError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Message
}

// Status returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

//nolint:ireturn
func NewError(code int, message string) Error {
	return StatusError{
		Code:    code,
		Message: message,
	}
}
