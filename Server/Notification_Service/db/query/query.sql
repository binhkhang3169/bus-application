-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    type,
    title,
    message
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetNotificationsByUserID :many
SELECT * FROM notifications
WHERE user_id = $1 OR user_id IS NULL -- Lấy cả thông báo chung và riêng cho user
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetBroadcastNotifications :many
SELECT * FROM notifications
WHERE user_id IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: MarkNotificationAsRead :one
UPDATE notifications
SET is_read = TRUE, updated_at = NOW()
WHERE id = $1 AND (user_id = $2 OR user_id IS NULL) -- Đảm bảo user chỉ mark read thông báo của họ hoặc broadcast
RETURNING *;

-- name: MarkAllUserNotificationsAsRead :many
UPDATE notifications
SET is_read = TRUE, updated_at = NOW()
WHERE user_id = $1 AND is_read = FALSE
RETURNING *;


-- ========= QUERIES MỚI CHO FCM TOKENS =========

-- name: RegisterFCMToken :one
-- Sử dụng ON CONFLICT để xử lý việc đăng ký lại token đã tồn tại
-- (ví dụ: người dùng đăng xuất rồi đăng nhập lại trên cùng thiết bị)
INSERT INTO fcm_tokens (user_id, token)
VALUES ($1, $2)
ON CONFLICT (token) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    created_at = NOW()
RETURNING *;


-- name: GetFCMTokensByUserID :many
-- Lấy tất cả token của một user cụ thể
SELECT token FROM fcm_tokens
WHERE user_id = $1;

-- name: GetAllFCMTokens :many
-- Lấy tất cả token trong database để gửi broadcast
SELECT token FROM fcm_tokens;

-- name: DeleteFCMToken :exec
-- Xóa token khi người dùng đăng xuất hoặc không muốn nhận thông báo nữa
DELETE FROM fcm_tokens
WHERE token = $1;