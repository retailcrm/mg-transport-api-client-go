package v1

import (
	"io"
	"net/http"
)

const MB = 1 << 20
const MaxSizeBody = MB * 0.5

func buildLimitedRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	limitReader := io.LimitReader(resp.Body, MaxSizeBody)
	body, err := io.ReadAll(limitReader)

	if err != nil {
		return body, err
	}

	return body, nil
}
