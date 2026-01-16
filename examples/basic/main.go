// Basic example of using libfizz-go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/visionik/libfizz-go/fizzy"
)

func main() {
	// Get credentials from environment
	token := os.Getenv("FIZZY_TOKEN")
	accountSlug := os.Getenv("FIZZY_ACCOUNT")

	if token == "" || accountSlug == "" {
		log.Fatal("FIZZY_TOKEN and FIZZY_ACCOUNT environment variables must be set")
	}

	// Create client
	client := fizzy.NewClient(token, accountSlug)

	ctx := context.Background()

	// Get identity
	identity, err := client.Identity.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get identity: %v", err)
	}
	if len(identity.Accounts) > 0 {
		account := identity.Accounts[0]
		if account.User != nil {
			fmt.Printf("Logged in as: %s (%s)\n", account.User.Name, *account.User.EmailAddress)
		}
		fmt.Printf("Account: %s\n", account.Name)
	}

	// List boards
	boards, err := client.Boards.List(ctx)
	if err != nil {
		log.Fatalf("Failed to list boards: %v", err)
	}

	fmt.Printf("\nBoards (%d):\n", len(boards))
	for _, board := range boards {
		fmt.Printf("  - %s (ID: %s)\n", board.Name, board.ID)
	}

	if len(boards) == 0 {
		fmt.Println("\nNo boards found. Create one first!")
		return
	}

	// List cards from first board
	result, err := client.Cards.List(ctx, &fizzy.CardListOptions{
		BoardID: boards[0].ID,
		Status:  "open",
	})
	if err != nil {
		log.Fatalf("Failed to list cards: %v", err)
	}

	fmt.Printf("\nOpen cards on '%s' (%d):\n", boards[0].Name, len(result.Items))
	for _, card := range result.Items {
		fmt.Printf("  - %s\n", card.Title)
	}
}
