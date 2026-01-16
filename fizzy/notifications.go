package fizzy

import (
	"context"
	"fmt"
)

// List retrieves all notifications for the authenticated user.
func (s *NotificationsService) List(ctx context.Context) ([]Notification, error) {
	path := "/my/notifications"

	var notifications []Notification
	if err := s.client.doRequest(ctx, "GET", path, nil, &notifications); err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	return notifications, nil
}

// Read marks a notification as read.
func (s *NotificationsService) Read(ctx context.Context, notificationID string) error {
	path := fmt.Sprintf("/my/notifications/%s/read", notificationID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}

// Unread marks a notification as unread.
func (s *NotificationsService) Unread(ctx context.Context, notificationID string) error {
	path := fmt.Sprintf("/my/notifications/%s/unread", notificationID)

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to mark notification as unread: %w", err)
	}

	return nil
}

// ReadAll marks all notifications as read.
func (s *NotificationsService) ReadAll(ctx context.Context) error {
	path := "/my/notifications/read_all"

	if err := s.client.doRequest(ctx, "POST", path, nil, nil); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}
