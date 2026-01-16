package fizzy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetry_SuccessOnFirstAttempt(t *testing.T) {
	var callCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.doRequest(context.Background(), "GET", "/test", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, int32(1), callCount.Load(), "Should not retry on success")
}

func TestRetry_ServerErrorWithRetry(t *testing.T) {
	var callCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := callCount.Add(1)
		if count < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "server error"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL), WithMaxRetries(3))
	err := client.doRequest(context.Background(), "GET", "/test", nil, nil)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, callCount.Load(), int32(3), "Should retry on 500 errors")
}

func TestRetry_MaxRetriesExceeded(t *testing.T) {
	var callCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error": "bad gateway"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL), WithMaxRetries(2))
	err := client.doRequest(context.Background(), "GET", "/test", nil, nil)

	require.Error(t, err)
	// Initial attempt + 2 retries = 3 total
	assert.Equal(t, int32(3), callCount.Load(), "Should try initial + maxRetries")
}

func TestRetry_ContextCancellationStopsRetries(t *testing.T) {
	var callCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL), WithMaxRetries(10))
	err := client.doRequest(ctx, "GET", "/test", nil, nil)

	require.Error(t, err)
	// Should have made some calls but not all 10 retries
	count := callCount.Load()
	assert.Greater(t, count, int32(0))
	assert.Less(t, count, int32(10), "Context cancellation should stop retries early")
}

func TestRetry_NonRetryableStatusCodes(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"Bad Request", http.StatusBadRequest},
		{"Unauthorized", http.StatusUnauthorized},
		{"Forbidden", http.StatusForbidden},
		{"Not Found", http.StatusNotFound},
		{"Unprocessable Entity", http.StatusUnprocessableEntity},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var callCount atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				callCount.Add(1)
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(`{"error": "error"}`))
			}))
			defer server.Close()

			client := NewClient("test-token", "test-account", WithBaseURL(server.URL), WithMaxRetries(3))
			err := client.doRequest(context.Background(), "GET", "/test", nil, nil)

			require.Error(t, err)
			assert.Equal(t, int32(1), callCount.Load(), "Should not retry 4xx errors")
		})
	}
}

func TestRetry_IsRetryableError(t *testing.T) {
	testCases := []struct {
		name      string
		err       error
		retryable bool
	}{
		{"Context Canceled", context.Canceled, false},
		{"Context Deadline Exceeded", context.DeadlineExceeded, false},
		{"Generic Error", errors.New("network error"), true},
		{"Nil Error", nil, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isRetryableError(tc.err)
			assert.Equal(t, tc.retryable, result)
		})
	}
}

func TestRetry_BackoffDelay(t *testing.T) {
	testCases := []struct {
		attempt      int
		expectedMin  time.Duration
		expectedMax  time.Duration
	}{
		{0, 1 * time.Second, 1 * time.Second},
		{1, 2 * time.Second, 2 * time.Second},
		{2, 4 * time.Second, 4 * time.Second},
		{3, 8 * time.Second, 8 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(string(rune('0'+tc.attempt)), func(t *testing.T) {
			delay := backoffDelay(tc.attempt)
			assert.GreaterOrEqual(t, delay, tc.expectedMin)
			assert.LessOrEqual(t, delay, tc.expectedMax)
		})
	}
}

func TestRetry_RateLimitWithRetryAfter(t *testing.T) {
	var callCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := callCount.Add(1)
		if count == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate limited"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL), WithMaxRetries(2))
	err := client.doRequest(context.Background(), "GET", "/test", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, int32(2), callCount.Load(), "Should retry after rate limit")
}

func TestRetry_ServiceUnavailable(t *testing.T) {
	var callCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := callCount.Add(1)
		if count <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "service unavailable"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL), WithMaxRetries(3))
	err := client.doRequest(context.Background(), "GET", "/test", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, int32(3), callCount.Load(), "Should retry 503 errors")
}

func TestRetry_GatewayTimeout(t *testing.T) {
	var callCount atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := callCount.Add(1)
		if count == 1 {
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write([]byte(`{"error": "gateway timeout"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL), WithMaxRetries(2))
	err := client.doRequest(context.Background(), "GET", "/test", nil, nil)

	require.NoError(t, err)
	assert.Equal(t, int32(2), callCount.Load(), "Should retry 504 errors")
}
