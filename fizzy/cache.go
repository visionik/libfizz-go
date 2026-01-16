package fizzy

import (
	"net/http"
	"sync"
)

// etagCache stores ETags for URLs to support HTTP caching.
type etagCache struct {
	mu    sync.RWMutex
	store map[string]string // URL -> ETag
}

// newETagCache creates a new ETag cache.
func newETagCache() *etagCache {
	return &etagCache{
		store: make(map[string]string),
	}
}

// get retrieves an ETag for a URL.
func (c *etagCache) get(url string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	etag, ok := c.store[url]
	return etag, ok
}

// set stores an ETag for a URL.
func (c *etagCache) set(url, etag string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[url] = etag
}

// clear removes all cached ETags.
func (c *etagCache) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]string)
}

// cacheRoundTripper is an HTTP RoundTripper that implements ETag caching.
type cacheRoundTripper struct {
	next  http.RoundTripper
	cache *etagCache
}

// RoundTrip implements http.RoundTripper with ETag caching support.
func (t *cacheRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only cache GET requests
	if req.Method != http.MethodGet {
		return t.next.RoundTrip(req)
	}

	url := req.URL.String()

	// Add If-None-Match header if we have a cached ETag
	if etag, ok := t.cache.get(url); ok {
		req.Header.Set("If-None-Match", etag)
	}

	resp, err := t.next.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	// Store ETag from response
	if etag := resp.Header.Get("ETag"); etag != "" {
		t.cache.set(url, etag)
	}

	return resp, nil
}
