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

type MGErrors struct {
	Errors []string
}

type httpClientError struct {
	ErrorMsg  string
	BaseError error
	Response  io.Reader
}

func (err *httpClientError) Unwrap() error {
	return err.BaseError
}

func (err *httpClientError) Error() string {
	message := defaultErrorMessage

	if err.BaseError != nil {
		message = fmt.Sprintf("%s: %s", defaultErrorMessage, err.BaseError.Error())
	}

	if len(err.ErrorMsg) > 0 {
		message = err.ErrorMsg
	}

	return message
}

func NewCriticalHTTPError(err error) error {
	return &httpClientError{BaseError: err}
}

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

	return &httpClientError{ErrorMsg: message}
}

func NewServerError(response *http.Response) error {
	var serverError *httpClientError

	body, _ := buildLimitedRawResponse(response)
	err := NewAPIClientError(body)

	if errors.As(err, &serverError) && len(body) > 0 {
		serverError.Response = bytes.NewBuffer(body)
		return serverError
	}

	return err
}
