package fizzy

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// List retrieves cards with optional filtering and pagination.
//
// Returns a single page of results. Use ListAll() to fetch all pages automatically.
//
// Example:
//
//	result, err := client.Cards.List(ctx, &fizzy.CardListOptions{
//		BoardID: "board-123",
//		Status:  "open",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, card := range result.Items {
//		fmt.Printf("Card: %s\n", card.Title)
//	}
func (s *CardsService) List(ctx context.Context, opts *CardListOptions) (*PagedResult[Card], error) {
	path := fmt.Sprintf("/%s/cards", s.client.accountSlug)

	// Build query parameters
	if opts != nil {
		query := url.Values{}
		if opts.BoardID != "" {
			query.Set("board_id", opts.BoardID)
		}
		if opts.Status != "" {
			query.Set("status", opts.Status)
		}
		if opts.ColumnID != "" {
			query.Set("column_id", opts.ColumnID)
		}
		if len(opts.TagIDs) > 0 {
			query.Set("tag_ids", strings.Join(opts.TagIDs, ","))
		}
		if opts.Page != "" {
			query.Set("page", opts.Page)
		}
		if len(query) > 0 {
			path += "?" + query.Encode()
		}
	}

	// API returns array directly, not nested object
	var cards []Card
	if err := s.client.doRequest(ctx, "GET", path, nil, &cards); err != nil {
		return nil, fmt.Errorf("failed to list cards: %w", err)
	}

	// TODO: Handle pagination - need to check response headers for next page
	return &PagedResult[Card]{
		Items:    cards,
		NextPage: "", // Pagination not yet implemented
	}, nil
}

// ListAll retrieves all cards across all pages.
//
// Example:
//
//	cards, err := client.Cards.ListAll(ctx, &fizzy.CardListOptions{
//		BoardID: "board-123",
//	})
func (s *CardsService) ListAll(ctx context.Context, opts *CardListOptions) ([]Card, error) {
	if opts == nil {
		opts = &CardListOptions{}
	}
	return listAll(ctx, func(ctx context.Context, page string) (*PagedResult[Card], error) {
		opts.Page = page
		return s.List(ctx, opts)
	})
}

// Get retrieves a specific card by ID.
//
// Example:
//
//	card, err := client.Cards.Get(ctx, "card-123")
func (s *CardsService) Get(ctx context.Context, cardID string) (*Card, error) {
	path := fmt.Sprintf("/%s/cards/%s", s.client.accountSlug, cardID)

	var card Card
	if err := s.client.doRequest(ctx, "GET", path, nil, &card); err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}

	return &card, nil
}

// Create creates a new card.
//
// Example:
//
//	card, err := client.Cards.Create(ctx, &fizzy.CardCreateOptions{
//		BoardID: "board-123",
//		Title:   "Bug: Login broken",
//		Body:    ptrString("Users can't log in with SSO"),
//	})
func (s *CardsService) Create(ctx context.Context, opts *CardCreateOptions) (*Card, error) {
	path := fmt.Sprintf("/%s/boards/%s/cards", s.client.accountSlug, opts.BoardID)

	// Wrap in card object per API spec
	body := map[string]interface{}{
		"card": map[string]interface{}{
			"title":       opts.Title,
			"description": opts.Body,
		},
	}

	var card Card
	if err := s.client.doRequest(ctx, "POST", path, body, &card); err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	return &card, nil
}

// Update updates an existing card.
//
// Example:
//
//	card, err := client.Cards.Update(ctx, "card-123", &fizzy.CardUpdateOptions{
//		Title: ptrString("Updated title"),
//	})
func (s *CardsService) Update(ctx context.Context, cardID string, opts *CardUpdateOptions) (*Card, error) {
	path := fmt.Sprintf("/%s/cards/%s", s.client.accountSlug, cardID)

	// Wrap in card object per API spec
	body := map[string]interface{}{
		"card": map[string]interface{}{},
	}
	if opts.Title != nil {
		body["card"].(map[string]interface{})["title"] = *opts.Title
	}
	if opts.Body != nil {
		body["card"].(map[string]interface{})["description"] = *opts.Body
	}

	var card Card
	if err := s.client.doRequest(ctx, "PATCH", path, body, &card); err != nil {
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	return &card, nil
}

// Delete deletes a card.
//
// Example:
//
//	err := client.Cards.Delete(ctx, "card-123")
func (s *CardsService) Delete(ctx context.Context, cardID string) error {
	path := fmt.Sprintf("/%s/cards/%s", s.client.accountSlug, cardID)

	if err := s.client.doRequest(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	return nil
}

// Close marks a card as closed.
//
// Example:
//
//	err := client.Cards.Close(ctx, "card-123")
func (s *CardsService) Close(ctx context.Context, cardID string) error {
	path := fmt.Sprintf("/%s/cards/%s/closure", s.client.accountSlug, cardID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to close card: %w", err)
	}

	return nil
}

// Reopen reopens a closed card.
//
// Example:
//
//	err := client.Cards.Reopen(ctx, "card-123")
func (s *CardsService) Reopen(ctx context.Context, cardID string) error {
	path := fmt.Sprintf("/%s/cards/%s/closure", s.client.accountSlug, cardID)

	if err := s.client.doRequest(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to reopen card: %w", err)
	}

	return nil
}

// Postpone moves a card to "Not Now" status.
//
// Example:
//
//	err := client.Cards.Postpone(ctx, "card-123")
func (s *CardsService) Postpone(ctx context.Context, cardID string) error {
	path := fmt.Sprintf("/%s/cards/%s/not_now", s.client.accountSlug, cardID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to postpone card: %w", err)
	}

	return nil
}

// Triage sends a card back to triage (Maybe column).
//
// Example:
//
//	err := client.Cards.Triage(ctx, "card-123")
func (s *CardsService) Triage(ctx context.Context, cardID string) error {
	path := fmt.Sprintf("/%s/cards/%s/triage", s.client.accountSlug, cardID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to triage card: %w", err)
	}

	return nil
}

// MoveToColumn moves a card to a specific column.
//
// Example:
//
//	err := client.Cards.MoveToColumn(ctx, "card-123", "column-456")
func (s *CardsService) MoveToColumn(ctx context.Context, cardID, columnID string) error {
	path := fmt.Sprintf("/%s/cards/%s/column", s.client.accountSlug, cardID)

	body := map[string]string{"column_id": columnID}
	if err := s.client.doRequest(ctx, "POST", path, body, nil); err != nil {
		return fmt.Errorf("failed to move card to column: %w", err)
	}

	return nil
}

// Assign toggles assignment of a user to a card.
// If the user is already assigned, they will be unassigned, and vice versa.
//
// Example:
//
//	err := client.Cards.Assign(ctx, "card-123", "user-456")
func (s *CardsService) Assign(ctx context.Context, cardID, userID string) error {
	path := fmt.Sprintf("/%s/cards/%s/assignments/%s/toggle", s.client.accountSlug, cardID, userID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to assign card: %w", err)
	}

	return nil
}

// Tag toggles a tag on a card.
// If the tag is already applied, it will be removed, and vice versa.
// If the tag doesn't exist, it will be created.
//
// Example:
//
//	err := client.Cards.Tag(ctx, "card-123", "bug")
func (s *CardsService) Tag(ctx context.Context, cardID, tagName string) error {
	path := fmt.Sprintf("/%s/cards/%s/tags/%s/toggle", s.client.accountSlug, cardID, tagName)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to tag card: %w", err)
	}

	return nil
}

// Watch subscribes to notifications for a card.
//
// Example:
//
//	err := client.Cards.Watch(ctx, "card-123")
func (s *CardsService) Watch(ctx context.Context, cardID string) error {
	path := fmt.Sprintf("/%s/cards/%s/watch", s.client.accountSlug, cardID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to watch card: %w", err)
	}

	return nil
}

// Unwatch unsubscribes from notifications for a card.
//
// Example:
//
//	err := client.Cards.Unwatch(ctx, "card-123")
func (s *CardsService) Unwatch(ctx context.Context, cardID string) error {
	path := fmt.Sprintf("/%s/cards/%s/unwatch", s.client.accountSlug, cardID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to unwatch card: %w", err)
	}

	return nil
}

// MarkGolden marks a card as golden (featured/important).
//
// Example:
//
//	err := client.Cards.MarkGolden(ctx, "card-123")
func (s *CardsService) MarkGolden(ctx context.Context, cardID string) error {
	path := fmt.Sprintf("/%s/cards/%s/golden", s.client.accountSlug, cardID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to mark card as golden: %w", err)
	}

	return nil
}

// UnmarkGolden removes golden status from a card.
//
// Example:
//
//	err := client.Cards.UnmarkGolden(ctx, "card-123")
func (s *CardsService) UnmarkGolden(ctx context.Context, cardID string) error {
	path := fmt.Sprintf("/%s/cards/%s/ungolden", s.client.accountSlug, cardID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to unmark card as golden: %w", err)
	}

	return nil
}
