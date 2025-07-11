package models

import (
	"database/sql"
	"ticket-service/internal/db"
	"time"
)

// TicketInput defines the structure for creating a new ticket, supporting both one-way and round-trip.
type TicketInput struct {
	TicketType     int     `json:"ticket_type"` // 0 for one-way, 1 for round-trip
	Price          float64 `json:"price"`
	Status         int16   `json:"status"`
	PaymentStatus  int16   `json:"payment_status"`
	BookingChannel int16   `json:"booking_channel"`
	PolicyID       int32   `json:"policy_id"`
	Phone          string  `json:"phone,omitempty"`
	Email          string  `json:"email,omitempty"`
	Name           string  `json:"name"`
	BookedBy       string  `json:"booked_by,omitempty"`

	// Outbound Trip Information
	TripIDBegin          string  `json:"trip_id_begin"`
	SeatIDBegin          []int32 `json:"seat_id_begin"`
	PickupLocationBegin  int32   `json:"pickup_location_begin"`
	DropoffLocationBegin int32   `json:"dropoff_location_begin"`

	// Return Trip Information (optional, used for round-trip)
	TripIDEnd          string  `json:"trip_id_end,omitempty"`
	SeatIDEnd          []int32 `json:"seat_id_end,omitempty"`
	PickupLocationEnd  int32   `json:"pickup_location_end,omitempty"`
	DropoffLocationEnd int32   `json:"dropoff_location_end,omitempty"`
}

type TicketInfoInput struct {
	Phone    string `json:"phone,omitempty"`
	TicketID string `json:"ticket_id"`
}

// TicketReturn defines the structure for returning ticket data to the client.
type TicketReturn struct {
	TicketID         string         `json:"ticket_id"`
	Type             int16          `json:"type"`
	TripIDBegin      string         `json:"trip_id_begin"`
	TripIDEnd        sql.NullString `json:"trip_id_end,omitempty"`
	CustomerID       sql.NullInt32  `json:"customer_id"`
	Phone            sql.NullString `json:"phone,omitempty"`
	Email            sql.NullString `json:"email,omitempty"`
	Name             sql.NullString `json:"name"`
	Price            float64        `json:"price"`
	Status           int16          `json:"status"`
	BookingTime      time.Time      `json:"booking_time"`
	PaymentStatus    int16          `json:"payment_status"`
	BookingChannel   int16          `json:"booking_channel"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	PolicyID         int32          `json:"policy_id"`
	BookedBy         sql.NullString `json:"booked_by"`
	Details          []db.TicketDetail
	SeatTicketsBegin []db.GetSeatTicketsByTicketIDRow
	SeatTicketsEnd   []db.GetSeatTicketsByTicketIDRow // Updated to use the generated row struct
	TripDetails      *TripInfo                        `json:"trip_details,omitempty"`
}

// VehicleInfo corresponds to the Java Vehicle entity (customize fields as needed)
type VehicleInfo struct {
	ID           int    `json:"id"`
	LicensePlate string `json:"licensePlate,omitempty"` // Example field
	Type         string `json:"type,omitempty"`         // Example field
	// Add other relevant fields from your Java Vehicle entity
}

// SpecialDayInfo corresponds to the Java SpecialDay entity (customize fields as needed)
type SpecialDayInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"` // Example field
	// Add other relevant fields
}

// RouteInfo corresponds to the Java Route entity (customize fields as needed)
type RouteInfo struct {
	ID            int    `json:"id"`
	Origin        string `json:"origin,omitempty"`        // Example field
	Destination   string `json:"destination,omitempty"`   // Example field
	DepartureStop string `json:"departureStop,omitempty"` // Example field
	ArrivalStop   string `json:"arrivalStop,omitempty"`   // Example field
	// Add other relevant fields
}

// TripInfo is the Go representation of the Trip data from trip-service
type TripInfo struct {
	ID            string          `json:"id"`            // Changed to string to match Trip_Id format
	DepartureDate string          `json:"departureDate"` // e.g., "YYYY-MM-DD"
	DepartureTime string          `json:"departureTime"` // e.g., "HH:MM:SS"
	ArrivalDate   string          `json:"arrivalDate"`
	ArrivalTime   string          `json:"arrivalTime"`
	Vehicle       *VehicleInfo    `json:"vehicle,omitempty"` // Use pointer if it can be null
	DriverID      *int            `json:"driverId,omitempty"`
	TotalSeats    int             `json:"total"` // Corresponds to 'total' in Java Trip
	StockSeats    int             `json:"stock"` // Corresponds to 'stock' in Java Trip
	Status        int             `json:"status"`
	Special       *SpecialDayInfo `json:"special,omitempty"`   // Use pointer if it can be null
	Route         *RouteInfo      `json:"route,omitempty"`     // Use pointer if it can be null
	CreatedAt     string          `json:"createdAt,omitempty"` // Expecting ISO 8601 string or timestamp string
	CreatedBy     *int            `json:"createdBy,omitempty"`
}

type PaginatedTickets struct {
	Tickets []*TicketReturn `json:"tickets"`
	Total   int64           `json:"total"`
	Page    int             `json:"page"`
	Limit   int             `json:"limit"`
}
