package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"ticket-service/domain/models" // For domain-specific models/constants
	"ticket-service/internal/db"   // sqlc generated package
	"ticket-service/pkg/utils"

	"github.com/redis/go-redis/v9"
)

type UpdateStatusTransactionParams struct {
	TicketID            string
	PaymentStatus       int16
	GeneralTicketStatus int16
	SeatTicketStatus    int16
	OutboxEvents        []db.CreateOutboxEventParams
}

type ManagerTicketInterface interface {
	CreateManagerTicket(ctx context.Context, seat *db.Seat) error // Use db.Seat
	UpdateStatusByTicketID(ctx context.Context, ticketID string, statusCode string) error
	GetSeatsByTripID(ctx context.Context, tripID string) ([]db.Seat, error) // Use db.Seat
	UpdateStatusInTransaction(ctx context.Context, params UpdateStatusTransactionParams) error
}

type ManagerTicketRepository struct {
	sqlDB      *sql.DB // For BeginTx
	q          *db.Queries
	redis      *redis.Client
	ticketRepo TicketRepositoryInterface // Changed to interface type for flexibility
	logger     utils.Logger
}

func NewManagerTicket(sqlDB *sql.DB, redis *redis.Client, ticketRepo TicketRepositoryInterface, logger utils.Logger) ManagerTicketInterface {
	return &ManagerTicketRepository{
		sqlDB:      sqlDB,
		q:          db.New(sqlDB),
		redis:      redis,
		ticketRepo: ticketRepo,
		logger:     logger,
	}
}

func (m *ManagerTicketRepository) CreateManagerTicket(ctx context.Context, seat *db.Seat) error {
	// The db.CreateSeatParams takes TripID and SeatName.
	// The input `seat *db.Seat` should have these fields populated.
	// The ID, CreatedAt, UpdatedAt will be set by the database and returned.
	createdSeat, err := m.q.CreateSeat(ctx, db.CreateSeatParams{
		TripID:   seat.TripID,
		SeatName: seat.SeatName,
	})
	if err != nil {
		return err
	}
	// Update the input seat object with the generated ID and timestamps
	seat.ID = createdSeat.ID
	seat.CreatedAt = createdSeat.CreatedAt
	seat.UpdatedAt = createdSeat.UpdatedAt
	return nil
}

func (m *ManagerTicketRepository) GetSeatsByTripID(ctx context.Context, tripID string) ([]db.Seat, error) {
	return m.q.GetSeatsByTripID(ctx, tripID)
}

func (m *ManagerTicketRepository) UpdateStatusInTransaction(ctx context.Context, params UpdateStatusTransactionParams) error {
	tx, err := m.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := m.q.WithTx(tx)

	// 1. Update ticket's main status
	updateTicketParams := db.UpdateTicketPaymentStatusParams{
		TicketID:      params.TicketID,
		PaymentStatus: params.PaymentStatus,
		Status:        params.GeneralTicketStatus,
	}
	if _, err := qtx.UpdateTicketPaymentStatus(ctx, updateTicketParams); err != nil {
		return fmt.Errorf("failed to update ticket payment status: %w", err)
	}

	// 2. Update status of associated seat_tickets
	updateSeatStatusParams := db.UpdateSeatTicketStatusByTicketIDParams{
		TicketID: params.TicketID,
		Status:   params.SeatTicketStatus,
	}
	if _, err := qtx.UpdateSeatTicketStatusByTicketID(ctx, updateSeatStatusParams); err != nil {
		return fmt.Errorf("failed to update seat_tickets status: %w", err)
	}

	// 3. Create all required outbox events
	for _, event := range params.OutboxEvents {
		if err := qtx.CreateOutboxEvent(ctx, event); err != nil {
			return fmt.Errorf("failed to create outbox event for topic %s: %w", event.Topic, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// Cập nhật trạng thái vé theo payment callback
func (m *ManagerTicketRepository) UpdateStatusByTicketID(ctx context.Context, ticketID string, statusCode string) error {
	var paymentStatus int16
	var generalTicketStatus int16 // Required for UpdateTicketPaymentStatus sqlc query

	if statusCode == "1" { // Payment successful
		paymentStatus = int16(models.PaymentStatusPaid)
		generalTicketStatus = int16(models.TicketStatusConfirmed) // Assuming successful payment confirms the ticket
	} else if statusCode == "2" { // Payment failed
		paymentStatus = int16(models.PaymentStatusFailed)
		generalTicketStatus = int16(models.TicketStatusCancelled) // Assuming failed payment cancels the ticket
	} else {
		return fmt.Errorf("invalid status code: %s", statusCode)
	}

	tx, err := m.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := m.q.WithTx(tx)

	// Cập nhật trạng thái thanh toán và trạng thái chung của vé
	// sqlc.UpdateTicketPaymentStatus updates both payment_status and general status.
	updateTicketParams := db.UpdateTicketPaymentStatusParams{
		TicketID:      ticketID,
		PaymentStatus: paymentStatus,
		Status:        generalTicketStatus, // General ticket status
	}
	_, err = qtx.UpdateTicketPaymentStatus(ctx, updateTicketParams)
	if err != nil {
		return fmt.Errorf("failed to update ticket payment status and general status: %w", err)
	}

	// Cập nhật trạng thái seat_tickets
	var seatTicketStatus int16
	if statusCode == "1" { // Payment successful
		seatTicketStatus = int16(models.SeatStatusConfirmed)
	} else if statusCode == "2" { // Payment failed
		seatTicketStatus = int16(models.SeatStatusCancelled)
	}
	// Only update seat tickets if status is 1 or 2
	if statusCode == "1" || statusCode == "2" {
		updateSeatStatusParams := db.UpdateSeatTicketStatusByTicketIDParams{
			TicketID: ticketID,
			Status:   seatTicketStatus,
		}
		_, err = qtx.UpdateSeatTicketStatusByTicketID(ctx, updateSeatStatusParams)
		if err != nil {
			return fmt.Errorf("failed to update seat_tickets status: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
