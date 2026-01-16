// Package fizzy provides a Go client library for the Fizzy API.
//
// Fizzy is a modern kanban board application for tracking bugs, issues, ideas, and small projects.
// This library provides complete coverage of the Fizzy HTTP API with idiomatic Go patterns.
//
// Basic usage:
//
//	client := fizzy.NewClient("your-token", "your-account-slug")
//	boards, err := client.Boards.List(context.Background())
//	if err != nil {
//		log.Fatal(err)
//	}
//
// With options:
//
//	client := fizzy.NewClient(token, accountSlug,
//		fizzy.WithTimeout(60*time.Second),
//		fizzy.WithCache(false),
//	)
package fizzy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultBaseURL    = "https://app.fizzy.do"
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 3
)

// Client is the main entrypoint for interacting with the Fizzy API.
type Client struct {
	token        string
	accountSlug  string
	baseURL      string
	timeout      time.Duration
	maxRetries   int
	cacheEnabled bool
	httpClient   *http.Client
	cache        *etagCache

	// Services
	Identity      *IdentityService
	Boards        *BoardsService
	Cards         *CardsService
	Comments      *CommentsService
	Reactions     *ReactionsService
	Steps         *StepsService
	Tags          *TagsService
	Columns       *ColumnsService
	Users         *UsersService
	Notifications *NotificationsService
	Uploads       *UploadsService
}

// NewClient creates a new Fizzy API client.
//
// The token parameter is your Personal Access Token from Fizzy.
// The accountSlug is your account identifier (e.g., "my-company").
//
// Optional configuration can be provided using ClientOption functions:
//
//	client := NewClient(token, accountSlug,
//		WithTimeout(60*time.Second),
//		WithCache(false),
//		WithMaxRetries(5),
//	)
func NewClient(token, accountSlug string, opts ...ClientOption) *Client {
	c := &Client{
		token:        token,
		accountSlug:  accountSlug,
		baseURL:      defaultBaseURL,
		timeout:      defaultTimeout,
		maxRetries:   defaultMaxRetries,
		cacheEnabled: true,
		cache:        newETagCache(),
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Build HTTP client with middleware
	transport := http.DefaultTransport

	// Add retry logic
	transport = &retryRoundTripper{
		next:       transport,
		maxRetries: c.maxRetries,
	}

	// Add caching if enabled
	if c.cacheEnabled {
		transport = &cacheRoundTripper{
			next:  transport,
			cache: c.cache,
		}
	}

	c.httpClient = &http.Client{
		Timeout:   c.timeout,
		Transport: transport,
	}

	// Initialize services
	c.Identity = &IdentityService{client: c}
	c.Boards = &BoardsService{client: c}
	c.Cards = &CardsService{client: c}
	c.Comments = &CommentsService{client: c}
	c.Reactions = &ReactionsService{client: c}
	c.Steps = &StepsService{client: c}
	c.Tags = &TagsService{client: c}
	c.Columns = &ColumnsService{client: c}
	c.Users = &UsersService{client: c}
	c.Notifications = &NotificationsService{client: c}
	c.Uploads = &UploadsService{client: c}

	return c
}

// doRequest performs an HTTP request and handles common response processing.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	// Build URL - concatenate base URL with path (which may include query string)
	fullURL := c.baseURL + path

	// Encode body if provided
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// For retry logic to work with body, we need GetBody
	if bodyReader != nil {
		bodyBytes, _ := json.Marshal(body)
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyBytes)), nil
		}
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle 304 Not Modified (cached response)
	if resp.StatusCode == http.StatusNotModified {
		// No new data, result stays unchanged
		return nil
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		message := string(bodyBytes)
		if message == "" {
			message = http.StatusText(resp.StatusCode)
		}
		requestID := resp.Header.Get("X-Request-Id")
		return newErrorFromResponse(resp.StatusCode, message, requestID)
	}

	// Handle 201 Created with Location header (no body)
	if resp.StatusCode == http.StatusCreated {
		location := resp.Header.Get("Location")
		if location != "" && result != nil {
			// For 201 Created with Location, need to GET the resource
			// This is a follow-up request to fetch the created resource
			return c.doRequest(ctx, "GET", location, nil, result)
		}
	}

	// Decode response if result is provided
	if result != nil && resp.StatusCode != http.StatusNoContent {
		// Read body first to check if it's empty
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}
		
		// Only decode if body is non-empty
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, result); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
		}
		// If body is empty, result stays zero-valued
	}

	return nil
}
