package v1

import (
	"encoding/json"
)

var defaultErrorMessage = "Internal http client error"
var internalServerError = "Internal server error"

type httpClientError struct {
	ErrorMsg  string
	BaseError error
}

func (err *httpClientError) Unwrap() error {
	return err.BaseError
}

func (err *httpClientError) Error() string {
	message := defaultErrorMessage

	if err.BaseError != nil {
		message = err.BaseError.Error()
	}

	if len([]rune(err.ErrorMsg)) > 0 {
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