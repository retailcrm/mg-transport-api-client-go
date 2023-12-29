package v1

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const MB = 1 << 20
const LimitResponse = 25 * MB

func buildLimitedRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	limitReader := io.LimitReader(resp.Body, LimitResponse)
	body, err := ioutil.ReadAll(limitReader)

	if err != nil {
		return body, err
	}

	return body, nil
}

// BoolPtr returns provided boolean as pointer. Can be used while editing the integration module activity.
func BoolPtr(v bool) *bool {
	return &v
}

// TimePtr returns provided time.Time's pointer.
func TimePtr(v time.Time) *time.Time {
	return &v
}
