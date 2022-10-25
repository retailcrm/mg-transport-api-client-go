package v1

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCriticalHTTPError(t *testing.T) {
	err := &url.Error{Op: "Get", URL: "http//example.com", Err: errors.New("EOF")}
	httpErr := NewCriticalHTTPError(err)

	assert.IsType(t, new(httpClientError), httpErr)
	assert.IsType(t, new(url.Error), errors.Unwrap(httpErr))
	assert.IsType(t, new(url.Error), errors.Unwrap(httpErr))
	assert.Equal(t, httpErr.Error(), err.Error())
}

func TestNewApiClientError(t *testing.T) {
	body := []byte(`{"errors" : ["Channel not found"]}`)
	httpErr := NewAPIClientError(body)

	assert.IsType(t, new(httpClientError), httpErr)
	assert.Equal(t, httpErr.Error(), "Channel not found")

	body = []byte{}
	httpErr = NewAPIClientError(body)

	assert.IsType(t, new(httpClientError), httpErr)
	assert.Equal(t, httpErr.Error(), internalServerError)
}