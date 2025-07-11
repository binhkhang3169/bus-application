-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Bảng lưu trữ thông báo
CREATE TABLE
    notifications (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id VARCHAR(255) NULL, -- NULL nếu là thông báo cho tất cả user
        type VARCHAR(50) NOT NULL,
        title VARCHAR(255) NOT NULL,
        message TEXT NOT NULL,
        is_read BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- Bảng mới để lưu FCM tokens
CREATE TABLE
    fcm_tokens (
        user_id VARCHAR(255) NOT NULL,
        token TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        -- Mỗi token chỉ tồn tại một lần
        PRIMARY KEY (token)
    );

-- Tạo các chỉ mục để tăng tốc độ truy vấn
CREATE INDEX idx_notifications_user_id ON notifications (user_id);
CREATE INDEX idx_fcm_tokens_user_id ON fcm_tokens (user_id);

-- Thêm comment để giải thích
COMMENT ON COLUMN notifications.user_id IS 'NULL for broadcast messages, user identifier for specific messages';
COMMENT ON TABLE fcm_tokens IS 'Stores FCM device tokens for users';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- +goose StatementEnd