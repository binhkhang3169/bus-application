// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: query.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createNotification = `-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    type,
    title,
    message
) VALUES (
    $1, $2, $3, $4
) RETURNING id, user_id, type, title, message, is_read, created_at, updated_at
`

type CreateNotificationParams struct {
	UserID  pgtype.Text `json:"user_id"`
	Type    string      `json:"type"`
	Title   string      `json:"title"`
	Message string      `json:"message"`
}

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error) {
	row := q.db.QueryRow(ctx, createNotification,
		arg.UserID,
		arg.Type,
		arg.Title,
		arg.Message,
	)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Type,
		&i.Title,
		&i.Message,
		&i.IsRead,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFCMToken = `-- name: DeleteFCMToken :exec
DELETE FROM fcm_tokens
WHERE token = $1
`

// Xóa token khi người dùng đăng xuất hoặc không muốn nhận thông báo nữa
func (q *Queries) DeleteFCMToken(ctx context.Context, token string) error {
	_, err := q.db.Exec(ctx, deleteFCMToken, token)
	return err
}

const getAllFCMTokens = `-- name: GetAllFCMTokens :many
SELECT token FROM fcm_tokens
`

// Lấy tất cả token trong database để gửi broadcast
func (q *Queries) GetAllFCMTokens(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, getAllFCMTokens)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		items = append(items, token)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBroadcastNotifications = `-- name: GetBroadcastNotifications :many
SELECT id, user_id, type, title, message, is_read, created_at, updated_at FROM notifications
WHERE user_id IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type GetBroadcastNotificationsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetBroadcastNotifications(ctx context.Context, arg GetBroadcastNotificationsParams) ([]Notification, error) {
	rows, err := q.db.Query(ctx, getBroadcastNotifications, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Notification{}
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Title,
			&i.Message,
			&i.IsRead,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFCMTokensByUserID = `-- name: GetFCMTokensByUserID :many
SELECT token FROM fcm_tokens
WHERE user_id = $1
`

// Lấy tất cả token của một user cụ thể
func (q *Queries) GetFCMTokensByUserID(ctx context.Context, userID string) ([]string, error) {
	rows, err := q.db.Query(ctx, getFCMTokensByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		items = append(items, token)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNotificationsByUserID = `-- name: GetNotificationsByUserID :many
SELECT id, user_id, type, title, message, is_read, created_at, updated_at FROM notifications
WHERE user_id = $1 OR user_id IS NULL -- Lấy cả thông báo chung và riêng cho user
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type GetNotificationsByUserIDParams struct {
	UserID pgtype.Text `json:"user_id"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

func (q *Queries) GetNotificationsByUserID(ctx context.Context, arg GetNotificationsByUserIDParams) ([]Notification, error) {
	rows, err := q.db.Query(ctx, getNotificationsByUserID, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Notification{}
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Title,
			&i.Message,
			&i.IsRead,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markAllUserNotificationsAsRead = `-- name: MarkAllUserNotificationsAsRead :many
UPDATE notifications
SET is_read = TRUE, updated_at = NOW()
WHERE user_id = $1 AND is_read = FALSE
RETURNING id, user_id, type, title, message, is_read, created_at, updated_at
`

func (q *Queries) MarkAllUserNotificationsAsRead(ctx context.Context, userID pgtype.Text) ([]Notification, error) {
	rows, err := q.db.Query(ctx, markAllUserNotificationsAsRead, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Notification{}
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Title,
			&i.Message,
			&i.IsRead,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markNotificationAsRead = `-- name: MarkNotificationAsRead :one
UPDATE notifications
SET is_read = TRUE, updated_at = NOW()
WHERE id = $1 AND (user_id = $2 OR user_id IS NULL) -- Đảm bảo user chỉ mark read thông báo của họ hoặc broadcast
RETURNING id, user_id, type, title, message, is_read, created_at, updated_at
`

type MarkNotificationAsReadParams struct {
	ID     pgtype.UUID `json:"id"`
	UserID pgtype.Text `json:"user_id"`
}

func (q *Queries) MarkNotificationAsRead(ctx context.Context, arg MarkNotificationAsReadParams) (Notification, error) {
	row := q.db.QueryRow(ctx, markNotificationAsRead, arg.ID, arg.UserID)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Type,
		&i.Title,
		&i.Message,
		&i.IsRead,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const registerFCMToken = `-- name: RegisterFCMToken :one

INSERT INTO fcm_tokens (user_id, token)
VALUES ($1, $2)
ON CONFLICT (token) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    created_at = NOW()
RETURNING user_id, token, created_at
`

type RegisterFCMTokenParams struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// ========= QUERIES MỚI CHO FCM TOKENS =========
// Sử dụng ON CONFLICT để xử lý việc đăng ký lại token đã tồn tại
// (ví dụ: người dùng đăng xuất rồi đăng nhập lại trên cùng thiết bị)
func (q *Queries) RegisterFCMToken(ctx context.Context, arg RegisterFCMTokenParams) (FcmToken, error) {
	row := q.db.QueryRow(ctx, registerFCMToken, arg.UserID, arg.Token)
	var i FcmToken
	err := row.Scan(&i.UserID, &i.Token, &i.CreatedAt)
	return i, err
}
