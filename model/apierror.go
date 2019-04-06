package model

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type ApiResponseError struct {
	error   `json:"-"` // internal error/response
	Code    int        `json:"code"`              // error code
	Type    string     `json:"type,omitempty"`    // unused?
	Message string     `json:"message,omitempty"` // external error/response
}

const (
	ApiErrorNotFound   = http.StatusNotFound
	ApiErrorBadRequest = http.StatusBadRequest
	ApiErrorInternal   = http.StatusInternalServerError

	ApiErrorCount  = http.StatusUnprocessableEntity
	ApiErrorExists = http.StatusExpectationFailed
)

// NewApiError is used to respond with a message and code
func NewApiError(code int, message string, internal error) *ApiResponseError {
	res := ApiResponseError{Code: code, Message: message}
	if res.Message == "" {
		res.Message = http.StatusText(code)
	}
	if internal == nil {
		internal = errors.New(res.Message)
	}
	res.error = internal

	return &res
}

// Error makes it compatible with error interface.
func (res *ApiResponseError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", res.Code, res.Message)
}

func (res *ApiResponseError) Internal() error {
	return res.error
}
