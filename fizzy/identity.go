package fizzy

import (
	"context"
	"fmt"
)

// Get retrieves the current user's identity and account information.
//
// This endpoint returns information about the authenticated user including
// their name, email, and the list of accounts they have access to.
//
// Example:
//
//	identity, err := client.Identity.Get(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Logged in as %s\n", identity.Name)
func (s *IdentityService) Get(ctx context.Context) (*Identity, error) {
	path := "/my/identity"

	var identity Identity
	if err := s.client.doRequest(ctx, "GET", path, nil, &identity); err != nil {
		return nil, fmt.Errorf("failed to get identity: %w", err)
	}

	return &identity, nil
}
