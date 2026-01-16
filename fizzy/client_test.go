package fizzy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-token", "test-account")

	assert.NotNil(t, client)
	assert.Equal(t, "test-account", client.accountSlug)
	assert.Equal(t, defaultBaseURL, client.baseURL)
	assert.Equal(t, defaultTimeout, client.httpClient.Timeout)
	assert.NotNil(t, client.Boards)
	assert.NotNil(t, client.Cards)
	assert.NotNil(t, client.Identity)
}

func TestWithBaseURL(t *testing.T) {
	customURL := "https://custom.fizzy.do"
	client := NewClient("test-token", "test-account", WithBaseURL(customURL))

	assert.Equal(t, customURL, client.baseURL)
}

func TestWithTimeout(t *testing.T) {
	timeout := 5 * time.Second
	client := NewClient("test-token", "test-account", WithTimeout(timeout))

	assert.Equal(t, timeout, client.httpClient.Timeout)
}

func TestWithMaxRetries(t *testing.T) {
	client := NewClient("test-token", "test-account", WithMaxRetries(5))

	// We can't easily check the maxRetries value directly, but we can verify the client was created
	assert.NotNil(t, client)
}

func TestWithCache(t *testing.T) {
	client := NewClient("test-token", "test-account", WithCache(false))

	assert.NotNil(t, client)
	assert.False(t, client.cacheEnabled)

	client2 := NewClient("test-token", "test-account", WithCache(true))
	assert.True(t, client2.cacheEnabled)
}

func TestClient_Authentication(t *testing.T) {
	var capturedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient("my-secret-token", "test-account", WithBaseURL(server.URL))
	client.doRequest(context.Background(), "GET", "/test", nil, nil)

	assert.Equal(t, "Bearer my-secret-token", capturedAuth)
}

func TestClient_UserAgent(t *testing.T) {
	var capturedUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	client.doRequest(context.Background(), "GET", "/test", nil, nil)

	// Just verify a user agent was set
	assert.NotEmpty(t, capturedUA)
}

func TestClient_ContentType(t *testing.T) {
	var capturedContentType string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedContentType = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	client.doRequest(context.Background(), "POST", "/test", map[string]string{"key": "value"}, nil)

	assert.Equal(t, "application/json", capturedContentType)
}

func TestClient_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.doRequest(context.Background(), "GET", "/nonexistent", nil, nil)

	require.Error(t, err)
	var notFoundErr *NotFoundError
	if assert.ErrorAs(t, err, &notFoundErr) {
		assert.Equal(t, 404, notFoundErr.StatusCode)
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := client.doRequest(ctx, "GET", "/test", nil, nil)
	assert.Error(t, err)
}

func TestListAll(t *testing.T) {
	// Mock a paginated API
	page1Items := []Board{{ID: "1", Name: "Board 1"}, {ID: "2", Name: "Board 2"}}
	page2Items := []Board{{ID: "3", Name: "Board 3"}}

	currentPage := 1
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if currentPage == 1 {
			json.NewEncoder(w).Encode(page1Items)
			currentPage++
		} else {
			json.NewEncoder(w).Encode(page2Items)
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))

	// Test listAll with a fetcher that simulates pagination
	fetcher := func(ctx context.Context, page string) (*PagedResult[Board], error) {
		path := "/test-account/boards"
		var boards []Board
		if err := client.doRequest(ctx, "GET", path, nil, &boards); err != nil {
			return nil, err
		}

		nextPage := ""
		if len(boards) == 2 {
			nextPage = "page2"
		}

		return &PagedResult[Board]{
			Items:    boards,
			NextPage: nextPage,
		}, nil
	}

	ctx := context.Background()
	result, err := listAll(ctx, fetcher)

	require.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "Board 1", result[0].Name)
	assert.Equal(t, "Board 3", result[2].Name)
}

func TestClient_JSONDecoding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Test"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))

	var result struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	err := client.doRequest(context.Background(), "GET", "/test", nil, &result)

	require.NoError(t, err)
	assert.Equal(t, "123", result.ID)
	assert.Equal(t, "Test", result.Name)
}

func TestClient_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.doRequest(context.Background(), "DELETE", "/test", nil, nil)

	require.NoError(t, err)
}
