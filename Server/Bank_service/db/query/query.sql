-- name: CreateAccount :one
INSERT INTO accounts (
  owner_name,
  balance,
  currency,
  status
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE; -- Để tránh deadlock khi cập nhật balance

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateAccountStatus :one
UPDATE accounts
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :exec
-- Thực tế không xóa, chỉ dùng để minh họa, chúng ta sẽ dùng UpdateAccountStatus
DELETE FROM accounts
WHERE id = $1;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount), updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: CreateTransactionHistory :one
INSERT INTO transaction_history (
  account_id,
  transaction_type,
  amount,
  currency,
  description
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListTransactionHistoryByAccountID :many
SELECT * FROM transaction_history
WHERE account_id = $1
ORDER BY created_at DESC -- Show newest first
LIMIT $2
OFFSET $3;