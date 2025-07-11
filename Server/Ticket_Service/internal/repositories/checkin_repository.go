package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"ticket-service/internal/db" // sqlc generated package

	// For domain-specific models/constants not in db
	"ticket-service/pkg/utils"
	// "github.com/jackc/pgx/v4" // No longer needed for pgx.ErrNoRows directly if using sql.ErrNoRows
)

type CheckinRepositoryInterface interface {
	// Note: The returned db.SeatTicket and db.Seat might have fewer fields (e.g. CreatedAt, UpdatedAt will be zero)
	// if GetSeatTicketAndSeatInfoByTicketID sqlc query doesn't select them.
	// Also, the sqlc query GetSeatTicketAndSeatInfoByTicketID has an additional condition `st.status = 1`.
	GetSeatTicketAndSeatInfoByTicketID(ctx context.Context, ticketID string) (*db.SeatTicket, *db.Seat, error)
	GetSeatTicketDetails(ctx context.Context, seatTicketID int32, ticketId string) (db.GetSeatTicketByIDRow, error)
	PerformCheckin(ctx context.Context, seatTicketID int32, ticketID string, tripID string, seatName sql.NullString, note string, newSeatTicketStatus int16, newTicketStatus int16) (*db.Checkin, error)
	// Mới: Thêm phương thức lấy tất cả check-in theo tripID
	GetAllCheckinsByTripID(ctx context.Context, tripID string) ([]db.Checkin, error)
}

type CheckinRepository struct {
	sqlDB  *sql.DB // For BeginTx
	q      *db.Queries
	logger utils.Logger
}

func NewCheckinRepository(sqlDB *sql.DB, logger utils.Logger) CheckinRepositoryInterface {
	return &CheckinRepository{
		sqlDB:  sqlDB,
		q:      db.New(sqlDB),
		logger: logger,
	}
}

// GetSeatTicketAndSeatInfoByTicketID retrieves seat_ticket and associated seat details.
// Note: This implementation uses the sqlc generated GetSeatTicketAndSeatInfoByTicketID.
// The sqlc query fetches st.id as seat_ticket_id, st.status as seat_ticket_status, s.id as seat_table_id.
// It does NOT fetch created_at/updated_at for seat_tickets or seats.
// The sqlc query also has a condition `AND st.status = 1` which was not in the original manual query.
func (r *CheckinRepository) GetSeatTicketAndSeatInfoByTicketID(ctx context.Context, ticketID string) (*db.SeatTicket, *db.Seat, error) {
	row, err := r.q.GetSeatTicketAndSeatInfoByTicketID(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows { // pgx.ErrNoRows is mapped to sql.ErrNoRows by the driver
			return nil, nil, fmt.Errorf("no seat_ticket found for ticket_id %s with status 1: %w", ticketID, err)
		}
		r.logger.Error("Error getting seat_ticket and seat info by ticket ID %s: %v", ticketID, err)
		return nil, nil, fmt.Errorf("database error: %w", err)
	}

	seatTicket := &db.SeatTicket{
		ID:       row.SeatTicketID,
		SeatID:   row.SeatID,
		TicketID: row.TicketID,
		Status:   row.SeatTicketStatus,
	}
	seat := &db.Seat{
		ID:       row.SeatTableID,
		TripID:   row.TripID,
		SeatName: row.SeatName,
	}

	return seatTicket, seat, nil
}

func (r *CheckinRepository) GetSeatTicketDetails(ctx context.Context, seatTicketID int32, ticketId string) (db.GetSeatTicketByIDRow, error) {
	arg := db.GetSeatTicketByIDParams{
		SeatID:   seatTicketID,
		TicketID: ticketId,
	}
	return r.q.GetSeatTicketByID(ctx, arg)
}

// Mới: Implement phương thức lấy danh sách check-in
func (r *CheckinRepository) GetAllCheckinsByTripID(ctx context.Context, tripID string) ([]db.Checkin, error) {
	checkins, err := r.q.GetAllCheckinsByTripID(ctx, tripID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []db.Checkin{}, nil // Trả về slice rỗng nếu không có kết quả
		}
		r.logger.Error("Error getting checkins by trip ID %s: %v", tripID, err)
		return nil, fmt.Errorf("database error when fetching checkins for trip %s: %w", tripID, err)
	}
	return checkins, nil
}

// PerformCheckin creates a checkin record and updates statuses within a transaction.
func (r *CheckinRepository) PerformCheckin(
	ctx context.Context,
	seatTicketID int32, // sqlc uses int32 for ID
	ticketID string,
	tripID string,
	seatName sql.NullString, // sqlc uses sql.NullString
	note string,
	newSeatTicketStatus int16, // sqlc uses int16 for status
	newTicketStatus int16,
) (*db.Checkin, error) {
	tx, err := r.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	checkinParams := db.CreateCheckinParams{
		SeatTicketID: seatTicketID,
		TicketID:     ticketID,
		TripID:       tripID,
		SeatName:     seatName,
		Note:         note,
	}
	createdCheckin, err := qtx.CreateCheckin(ctx, checkinParams)
	if err != nil {
		r.logger.Error("Failed to insert into checkins: %v", err)
		return nil, fmt.Errorf("failed to create checkin record: %w", err)
	}

	updateSeatTicketParams := db.UpdateSeatTicketStatusAfterCheckinParams{
		ID:     createdCheckin.ID,
		Status: newSeatTicketStatus,
	}
	_, err = qtx.UpdateSeatTicketStatusAfterCheckin(ctx, updateSeatTicketParams)
	if err != nil {
		r.logger.Error("Failed to update seat_tickets status: %v", err)
		return nil, fmt.Errorf("failed to update seat_ticket status: %w", err)
	}

	updateTicketParams := db.UpdateTicketStatusAfterCheckinParams{
		TicketID: ticketID,
		Status:   newTicketStatus,
	}
	_, err = qtx.UpdateTicketStatusAfterCheckin(ctx, updateTicketParams)
	if err != nil {
		r.logger.Error("Failed to update Ticket status: %v", err)
		return nil, fmt.Errorf("failed to update ticket status: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &createdCheckin, nil
}
