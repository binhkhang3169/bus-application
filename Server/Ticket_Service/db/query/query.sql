-- name: CreateSeat :one
-- Inserts a new seat for a trip.
INSERT INTO seats (trip_id, seat_name)
VALUES ($1, $2)
RETURNING *;

-- name: GetSeatByID :one
-- Retrieves a specific seat by its ID.
SELECT * FROM seats
WHERE id = $1;

-- name: GetSeatsByTripID :many
-- Retrieves all seats for a given trip_id.
SELECT * FROM seats
WHERE trip_id = $1
ORDER BY seat_name;

-- name: CreateTicket :one
-- Inserts a new ticket record for both one-way and round-trip.
INSERT INTO Ticket (
    Ticket_Id, Trip_Id_Begin, Trip_Id_End, Type, Customer_Id, Phone, Email, Name, Price, Status,
    Booking_Time, Payment_Status, Booking_Channel, Policy_Id, Booked_By
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;

-- name: CreateTicketDetails :one
-- Inserts details for a ticket.
INSERT INTO Ticket_Details (
    Ticket_Id, Pickup_Location_Begin, Dropoff_Location_Begin, Pickup_Location_End, Dropoff_Location_End
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: CreateSeatTicket :one
-- Links a seat to a ticket, now includes trip_id directly.
INSERT INTO seat_tickets (seat_id, ticket_id, status, trip_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTicketCore :one
-- Retrieves core ticket information by Ticket_Id.
SELECT * FROM Ticket
WHERE Ticket_Id = $1;

-- name: GetTicketDetailsByTicketID :many
-- Retrieves all details for a given Ticket_Id.
SELECT * FROM Ticket_Details
WHERE Ticket_Id = $1;

-- name: GetSeatTicketsByTicketID :many
-- Retrieves all seat_ticket entries for a given Ticket_Id, no longer needs JOIN with seats.
SELECT st.*, s.seat_name
FROM seat_tickets st
JOIN seats s ON st.seat_id = s.id
WHERE st.ticket_id = $1;

-- name: GetTicketByPhoneAndIDCore :one
-- Retrieves core ticket information by Ticket_Id and Phone.
SELECT * FROM Ticket
WHERE Ticket_Id = $1 AND Phone = $2;

-- name: GetTicketsByCustomerIDCore :many
-- Retrieves all core ticket information for a given Customer_Id.
SELECT * FROM Ticket
WHERE Customer_Id = $1
ORDER BY Booking_Time DESC;

-- name: IsSeatBookedOnTrip :one
-- Checks if a specific seat_id is booked on a specific trip - now uses trip_id directly.
SELECT EXISTS (
    SELECT 1
    FROM seat_tickets st
    WHERE st.seat_id = $1 AND st.trip_id = $2 AND st.status IN (0, 1)
);

-- name: IsSeatGenerallyBooked :one
-- Checks if a specific seat_id is currently booked or pending (status 0 or 1).
SELECT EXISTS (
    SELECT 1
    FROM seat_tickets
    WHERE seat_id = $1 AND status IN (0, 1)
);

-- name: UpdateTicketStatus :one
-- Updates the status of a ticket.
UPDATE Ticket
SET Status = $2, Updated_At = CURRENT_TIMESTAMP
WHERE Ticket_Id = $1
RETURNING *;

-- name: UpdateTicketPaymentStatus :one
-- Updates the payment_status and general status of a ticket.
UPDATE Ticket
SET Payment_Status = $2, Status = $3, Updated_At = CURRENT_TIMESTAMP
WHERE Ticket_Id = $1
RETURNING *;

-- name: UpdateSeatTicketStatus :one
-- Updates the status of a seat_ticket entry by its ID.
UPDATE seat_tickets
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateSeatTicketStatusByTicketID :many
-- Updates the status of all seat_ticket entries for a given ticket_id.
UPDATE seat_tickets
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE ticket_id = $1
RETURNING *;

-- name: GetSeatIDsByTicketID :many
-- Retrieves seat_ids associated with a ticket_id.
SELECT seat_id FROM seat_tickets
WHERE ticket_id = $1;

-- name: ListAvailableSeatsByTripID :many
-- Lists all seats for a trip_id that are not in seat_tickets or have status 2 (cancelled).
SELECT s.id, s.trip_id, s.seat_name
FROM seats s
WHERE s.trip_id = $1
AND NOT EXISTS (
    SELECT 1
    FROM seat_tickets st
    WHERE st.seat_id = s.id AND st.status IN (0, 1, 3) -- 0: pending, 1: paid, 3: checked-in
)
ORDER BY s.seat_name;

-- name: CreateTicketLog :one
-- Inserts a log entry for a ticket action.
INSERT INTO Ticket_Logs (Ticket_Id, Action)
VALUES ($1, $2)
RETURNING *;

-- name: GetSeatTicketAndSeatInfoByTicketID :one
-- Retrieves seat_ticket and associated seat details for a given ticket_id.
SELECT
    st.id as seat_ticket_id, st.seat_id, st.ticket_id, st.status as seat_ticket_status, st.trip_id,
    s.id as seat_table_id, s.seat_name
FROM seat_tickets st
JOIN seats s ON st.seat_id = s.id
WHERE st.ticket_id = $1 AND st.status = 1; -- Typically checkin for confirmed/paid tickets

-- name: GetSeatTicketByID :one
-- Retrieves a specific seat_ticket by its ID - now uses trip_id directly.
SELECT st.*, s.seat_name
FROM seat_tickets st
JOIN seats s ON st.seat_id = s.id
WHERE st.seat_id = $1 and st.ticket_id = $2 and st.status  = 1; -- Typically checkin for confirmed/paid tickets;


-- name: CreateCheckin :one
-- Inserts a new checkin record - can now get trip_id from seat_tickets directly.
INSERT INTO checkins (seat_ticket_id, ticket_id, trip_id, seat_name, note)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: UpdateTicketStatusAfterCheckin :one
-- Updates the ticket's main status to 'used'.
UPDATE Ticket
SET Status = $2, Updated_At = CURRENT_TIMESTAMP -- $2 would be the 'used' status
WHERE Ticket_Id = $1
RETURNING *;

-- name: UpdateSeatTicketStatusAfterCheckin :one
-- Updates the seat_ticket status to 'checked-in'.
UPDATE seat_tickets
SET status = $2, updated_at = CURRENT_TIMESTAMP -- $2 would be the 'checked-in' status
WHERE id = $1 -- seat_ticket.id
RETURNING *;

-- name: GetTicketStatus :one
-- Retrieves just the status of a ticket.
SELECT status FROM Ticket
WHERE Ticket_Id = $1;

-- name: GetSeatTicketStatus :one
-- Retrieves just the status of a seat_ticket.
SELECT status FROM seat_tickets
WHERE id = $1;

-- name: AreSeatsAvailable :many
-- Checks if a list of seats are available (not booked or held).
SELECT s.id,
       EXISTS(
           SELECT 1
           FROM seat_tickets st
           WHERE st.seat_id = s.id AND st.status IN (0, 1, 3) -- 0: pending, 1: paid, 3: checked-in
       ) as is_booked
FROM seats s
WHERE s.id = ANY(@seat_ids::int[]);


-- name: CreateOutboxEvent :exec
-- For Transactional Outbox Pattern
INSERT INTO outbox_events (id, topic, key, payload)
VALUES ($1, $2, $3, $4);

-- name: GetOutboxEvents :many
-- For the Outbox Poller/Relay
SELECT * FROM outbox_events
ORDER BY created_at
LIMIT $1;

-- name: DeleteOutboxEvents :exec
-- For the Outbox Poller/Relay
DELETE FROM outbox_events
WHERE id = ANY(@event_ids::uuid[]);

-- name: GetAllCheckinsByTripID :many
SELECT * FROM checkins
WHERE trip_id = $1
ORDER BY checked_in_at DESC;

-- name: GetAllTickets :many
-- Retrieves a paginated list of all tickets, ordered by booking time.
SELECT * FROM Ticket
ORDER BY Booking_Time DESC
LIMIT $1
OFFSET $2;

-- name: GetTotalTicketCount :one
-- Retrieves the total number of tickets.
SELECT COUNT(*) FROM Ticket;

-- name: GetSeatTicketsByTicketIDs :many
-- Retrieves all seat_ticket entries for a given list of Ticket_Ids.
SELECT st.*, s.seat_name
FROM seat_tickets st
JOIN seats s ON st.seat_id = s.id
WHERE st.ticket_id = ANY(@ticket_ids::varchar[]);

-- name: GetTicketDetailsByTicketIDs :many
-- Retrieves all details for a given list of Ticket_Ids.
SELECT * FROM Ticket_Details
WHERE Ticket_Id = ANY(@ticket_ids::varchar[]);

