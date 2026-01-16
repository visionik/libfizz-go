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

// Identity Service Tests
func TestIdentityService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/my/identity", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode(Identity{
			Accounts: []Account{
				{
					ID:   "123",
					Name: "Test Account",
					User: &User{
						ID:   "user1",
						Name: "Test User",
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	identity, err := client.Identity.Get(context.Background())

	require.NoError(t, err)
	assert.Len(t, identity.Accounts, 1)
	assert.Equal(t, "123", identity.Accounts[0].ID)
	assert.Equal(t, "Test Account", identity.Accounts[0].Name)
}

// Tags Service Tests
func TestTagsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/tags", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode([]Tag{
			{ID: "1", Name: "bug", Color: "red"},
			{ID: "2", Name: "feature", Color: "blue"},
		})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	tags, err := client.Tags.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, "bug", tags[0].Name)
	assert.Equal(t, "feature", tags[1].Name)
}

func TestTagsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/tags", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var opts TagCreateOptions
		require.NoError(t, json.NewDecoder(r.Body).Decode(&opts))
		assert.Equal(t, "new-tag", opts.Name)

		json.NewEncoder(w).Encode(Tag{ID: "456", Name: opts.Name})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	tag, err := client.Tags.Create(context.Background(), &TagCreateOptions{
		Name: "new-tag",
	})

	require.NoError(t, err)
	assert.Equal(t, "456", tag.ID)
	assert.Equal(t, "new-tag", tag.Name)
}

// Comments Service Tests
func TestCommentsService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/comments", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode([]Comment{
			{ID: "1", Body: "Comment 1"},
			{ID: "2", Body: "Comment 2"},
		})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	comments, err := client.Comments.List(context.Background(), "card-123")

	require.NoError(t, err)
	assert.Len(t, comments, 2)
	assert.Equal(t, "Comment 1", comments[0].Body)
	assert.Equal(t, "Comment 2", comments[1].Body)
}

func TestCommentsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/comments/456", r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)
		
		// Verify wrapped request body
		var payload map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Contains(t, payload, "comment")
		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Comment{ID: "456", Body: "Updated"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	comment, err := client.Comments.Update(context.Background(), "card-123", "456", &CommentUpdateOptions{
		Body: "Updated",
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated", comment.Body)
}

func TestCommentsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/comments/456", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Comments.Delete(context.Background(), "card-123", "456")

	require.NoError(t, err)
}

func TestCommentsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/card-123/comments", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Verify wrapped request body
		var payload map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Contains(t, payload, "comment")
		commentData := payload["comment"].(map[string]interface{})
		assert.Equal(t, "New comment", commentData["body"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Comment{ID: "456", Body: "New comment"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	comment, err := client.Comments.Create(context.Background(), "card-123", &CommentCreateOptions{
		Body: "New comment",
	})

	require.NoError(t, err)
	assert.Equal(t, "456", comment.ID)
	assert.Equal(t, "New comment", comment.Body)
}


// Users Service Tests
func TestUsersService_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/users", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode([]User{
			{ID: "1", Name: "User 1"},
			{ID: "2", Name: "User 2"},
		})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	users, err := client.Users.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "User 1", users[0].Name)
	assert.Equal(t, "User 2", users[1].Name)
}

