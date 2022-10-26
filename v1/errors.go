package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var defaultErrorMessage = "Internal http client error"
var internalServerError = "Internal server error"

type httpClientError struct {
	ErrorMsg     string
	BaseError    error
	ResponseBody io.ReadCloser
}

func (err *httpClientError) Unwrap() error {
	return err.BaseError
}

func (err *httpClientError) Error() string {
	message := defaultErrorMessage

	if err.BaseError != nil {
		message = fmt.Sprintf("%s - %s", defaultErrorMessage, err.BaseError.Error())
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
	var data map[string]interface{}
	var message string

	if len(responseBody) == 0 {
		message = internalServerError
	} else {
		if err := json.Unmarshal(responseBody, &data); err != nil {
			return err
		}

		values := data["errors"].([]interface{})
		message = values[0].(string)
	}

	return &httpClientError{ErrorMsg: message}
}

func NewServerError(response *http.Response) error {
	var data []byte
	body, err := buildRawResponse(response)
	if err == nil {
		data = body
	}

	err = NewAPIClientError(data)
	var serverError *httpClientError

	if errors.As(err, &serverError) {
		serverError.ResponseBody = response.Body
		return serverError
	}

	return err
}
