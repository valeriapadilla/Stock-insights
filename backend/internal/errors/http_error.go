package errors

import "net/http"

type HTTPError struct {
	Code int
	Message string
	Err error
}

func (e *HTTPError) Error() string {
	return e.Message
}

func New(code int, msg string) *HTTPError {
	return &HTTPError{
		Code: code,
		Message: msg,
	}
}

func Wrap(code int, msg string, err error) *HTTPError {
	return &HTTPError{Code: code, Message: msg, Err: err}
}

var (
	ErrBadRequest     = New(http.StatusBadRequest, "bad request")
	ErrNotFound       = New(http.StatusNotFound, "resource not found")
	ErrInternalServer = New(http.StatusInternalServerError, "internal server error")
)
