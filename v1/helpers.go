package v1

import (
	"io"
	"io/ioutil"
	"net/http"
)

const MB = 1 << 20

func buildLimitedRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	limitReader := io.LimitReader(resp.Body, MB)
	body, err := ioutil.ReadAll(limitReader)

	if err != nil {
		return body, err
	}

	return body, nil
}
