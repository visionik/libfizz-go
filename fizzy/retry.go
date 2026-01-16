package fizzy

import (
	"context"
	"io"
	"math"
	"net/http"
	"time"
)

// retryRoundTripper implements automatic retry logic with exponential backoff.
type retryRoundTripper struct {
	next       http.RoundTripper
	maxRetries int
}

// RoundTrip implements http.RoundTripper with retry logic.
func (t *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// Get context from request
	ctx := req.Context()

	for attempt := 0; attempt <= t.maxRetries; attempt++ {
		// Check context cancellation before each attempt
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		// Clone request for retry (body might have been consumed)
		reqClone := req.Clone(ctx)
		if req.Body != nil {
			// For requests with body, we need to be able to re-read it
			// This is handled by the caller ensuring GetBody is set
			if req.GetBody != nil {
				body, bodyErr := req.GetBody()
				if bodyErr != nil {
					return nil, bodyErr
				}
				reqClone.Body = body
			}
		}

		resp, err = t.next.RoundTrip(reqClone)

		// If request failed, check if we should retry
		if err != nil {
			if attempt < t.maxRetries && isRetryableError(err) {
				delay := backoffDelay(attempt)
				select {
				case <-time.After(delay):
					continue
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
			return nil, err
		}

		// If we got a response, check status code
		statusCode := resp.StatusCode

		// Success - return immediately
		if statusCode < 400 {
			return resp, nil
		}

		// Read and close response body for error cases
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Retry on 429 (rate limit) or 5xx (server errors)
		shouldRetry := (statusCode == 429 || statusCode >= 500) && attempt < t.maxRetries

		if !shouldRetry {
			// Not retryable or out of retries - return error
			message := string(body)
			if message == "" {
				message = http.StatusText(statusCode)
			}
			requestID := resp.Header.Get("X-Request-Id")
			return nil, newErrorFromResponse(statusCode, message, requestID)
		}

		// Calculate delay and wait
		delay := backoffDelay(attempt)
		select {
		case <-time.After(delay):
			// Continue to next retry
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Should not reach here, but return last error if we do
	return resp, err
}

// backoffDelay calculates exponential backoff delay: 1s, 2s, 4s
func backoffDelay(attempt int) time.Duration {
	seconds := math.Pow(2, float64(attempt))
	return time.Duration(seconds) * time.Second
}

// isRetryableError checks if an error is retryable (network errors, timeouts, etc.)
func isRetryableError(err error) bool {
	// For now, we retry on any error that's not context cancellation
	// In production, you might want to be more selective
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}
	return true
}
