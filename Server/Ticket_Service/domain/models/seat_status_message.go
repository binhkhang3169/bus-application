package models

import "time"

// Add to models package (maybe create seat_models.go)
type SeatStatusMessage struct {
	EventType string    `json:"event_type"` // "seat_reserved", "seat_released", "payment_completed", "payment_failed"
	TicketID  string    `json:"ticket_id"`
	SeatID    int32     `json:"seat_id"`
	TripID    string    `json:"trip_id"`
	Timestamp time.Time `json:"timestamp"`
}
