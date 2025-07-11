-- name: CreateInvoice :one
INSERT INTO invoices (
    invoice_id,
    invoice_number,
    invoice_type,
    customer_id,
    ticket_id,
    total_amount,
    discount_amount,
    tax_amount,
    final_amount,
    currency,
    payment_status,
    payment_method,
    issue_date,
    notes,
    -- vnpay fields
    vnpay_txn_ref,
    vnpay_bank_code,
    vnpay_txn_no,
    vnpay_pay_date,
    -- stripe fields
    stripe_payment_intent_id,
    stripe_charge_id,
    stripe_customer_id,
    stripe_payment_method_details,
    -- bank transfer fields
    bank_transfer_code,
    bank_payment_details -- Other bank fields like account_name, account_number, bank_name might be updated later upon confirmation
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24
) RETURNING *;

-- name: GetInvoiceByID :one
SELECT * FROM invoices
WHERE invoice_id = $1 LIMIT 1;

-- name: GetInvoiceByVNPayTxnRef :one
SELECT * FROM invoices
WHERE vnpay_txn_ref = $1 LIMIT 1;

-- name: GetInvoiceByStripePaymentIntentID :one
SELECT * FROM invoices
WHERE stripe_payment_intent_id = $1 LIMIT 1;

-- name: GetInvoiceByBankTransferCode :one
SELECT * FROM invoices
WHERE bank_transfer_code = $1 LIMIT 1;

-- name: GetLatestCompletedInvoiceByTicketID :one
SELECT * FROM invoices
WHERE ticket_id = $1 AND payment_status = 'COMPLETED'
ORDER BY created_at DESC
LIMIT 1;

-- name: ListInvoicesByCustomerID :many
SELECT * FROM invoices
WHERE customer_id = $1
ORDER BY created_at DESC;

-- name: UpdateInvoiceVNPayStatus :one
UPDATE invoices
SET
    payment_status = $2,
    vnpay_bank_code = $3,
    vnpay_txn_no = $4,
    vnpay_pay_date = $5,
    updated_at = NOW()
WHERE vnpay_txn_ref = $1
RETURNING *;

-- name: UpdateInvoiceStripePaymentSuccess :one
UPDATE invoices
SET
    payment_status = $2,
    stripe_charge_id = $3,
    stripe_payment_method_details = $4,
    updated_at = NOW()
WHERE stripe_payment_intent_id = $1
RETURNING *;

-- name: UpdateInvoiceStripePaymentIntent :one
-- Used when creating Payment Intent and need to update invoice with PI ID
UPDATE invoices
SET
    stripe_payment_intent_id = $2,
    payment_method = $3, -- 'STRIPE'
    payment_status = $4, -- 'PENDING' or 'REQUIRES_PAYMENT_METHOD'
    updated_at = NOW()
WHERE invoice_id = $1
RETURNING *;

-- name: UpdateInvoiceBankPaymentRequest :one
-- Used when creating a bank payment request (invoice is PENDING or AWAITING_CONFIRMATION)
UPDATE invoices
SET
    payment_method = $2, -- 'BANK'
    payment_status = $3, -- 'AWAITING_CONFIRMATION' or 'PENDING'
    bank_transfer_code = $4, -- The code the user should use for the transfer
    notes = $5, -- Instructions for bank payment
    updated_at = NOW()
WHERE invoice_id = $1
RETURNING *;

-- name: UpdateInvoiceBankPaymentConfirmation :one
-- Used when an admin/system confirms a bank payment
UPDATE invoices
SET
    payment_status = $2, -- 'COMPLETED'
    bank_account_name = $3,
    bank_account_number = $4,
    bank_name = $5,
    bank_transaction_id = $6,
    bank_payment_details = $7, -- Append confirmation details
    notes = $8, -- Append confirmation notes
    updated_at = NOW()
WHERE invoice_id = $1 -- Could also be bank_transfer_code
RETURNING *;


-- name: UpdateInvoicePaymentFailed :one
-- Update status when payment fails (generic for Stripe, VNPay, or Bank based on invoice_id)
UPDATE invoices
SET
    payment_status = $2, -- 'FAILED'
    notes = $3, -- Reason for failure
    updated_at = NOW()
WHERE invoice_id = $1
RETURNING *;

-- name: UpdateInvoiceStatusGeneral :one
-- Update general status (e.g., REFUNDED, CANCELED)
UPDATE invoices
SET
    payment_status = $2,
    notes = $3, -- Notes for refund, cancellation, etc.
    updated_at = NOW()
WHERE invoice_id = $1
RETURNING *;
