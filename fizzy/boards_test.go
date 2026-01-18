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

func TestBoardsService_List(t *testing.T) {
	tests := []struct {
		name         string
		responseCode int
		responseBody string
		wantErr      bool
		errType      error
	}{
		{
			"success",
			http.StatusOK,
			`[{"id":"board-1","name":"Engineering"},{"id":"board-2","name":"Design"}]`,
			false,
			nil,
		},
		{
			"unauthorized",
			http.StatusUnauthorized,
			`{"error":"Unauthorized"}`,
			true,
			&AuthenticationError{},
		},
		{
			"server error",
			http.StatusInternalServerError,
			`{"error":"Internal server error"}`,
			true,
			&ServerError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/test-account/boards.json", r.URL.Path)
				assert.Equal(t, "GET", r.Method)
				assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

				w.WriteHeader(tt.responseCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
			boards, err := client.Boards.List(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorAs(t, err, &tt.errType)
				}
				return
			}

			require.NoError(t, err)
			assert.Len(t, boards, 2)
			assert.Equal(t, "board-1", boards[0].ID)
			assert.Equal(t, "Engineering", boards[0].Name)
		})
	}
}

func TestBoardsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/boards/board-123.json", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"board-123","name":"Engineering","description":"Dev tasks"}`))
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	board, err := client.Boards.Get(context.Background(), "board-123")

	require.NoError(t, err)
	assert.Equal(t, "board-123", board.ID)
	assert.Equal(t, "Engineering", board.Name)
	assert.NotNil(t, board.Description)
	assert.Equal(t, "Dev tasks", *board.Description)
}

func TestBoardsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test-account/boards" {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			
			// Verify wrapped request body
			var payload map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			assert.Contains(t, payload, "board")
			board := payload["board"].(map[string]interface{})
			assert.Equal(t, "New Board", board["name"])
			
			// Return 201 with board data
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"id":"new-board","name":"New Board"}`))
		} else if r.URL.Path == "/test-account/boards/new-board.json" {
			// Follow-up GET request
			assert.Equal(t, "GET", r.Method)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"new-board","name":"New Board"}`))
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	
	desc := "Board description"
	board, err := client.Boards.Create(context.Background(), &BoardCreateOptions{
		Name:        "New Board",
		Description: &desc,
	})

	require.NoError(t, err)
	assert.Equal(t, "new-board", board.ID)
	assert.Equal(t, "New Board", board.Name)
}

func TestBoardsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			assert.Equal(t, "/test-account/boards/board-123", r.URL.Path)
			
			// Verify wrapped request body
			var payload map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			assert.Contains(t, payload, "board")
			
			// Return 204 No Content
			w.WriteHeader(http.StatusNoContent)
		} else if r.Method == "GET" {
			// Follow-up GET request after 204
			assert.Equal(t, "/test-account/boards/board-123.json", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id":"board-123","name":"Updated Name"}`))
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	
	newName := "Updated Name"
	board, err := client.Boards.Update(context.Background(), "board-123", &BoardUpdateOptions{
		Name: &newName,
	})

	require.NoError(t, err)
	assert.Equal(t, "board-123", board.ID)
	assert.Equal(t, "Updated Name", board.Name)
}

func TestBoardsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/boards/board-123", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Boards.Delete(context.Background(), "board-123")

	assert.NoError(t, err)
}
