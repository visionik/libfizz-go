package fizzy

import (
	"context"
	"fmt"
)

// List retrieves all comments for a card.
//
// Example:
//
//	comments, err := client.Comments.List(ctx, "card-123")
func (s *CommentsService) List(ctx context.Context, cardID string) ([]Comment, error) {
	path := fmt.Sprintf("/%s/cards/%s/comments", s.client.accountSlug, cardID)

	var comments []Comment
	if err := s.client.doRequest(ctx, "GET", path, nil, &comments); err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	return comments, nil
}

// Create creates a new comment on a card.
//
// Example:
//
//	comment, err := client.Comments.Create(ctx, "card-123", &fizzy.CommentCreateOptions{
//		Body: "This looks good!",
//	})
func (s *CommentsService) Create(ctx context.Context, cardID string, opts *CommentCreateOptions) (*Comment, error) {
	path := fmt.Sprintf("/%s/cards/%s/comments", s.client.accountSlug, cardID)

	// Wrap in comment object per API spec
	body := map[string]interface{}{"comment": opts}

	var comment Comment
	if err := s.client.doRequest(ctx, "POST", path, body, &comment); err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &comment, nil
}

// Update updates an existing comment.
//
// Example:
//
//	comment, err := client.Comments.Update(ctx, "card-123", "comment-456", &fizzy.CommentUpdateOptions{
//		Body: "Updated comment text",
//	})
func (s *CommentsService) Update(ctx context.Context, cardID, commentID string, opts *CommentUpdateOptions) (*Comment, error) {
	path := fmt.Sprintf("/%s/cards/%s/comments/%s", s.client.accountSlug, cardID, commentID)

	// Wrap in comment object per API spec
	body := map[string]interface{}{"comment": opts}

	var comment Comment
	if err := s.client.doRequest(ctx, "PATCH", path, body, &comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}

	return &comment, nil
}

// Delete deletes a comment.
//
// Example:
//
//	err := client.Comments.Delete(ctx, "card-123", "comment-456")
func (s *CommentsService) Delete(ctx context.Context, cardID, commentID string) error {
	path := fmt.Sprintf("/%s/cards/%s/comments/%s", s.client.accountSlug, cardID, commentID)

	if err := s.client.doRequest(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}
