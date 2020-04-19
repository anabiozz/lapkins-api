package users

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	errUserNotFound = errForbidden("user not found")
)

// Error is a user-info service error.
type Error struct {
	err       string
	RequestID string
	ErrorCode int
}

func (e *Error) MarshalJSON() ([]byte, error) {
	tempErr := struct {
		RequestID string
		ErrorCode int
	}{
		RequestID: e.RequestID,
		ErrorCode: e.ErrorCode,
	}

	return json.Marshal(tempErr)
}

// Error returns a text message corresponding to the given error.
func (e *Error) Error() string {
	return e.err
}

// StatusCode returns an HTTP status code corresponding to the given error.
func (e *Error) StatusCode() int {
	return e.ErrorCode
}

func errBadRequest(format string, v ...interface{}) error {
	return &Error{
		ErrorCode: http.StatusBadRequest,
		err:       fmt.Sprintf(format, v...),
	}
}

func errConflict(format string, v ...interface{}) error {
	return &Error{
		ErrorCode: http.StatusConflict,
		err:       fmt.Sprintf(format, v...),
	}
}

func errNotFound(format string, v ...interface{}) error {
	return &Error{
		ErrorCode: http.StatusNotFound,
		err:       fmt.Sprintf(format, v...),
	}
}

func errForbidden(format string, v ...interface{}) error {
	return &Error{
		ErrorCode: http.StatusForbidden,
		err:       fmt.Sprintf(format, v...),
	}
}

func errUnauthorized(format string, v ...interface{}) error {
	return &Error{
		ErrorCode: http.StatusUnauthorized,
		err:       fmt.Sprintf(format, v...),
	}
}

func errInternal(format string, v ...interface{}) error {
	return &Error{
		ErrorCode: http.StatusInternalServerError,
		err:       fmt.Sprintf(format, v...),
	}
}

func errIsNotFound(err error) bool {
	type notFound interface {
		NotFound() bool
	}
	e, ok := err.(notFound)
	return ok && e.NotFound()
}
