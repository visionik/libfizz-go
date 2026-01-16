package fizzy

import (
	"context"
	"fmt"
)

// List retrieves all columns for a board.
func (s *ColumnsService) List(ctx context.Context, boardID string) ([]Column, error) {
	path := fmt.Sprintf("/%s/boards/%s/columns", s.client.accountSlug, boardID)

	var columns []Column
	if err := s.client.doRequest(ctx, "GET", path, nil, &columns); err != nil {
		return nil, fmt.Errorf("failed to list columns: %w", err)
	}

	return columns, nil
}

// Get retrieves a specific column.
func (s *ColumnsService) Get(ctx context.Context, boardID, columnID string) (*Column, error) {
	path := fmt.Sprintf("/%s/boards/%s/columns/%s", s.client.accountSlug, boardID, columnID)

	var column Column
	if err := s.client.doRequest(ctx, "GET", path, nil, &column); err != nil {
		return nil, fmt.Errorf("failed to get column: %w", err)
	}

	return &column, nil
}

// Create creates a new column on a board.
func (s *ColumnsService) Create(ctx context.Context, boardID string, opts *ColumnCreateOptions) (*Column, error) {
	path := fmt.Sprintf("/%s/boards/%s/columns", s.client.accountSlug, boardID)

	// Wrap in column object per API spec
	body := map[string]interface{}{"column": opts}

	var column Column
	if err := s.client.doRequest(ctx, "POST", path, body, &column); err != nil {
		return nil, fmt.Errorf("failed to create column: %w", err)
	}

	return &column, nil
}

// Update updates an existing column.
func (s *ColumnsService) Update(ctx context.Context, boardID, columnID string, opts *ColumnUpdateOptions) (*Column, error) {
	path := fmt.Sprintf("/%s/boards/%s/columns/%s", s.client.accountSlug, boardID, columnID)

	// Wrap in column object per API spec
	body := map[string]interface{}{"column": opts}

	var column Column
	if err := s.client.doRequest(ctx, "PATCH", path, body, &column); err != nil {
		return nil, fmt.Errorf("failed to update column: %w", err)
	}

	return &column, nil
}

// Delete deletes a column.
func (s *ColumnsService) Delete(ctx context.Context, boardID, columnID string) error {
	path := fmt.Sprintf("/%s/boards/%s/columns/%s", s.client.accountSlug, boardID, columnID)

	if err := s.client.doRequest(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete column: %w", err)
	}

	return nil
}
