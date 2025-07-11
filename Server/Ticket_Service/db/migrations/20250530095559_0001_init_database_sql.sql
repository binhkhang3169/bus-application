-- +goose Up
-- +goose StatementBegin

CREATE TABLE Ticket (
  Ticket_Id VARCHAR(6) NOT NULL PRIMARY KEY,
  Trip_Id_Begin VARCHAR NOT NULL,
  Trip_Id_End VARCHAR,                             -- NULL for one-way tickets
  Type SMALLINT NOT NULL,                          -- 0: one-way, 1: round-trip
  Customer_Id INT,
  Phone VARCHAR(15),
  Email VARCHAR(100),
  Name VARCHAR(100),
  Price DECIMAL(10,2) NOT NULL,
  Status SMALLINT NOT NULL,                        -- E.g., 0: wait for payment, 1: active, 2: cancel
  Booking_Time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  Payment_Status SMALLINT NOT NULL,                -- 0: unpaid, 1: paid
  Booking_Channel SMALLINT NOT NULL,               -- 0: web, 1: app, 2: offline.
  Created_At TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  Updated_At TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  Policy_Id INT NOT NULL,
  Booked_By VARCHAR(50) -- Stores staff ID or 'customer'
);

CREATE TABLE Ticket_Details (
  Detail_Id SERIAL PRIMARY KEY,
  Ticket_Id VARCHAR(6) NOT NULL,
  Pickup_Location_Begin INT,
  Dropoff_Location_Begin INT,
  Pickup_Location_End INT,                         -- NULL for one-way tickets
  Dropoff_Location_End INT,                        -- NULL for one-way tickets
  Created_At TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  Updated_At TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (Ticket_Id) REFERENCES Ticket(Ticket_Id) ON DELETE CASCADE
);

CREATE TABLE Ticket_Logs (
  Log_Id SERIAL PRIMARY KEY,
  Ticket_Id VARCHAR(6) NOT NULL,
  Action VARCHAR(255) NOT NULL,
  Created_At TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (Ticket_Id) REFERENCES Ticket(Ticket_Id) ON DELETE CASCADE
);

CREATE TABLE seats (
    id SERIAL PRIMARY KEY,
    trip_id VARCHAR NOT NULL,
    seat_name VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(trip_id, seat_name)
);

CREATE TABLE seat_tickets (
    id SERIAL PRIMARY KEY,
    seat_id INT NOT NULL,
    ticket_id VARCHAR(6) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0,            -- 0: Pending, 1: Confirmed, 2: Cancelled, 3: Checked-in
    trip_id VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(ticket_id, seat_id),
    FOREIGN KEY (seat_id) REFERENCES seats(id) ON DELETE CASCADE,
    FOREIGN KEY (ticket_id) REFERENCES Ticket(Ticket_Id) ON DELETE CASCADE
);

CREATE TABLE checkins (
    id SERIAL PRIMARY KEY,
    seat_ticket_id INT NOT NULL,
    ticket_id VARCHAR(6) NOT NULL,
    trip_id VARCHAR NOT NULL,
    seat_name VARCHAR,
    checked_in_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    note TEXT,
    FOREIGN KEY (seat_ticket_id) REFERENCES seats(id) ON DELETE CASCADE,
    FOREIGN KEY (ticket_id) REFERENCES Ticket(Ticket_Id) ON DELETE CASCADE
);

CREATE TABLE outbox_events (
    id UUID PRIMARY KEY,
    topic VARCHAR(255) NOT NULL,
    key VARCHAR(255) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ticket_customer_id ON Ticket(Customer_Id);
CREATE INDEX idx_ticket_phone ON Ticket(Phone);
CREATE INDEX idx_ticket_details_ticket_id ON Ticket_Details(Ticket_Id);
CREATE INDEX idx_ticket_logs_ticket_id ON Ticket_Logs(Ticket_Id);
CREATE INDEX idx_seats_trip_id ON seats(trip_id);
CREATE INDEX idx_seat_tickets_seat_id ON seat_tickets(seat_id);
CREATE INDEX idx_seat_tickets_ticket_id ON seat_tickets(ticket_id);
CREATE INDEX idx_checkins_ticket_id ON checkins(ticket_id);
CREATE INDEX idx_checkins_trip_id ON checkins(trip_id);
CREATE INDEX idx_checkins_seat_ticket_id ON checkins(seat_ticket_id);
CREATE INDEX idx_ticket_trip_id_begin ON Ticket(Trip_Id_Begin);
CREATE INDEX idx_ticket_trip_id_end ON Ticket(Trip_Id_End);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.Updated_At = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_ticket_timestamp BEFORE UPDATE ON Ticket FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_ticket_details_timestamp BEFORE UPDATE ON Ticket_Details FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_seats_timestamp BEFORE UPDATE ON seats FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
CREATE TRIGGER set_seat_tickets_timestamp BEFORE UPDATE ON seat_tickets FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_seat_tickets_timestamp ON seat_tickets;
DROP TRIGGER IF EXISTS set_seats_timestamp ON seats;
DROP TRIGGER IF EXISTS set_ticket_details_timestamp ON Ticket_Details;
DROP TRIGGER IF EXISTS set_ticket_timestamp ON Ticket;
DROP FUNCTION IF EXISTS trigger_set_timestamp();

DROP INDEX IF EXISTS idx_checkins_seat_ticket_id;
DROP INDEX IF EXISTS idx_checkins_trip_id;
DROP INDEX IF EXISTS idx_checkins_ticket_id;
DROP INDEX IF EXISTS idx_seat_tickets_ticket_id;
DROP INDEX IF EXISTS idx_seat_tickets_seat_id;
DROP INDEX IF EXISTS idx_seats_trip_id;
DROP INDEX IF EXISTS idx_ticket_logs_ticket_id;
DROP INDEX IF EXISTS idx_ticket_details_ticket_id;
DROP INDEX IF EXISTS idx_ticket_phone;
DROP INDEX IF EXISTS idx_ticket_customer_id;
DROP INDEX IF EXISTS idx_ticket_trip_id_begin;
DROP INDEX IF EXISTS idx_ticket_trip_id_end;

DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS checkins;
DROP TABLE IF EXISTS seat_tickets;
DROP TABLE IF EXISTS seats;
DROP TABLE IF EXISTS Ticket_Logs;
DROP TABLE IF EXISTS Ticket_Details;
DROP TABLE IF EXISTS Ticket;

-- +goose StatementEnd
