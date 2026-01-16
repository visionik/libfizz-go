//go:build integration

package fizzy

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests require real Fizzy API credentials
// Run with: go test -tags=integration -v ./fizzy/...
// Required environment variables:
//   FIZZY_TOKEN - Your Personal Access Token
//   FIZZY_ACCOUNT - Your account slug

func getTestClient(t *testing.T) *Client {
	token := os.Getenv("FIZZY_TOKEN")
	account := os.Getenv("FIZZY_ACCOUNT")

	if token == "" || account == "" {
		t.Skip("Skipping integration test: FIZZY_TOKEN and FIZZY_ACCOUNT must be set")
	}

	return NewClient(token, account, WithTimeout(30*time.Second))
}

func TestIntegration_Identity(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("Get identity", func(t *testing.T) {
		identity, err := client.Identity.Get(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, identity.Accounts)
		
		if len(identity.Accounts) > 0 {
			account := identity.Accounts[0]
			assert.NotEmpty(t, account.ID)
			assert.NotEmpty(t, account.Name)
			if account.User != nil {
				t.Logf("Authenticated as: %s (%s)", account.User.Name, *account.User.EmailAddress)
				t.Logf("Account: %s", account.Name)
			}
		}
	})
}

func TestIntegration_Boards(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("List boards", func(t *testing.T) {
		boards, err := client.Boards.List(ctx)
		require.NoError(t, err)
		t.Logf("Found %d boards", len(boards))

		if len(boards) > 0 {
			t.Logf("First board: %s (ID: %s)", boards[0].Name, boards[0].ID)
		}
	})

	var createdBoardID string

	t.Run("Create board", func(t *testing.T) {
		desc := "Test board created by integration tests"
		board, err := client.Boards.Create(ctx, &BoardCreateOptions{
			Name:        "Integration Test Board",
			Description: &desc,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, board.ID)
		assert.Equal(t, "Integration Test Board", board.Name)
		createdBoardID = board.ID
		t.Logf("Created board: %s", board.ID)
	})

	t.Run("Get board", func(t *testing.T) {
		if createdBoardID == "" {
			t.Skip("No board created")
		}

		board, err := client.Boards.Get(ctx, createdBoardID)
		require.NoError(t, err)
		assert.Equal(t, createdBoardID, board.ID)
		assert.Equal(t, "Integration Test Board", board.Name)
	})

	t.Run("Update board", func(t *testing.T) {
		if createdBoardID == "" {
			t.Skip("No board created")
		}

		newName := "Updated Integration Test Board"
		board, err := client.Boards.Update(ctx, createdBoardID, &BoardUpdateOptions{
			Name: &newName,
		})
		require.NoError(t, err)
		assert.Equal(t, newName, board.Name)
	})

	t.Run("Delete board", func(t *testing.T) {
		if createdBoardID == "" {
			t.Skip("No board created")
		}

		err := client.Boards.Delete(ctx, createdBoardID)
		require.NoError(t, err)
		t.Logf("Deleted board: %s", createdBoardID)
	})
}

func TestIntegration_Cards(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Get a board to work with
	boards, err := client.Boards.List(ctx)
	require.NoError(t, err)
	if len(boards) == 0 {
		t.Skip("No boards available for card tests")
	}
	boardID := boards[0].ID
	t.Logf("Using board: %s (%s)", boards[0].Name, boardID)

	var createdCardNumber string

	t.Run("Create card", func(t *testing.T) {
		body := "This is a test card created by integration tests"
		card, err := client.Cards.Create(ctx, &CardCreateOptions{
			BoardID: boardID,
			Title:   "Integration Test Card",
			Body:    &body,
		})
		require.NoError(t, err)
		assert.NotEmpty(t, card.ID)
		assert.NotZero(t, card.Number)
		assert.Equal(t, "Integration Test Card", card.Title)
		createdCardNumber = fmt.Sprintf("%d", card.Number)
		t.Logf("Created card: %s (number: %d)", card.ID, card.Number)
	})

	t.Run("Get card", func(t *testing.T) {
		if createdCardNumber == "" {
			t.Skip("No card created")
		}

		card, err := client.Cards.Get(ctx, createdCardNumber)
		require.NoError(t, err)
		assert.Equal(t, "Integration Test Card", card.Title)
	})

	t.Run("List cards", func(t *testing.T) {
		result, err := client.Cards.List(ctx, &CardListOptions{
			BoardID: boardID,
			Status:  "open",
		})
		require.NoError(t, err)
		assert.NotNil(t, result.Items)
		t.Logf("Found %d open cards", len(result.Items))
	})

	t.Run("Update card", func(t *testing.T) {
		if createdCardNumber == "" {
			t.Skip("No card created")
		}

		newTitle := "Updated Integration Test Card"
		card, err := client.Cards.Update(ctx, createdCardNumber, &CardUpdateOptions{
			Title: &newTitle,
		})
		require.NoError(t, err)
		assert.Equal(t, newTitle, card.Title)
	})

	t.Run("Close card", func(t *testing.T) {
		if createdCardNumber == "" {
			t.Skip("No card created")
		}

		err := client.Cards.Close(ctx, createdCardNumber)
		require.NoError(t, err)
		t.Logf("Closed card: %s", createdCardNumber)
	})

	t.Run("Reopen card", func(t *testing.T) {
		if createdCardNumber == "" {
			t.Skip("No card created")
		}

		err := client.Cards.Reopen(ctx, createdCardNumber)
		require.NoError(t, err)
		t.Logf("Reopened card: %s", createdCardNumber)
	})

	t.Run("Delete card", func(t *testing.T) {
		if createdCardNumber == "" {
			t.Skip("No card created")
		}

		err := client.Cards.Delete(ctx, createdCardNumber)
		require.NoError(t, err)
		t.Logf("Deleted card: %s", createdCardNumber)
	})
}

func TestIntegration_Tags(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("List tags", func(t *testing.T) {
		tags, err := client.Tags.List(ctx)
		require.NoError(t, err)
		t.Logf("Found %d tags", len(tags))
	})
}

func TestIntegration_Users(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("List users", func(t *testing.T) {
		users, err := client.Users.List(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, users)
		t.Logf("Found %d users", len(users))
	})
}

func TestIntegration_Notifications(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("List notifications", func(t *testing.T) {
		notifications, err := client.Notifications.List(ctx)
		require.NoError(t, err)
		t.Logf("Found %d notifications", len(notifications))
	})
}

func TestIntegration_ErrorHandling(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("Not found error", func(t *testing.T) {
		_, err := client.Cards.Get(ctx, "invalid-card-id-12345")
		require.Error(t, err)
		
		var notFound *NotFoundError
		assert.ErrorAs(t, err, &notFound)
		t.Logf("Got expected NotFoundError: %v", err)
	})

	t.Run("Invalid board ID", func(t *testing.T) {
		_, err := client.Boards.Get(ctx, "invalid-board-id-67890")
		require.Error(t, err)
		t.Logf("Got error for invalid board: %v", err)
	})
}
