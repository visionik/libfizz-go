package fizzy

import "fmt"

// FizzyError is the base error type for all Fizzy API errors.
type FizzyError struct {
	StatusCode int
	Message    string
	RequestID  string
}

func (e *FizzyError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("fizzy: %s (status %d, request %s)", e.Message, e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("fizzy: %s (status %d)", e.Message, e.StatusCode)
}

// AuthenticationError represents a 401 Unauthorized error.
type AuthenticationError struct {
	FizzyError
}

// ForbiddenError represents a 403 Forbidden error.
type ForbiddenError struct {
	FizzyError
}

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct {
	FizzyError
}

// BadRequestError represents a 400 Bad Request error.
type BadRequestError struct {
	FizzyError
}

// UnprocessableEntityError represents a 422 Unprocessable Entity error.
type UnprocessableEntityError struct {
	FizzyError
}

// RateLimitError represents a 429 Too Many Requests error.
type RateLimitError struct {
	FizzyError
	RetryAfter int // Seconds to wait before retry
}

// ServerError represents a 5xx server error.
type ServerError struct {
	FizzyError
}

// newErrorFromResponse creates an appropriate error type based on HTTP status code.
func newErrorFromResponse(statusCode int, message, requestID string) error {
	baseErr := FizzyError{
		StatusCode: statusCode,
		Message:    message,
		RequestID:  requestID,
	}

	switch statusCode {
	case 400:
		return &BadRequestError{FizzyError: baseErr}
	case 401:
		return &AuthenticationError{FizzyError: baseErr}
	case 403:
		return &ForbiddenError{FizzyError: baseErr}
	case 404:
		return &NotFoundError{FizzyError: baseErr}
	case 422:
		return &UnprocessableEntityError{FizzyError: baseErr}
	case 429:
		return &RateLimitError{FizzyError: baseErr}
	default:
		if statusCode >= 500 {
			return &ServerError{FizzyError: baseErr}
		}
		return &baseErr
	}
}
