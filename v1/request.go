package v1

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const MaxRPS = 100

var prefix = "/api/transport/v1"

// GetRequest performs GET request to the provided route.
func (c *MgClient) GetRequest(url string, parameters []byte) ([]byte, int, error) {
	return makeRequest(
		"GET",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		bytes.NewBuffer(parameters),
		c,
	)
}

// PostRequest performs POST request to the provided route.
func (c *MgClient) PostRequest(url string, parameters io.Reader) ([]byte, int, error) {
	return makeRequest(
		"POST",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		parameters,
		c,
	)
}

// PutRequest performs PUT request to the provided route.
func (c *MgClient) PutRequest(url string, parameters []byte) ([]byte, int, error) {
	return makeRequest(
		"PUT",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
		bytes.NewBuffer(parameters),
		c,
	)
}

// DeleteRequest performs DELETE request to the provided route.
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

	c.mux.Lock()
	defer c.mux.Unlock()

	attempt := 0
tryAgain:
	sleepTime := time.Second - time.Since(c.lastTime)
	if sleepTime < 0 {
		c.lastTime = time.Now()
		c.rps = 0
	} else if c.rps == MaxRPS {
		time.Sleep(sleepTime)
		c.lastTime = time.Now()
		c.rps = 0
	}
	c.rps++

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return res, 0, NewCriticalHTTPError(err)
	}

	if resp.StatusCode == http.StatusTooManyRequests && attempt < 3 {
		attempt++
		goto tryAgain
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
