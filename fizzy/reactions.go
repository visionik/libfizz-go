package fizzy

import (
	"context"
	"fmt"
)

// List retrieves all reactions for a comment.
func (s *ReactionsService) List(ctx context.Context, cardID, commentID string) ([]Reaction, error) {
	path := fmt.Sprintf("/%s/cards/%s/comments/%s/reactions", s.client.accountSlug, cardID, commentID)

	var reactions []Reaction
	if err := s.client.doRequest(ctx, "GET", path, nil, &reactions); err != nil {
		return nil, fmt.Errorf("failed to list reactions: %w", err)
	}

	return reactions, nil
}

// Create adds a reaction to a comment.
func (s *ReactionsService) Create(ctx context.Context, cardID, commentID string, opts *ReactionCreateOptions) (*Reaction, error) {
	path := fmt.Sprintf("/%s/cards/%s/comments/%s/reactions", s.client.accountSlug, cardID, commentID)

	// Wrap in reaction object per API spec
	body := map[string]interface{}{"reaction": opts}

	var reaction Reaction
	if err := s.client.doRequest(ctx, "POST", path, body, &reaction); err != nil {
		return nil, fmt.Errorf("failed to create reaction: %w", err)
	}

	return &reaction, nil
}

// Delete removes a reaction from a comment.
func (s *ReactionsService) Delete(ctx context.Context, cardID, commentID, reactionID string) error {
	path := fmt.Sprintf("/%s/cards/%s/comments/%s/reactions/%s", s.client.accountSlug, cardID, commentID, reactionID)

	if err := s.client.doRequest(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete reaction: %w", err)
	}

	return nil
}
