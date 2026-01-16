package fizzy

import (
	"context"
	"fmt"
)

// List retrieves all steps for a card.
func (s *StepsService) List(ctx context.Context, cardID string) ([]Step, error) {
	path := fmt.Sprintf("/%s/cards/%s/steps", s.client.accountSlug, cardID)

	var steps []Step
	if err := s.client.doRequest(ctx, "GET", path, nil, &steps); err != nil {
		return nil, fmt.Errorf("failed to list steps: %w", err)
	}

	return steps, nil
}

// Get retrieves a specific step.
func (s *StepsService) Get(ctx context.Context, cardID, stepID string) (*Step, error) {
	path := fmt.Sprintf("/%s/cards/%s/steps/%s", s.client.accountSlug, cardID, stepID)

	var step Step
	if err := s.client.doRequest(ctx, "GET", path, nil, &step); err != nil {
		return nil, fmt.Errorf("failed to get step: %w", err)
	}

	return &step, nil
}

// Create creates a new step on a card.
func (s *StepsService) Create(ctx context.Context, cardID string, opts *StepCreateOptions) (*Step, error) {
	path := fmt.Sprintf("/%s/cards/%s/steps", s.client.accountSlug, cardID)

	// Wrap in step object per API spec
	body := map[string]interface{}{"step": opts}

	var step Step
	if err := s.client.doRequest(ctx, "POST", path, body, &step); err != nil {
		return nil, fmt.Errorf("failed to create step: %w", err)
	}

	return &step, nil
}

// Update updates an existing step.
func (s *StepsService) Update(ctx context.Context, cardID, stepID string, opts *StepUpdateOptions) (*Step, error) {
	path := fmt.Sprintf("/%s/cards/%s/steps/%s", s.client.accountSlug, cardID, stepID)

	// Wrap in step object per API spec
	body := map[string]interface{}{"step": opts}

	var step Step
	if err := s.client.doRequest(ctx, "PATCH", path, body, &step); err != nil {
		return nil, fmt.Errorf("failed to update step: %w", err)
	}

	return &step, nil
}

// Delete deletes a step.
func (s *StepsService) Delete(ctx context.Context, cardID, stepID string) error {
	path := fmt.Sprintf("/%s/cards/%s/steps/%s", s.client.accountSlug, cardID, stepID)

	if err := s.client.doRequest(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("failed to delete step: %w", err)
	}

	return nil
}
