# libfizz-go

A simple, idiomatic Go client library for the [Fizzy API](https://fizzy.do).

[![Go Reference](https://pkg.go.dev/badge/github.com/visionik/libfizz-go/fizzy.svg)](https://pkg.go.dev/github.com/visionik/libfizz-go/fizzy)

## Features

- **Complete API Coverage**: All Fizzy API resources (Identity, Boards, Cards, Comments, Reactions, Steps, Tags, Columns, Users, Notifications, Uploads)
- **Idiomatic Go**: Follows Go best practices with context support, functional options, and typed errors
- **Automatic Retries**: Exponential backoff for rate limits (429) and server errors (5xx)
- **ETag Caching**: Built-in HTTP caching to reduce bandwidth and improve performance
- **Thread-Safe**: Client can be safely shared across goroutines
- **Comprehensive Testing**: 86.6% test coverage with 160 unit tests + integration tests
- **Production Ready**: All write operations tested against live Fizzy API

## Installation

```bash
go get github.com/visionik/libfizz-go/fizzy
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/visionik/libfizz-go/fizzy"
)

func main() {
	// Create client with Personal Access Token
	client := fizzy.NewClient("your-token", "your-account-slug")

	ctx := context.Background()

	// Get your identity
	identity, err := client.Identity.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(identity.Accounts) > 0 && identity.Accounts[0].User != nil {
		fmt.Printf("Logged in as: %s\n", identity.Accounts[0].User.Name)
	}

	// List boards
	boards, err := client.Boards.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, board := range boards {
		fmt.Printf("Board: %s\n", board.Name)
	}

	// Create a card
	card, err := client.Cards.Create(ctx, &fizzy.CardCreateOptions{
		BoardID: boards[0].ID,
		Title:   "New issue",
		Body:    ptrString("Description goes here"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created card: %s\n", card.Title)
}

func ptrString(s string) *string {
	return &s
}
```

## Configuration

Customize the client with functional options:

```go
import "time"

client := fizzy.NewClient(token, accountSlug,
	fizzy.WithTimeout(60*time.Second),      // Default: 30s
	fizzy.WithBaseURL("https://custom.url"), // Default: https://app.fizzy.do
	fizzy.WithCache(false),                  // Default: true (enabled)
	fizzy.WithMaxRetries(5),                 // Default: 3
)
```

## Usage Examples

### Boards

```go
// List all boards
boards, err := client.Boards.List(ctx)

// Get a specific board
board, err := client.Boards.Get(ctx, "board-id")

// Create a board
board, err := client.Boards.Create(ctx, &fizzy.BoardCreateOptions{
	Name:        "Engineering",
	Description: ptrString("Development tasks"),
})

// Update a board
board, err := client.Boards.Update(ctx, "board-id", &fizzy.BoardUpdateOptions{
	Name: ptrString("Updated Name"),
})

// Delete a board
err := client.Boards.Delete(ctx, "board-id")
```

### Cards

```go
// List cards (manual pagination)
result, err := client.Cards.List(ctx, &fizzy.CardListOptions{
	BoardID: "board-id",
	Status:  "open",
})
for _, card := range result.Items {
	fmt.Println(card.Title)
}

// Fetch next page if exists
if result.NextPage != "" {
	result, err = client.Cards.List(ctx, &fizzy.CardListOptions{
		Page: result.NextPage,
	})
}

// List all cards (automatic pagination)
cards, err := client.Cards.ListAll(ctx, &fizzy.CardListOptions{
	BoardID: "board-id",
})

// Create a card
card, err := client.Cards.Create(ctx, &fizzy.CardCreateOptions{
	BoardID: "board-id",
	Title:   "Bug: Fix login issue",
	Body:    ptrString("Users cannot log in with SSO"),
})

// Update a card (use card number, not ID)
card, err := client.Cards.Update(ctx, "card-number", &fizzy.CardUpdateOptions{
	Title: ptrString("Updated title"),
})

// Card actions
err = client.Cards.Close(ctx, "card-number")
err = client.Cards.Reopen(ctx, "card-number")
err = client.Cards.Postpone(ctx, "card-number")
err = client.Cards.Assign(ctx, "card-number", "user-id")
err = client.Cards.Tag(ctx, "card-number", "tag-name")
```

### Error Handling

The library provides typed errors for different HTTP status codes:

```go
card, err := client.Cards.Get(ctx, "invalid-id")
if err != nil {
	var notFound *fizzy.NotFoundError
	if errors.As(err, &notFound) {
		fmt.Println("Card not found")
		return
	}

	var authErr *fizzy.AuthenticationError
	if errors.As(err, &authErr) {
		fmt.Println("Authentication failed - check your token")
		return
	}

	var rateLimit *fizzy.RateLimitError
	if errors.As(err, &rateLimit) {
		fmt.Println("Rate limited - already retried 3 times")
		return
	}

	// Generic error handling
	log.Fatal(err)
}
```

## API Coverage

✅ **Fully Implemented & Tested:**

| Service | Operations | Status |
|---------|-----------|--------|
| **Identity** | Get | ✅ Complete |
| **Boards** | List, Get, Create, Update, Delete | ✅ Complete |
| **Cards** | List, ListAll, Get, Create, Update, Delete<br>Close, Reopen, Postpone, Triage<br>Assign, Tag, Watch, MoveToColumn<br>MarkGolden, UnmarkGolden | ✅ Complete (16 methods) |
| **Comments** | List, Create, Update, Delete | ✅ Complete |
| **Reactions** | List, Create, Delete | ✅ Complete |
| **Steps** | List, Get, Create, Update, Delete | ✅ Complete |
| **Tags** | List, Create | ✅ Complete |
| **Columns** | List, Get, Create, Update, Delete | ✅ Complete |
| **Users** | List | ✅ Complete |
| **Notifications** | List, Read, Unread, ReadAll | ✅ Complete |
| **Uploads** | RequestUpload, UploadFile | ✅ Complete |

**Total:** 11 services, 60+ API methods, 160 unit tests, 86.6% coverage

## Development

### Prerequisites

- Go 1.21+
- [Task](https://taskfile.dev) (optional, for running tasks)

### Building

```bash
go build ./fizzy
```

### Testing

```bash
# Run unit tests
go test ./fizzy

# Run with coverage (86.6% coverage achieved)
go test -coverprofile=coverage.out ./fizzy
go tool cover -html=coverage.out

# Run integration tests (requires FIZZY_TOKEN and FIZZY_ACCOUNT env vars)
export FIZZY_TOKEN="your-token"
export FIZZY_ACCOUNT="your-account-slug"
go test -tags=integration -v ./fizzy

# Or use Task
task test
task test:coverage
task test:integration
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
go vet ./fizzy/...

# Or use Task for pre-commit checks
task check
```

## Notes

**Card IDs vs Numbers:** The Fizzy API uses card **numbers** (not IDs) in most endpoint URLs. When you create a card, both `ID` and `Number` fields are populated. Use the `Number` field for subsequent operations (Get, Update, Delete, etc.).

```go
card, _ := client.Cards.Create(ctx, opts)
fmt.Printf("ID: %s, Number: %d\n", card.ID, card.Number)

// Use Number for operations
cardNumberStr := fmt.Sprintf("%d", card.Number)
client.Cards.Get(ctx, cardNumberStr)
```

## Contributing

Contributions are welcome! Please:

1. Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification
2. Run `task check` before submitting (formats, lints, and tests)
3. Maintain ≥85% test coverage for new code
4. Add godoc comments for all exported types and methods
5. Test against live API when possible (integration tests)

## License

MIT

## Resources

- [Fizzy API Documentation](https://github.com/basecamp/fizzy/blob/main/docs/API.md)
- [Fizzy Website](https://fizzy.do)
- [Specification](./SPECIFICATION.md)
