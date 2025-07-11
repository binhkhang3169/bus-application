
-- name: CreateNews :one
INSERT INTO news (
    title,
    image_url,
    content,
    created_by
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetNewsByID :one
SELECT * FROM news
WHERE id = $1 LIMIT 1;

-- name: ListNews :many
SELECT * FROM news
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: UpdateNews :one
UPDATE news
SET
    title = $2,
    image_url = $3,
    content = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteNews :exec
DELETE FROM news
WHERE id = $1;

-- name: CountNews :one
SELECT COUNT(*) FROM news;