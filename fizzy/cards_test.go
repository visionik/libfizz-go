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

func TestCardsService_List(t *testing.T) {
	tests := []struct {
		name       string
		response   interface{}
		statusCode int
		wantErr    bool
	}{
		{
			name: "success",
			response: []Card{
				{ID: "1", Title: "Card 1", Status: "open"},
				{ID: "2", Title: "Card 2", Status: "closed"},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "unauthorized",
			response:   map[string]string{"error": "Unauthorized"},
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/test-account/cards", r.URL.Path)
				assert.Equal(t, "GET", r.Method)
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
			ctx := context.Background()

			result, err := client.Cards.List(ctx, &CardListOptions{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, result.Items, 2)
			assert.Equal(t, "Card 1", result.Items[0].Title)
		})
	}
}

func TestCardsService_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		json.NewEncoder(w).Encode(Card{ID: "123", Title: "Test Card"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	card, err := client.Cards.Get(context.Background(), "123")

	require.NoError(t, err)
	assert.Equal(t, "123", card.ID)
	assert.Equal(t, "Test Card", card.Title)
}

func TestCardsService_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test-account/boards/board123/cards" {
			assert.Equal(t, "POST", r.Method)
			
			// Verify wrapped request body
			var payload map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			assert.Contains(t, payload, "card")
			cardData := payload["card"].(map[string]interface{})
			assert.Equal(t, "New Card", cardData["title"])
			
			// Return 201 with Location header
			w.Header().Set("Location", "/test-account/cards/456")
			w.WriteHeader(http.StatusCreated)
		} else if r.URL.Path == "/test-account/cards/456" {
			// Follow-up GET request
			assert.Equal(t, "GET", r.Method)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Card{ID: "456", Number: 1, Title: "New Card"})
		}
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	card, err := client.Cards.Create(context.Background(), &CardCreateOptions{
		Title:   "New Card",
		BoardID: "board123",
	})

	require.NoError(t, err)
	assert.Equal(t, "456", card.ID)
	assert.Equal(t, "New Card", card.Title)
}

func TestCardsService_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123", r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)

		// Verify wrapped request body
		var payload map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Contains(t, payload, "card")
		cardData := payload["card"].(map[string]interface{})
		assert.Equal(t, "Updated Title", cardData["title"])

		// API returns 200 OK with body for PATCH
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Card{ID: "123", Title: "Updated Title"})
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	title := "Updated Title"
	card, err := client.Cards.Update(context.Background(), "123", &CardUpdateOptions{
		Title: &title,
	})

	require.NoError(t, err)
	assert.Equal(t, "123", card.ID)
	assert.Equal(t, "Updated Title", card.Title)
}

func TestCardsService_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Delete(context.Background(), "123")

	require.NoError(t, err)
}

func TestCardsService_Close(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/closure", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Close(context.Background(), "123")

	require.NoError(t, err)
}

func TestCardsService_Tag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/tags/bug/toggle", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Tag(context.Background(), "123", "bug")

	require.NoError(t, err)
}

func TestCardsService_Reopen(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/closure", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Reopen(context.Background(), "123")

	require.NoError(t, err)
}

func TestCardsService_Postpone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/not_now", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Postpone(context.Background(), "123")

	require.NoError(t, err)
}

func TestCardsService_Triage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/triage", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Triage(context.Background(), "123")

	require.NoError(t, err)
}

func TestCardsService_MoveToColumn(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/column", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.MoveToColumn(context.Background(), "123", "col456")

	require.NoError(t, err)
}

func TestCardsService_Assign(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/assignments/user456/toggle", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Assign(context.Background(), "123", "user456")

	require.NoError(t, err)
}

func TestCardsService_Watch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/watch", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Watch(context.Background(), "123")

	require.NoError(t, err)
}

func TestCardsService_Unwatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/unwatch", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.Unwatch(context.Background(), "123")

	require.NoError(t, err)
}

func TestCardsService_MarkGolden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/golden", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.MarkGolden(context.Background(), "123")

	require.NoError(t, err)
}

func TestCardsService_ListAll(t *testing.T) {
	cards := []Card{{ID: "1", Title: "Card 1"}, {ID: "2", Title: "Card 2"}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(cards)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	result, err := client.Cards.ListAll(context.Background(), &CardListOptions{})

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCardsService_UnmarkGolden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-account/cards/123/ungolden", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token", "test-account", WithBaseURL(server.URL))
	err := client.Cards.UnmarkGolden(context.Background(), "123")

	require.NoError(t, err)
}
