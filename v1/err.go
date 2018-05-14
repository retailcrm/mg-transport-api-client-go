package v1

import (
	"encoding/json"
	"fmt"
)

// Error implements generic error interface
type Error interface {
	error
	ApiError() string
	ApiErrors() map[string]string
}

// Failure struct implode runtime & api errors
type Failure struct {
	RuntimeErr error
	ApiErr     string
	ApiErrs    map[string]string
}

// FailureResponse convert json error response into object
type FailureResponse struct {
	ErrorMsg string            `json:"errorMsg,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

// Error returns the string representation of the error and satisfies the error interface.
func (f *Failure) Error() string {
	return f.RuntimeErr.Error()
}

// ApiError returns formatted string representation of the API error
func (f *Failure) ApiError() string {
	return fmt.Sprintf("%v", f.ApiErr)
}

// ApiErrors returns array of formatted strings that represents API errors
func (f *Failure) ApiErrors() map[string]string {
	return f.ApiErrs
}

// ErrorResponse method
func ErrorResponse(data []byte) (FailureResponse, error) {
	var resp FailureResponse
	err := json.Unmarshal(data, &resp)

	return resp, err
}
