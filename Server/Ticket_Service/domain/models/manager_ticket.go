package models

type ManagerTicket struct {
	TripID   string  `json:"trip_id"`
	SeatID   string  `json:"seat_id"`
	TicketID *string `json:"ticket_id"` // null khi chưa có vé
}
