package fizzy

import "context"

// PagedResult wraps a list result with pagination information.
type PagedResult[T any] struct {
	Items    []T
	NextPage string
}

// listAllFunc is a function that lists a page of results.
type listAllFunc[T any] func(ctx context.Context, page string) (*PagedResult[T], error)

// listAll fetches all pages of results using the provided list function.
// This is a helper for implementing ListAll methods across services.
func listAll[T any](ctx context.Context, fn listAllFunc[T]) ([]T, error) {
	var allItems []T
	page := ""

	for {
		result, err := fn(ctx, page)
		if err != nil {
			return nil, err
		}

		allItems = append(allItems, result.Items...)

		// No more pages
		if result.NextPage == "" {
			break
		}

		page = result.NextPage
	}

	return allItems, nil
}
