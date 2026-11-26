package v1

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequest429Retry(t *testing.T) {
	const payload = `{"message":{"externalId":"1"},"channelId":2}`
	const rateLimitResponse = `{"errors":["Too many requests"]}`
	const successResponse = `{}`

	cases := []struct {
		name         string
		body         io.Reader
		requestBody  string
		limiter      Limiter
		rateLimited  int32
		requestCount int
		status       int
		responseBody string
	}{
		{
			name: "replay POST on each retry", body: strings.NewReader(payload), requestBody: payload,
			limiter: NoopLimiter, rateLimited: 3, requestCount: 4, status: http.StatusOK, responseBody: successResponse,
		},
		{
			name: "return final 429", body: strings.NewReader(payload), requestBody: payload,
			limiter: NoopLimiter, rateLimited: 4, requestCount: 4,
			status: http.StatusTooManyRequests, responseBody: rateLimitResponse,
		},
		{
			name: "one retry without limiter", body: strings.NewReader(payload), requestBody: payload,
			rateLimited: 2, requestCount: 2, status: http.StatusTooManyRequests, responseBody: rateLimitResponse,
		},
		{
			name: "retry without a body", limiter: NoopLimiter, rateLimited: 1,
			requestCount: 2, status: http.StatusOK, responseBody: successResponse,
		},
		{
			name: "return 429 for streaming body",
			body: struct{ io.Reader }{strings.NewReader(payload)}, requestBody: payload,
			limiter: NoopLimiter, rateLimited: 1, requestCount: 1,
			status: http.StatusTooManyRequests, responseBody: rateLimitResponse,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			type observedRequest struct {
				body string
				addr string
			}
			requests := make(chan observedRequest, tc.requestCount)
			var attempts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					return
				}
				select {
				case requests <- observedRequest{body: string(body), addr: r.RemoteAddr}:
				default:
					t.Error("too many requests")
				}
				if attempts.Add(1) <= tc.rateLimited {
					w.WriteHeader(http.StatusTooManyRequests)
					_, _ = io.WriteString(w, rateLimitResponse)
					return
				}
				_, _ = io.WriteString(w, successResponse)
			}))
			defer server.Close()

			httpClient := server.Client()
			httpClient.Timeout = 3 * time.Second
			client := NewWithClient(server.URL, "test-token", httpClient).WithLimiter(tc.limiter)
			body, status, err := client.PostRequest("/messages/read", tc.body)
			require.NoError(t, err)
			require.Equal(t, tc.status, status)
			require.Equal(t, tc.responseBody, string(body))
			require.Len(t, requests, tc.requestCount)

			first := <-requests
			require.Equal(t, tc.requestBody, first.body)
			for request := 1; request < tc.requestCount; request++ {
				require.Equal(t, first, <-requests, "retry should replay the body and reuse the connection")
			}
		})
	}
}
