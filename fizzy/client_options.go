package fizzy

import "time"

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithTimeout sets the HTTP client timeout.
// Default is 30 seconds.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithBaseURL sets the base URL for the Fizzy API.
// Default is "https://app.fizzy.do".
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithCache enables or disables ETag caching.
// Default is true (enabled).
func WithCache(enabled bool) ClientOption {
	return func(c *Client) {
		c.cacheEnabled = enabled
	}
}

// WithMaxRetries sets the maximum number of retry attempts.
// Default is 3.
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}
