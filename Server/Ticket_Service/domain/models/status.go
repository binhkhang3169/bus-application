package models

import (
	"database/sql"
	"time"
)

// ... other existing constants

// SeatTicket Statuses
const (
	SeatStatusPendingPayment = 0 // Đang đợi thanh toán
	SeatStatusCheckedIn      = 3 // Đã check-in
	SeatStatusMissed         = 4 // Bỏ lỡ chuyến (No-show)
	TicketStatusConfirmed    = 1 // Đã xác nhận
	SeatStatusAvailable      = 1 // Trống
)

// Ticket Statuses (adjust as needed for overall ticket lifecycle)

// Payment Statuses (already defined in your example)
// const (
//  PaymentStatusUnpaid = 0
//  PaymentStatusPaid   = 1
//  PaymentStatusFailed = 2
// )

// Checkin model
type Checkin struct {
	ID           int       `json:"id"`
	SeatTicketID int       `json:"seat_ticket_id"`
	TicketID     string    `json:"ticket_id"`
	TripID       string    `json:"trip_id"`
	SeatName     string    `json:"seat_name"`
	CheckedInAt  time.Time `json:"checked_in_at"`
	Note         string    `json:"note"`
}

// CheckinResponse for API
type CheckinResponse struct {
	TicketID    string         `json:"ticket_id"`
	TripID      string         `json:"trip_id"`
	SeatName    sql.NullString `json:"seat_name"`
	CheckedInAt sql.NullTime   `json:"checked_in_at"`
	CheckinNote string         `json:"checkin_note"`
	Message     string         `json:"message"`
}
