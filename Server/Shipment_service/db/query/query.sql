

-- name: CreateShipment :one
INSERT INTO shipments (
    trip_id, sender_name, receiver_name, item_name, item_type,
    weight, length, width, height, volume, price, payer_type, note, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW()
) RETURNING *;

-- name: GetShipmentByID :one
SELECT * FROM shipments
WHERE id = $1 LIMIT 1;

-- name: ListShipments :many
SELECT * FROM shipments
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: ListShipmentsByTripID :many
SELECT * FROM shipments
WHERE trip_id = $1
ORDER BY created_at DESC;

-- name: CreateInvoice :one
INSERT INTO invoices (
    shipment_id, amount, issued_at, created_at
) VALUES (
    $1, $2, $3, NOW()
) RETURNING *;

-- name: GetInvoiceByID :one
SELECT * FROM invoices
WHERE id = $1 LIMIT 1;

-- name: GetInvoiceByShipmentID :one
SELECT * FROM invoices
WHERE shipment_id = $1 LIMIT 1;

-- name: ListInvoices :many
SELECT * FROM invoices
ORDER BY issued_at DESC
LIMIT $1
OFFSET $2;

-- name: ListInvoicesByTripID :many
SELECT inv.* FROM invoices inv
JOIN shipments s ON inv.shipment_id = s.id
WHERE s.trip_id = $1
ORDER BY inv.issued_at DESC;