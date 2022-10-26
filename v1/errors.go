package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var defaultErrorMessage = "internal http client error"
var internalServerError = "internal server error"

type httpClientError struct {
	ErrorMsg        string
	BaseError       error
	LimitedResponse io.Reader
}

func (err *httpClientError) Unwrap() error {
	return err.BaseError
}

func (err *httpClientError) Is(target error) bool {
	return errors.As(target, &err)
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
	body, err := buildLimitedRawResponse(response)
	if err == nil {
		data = body
	}

	err = NewAPIClientError(data)
	var serverError *httpClientError

	if errors.As(err, &serverError) {
		serverError.LimitedResponse = bytes.NewBuffer(body)
		return serverError
	}

	return err
}
