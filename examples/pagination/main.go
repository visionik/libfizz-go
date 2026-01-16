// Example demonstrating pagination
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/visionik/libfizz-go/fizzy"
)

func main() {
	token := os.Getenv("FIZZY_TOKEN")
	accountSlug := os.Getenv("FIZZY_ACCOUNT")
	boardID := os.Getenv("FIZZY_BOARD_ID")

	if token == "" || accountSlug == "" {
		log.Fatal("FIZZY_TOKEN and FIZZY_ACCOUNT must be set")
	}
	if boardID == "" {
		log.Fatal("FIZZY_BOARD_ID must be set")
	}

	client := fizzy.NewClient(token, accountSlug)
	ctx := context.Background()

	fmt.Println("=== Manual Pagination ===")
	opts := &fizzy.CardListOptions{BoardID: boardID}
	pageNum := 1

	for {
		result, err := client.Cards.List(ctx, opts)
		if err != nil {
			log.Fatalf("Failed to list cards: %v", err)
		}

		fmt.Printf("\nPage %d (%d cards):\n", pageNum, len(result.Items))
		for _, card := range result.Items {
			fmt.Printf("  - %s\n", card.Title)
		}

		if result.NextPage == "" {
			break
		}

		opts.Page = result.NextPage
		pageNum++
	}

	fmt.Println("\n=== Automatic Pagination (ListAll) ===")
	allCards, err := client.Cards.ListAll(ctx, &fizzy.CardListOptions{
		BoardID: boardID,
	})
	if err != nil {
		log.Fatalf("Failed to list all cards: %v", err)
	}

	fmt.Printf("\nTotal cards: %d\n", len(allCards))
	for i, card := range allCards {
		fmt.Printf("  %d. %s\n", i+1, card.Title)
	}
}
