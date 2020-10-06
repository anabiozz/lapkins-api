package products

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type serviceError struct {
	code    int
	message string
}

func (e *serviceError) UnmarshalJSON(data []byte) error {
	tmp := &struct {
		Err string `json:"error"`
	}{}

	err := json.Unmarshal(data, tmp)
	if err != nil {
		return err
	}

	e.message = tmp.Err

	return nil
}

func (e *serviceError) MarshalJSON() ([]byte, error) {
	tmp := &struct {
		Err string `json:"error"`
	}{
		Err: e.message,
	}

	return json.Marshal(tmp)
}

func (e *serviceError) StatusCode() int {
	return e.code
}

func (e *serviceError) Error() string {
	return e.message
}

func errBadRequest(format string, v ...interface{}) error {
	return &serviceError{
		code:    http.StatusBadRequest,
		message: fmt.Sprintf(format, v...),
	}
}

func errNotFound(format string, v ...interface{}) error {
	return &serviceError{
		code:    http.StatusNotFound,
		message: fmt.Sprintf(format, v...),
	}
}

func errInternal(format string, v ...interface{}) error {
	return &serviceError{
		code:    http.StatusInternalServerError,
		message: fmt.Sprintf(format, v...),
	}
}
