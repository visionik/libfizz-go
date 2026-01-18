package fizzy

import (
	"context"
	"fmt"
)

// List retrieves all boards for the authenticated account.
//
// Example:
//
//	boards, err := client.Boards.List(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, board := range boards {
//		fmt.Printf("Board: %s\n", board.Name)
//	}
func (s *BoardsService) List(ctx context.Context) ([]Board, error) {
	path := fmt.Sprintf("/%s/boards.json", s.client.accountSlug)

	var boards []Board
	if err := s.client.doRequest(ctx, "GET", path, nil, &boards); err != nil {
		return nil, fmt.Errorf("failed to list boards: %w", err)
	}

	return boards, nil
}

// Get retrieves a specific board by ID.
//
// Example:
//
//	board, err := client.Boards.Get(ctx, "board-123")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Board: %s\n", board.Name)
func (s *BoardsService) Get(ctx context.Context, boardID string) (*Board, error) {
	path := fmt.Sprintf("/%s/boards/%s.json", s.client.accountSlug, boardID)

	var board Board
	if err := s.client.doRequest(ctx, "GET", path, nil, &board); err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	return &board, nil
}

// Create creates a new board.
//
// Example:
//
//	opts := &fizzy.BoardCreateOptions{
//		Name: "Engineering",
//		Description: ptrString("Development tasks"),
//	}
//	board, err := client.Boards.Create(ctx, opts)
//	if err != nil {
//		log.Fatal(err)
//	}
func (s *BoardsService) Create(ctx context.Context, opts *BoardCreateOptions) (*Board, error) {
	path := fmt.Sprintf("/%s/boards", s.client.accountSlug)

	// Wrap in board object per API spec
	body := map[string]interface{}{"board": opts}

	var board Board
	if err := s.client.doRequest(ctx, "POST", path, body, &board); err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}

	return &board, nil
}

// Update updates an existing board.
//
// Example:
//
//	opts := &fizzy.BoardUpdateOptions{
//		Name: ptrString("New Name"),
//	}
//	board, err := client.Boards.Update(ctx, "board-123", opts)
//	if err != nil {
//		log.Fatal(err)
//	}
func (s *BoardsService) Update(ctx context.Context, boardID string, opts *BoardUpdateOptions) (*Board, error) {
	path := fmt.Sprintf("/%s/boards/%s", s.client.accountSlug, boardID)

	// Wrap in board object per API spec
	body := map[string]interface{}{"board": opts}

	// API returns 204 No Content, so we need to fetch the board after update
	if err := s.client.doRequest(ctx, "PATCH", path, body, nil); err != nil {
		return nil, fmt.Errorf("failed to update board: %w", err)
	}

	// Fetch the updated board
	return s.Get(ctx, boardID)
}

// Delete deletes a board.
//
// Example:
//
//	err := client.Boards.Delete(ctx, "board-123")
//	if err != nil {
//		log.Fatal(err)
//	}
func (s *BoardsService) Delete(ctx context.Context, boardID string) error {
	path := fmt.Sprintf("/%s/boards/%s", s.client.accountSlug, boardID)

	if err := s.client.doRequest(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	return nil
}
