-- +goose Up
-- +goose StatementBegin
-- Create extension for UUID generation if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create invoices table
CREATE TABLE
    IF NOT EXISTS invoices (
        invoice_id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        invoice_number VARCHAR(50) UNIQUE NOT NULL,
        invoice_type VARCHAR(50),
        customer_id VARCHAR(100) NOT NULL,
        ticket_id VARCHAR(100) NOT NULL,
        total_amount DECIMAL(15, 2) NOT NULL,
        discount_amount DECIMAL(15, 2) DEFAULT 0.00,
        tax_amount DECIMAL(15, 2) DEFAULT 0.00,
        final_amount DECIMAL(15, 2) NOT NULL,
        currency VARCHAR(10),
        payment_status VARCHAR(20) DEFAULT 'PENDING', -- PENDING, COMPLETED, FAILED, REFUNDED, AWAITING_CONFIRMATION (for bank)
        payment_method VARCHAR(50), -- VNPAY, STRIPE, BANK
        issue_date TIMESTAMP,
        notes TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        -- VNPay specific fields
        vnpay_txn_ref VARCHAR(100) UNIQUE,
        vnpay_bank_code VARCHAR(50),
        vnpay_txn_no VARCHAR(100),
        vnpay_pay_date VARCHAR(50),
        -- Stripe specific fields
        stripe_payment_intent_id VARCHAR(255) UNIQUE,
        stripe_charge_id VARCHAR(255),
        stripe_customer_id VARCHAR(255),
        stripe_payment_method_details TEXT,
        -- Bank Transfer specific fields
        bank_transfer_code VARCHAR(100) UNIQUE, -- Unique code for user to reference or system generated
        bank_account_name VARCHAR(255), -- Name of account holder who paid (optional, if user provides)
        bank_account_number VARCHAR(100), -- Account number from which payment was made (optional)
        bank_name VARCHAR(100), -- Name of the bank (optional)
        bank_transaction_id VARCHAR(255), -- Bank's transaction ID after confirmation (optional)
        bank_payment_details TEXT -- Additional details, e.g., screenshot URL, user notes for payment
    );

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_invoices_customer_id ON invoices (customer_id);

CREATE INDEX IF NOT EXISTS idx_invoices_ticket_id ON invoices (ticket_id);

CREATE INDEX IF NOT EXISTS idx_invoices_payment_status ON invoices (payment_status);

CREATE INDEX IF NOT EXISTS idx_invoices_payment_method ON invoices (payment_method);

CREATE INDEX IF NOT EXISTS idx_invoices_created_at ON invoices (created_at);

CREATE INDEX IF NOT EXISTS idx_invoices_vnpay_txn_ref ON invoices (vnpay_txn_ref);

CREATE INDEX IF NOT EXISTS idx_invoices_stripe_payment_intent_id ON invoices (stripe_payment_intent_id);

CREATE INDEX IF NOT EXISTS idx_invoices_stripe_charge_id ON invoices (stripe_charge_id);

CREATE INDEX IF NOT EXISTS idx_invoices_bank_transfer_code ON invoices (bank_transfer_code);

CREATE INDEX IF NOT EXISTS idx_invoices_bank_transaction_id ON invoices (bank_transaction_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invoices;

DROP EXTENSION IF EXISTS "uuid-ossp";

-- +goose StatementEnd