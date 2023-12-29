package v1

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCriticalHTTPError(t *testing.T) {
	err := &url.Error{Op: "Get", URL: "http//example.com", Err: errors.New("EOF")}
	httpErr := NewCriticalHTTPError(err)

	assert.IsType(t, new(HTTPClientError), httpErr)
	assert.IsType(t, new(url.Error), errors.Unwrap(httpErr))
	assert.IsType(t, new(url.Error), errors.Unwrap(httpErr))
	assert.Equal(t, httpErr.Error(), fmt.Sprintf("%s: %s", defaultErrorMessage, err.Error()))
}

func TestNewApiClientError(t *testing.T) {
	body := []byte(`{"errors" : ["Channel not found"]}`)
	httpErr := NewAPIClientError(body)

	assert.IsType(t, new(HTTPClientError), httpErr)
	assert.Equal(t, httpErr.Error(), "Channel not found")

	body = []byte{}
	httpErr = NewAPIClientError(body)

	assert.IsType(t, new(HTTPClientError), httpErr)
	assert.Equal(t, httpErr.Error(), internalServerError)
}

func TestNewServerError(t *testing.T) {
	body := []byte(`{"errors" : ["Something went wrong"]}`)
	response := new(http.Response)
	response.Body = ioutil.NopCloser(bytes.NewReader(body))
	serverErr := NewServerError(response)

	assert.IsType(t, new(HTTPClientError), serverErr)
	assert.Equal(t, serverErr.Error(), "Something went wrong")

	var err *HTTPClientError
	if errors.As(serverErr, &err) {
		assert.NotNil(t, err.Response)
	} else {
		t.Fatal("Unexpected type of error")
	}

	body = []byte(`{"invalid_json"`)
	response = new(http.Response)
	response.Body = ioutil.NopCloser(bytes.NewReader(body))
	serverErr = NewServerError(response)

	assert.IsType(t, new(HTTPClientError), serverErr)
	assert.Equal(t, serverErr.Error(), marshalError)

	if errors.As(serverErr, &err) {
		assert.NotNil(t, err.Response)
	}
}

func TestAsClientError(t *testing.T) {
	assert.Nil(t, AsClientError(nil))
	assert.Nil(t, AsClientError(errors.New("arbitrary")))
	assert.NotNil(t, AsClientError(NewCriticalHTTPError(errors.New("arbitrary"))))
}
