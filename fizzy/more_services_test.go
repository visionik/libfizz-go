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

// Steps Service Tests
func TestStepsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/steps", r.URL.Path)
		json.NewEncoder(w).Encode([]Step{
			{ID: "1", Content: "Step 1"},
			{ID: "2", Content: "Step 2"},
		})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	steps, err := client.Steps.List(context.Background(), "card-123")

	require.NoError(t, err)
	assert.Len(t, steps, 2)
}

func TestStepsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/steps", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		
		// Verify wrapped request body
		var payload map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Contains(t, payload, "step")
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Step{ID: "456", Content: "New Step"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	step, err := client.Steps.Create(context.Background(), "card-123", &StepCreateOptions{
		Content: "New Step",
	})

	require.NoError(t, err)
	assert.Equal(t, "456", step.ID)
}

func TestStepsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/steps/step-456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode(Step{ID: "step-456", Content: "Test Step"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	step, err := client.Steps.Get(context.Background(), "card-123", "step-456")

	require.NoError(t, err)
	assert.Equal(t, "step-456", step.ID)
}

func TestStepsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/steps/123", r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)
		
		// Verify wrapped request body
		var payload map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Contains(t, payload, "step")
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Step{ID: "123", Content: "Updated"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	content := "Updated"
	step, err := client.Steps.Update(context.Background(), "card-123", "123", &StepUpdateOptions{
		Content: &content,
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated", step.Content)
}

func TestStepsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/steps/123", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Steps.Delete(context.Background(), "card-123", "123")

	require.NoError(t, err)
}

// Columns Service Tests
func TestColumnsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/boards/board-123/columns", r.URL.Path)
		json.NewEncoder(w).Encode([]Column{
			{ID: "1", Name: "To Do"},
			{ID: "2", Name: "Done"},
		})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	columns, err := client.Columns.List(context.Background(), "board-123")

	require.NoError(t, err)
	assert.Len(t, columns, 2)
}

func TestColumnsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/boards/board-123/columns", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		
		// Verify wrapped request body
		var payload map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Contains(t, payload, "column")
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Column{ID: "456", Name: "In Progress"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	column, err := client.Columns.Create(context.Background(), "board-123", &ColumnCreateOptions{
		Name: "In Progress",
	})

	require.NoError(t, err)
	assert.Equal(t, "456", column.ID)
}

func TestColumnsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/boards/board-123/columns/col-456", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode(Column{ID: "col-456", Name: "In Progress"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	column, err := client.Columns.Get(context.Background(), "board-123", "col-456")

	require.NoError(t, err)
	assert.Equal(t, "col-456", column.ID)
	assert.Equal(t, "In Progress", column.Name)
}

func TestColumnsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/boards/board-123/columns/col-456", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Columns.Delete(context.Background(), "board-123", "col-456")

	require.NoError(t, err)
}

func TestColumnsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/boards/board-123/columns/col-123", r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)
		
		// Verify wrapped request body
		var payload map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Contains(t, payload, "column")
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Column{ID: "col-123", Name: "Updated", Position: 2})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	position := 2
	column, err := client.Columns.Update(context.Background(), "board-123", "col-123", &ColumnUpdateOptions{
		Position: &position,
	})

	require.NoError(t, err)
	assert.Equal(t, 2, column.Position)
}

// Reactions Service Tests
func TestReactionsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/comments/comment-456/reactions", r.URL.Path)
		json.NewEncoder(w).Encode([]Reaction{
			{ID: "1", Content: "👍"},
			{ID: "2", Content: "❤️"},
		})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	reactions, err := client.Reactions.List(context.Background(), "card-123", "comment-456")

	require.NoError(t, err)
	assert.Len(t, reactions, 2)
}

func TestReactionsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/comments/comment-456/reactions", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		
		// Verify wrapped request body
		var payload map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Contains(t, payload, "reaction")
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Reaction{ID: "789", Content: "🎉"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	reaction, err := client.Reactions.Create(context.Background(), "card-123", "comment-456", &ReactionCreateOptions{
		Content: "🎉",
	})

	require.NoError(t, err)
	assert.Equal(t, "789", reaction.ID)
}

func TestReactionsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/comments/comment-456/reactions/123", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Reactions.Delete(context.Background(), "card-123", "comment-456", "123")

	require.NoError(t, err)
}

// Notifications Service Tests
func TestNotificationsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/my/notifications", r.URL.Path)
		json.NewEncoder(w).Encode([]Notification{
			{ID: "1", Type: "mention"},
			{ID: "2", Type: "assignment"},
		})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	notifications, err := client.Notifications.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, notifications, 2)
}

func TestNotificationsService_Read(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/my/notifications/123/read", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Notifications.Read(context.Background(), "123")

	require.NoError(t, err)
}

func TestNotificationsService_Unread(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/my/notifications/123/unread", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Notifications.Unread(context.Background(), "123")

	require.NoError(t, err)
}

func TestNotificationsService_ReadAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/my/notifications/read_all", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Notifications.ReadAll(context.Background())

	require.NoError(t, err)
}

// Uploads Service Tests
func TestUploadsService_CreateDirectUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/rails/active_storage/direct_uploads", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		json.NewEncoder(w).Encode(DirectUploadResponse{
			DirectUploadURL: "https://s3.example.com/upload",
			BlobID:          "blob-123",
			Headers:         map[string]string{"Content-Type": "image/png"},
		})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	resp, err := client.Uploads.CreateDirectUpload(context.Background(), &DirectUploadRequest{
		Blob: DirectUploadBlob{
			Filename:    "test.png",
			ByteSize:    1024,
			Checksum:    "abc123",
			ContentType: "image/png",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "blob-123", resp.BlobID)
	assert.Equal(t, "https://s3.example.com/upload", resp.DirectUploadURL)
}
