package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var defaultErrorMessage = "http client error"
var internalServerError = "internal server error"
var marshalError = "cannot unmarshal response body"

// MGErrors contains a list of errors as sent by MessageGateway.
type MGErrors struct {
	Errors []string
}

// HTTPClientError is a common error type used in the client.
type HTTPClientError struct {
	ErrorMsg  string
	BaseError error
	Response  io.Reader
}

// Unwrap returns underlying error. Its presence usually indicates a problem with the network.
func (err *HTTPClientError) Unwrap() error {
	return err.BaseError
}

// Error message will contain either an error from MG or underlying error message.
func (err *HTTPClientError) Error() string {
	message := defaultErrorMessage

	if err.BaseError != nil {
		message = fmt.Sprintf("%s: %s", defaultErrorMessage, err.BaseError.Error())
	} else if len(err.ErrorMsg) > 0 {
		message = err.ErrorMsg
	}

	return message
}

// NewCriticalHTTPError wraps *http.Client error.
func NewCriticalHTTPError(err error) error {
	return &HTTPClientError{BaseError: err}
}

// NewAPIClientError wraps MG error.
func NewAPIClientError(responseBody []byte) error {
	var data MGErrors
	var message string

	if len(responseBody) == 0 {
		message = internalServerError
	} else {
		err := json.Unmarshal(responseBody, &data)

		if err != nil {
			message = marshalError
		} else if len(data.Errors) > 0 {
			message = data.Errors[0]
		}
	}

	return &HTTPClientError{ErrorMsg: message}
}

// NewServerError wraps an unexpected API error (e.g. 5xx).
func NewServerError(response *http.Response) error {
	var serverError *HTTPClientError

	body, _ := buildLimitedRawResponse(response)
	err := NewAPIClientError(body)

	if errors.As(err, &serverError) && len(body) > 0 {
		serverError.Response = bytes.NewBuffer(body)
		return serverError
	}

	return err
}

func AsClientError(err error) *HTTPClientError {
	for {
		if err == nil {
			return nil
		}
		if typed, ok := err.(*HTTPClientError); ok {
			return typed
		}
		err = errors.Unwrap(err)
	}
}
