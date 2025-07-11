CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Cho phép sử dụng uuid_generate_v4()
CREATE TABLE
    news (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        title VARCHAR(255) NOT NULL,
        image_url TEXT,
        content TEXT NOT NULL,
        created_by VARCHAR(100) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_news_created_at ON news (created_at DESC);