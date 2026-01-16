package fizzy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestETagCache(t *testing.T) {
	cache := newETagCache()

	t.Run("get empty cache", func(t *testing.T) {
		_, ok := cache.get("http://example.com")
		assert.False(t, ok)
	})

	t.Run("set and get", func(t *testing.T) {
		cache.set("http://example.com", "etag-123")
		etag, ok := cache.get("http://example.com")
		assert.True(t, ok)
		assert.Equal(t, "etag-123", etag)
	})

	t.Run("clear cache", func(t *testing.T) {
		cache.set("http://example.com", "etag-123")
		cache.clear()
		_, ok := cache.get("http://example.com")
		assert.False(t, ok)
	})
}

func TestCacheRoundTripper(t *testing.T) {
	t.Run("caches GET request ETags", func(t *testing.T) {
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.Header().Set("ETag", "etag-abc")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"result":"ok"}`))
		}))
		defer server.Close()

		cache := newETagCache()
		rt := &cacheRoundTripper{
			next:  http.DefaultTransport,
			cache: cache,
		}

		client := &http.Client{Transport: rt}

		// First request - no ETag header
		req, _ := http.NewRequest("GET", server.URL, nil)
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Equal(t, 1, requestCount)
		etag, ok := cache.get(server.URL)
		assert.True(t, ok)
		assert.Equal(t, "etag-abc", etag)

		// Second request - should include If-None-Match
		req2, _ := http.NewRequest("GET", server.URL, nil)
		resp2, err := client.Do(req2)
		require.NoError(t, err)
		resp2.Body.Close()

		assert.Equal(t, 2, requestCount)
	})

	t.Run("does not cache non-GET requests", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("ETag", "etag-post")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		cache := newETagCache()
		rt := &cacheRoundTripper{
			next:  http.DefaultTransport,
			cache: cache,
		}

		client := &http.Client{Transport: rt}

		req, _ := http.NewRequest("POST", server.URL, nil)
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()

		_, ok := cache.get(server.URL)
		assert.False(t, ok, "POST requests should not be cached")
	})

	t.Run("adds If-None-Match header on cached URL", func(t *testing.T) {
		receivedHeader := ""
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receivedHeader = r.Header.Get("If-None-Match")
			w.Header().Set("ETag", "etag-xyz")
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		cache := newETagCache()
		cache.set(server.URL, "cached-etag")

		rt := &cacheRoundTripper{
			next:  http.DefaultTransport,
			cache: cache,
		}

		client := &http.Client{Transport: rt}

		req, _ := http.NewRequest("GET", server.URL, nil)
		resp, err := client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()

		assert.Equal(t, "cached-etag", receivedHeader)
	})
}
