package fizzy

import "time"

// Identity represents the current user's identity and account access.
type Identity struct {
	Accounts []Account `json:"accounts"`
}

// Account represents a Fizzy account.
type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

// User represents a user in Fizzy.
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Role         string    `json:"role,omitempty"`
	Active       bool      `json:"active"`
	EmailAddress *string   `json:"email_address,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	URL          string    `json:"url,omitempty"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
}

// Board represents a Fizzy board.
type Board struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	AllAccess   bool      `json:"all_access,omitempty"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	URL         string    `json:"url,omitempty"`
	Creator     *User     `json:"creator,omitempty"`
}

// BoardCreateOptions contains options for creating a board.
type BoardCreateOptions struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// BoardUpdateOptions contains options for updating a board.
type BoardUpdateOptions struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Position    *int    `json:"position,omitempty"`
}

// Card represents a Fizzy card.
type Card struct {
	ID                 string     `json:"id"`
	Number             int        `json:"number"` // Card number used in URLs
	BoardID            string     `json:"board_id,omitempty"`
	Title              string     `json:"title"`
	Description        *string    `json:"description,omitempty"`        // Plain text description
	DescriptionHTML    *string    `json:"description_html,omitempty"`   // HTML formatted description
	Body               *string    `json:"body,omitempty"`               // Alias for Description (deprecated, use Description)
	ImageURL           *string    `json:"image_url,omitempty"`          // Attached image URL
	Status             string     `json:"status"`                        // "published", "closed", "maybe", "not_now"
	Closed             bool       `json:"closed,omitempty"`             // Whether card is closed
	ColumnID           *string    `json:"column_id,omitempty"`
	Position           int        `json:"position"`
	Golden             bool       `json:"golden"`
	LastActiveAt       *time.Time `json:"last_active_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at,omitempty"`
	ClosedAt           *time.Time `json:"closed_at,omitempty"`
	Board              *Board     `json:"board,omitempty"`              // Full board object
	Creator            *User      `json:"creator,omitempty"`
	Assignees          []User     `json:"assignees,omitempty"`
	HasMoreAssignees   bool       `json:"has_more_assignees,omitempty"`
	Tags               []Tag      `json:"tags,omitempty"`
	URL                string     `json:"url"`
	CommentsURL        string     `json:"comments_url,omitempty"`
	CommentsCount      int        `json:"comments_count,omitempty"`
}

// CardCreateOptions contains options for creating a card.
type CardCreateOptions struct {
	BoardID string  `json:"board_id"`
	Title   string  `json:"title"`
	Body    *string `json:"body,omitempty"`
}

// CardUpdateOptions contains options for updating a card.
type CardUpdateOptions struct {
	Title *string `json:"title,omitempty"`
	Body  *string `json:"body,omitempty"`
}

// CardListOptions contains options for listing cards.
type CardListOptions struct {
	BoardID  string   `json:"board_id,omitempty"`
	Status   string   `json:"status,omitempty"`
	TagIDs   []string `json:"tag_ids,omitempty"`
	ColumnID string   `json:"column_id,omitempty"`
	Page     string   `json:"page,omitempty"`
}

// Comment represents a comment on a card.
type Comment struct {
	ID        string    `json:"id"`
	CardID    string    `json:"card_id"`
	Body      string    `json:"body"`
	PlainText string    `json:"plain_text"`
	HTML      string    `json:"html"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Creator   *User     `json:"creator,omitempty"`
}

// CommentCreateOptions contains options for creating a comment.
type CommentCreateOptions struct {
	Body string `json:"body"`
}

// CommentUpdateOptions contains options for updating a comment.
type CommentUpdateOptions struct {
	Body string `json:"body"`
}

// Reaction represents a reaction to a comment.
type Reaction struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"` // Emoji, max 16 chars
	CreatedAt time.Time `json:"created_at"`
	Creator   *User     `json:"creator,omitempty"`
}

// ReactionCreateOptions contains options for creating a reaction.
type ReactionCreateOptions struct {
	Content string `json:"content"` // Emoji, max 16 chars
}

// Step represents a checklist step on a card.
type Step struct {
	ID        string    `json:"id"`
	CardID    string    `json:"card_id"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StepCreateOptions contains options for creating a step.
type StepCreateOptions struct {
	Content   string `json:"content"`
	Completed *bool  `json:"completed,omitempty"`
}

// StepUpdateOptions contains options for updating a step.
type StepUpdateOptions struct {
	Content   *string `json:"content,omitempty"`
	Completed *bool   `json:"completed,omitempty"`
}

// Tag represents a tag that can be applied to cards.
type Tag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// TagCreateOptions contains options for creating a tag.
type TagCreateOptions struct {
	Name  string  `json:"name"`
	Color *string `json:"color,omitempty"`
}

// Column represents a column on a board.
type Column struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ColumnCreateOptions contains options for creating a column.
type ColumnCreateOptions struct {
	Name string `json:"name"`
}

// ColumnUpdateOptions contains options for updating a column.
type ColumnUpdateOptions struct {
	Name     *string `json:"name,omitempty"`
	Position *int    `json:"position,omitempty"`
}

// Notification represents a notification.
type Notification struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	CardID    string     `json:"card_id"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// DirectUploadRequest contains parameters for requesting a direct upload URL.
type DirectUploadRequest struct {
	Blob DirectUploadBlob `json:"blob"`
}

// DirectUploadBlob contains file metadata for direct upload.
type DirectUploadBlob struct {
	Filename    string `json:"filename"`
	ByteSize    int64  `json:"byte_size"`
	Checksum    string `json:"checksum"` // Base64-encoded MD5 hash
	ContentType string `json:"content_type"`
}

// DirectUploadResponse contains the response from a direct upload request.
type DirectUploadResponse struct {
	DirectUploadURL string            `json:"direct_upload_url"`
	Headers         map[string]string `json:"headers"`
	BlobID          string            `json:"blob_id"`
}
