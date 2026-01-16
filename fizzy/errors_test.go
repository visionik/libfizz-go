package fizzy

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorFromResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		message    string
		requestID  string
		wantType   error
	}{
		{"bad request", 400, "Bad request", "req-123", &BadRequestError{}},
		{"unauthorized", 401, "Unauthorized", "req-456", &AuthenticationError{}},
		{"forbidden", 403, "Forbidden", "req-789", &ForbiddenError{}},
		{"not found", 404, "Not found", "req-abc", &NotFoundError{}},
		{"unprocessable", 422, "Invalid data", "req-def", &UnprocessableEntityError{}},
		{"rate limit", 429, "Too many requests", "req-ghi", &RateLimitError{}},
		{"server error", 500, "Internal error", "req-jkl", &ServerError{}},
		{"bad gateway", 502, "Bad gateway", "req-mno", &ServerError{}},
		{"generic error", 418, "I'm a teapot", "req-pqr", &FizzyError{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := newErrorFromResponse(tt.statusCode, tt.message, tt.requestID)
			assert.Error(t, err)
			
			// Check error type using type switch
			switch tt.wantType.(type) {
			case *BadRequestError:
				var target *BadRequestError
				assert.True(t, errors.As(err, &target))
			case *AuthenticationError:
				var target *AuthenticationError
				assert.True(t, errors.As(err, &target))
			case *ForbiddenError:
				var target *ForbiddenError
				assert.True(t, errors.As(err, &target))
			case *NotFoundError:
				var target *NotFoundError
				assert.True(t, errors.As(err, &target))
			case *UnprocessableEntityError:
				var target *UnprocessableEntityError
				assert.True(t, errors.As(err, &target))
			case *RateLimitError:
				var target *RateLimitError
				assert.True(t, errors.As(err, &target))
			case *ServerError:
				var target *ServerError
				assert.True(t, errors.As(err, &target))
			case *FizzyError:
				var target *FizzyError
				assert.True(t, errors.As(err, &target))
			}

			// Check error message contains key info
			errMsg := err.Error()
			assert.Contains(t, errMsg, tt.message)
			assert.Contains(t, errMsg, tt.requestID)
		})
	}
}

func TestFizzyError_Error(t *testing.T) {
	tests := []struct {
		name      string
		err       FizzyError
		wantMatch string
	}{
		{
			"with request ID",
			FizzyError{StatusCode: 404, Message: "Not found", RequestID: "req-123"},
			"fizzy: Not found (status 404, request req-123)",
		},
		{
			"without request ID",
			FizzyError{StatusCode: 500, Message: "Server error", RequestID: ""},
			"fizzy: Server error (status 500)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantMatch, tt.err.Error())
		})
	}
}

func TestErrorTypes_AsCheck(t *testing.T) {
	// Test that errors.As works correctly for specific error types
	err := newErrorFromResponse(404, "Not found", "req-123")

	var notFound *NotFoundError
	assert.True(t, errors.As(err, &notFound))

	var authErr *AuthenticationError
	assert.False(t, errors.As(err, &authErr))
}
