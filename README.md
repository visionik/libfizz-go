# libfizz-go

A simple, idiomatic Go client library for the [Fizzy API](https://fizzy.do).

[![Go Reference](https://pkg.go.dev/badge/github.com/visionik/libfizz-go/fizzy.svg)](https://pkg.go.dev/github.com/visionik/libfizz-go/fizzy)

## Features

- **Complete API Coverage**: All Fizzy API resources (Identity, Boards, Cards, Comments, Reactions, Steps, Tags, Columns, Users, Notifications, File Uploads)
- **Idiomatic Go**: Follows Go best practices with context support, functional options, and typed errors
- **Automatic Retries**: Exponential backoff for rate limits (429) and server errors (5xx)
- **ETag Caching**: Built-in HTTP caching to reduce bandwidth and improve performance
- **Thread-Safe**: Client can be safely shared across goroutines
- **Comprehensive Testing**: ≥85% test coverage

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
	fmt.Printf("Logged in as: %s\n", identity.Name)

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

// Update a card
card, err := client.Cards.Update(ctx, "card-id", &fizzy.CardUpdateOptions{
	Title: ptrString("Updated title"),
})

// Card actions
err = client.Cards.Close(ctx, "card-id")
err = client.Cards.Reopen(ctx, "card-id")
err = client.Cards.Assign(ctx, "card-id", "user-id")
err = client.Cards.Tag(ctx, "card-id", "tag-id")
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

✅ **Implemented:**
- Identity (Get)
- Boards (List, Get, Create, Update, Delete)

🚧 **In Progress** (see `IMPLEMENTATION_STATUS.md`):
- Cards (CRUD + Actions)
- Comments (List, Create, Update, Delete)
- Reactions (List, Create, Delete)
- Steps (List, Get, Create, Update, Delete)
- Tags (List, Create)
- Columns (List, Get, Create, Update, Delete)
- Users (List)
- Notifications (List, Read, Unread, ReadAll)
- Uploads (Direct upload, convenience wrapper)

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
# Run tests
go test ./fizzy/...

# Run with coverage
go test -coverprofile=coverage.out ./fizzy/...
go tool cover -html=coverage.out

# Or use Task
task test
task test:coverage
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

## Contributing

Contributions are welcome! Please:

1. Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification
2. Run `task check` before submitting
3. Ensure ≥85% test coverage for new code
4. Add godoc comments for all exported types and methods

## License

MIT

## Resources

- [Fizzy API Documentation](https://github.com/basecamp/fizzy/blob/main/docs/API.md)
- [Fizzy Website](https://fizzy.do)
- [Specification](./SPECIFICATION.md)
