// notification-service/internal/model/notification.go
package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Notification represents a notification message.
// Chúng ta có thể sử dụng trực tiếp repository.Notification do sqlc tạo ra.
// Tuy nhiên, nếu cần custom fields hoặc validation riêng ở tầng này, bạn có thể định nghĩa ở đây.
// Ví dụ:
type Notification struct {
	ID        pgtype.UUID        `json:"id"`
	UserID    pgtype.Text        `json:"user_id,omitempty"` // Sử dụng pgtype.Text cho nullable string
	Type      string             `json:"type"`
	Title     string             `json:"title"`
	Message   string             `json:"message"`
	IsRead    pgtype.Bool        `json:"is_read"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

// CreateNotificationRequest defines the structure for creating a new notification.
type CreateNotificationRequest struct {
	UserID  *string `json:"user_id"` // Pointer để có thể là nil (broadcast)
	Type    string  `json:"type" binding:"required"`
	Title   string  `json:"title" binding:"required"`
	Message string  `json:"message" binding:"required"`
}

// Struct mới cho request đăng ký FCM token
type RegisterFCMTokenRequest struct {
	Token string `json:"token" binding:"required"`
}
