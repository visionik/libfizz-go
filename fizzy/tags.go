package fizzy

import (
	"context"
	"fmt"
)

// List retrieves all tags for the account.
func (s *TagsService) List(ctx context.Context) ([]Tag, error) {
	path := fmt.Sprintf("/%s/tags", s.client.accountSlug)

	var tags []Tag
	if err := s.client.doRequest(ctx, "GET", path, nil, &tags); err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return tags, nil
}

// Create creates a new tag.
func (s *TagsService) Create(ctx context.Context, opts *TagCreateOptions) (*Tag, error) {
	path := fmt.Sprintf("/%s/tags", s.client.accountSlug)

	var tag Tag
	if err := s.client.doRequest(ctx, "POST", path, opts, &tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return &tag, nil
}
