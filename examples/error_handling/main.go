// Example demonstrating error handling with typed errors
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/visionik/libfizz-go/fizzy"
)

func main() {
	token := os.Getenv("FIZZY_TOKEN")
	accountSlug := os.Getenv("FIZZY_ACCOUNT")

	if token == "" || accountSlug == "" {
		log.Fatal("FIZZY_TOKEN and FIZZY_ACCOUNT must be set")
	}

	client := fizzy.NewClient(token, accountSlug)
	ctx := context.Background()

	// Example 1: Handle NotFoundError
	fmt.Println("=== Handling Not Found Errors ===")
	_, err := client.Cards.Get(ctx, "invalid-card-id")
	if err != nil {
		var notFound *fizzy.NotFoundError
		if errors.As(err, &notFound) {
			fmt.Printf("✓ Card not found (expected): %v\n", err)
		} else {
			log.Fatalf("Unexpected error: %v", err)
		}
	}

	// Example 2: Handle AuthenticationError
	fmt.Println("\n=== Handling Authentication Errors ===")
	badClient := fizzy.NewClient("invalid-token", accountSlug)
	_, err = badClient.Identity.Get(ctx)
	if err != nil {
		var authErr *fizzy.AuthenticationError
		if errors.As(err, &authErr) {
			fmt.Printf("✓ Authentication failed (expected): %v\n", err)
		} else {
			log.Fatalf("Unexpected error: %v", err)
		}
	}

	// Example 3: Successful request
	fmt.Println("\n=== Successful Request ===")
	identity, err := client.Identity.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to get identity: %v", err)
	}
	fmt.Printf("✓ Successfully authenticated as: %s\n", identity.Name)

	// Example 4: Generic error handling
	fmt.Println("\n=== Generic Error Handling ===")
	_, err = client.Boards.Get(ctx, "maybe-invalid-board")
	if err != nil {
		var fizzErr *fizzy.FizzyError
		if errors.As(err, &fizzErr) {
			fmt.Printf("Fizzy API error (status %d): %s\n", fizzErr.StatusCode, fizzErr.Message)
		} else {
			fmt.Printf("Non-API error: %v\n", err)
		}
	}
}
