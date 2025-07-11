-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    shipments (
        id SERIAL PRIMARY KEY,
        trip_id INTEGER NOT NULL,
        sender_name VARCHAR(255) NOT NULL,
        receiver_name VARCHAR(255) NOT NULL,
        item_name VARCHAR(255) NOT NULL,
        item_type VARCHAR(50) NOT NULL,
        weight NUMERIC(10, 2) NOT NULL,
        length NUMERIC(10, 2) NOT NULL,
        width NUMERIC(10, 2) NOT NULL,
        height NUMERIC(10, 2) NOT NULL,
        volume NUMERIC(10, 2) NOT NULL,
        price NUMERIC(10, 2) NOT NULL, -- Giá cước do frontend cung cấp
        payer_type VARCHAR(50) NOT NULL CHECK (payer_type IN ('sender', 'receiver')), -- Người trả phí
        note TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_shipments_trip_id ON shipments (trip_id);

CREATE TABLE
    invoices (
        id SERIAL PRIMARY KEY,
        shipment_id INTEGER NOT NULL UNIQUE REFERENCES shipments (id) ON DELETE CASCADE,
        amount NUMERIC(10, 2) NOT NULL,
        issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_invoices_shipment_id ON invoices (shipment_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invoices;

DROP TABLE IF EXISTS shipments;

-- +goose StatementEnd