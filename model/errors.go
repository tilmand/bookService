package model

import "net/http"

var (
	ErrInternalServerError = NewError(http.StatusInternalServerError, "something went wrong")
	ErrInvalidBody         = NewError(http.StatusBadRequest, "request invalid body")
	ErrUnauthorized        = NewError(http.StatusUnauthorized, "user unauthorized")
	ErrRefreshExpired      = NewError(http.StatusUnauthorized, "refresh")
	ErrForbidden           = NewError(http.StatusForbidden, "forbidden")
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (se StatusError) Error() string {
	return se.Message
}

func (se StatusError) Status() int {
	return se.Code
}

func NewError(code int, message string) Error {
	return StatusError{
		Code:    code,
		Message: message,
	}
}
