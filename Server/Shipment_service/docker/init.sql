CREATE TABLE
    shipments (
        id SERIAL PRIMARY KEY,
        trip_id INTEGER NOT NULL, -- Sẽ liên kết với bảng trips (nếu có)
        sender_name VARCHAR(255) NOT NULL,
        receiver_name VARCHAR(255) NOT NULL,
        item_name VARCHAR(255) NOT NULL,
        item_type VARCHAR(50) NOT NULL, -- vd: "document", "electronics", "furniture"
        weight NUMERIC(10, 2) NOT NULL, -- Khối lượng tính bằng kg
        length NUMERIC(10, 2) NOT NULL, -- Kích thước tính bằng cm
        width NUMERIC(10, 2) NOT NULL, -- Kích thước tính bằng cm
        height NUMERIC(10, 2) NOT NULL, -- Kích thước tính bằng cm
        volume NUMERIC(10, 2) NOT NULL, -- Thể tích thực tế tính bằng cm^3
        price NUMERIC(10, 2) NOT NULL, -- Giá cước vận chuyển
        note TEXT,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_shipments_trip_id ON shipments (trip_id);

CREATE TABLE
    invoices (
        id SERIAL PRIMARY KEY,
        shipment_id INTEGER NOT NULL UNIQUE REFERENCES shipments (id) ON DELETE CASCADE,
        amount NUMERIC(10, 2) NOT NULL, -- Số tiền hóa đơn (bằng giá của lô hàng)
        issued_at TIMESTAMPTZ NOT NULL DEFAULT NOW (), -- Ngày phát hành hóa đơn
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
        -- updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW() -- Nếu hóa đơn có thể cập nhật
    );

CREATE INDEX idx_invoices_shipment_id ON invoices (shipment_id);