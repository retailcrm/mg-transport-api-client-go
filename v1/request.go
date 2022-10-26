package v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var prefix = "/api/transport/v1"

// GetRequest implements GET Request.
func (c *MgClient) GetRequest(url string, parameters []byte) ([]byte, int, error) {
	return makeRequest(
		"GET",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		bytes.NewBuffer(parameters),
		c,
	)
}

// PostRequest implements POST Request.
func (c *MgClient) PostRequest(url string, parameters io.Reader) ([]byte, int, error) {
	return makeRequest(
		"POST",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		parameters,
		c,
	)
}

// PutRequest implements PUT Request.
func (c *MgClient) PutRequest(url string, parameters []byte) ([]byte, int, error) {
	return makeRequest(
		"PUT",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		bytes.NewBuffer(parameters),
		c,
	)
}

// DeleteRequest implements DELETE Request.
func (c *MgClient) DeleteRequest(url string, parameters []byte) ([]byte, int, error) {
	return makeRequest(
		"DELETE",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		bytes.NewBuffer(parameters),
		c,
	)
}

func makeRequest(reqType, url string, buf io.Reader, c *MgClient) ([]byte, int, error) {
	var res []byte
	req, err := http.NewRequest(reqType, url, buf)
	if err != nil {
		return res, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transport-Token", c.Token)

	if c.Debug {
		if strings.Contains(url, "/files/upload") {
			c.writeLog("MG TRANSPORT API Request: %s %s %s [file data]", reqType, url, c.Token)
		} else {
			c.writeLog("MG TRANSPORT API Request: %s %s %s %v", reqType, url, c.Token, buf)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return res, 0, NewCriticalHTTPError(err)
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		err = NewServerError(resp)
		return res, resp.StatusCode, err
	}

	res, err = buildLimitedRawResponse(resp)
	if err != nil {
		return res, 0, err
	}

	if c.Debug {
		c.writeLog("MG TRANSPORT API Response: %s", res)
	}

	return res, resp.StatusCode, err
}
