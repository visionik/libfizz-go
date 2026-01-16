package fizzy

import (
	"context"
	"fmt"
)

// List retrieves all users in the account.
func (s *UsersService) List(ctx context.Context) ([]User, error) {
	path := fmt.Sprintf("/%s/users", s.client.accountSlug)

	var users []User
	if err := s.client.doRequest(ctx, "GET", path, nil, &users); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}
