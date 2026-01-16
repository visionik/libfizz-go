# libfizz-go Specification

**Generated**: 2026-01-16
**Status**: Ready for Implementation
**Version**: 1.0.0

## Overview

libfizz-go is a simple, idiomatic Go client library for the Fizzy API (https://fizzy.do). Fizzy is a modern kanban board application for tracking bugs, issues, ideas, and small projects. This library provides full coverage of the Fizzy HTTP API with a focus on simplicity, Go best practices, and developer experience.

### Goals
- Complete coverage of all Fizzy API resources (Identity, Boards, Cards, Comments, Reactions, Steps, Tags, Columns, Users, Notifications, File Uploads)
- Idiomatic Go API design following standard library and popular SDK patterns
- Simple, safe defaults with configurability for advanced use cases
- Comprehensive error handling with typed errors
- Built-in resilience (automatic retries, ETag caching)
- ≥85% test coverage per warping Go standards

### Non-Goals
- Magic Link authentication (Personal Access Token only for v1.0)
- Webhook handling (separate package if needed)
- CLI tool (library only)
- Real-time updates (polling only)

## Requirements

### Functional Requirements

**FR-1: Authentication**
- Support Personal Access Token authentication via Bearer token
- Token provided at client initialization
- All requests include `Authorization: Bearer <token>` header

**FR-2: API Resource Coverage**
All Fizzy API resources must be supported:
- Identity: Get current user identity and accounts
- Boards: List, Get, Create, Update, Delete
- Cards: List, Get, Create, Update, Delete, Close, Reopen, Postpone, Triage, Assign, Tag, Watch, Golden
- Comments: List, Create, Update, Delete
- Reactions: List, Create, Delete
- Steps: List, Get, Create, Update, Delete
- Tags: List, Create
- Columns: List, Get, Create, Update, Delete
- Users: List
- Notifications: List, Read, Unread, ReadAll
- File Uploads: Direct upload for rich text attachments and card header images

**FR-3: Pagination**
- Manual pagination: Return results + next page cursor
- Automatic pagination: `ListAll()` helper methods to fetch all pages
- Caller controls which approach to use

**FR-4: Error Handling**
- Typed errors for each HTTP status (400, 401, 403, 404, 422, 429, 5xx)
- Automatic retry with exponential backoff for 429 (rate limit) and 5xx errors
- Maximum 3 retry attempts (configurable)
- Errors implement standard Go error interface and support `errors.Is/As`

**FR-5: HTTP Caching**
- Built-in ETag support for GET requests
- Automatic If-None-Match header on cached resources
- Handle 304 Not Modified responses
- In-memory cache, enabled by default, configurable

**FR-6: Context Support**
- All API methods accept `context.Context` as first parameter
- Support cancellation and timeouts via context
- Propagate context through retry logic

### Non-Functional Requirements

**NFR-1: Performance**
- Default timeout: 30 seconds per request
- ETag caching reduces bandwidth and improves response times
- Retry delays: 1s, 2s, 4s (exponential backoff)

**NFR-2: Reliability**
- Automatic retries for transient failures (429, 5xx)
- Thread-safe client (can be shared across goroutines)
- No panics in library code

**NFR-3: Usability**
- Zero-config default: `NewClient(token, accountSlug)` works immediately
- Functional options for advanced configuration
- Namespaced services for clear API organization
- Comprehensive documentation and examples

**NFR-4: Testing**
- ≥85% test coverage (per warping Go standards)
- Unit tests with httptest for all endpoints
- Table-driven tests for consistency
- Integration tests (optional, require live API token)

**NFR-5: Maintainability**
- Follow warping Go standards (see warping/languages/go.md)
- Clear separation: client, services, models, errors
- Conventional Commits for version control

## Architecture

### Package Structure

```
libfizz-go/
├── fizzy/                  # Main package
│   ├── client.go          # Client initialization, config, HTTP transport
│   ├── client_options.go  # Functional options (WithTimeout, WithCache, etc.)
│   ├── errors.go          # Typed error definitions
│   ├── retry.go           # Retry logic with exponential backoff
│   ├── cache.go           # ETag cache implementation
│   ├── models.go          # Shared types (Identity, Account, etc.)
│   ├── boards.go          # BoardsService
│   ├── cards.go           # CardsService
│   ├── comments.go        # CommentsService
│   ├── reactions.go       # ReactionsService
│   ├── steps.go           # StepsService
│   ├── tags.go            # TagsService
│   ├── columns.go         # ColumnsService
│   ├── users.go           # UsersService
│   ├── notifications.go   # NotificationsService
│   ├── uploads.go         # UploadsService
│   └── pagination.go      # Pagination helpers
├── examples/              # Usage examples
│   ├── basic/
│   ├── pagination/
│   └── error_handling/
├── Taskfile.yml          # Build, test, lint tasks
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

### Core Components

**Client**
- Holds configuration (token, account slug, base URL, timeout)
- Manages HTTP client with retry and caching middleware
- Exposes service instances (Boards, Cards, etc.)
- Thread-safe

**Services**
- Each resource has its own service (BoardsService, CardsService, etc.)
- Services hold reference to parent client
- Methods follow pattern: `Method(ctx context.Context, params...) (result, error)`
- Handle request construction, response parsing, error mapping

**Models**
- Go structs for API request/response types
- Use `json` struct tags for serialization
- Pointers for optional fields (distinguish nil from zero value)
- Time fields use `time.Time`

**Errors**
- Base error type: `FizzyError` with status code, message, request ID
- Specific types: `AuthenticationError`, `NotFoundError`, `RateLimitError`, etc.
- Implement `error` interface and support `errors.Is/As`

**Retry Logic**
- Intercepts HTTP responses
- Retries on 429 and 5xx with exponential backoff
- Respects context cancellation
- Configurable max attempts

**Cache**
- In-memory map: URL → ETag
- Thread-safe (mutex)
- Adds If-None-Match header on GET requests
- Returns cached response on 304

### API Design Patterns

**Client Initialization**
```go
// Simple
client := fizzy.NewClient(token, accountSlug)

// With options
client := fizzy.NewClient(token, accountSlug,
    fizzy.WithTimeout(60*time.Second),
    fizzy.WithBaseURL("https://custom.fizzy.do"),
    fizzy.WithCache(false),
    fizzy.WithMaxRetries(5),
)
```

**Service Usage**
```go
// Get identity
identity, err := client.Identity.Get(ctx)

// List boards
boards, err := client.Boards.List(ctx)

// Create card
card, err := client.Cards.Create(ctx, &fizzy.CardCreateOptions{
    BoardID: "board-123",
    Title:   "New issue",
    Body:    "Description...",
})

// Manual pagination
result, err := client.Cards.List(ctx, &fizzy.CardListOptions{
    BoardID: "board-123",
})
for _, card := range result.Cards {
    // process
}
if result.NextPage != "" {
    result, err = client.Cards.List(ctx, &fizzy.CardListOptions{
        Page: result.NextPage,
    })
}

// Automatic pagination
allCards, err := client.Cards.ListAll(ctx, &fizzy.CardListOptions{
    BoardID: "board-123",
})
```

**Error Handling**
```go
card, err := client.Cards.Get(ctx, "invalid-id")
if err != nil {
    var notFound *fizzy.NotFoundError
    if errors.As(err, &notFound) {
        // Handle not found
    }
    
    var rateLimit *fizzy.RateLimitError
    if errors.As(err, &rateLimit) {
        // Already retried 3 times, back off longer
    }
    
    return err
}
```

## Implementation Plan

### Phase 1: Foundation
**Dependencies**: None

#### Subphase 1.1: Project Setup
**Dependencies**: None

- Task 1.1.1: Initialize Go module
  - Run `go mod init github.com/yourusername/libfizz-go`
  - Create directory structure
  - **Acceptance**: `go.mod` exists, directories created

- Task 1.1.2: Create Taskfile.yml
  - Add tasks: `build`, `test`, `test:coverage`, `fmt`, `lint`, `check`
  - Follow warping taskfile standards
  - **Acceptance**: `task check` runs successfully

- Task 1.1.3: Add dependencies
  - Add `github.com/stretchr/testify` for testing
  - **Acceptance**: `go.mod` updated, `go mod tidy` succeeds

#### Subphase 1.2: Core Client Infrastructure
**Dependencies**: 1.1

- Task 1.2.1: Implement client.go
  - Client struct with config fields
  - `NewClient(token, accountSlug, ...opts)` constructor
  - HTTP client initialization
  - **Acceptance**: Can instantiate client, unit tests pass

- Task 1.2.2: Implement client_options.go
  - Functional options: `WithTimeout`, `WithBaseURL`, `WithCache`, `WithMaxRetries`
  - **Acceptance**: Options modify client config correctly, unit tests pass

- Task 1.2.3: Implement errors.go
  - Base `FizzyError` type
  - Specific error types for each HTTP status
  - `errors.Is/As` support
  - **Acceptance**: Can create and match errors, unit tests pass

- Task 1.2.4: Implement retry.go
  - HTTP RoundTripper for retry logic
  - Exponential backoff (1s, 2s, 4s)
  - Retry on 429 and 5xx
  - Respect context cancellation
  - **Acceptance**: Retries work correctly, unit tests with httptest pass

- Task 1.2.5: Implement cache.go
  - In-memory ETag cache with mutex
  - Add If-None-Match header
  - Handle 304 responses
  - **Acceptance**: Caching works, unit tests pass

#### Subphase 1.3: Models and Pagination
**Dependencies**: 1.2

- Task 1.3.1: Implement models.go
  - Core types: Identity, Account, User
  - Use json tags, pointer fields for optionals
  - **Acceptance**: Models marshal/unmarshal correctly, unit tests pass

- Task 1.3.2: Implement pagination.go
  - PagedResult wrapper with NextPage field
  - `ListAll()` helper logic (fetch all pages)
  - **Acceptance**: Pagination helpers work, unit tests pass

### Phase 2: API Resources - Core
**Dependencies**: Phase 1

#### Subphase 2.1: Identity and Boards
**Dependencies**: 1.3

- Task 2.1.1: Implement Identity service
  - `identity.go` with IdentityService
  - `Get(ctx)` method
  - Models: Identity, Account
  - Unit tests with httptest
  - **Acceptance**: Identity.Get works, ≥85% coverage

- Task 2.1.2: Implement Boards service
  - `boards.go` with BoardsService
  - Methods: List, Get, Create, Update, Delete
  - Models: Board, BoardCreateOptions, BoardUpdateOptions
  - Unit tests with httptest
  - **Acceptance**: All board methods work, ≥85% coverage

#### Subphase 2.2: Cards
**Dependencies**: 2.1

- Task 2.2.1: Implement Cards service - CRUD
  - `cards.go` with CardsService
  - Methods: List, ListAll, Get, Create, Update, Delete
  - Models: Card, CardCreateOptions, CardUpdateOptions, CardListOptions
  - Unit tests with httptest
  - **Acceptance**: CRUD methods work, ≥85% coverage

- Task 2.2.2: Implement Cards service - Actions
  - Methods: Close, Reopen, Postpone, Triage, Assign, Unassign, Tag, Untag, Watch, Unwatch, Golden, Ungolden
  - Unit tests with httptest
  - **Acceptance**: All action methods work, ≥85% coverage

### Phase 3: API Resources - Extended
**Dependencies**: Phase 2

#### Subphase 3.1: Comments and Reactions
**Dependencies**: 2.2

- Task 3.1.1: Implement Comments service
  - `comments.go` with CommentsService
  - Methods: List, Create, Update, Delete
  - Models: Comment, CommentCreateOptions, CommentUpdateOptions
  - Unit tests with httptest
  - **Acceptance**: All comment methods work, ≥85% coverage

- Task 3.1.2: Implement Reactions service
  - `reactions.go` with ReactionsService
  - Methods: List, Create, Delete
  - Models: Reaction, ReactionCreateOptions
  - Unit tests with httptest
  - **Acceptance**: All reaction methods work, ≥85% coverage

#### Subphase 3.2: Steps and Tags
**Dependencies**: 3.1

- Task 3.2.1: Implement Steps service
  - `steps.go` with StepsService
  - Methods: List, Get, Create, Update, Delete
  - Models: Step, StepCreateOptions, StepUpdateOptions
  - Unit tests with httptest
  - **Acceptance**: All step methods work, ≥85% coverage

- Task 3.2.2: Implement Tags service
  - `tags.go` with TagsService
  - Methods: List, Create
  - Models: Tag, TagCreateOptions
  - Unit tests with httptest
  - **Acceptance**: All tag methods work, ≥85% coverage

### Phase 4: API Resources - Management
**Dependencies**: Phase 3

#### Subphase 4.1: Columns, Users, Notifications
**Dependencies**: 3.2

- Task 4.1.1: Implement Columns service
  - `columns.go` with ColumnsService
  - Methods: List, Get, Create, Update, Delete
  - Models: Column, ColumnCreateOptions, ColumnUpdateOptions
  - Unit tests with httptest
  - **Acceptance**: All column methods work, ≥85% coverage

- Task 4.1.2: Implement Users service
  - `users.go` with UsersService
  - Methods: List
  - Models: User
  - Unit tests with httptest
  - **Acceptance**: User methods work, ≥85% coverage

- Task 4.1.3: Implement Notifications service
  - `notifications.go` with NotificationsService
  - Methods: List, Read, Unread, ReadAll
  - Models: Notification
  - Unit tests with httptest
  - **Acceptance**: All notification methods work, ≥85% coverage

#### Subphase 4.2: File Uploads
**Dependencies**: 4.1

- Task 4.2.1: Implement Uploads service
  - `uploads.go` with UploadsService
  - Methods: CreateDirectUpload, UploadFile (convenience wrapper)
  - Handle multipart uploads
  - Models: DirectUploadRequest, DirectUploadResponse
  - Unit tests with httptest
  - **Acceptance**: File uploads work, ≥85% coverage

### Phase 5: Polish and Documentation
**Dependencies**: Phase 4

#### Subphase 5.1: Examples and Documentation
**Dependencies**: 4.2

- Task 5.1.1: Create README.md
  - Installation instructions
  - Quick start example
  - API overview with links
  - Configuration options
  - Error handling examples
  - **Acceptance**: README complete and accurate

- Task 5.1.2: Create example programs
  - `examples/basic/` - simple card creation
  - `examples/pagination/` - manual and automatic pagination
  - `examples/error_handling/` - typed error handling
  - **Acceptance**: All examples run successfully

- Task 5.1.3: Add godoc comments
  - Package-level documentation
  - All exported types and methods documented
  - Follow go.dev/doc/comment standards
  - **Acceptance**: `go doc` output is comprehensive

#### Subphase 5.2: Final Validation
**Dependencies**: 5.1

- Task 5.2.1: Run full test suite
  - `task test:coverage` - ensure ≥85% coverage
  - `task lint` - no linting errors
  - `task fmt` - all code formatted
  - **Acceptance**: All tasks pass

- Task 5.2.2: Integration testing (optional)
  - Create integration test suite requiring live API token
  - Test against real Fizzy API (separate from unit tests)
  - **Acceptance**: Integration tests pass (if run)

- Task 5.2.3: Version and release
  - Tag v1.0.0
  - Create GitHub release with notes
  - **Acceptance**: Release published

## Testing Strategy

### Unit Tests
- Use `httptest.NewServer` for mocking Fizzy API
- Table-driven tests for each service method
- Test success cases and all error conditions
- Test retry logic (429, 5xx)
- Test caching (ETag, 304)
- Test pagination (manual and automatic)
- Use testify/assert and testify/require

### Coverage Requirements
- ≥85% overall coverage
- ≥85% per-package coverage
- Exclude: examples/, integration tests

### Test Organization
- `*_test.go` files alongside implementation
- Shared test helpers in `testing.go` (if needed)
- Integration tests in separate package with build tag: `//go:build integration`

### Example Test Structure
```go
func TestCardsService_Get(t *testing.T) {
    tests := []struct {
        name       string
        cardID     string
        response   string
        statusCode int
        wantErr    bool
        errType    error
    }{
        {"success", "card-123", `{"id":"card-123","title":"Test"}`, 200, false, nil},
        {"not found", "invalid", `{"error":"Not found"}`, 404, true, &fizzy.NotFoundError{}},
        {"unauthorized", "card-123", `{"error":"Unauthorized"}`, 401, true, &fizzy.AuthenticationError{}},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tt.statusCode)
                w.Write([]byte(tt.response))
            }))
            defer server.Close()
            
            client := fizzy.NewClient("token", "account", fizzy.WithBaseURL(server.URL))
            card, err := client.Cards.Get(context.Background(), tt.cardID)
            
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errType != nil {
                    assert.ErrorAs(t, err, &tt.errType)
                }
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.cardID, card.ID)
        })
    }
}
```

## Deployment

### Distribution
- Publish to GitHub: `github.com/yourusername/libfizz-go`
- Users install via: `go get github.com/yourusername/libfizz-go/fizzy`
- Follow semantic versioning (v1.0.0, v1.1.0, etc.)

### Release Process
1. Run `task check` - ensure all tests pass
2. Update version in README and godoc
3. Create git tag: `git tag v1.0.0`
4. Push tag: `git push origin v1.0.0`
5. Create GitHub release with changelog
6. Go module proxy automatically indexes new version

### Dependencies
- Go 1.21+ (for standard library features)
- github.com/stretchr/testify (testing only)
- No runtime dependencies beyond standard library

## Risks and Mitigations

### Risk 1: Fizzy API Changes
**Impact**: Breaking changes in Fizzy API could break library
**Mitigation**: 
- Version library semantically
- Monitor Fizzy GitHub repo for API changes
- Maintain backwards compatibility where possible
- Document API version compatibility in README

### Risk 2: Rate Limiting
**Impact**: Heavy usage could hit rate limits despite retries
**Mitigation**:
- Automatic retries with exponential backoff
- ETag caching reduces request volume
- Document rate limits in README
- Provide hooks for custom backoff strategies if needed later

### Risk 3: Test Coverage for All Endpoints
**Impact**: Many endpoints to test, risk of incomplete coverage
**Mitigation**:
- Systematic table-driven tests for each endpoint
- Enforce ≥85% coverage in CI
- Use httptest for deterministic testing
- Prioritize core resources (Boards, Cards) in Phase 2

### Risk 4: Context Cancellation Edge Cases
**Impact**: Context cancellation during retries could cause subtle bugs
**Mitigation**:
- Test context cancellation explicitly
- Use `context.Cause()` to distinguish cancellation vs timeout
- Respect context in all retry loops

## Appendix: Interview Q&A

**Q1: Authentication approach**
A: Personal Access Token only (Bearer token authentication)

**Q2: API coverage scope**
A: All resources from day one (comprehensive v1.0)

**Q3: Error handling strategy**
A: Typed errors + automatic retries (exponential backoff for 429/5xx, 3 attempts)

**Q4: HTTP client caching strategy**
A: Built-in ETag caching (automatic, enabled by default, configurable)

**Q5: Client configuration approach**
A: Functional options pattern (`WithTimeout`, `WithCache`, etc.)

**Q6: API client structure**
A: Namespaced services (`client.Boards.List()`, `client.Cards.Get()`)

**Q7: Pagination handling**
A: Helper + manual (return NextPage cursor, provide `ListAll()` helper)

## References

- Fizzy API Documentation: https://github.com/basecamp/fizzy/blob/main/docs/API.md
- Fizzy Website: https://fizzy.do
- Warping Go Standards: warping/languages/go.md
- Reference Implementations:
  - fizzy-mcp (TypeScript): https://github.com/Fabric-Pro/fizzy-mcp
  - fizzy-api-client (Python): https://libraries.io/pypi/fizzy-api-client
  - fizzy-cli (Ruby): https://github.com/robzolkos/fizzy-cli
