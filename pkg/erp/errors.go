package erp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Storage-related errors.
var (
	ErrNotFoundInStorage     = errors.New("not found in storage")
	ErrDuplicateKeyInStorage = errors.New("duplicate key in storage")
	ErrNotExist              = errors.New("not exist")
)

// ServiceError describes a web-service error.
type ServiceError struct {
	Code    int
	Message string
}

// Error returns a string representation of the error.
func (e *ServiceError) Error() string {
	return fmt.Sprintf("status %d: %s", e.Code, e.Message)
}

// StatusCode returns HTTP status code of the error.
func (e *ServiceError) StatusCode() int {
	return e.Code
}

// Encode encodes the error using the given HTTP response writer.
func (e *ServiceError) Encode(w http.ResponseWriter) {
	message := e.Message
	if e.Code == http.StatusInternalServerError {
		message = "internal error"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Decode decodes the error from the given HTTP response.
func (e *ServiceError) Decode(r *http.Response) {
	e.Code = r.StatusCode
	var res struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(r.Body).Decode(&res); err == nil && res.Error != "" {
		e.Message = res.Error
	} else {
		e.Message = http.StatusText(r.StatusCode)
	}
}

// ErrBadRequest creates a BadRequest service error.
func ErrBadRequest(format string, v ...interface{}) error {
	return &ServiceError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf(format, v...),
	}
}

// ErrUnauthorized creates a NotFound service error.
func ErrUnauthorized(format string, v ...interface{}) error {
	return &ServiceError{
		Code:    http.StatusUnauthorized,
		Message: fmt.Sprintf(format, v...),
	}
}

// ErrNotFound creates a NotFound service error.
func ErrNotFound(format string, v ...interface{}) error {
	return &ServiceError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf(format, v...),
	}
}

// ErrConflict creates a Conflict service error.
//func ErrConflict(format string, v ...interface{}) error {
//	return &ServiceError{
//		Code:    http.StatusConflict,
//		Message: fmt.Sprintf(format, v...),
//	}
//}

// ErrInternal creates an Internal service error.
func ErrInternal(format string, v ...interface{}) error {
	return &ServiceError{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf(format, v...),
	}
}
