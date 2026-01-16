package fizzy

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test error scenarios
func TestClient_BadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad request"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.doRequest(context.Background(), "POST", "/test", map[string]string{"bad": "data"}, nil)

	require.Error(t, err)
	var badReqErr *BadRequestError
	assert.ErrorAs(t, err, &badReqErr)
}

func TestClient_Forbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.doRequest(context.Background(), "GET", "/forbidden", nil, nil)

	require.Error(t, err)
	var forbiddenErr *ForbiddenError
	assert.ErrorAs(t, err, &forbiddenErr)
}

func TestClient_RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]string{"error": "Rate limited"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL), WithMaxRetries(0))
	err := client.doRequest(context.Background(), "GET", "/test", nil, nil)

	require.Error(t, err)
	var rateLimitErr *RateLimitError
	assert.ErrorAs(t, err, &rateLimitErr)
}

func TestClient_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Internal error"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL), WithMaxRetries(0))
	err := client.doRequest(context.Background(), "GET", "/test", nil, nil)

	require.Error(t, err)
	var serverErr *ServerError
	assert.ErrorAs(t, err, &serverErr)
}

// Test GET request methods
func TestBoardsService_GetNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	_, err := client.Boards.Get(context.Background(), "nonexistent")

	require.Error(t, err)
	var notFoundErr *NotFoundError
	assert.ErrorAs(t, err, &notFoundErr)
}

// Test update methods
func TestBoardsService_UpdateName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/boards/123", r.URL.Path)
		
		if r.Method == "PATCH" {
			// Verify wrapped request body
			var payload map[string]interface{}
			json.NewDecoder(r.Body).Decode(&payload)
			assert.Contains(t, payload, "board")
			board := payload["board"].(map[string]interface{})
			assert.NotNil(t, board["name"])
			assert.Equal(t, "Updated Name", board["name"])
			
			// Return 204 No Content
			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == "GET" {
			// Follow-up GET after 204
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Board{ID: "123", Name: "Updated Name"})
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	name := "Updated Name"
	board, err := client.Boards.Update(context.Background(), "123", &BoardUpdateOptions{
		Name: &name,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Name", board.Name)
}

// Test delete methods
func TestCardsService_DeleteError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot delete"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Delete(context.Background(), "123")

	require.Error(t, err)
}

// Test list methods for various services
func TestBoardsService_ListSuccess(t *testing.T) {
	boards := []Board{{ID: "1", Name: "Board"}, {ID: "2", Name: "Another"}}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(boards)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	result, err := client.Boards.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCardsService_ListWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.URL.RawQuery, "board_id=board1")
		assert.Contains(t, r.URL.RawQuery, "status=open")
		json.NewEncoder(w).Encode([]Card{{ID: "1", Title: "Card"}})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	result, err := client.Cards.List(context.Background(), &CardListOptions{
		BoardID: "board1",
		Status:  "open",
	})

	require.NoError(t, err)
	assert.Len(t, result.Items, 1)
}

