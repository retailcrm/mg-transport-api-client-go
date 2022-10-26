package v1

import (
	"io"
	"net/http"
)

func buildRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	return res, nil
}
