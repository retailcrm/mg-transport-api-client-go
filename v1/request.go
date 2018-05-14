package v1

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var prefix = "/api/v1"

// GetRequest implements GET Request
func (c *MgClient) GetRequest(urlWithParameters string) ([]byte, int, Failure) {
	var res []byte
	var cerr Failure

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.URL, prefix, urlWithParameters), nil)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	req.Header.Set("X-Transport-Token", c.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		cerr.ApiErr = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, cerr
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		cerr.RuntimeErr = err
	}

	return res, resp.StatusCode, cerr
}

// PostRequest implements POST Request
func (c *MgClient) PostRequest(url string, parameters []byte) ([]byte, int, Failure) {
	var res []byte
	var cerr Failure

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", c.URL, prefix, url), bytes.NewBuffer(parameters))
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transport-Token", c.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		cerr.ApiErr = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, cerr
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	return res, resp.StatusCode, cerr
}

// PutRequest implements PUT Request
func (c *MgClient) PutRequest(url string, parameters []byte) ([]byte, int, Failure) {
	var res []byte
	var cerr Failure

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s%s", c.URL, prefix, url), bytes.NewBuffer(parameters))
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transport-Token", c.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		cerr.ApiErr = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, cerr
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	return res, resp.StatusCode, cerr
}

// DeleteRequest implements DELETE Request
func (c *MgClient) DeleteRequest(url string) ([]byte, int, Failure) {
	var res []byte
	var buf []byte
	var cerr Failure

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", c.URL, prefix, url), bytes.NewBuffer(buf))
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Transport-Token", c.Token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		cerr.ApiErr = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, cerr
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	return res, resp.StatusCode, cerr
}

func buildRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	return res, nil
}

func buildErr(data []byte) Failure {
	var err = Failure{}

	eresp, errr := ErrorResponse(data)
	err.RuntimeErr = errr
	err.ApiErr = eresp.ErrorMsg

	if eresp.Errors != nil {
		err.ApiErrs = eresp.Errors
	}

	return err
}
