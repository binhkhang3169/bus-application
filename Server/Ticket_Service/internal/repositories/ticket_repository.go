package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"ticket-service/domain/models"
	"ticket-service/internal/db" // sqlc generated package
	"ticket-service/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
	// "github.com/jackc/pgx/v4" // No longer needed for pgx.ErrNoRows
)

// Các mã lỗi để xác định rõ loại lỗi
var (
	ErrRedisNotAvailable = errors.New("redis service is not available")
	ErrSeatAlreadyHeld   = errors.New("seat is already held by another user")
	ErrLockNotAcquired   = errors.New("could not acquire lock")
	ErrCacheMiss         = errors.New("cache miss")
)

// NEW: Define params for the new atomic creation method
type CreateTicketTransactionParams struct {
	Ticket       db.CreateTicketParams
	TicketDetail db.CreateTicketDetailsParams
	SeatIDsBegin []int32
	SeatIDsEnd   []int32
	OutboxEvents []db.CreateOutboxEventParams
}

// TicketRepositoryInterface defines the methods for ticket repository
// Note: Renamed original TicketRepository to TicketRepositoryInterface for clarity with implementation struct
type TicketRepositoryInterface interface {
	GetTicketByID(ctx context.Context, ticketID string) (*models.TicketReturn, error)
	GetTicketByCustomer(ctx context.Context, customerID sql.NullInt32) ([]*models.TicketReturn, error) // Use sql.NullInt32 from sqlc
	GetInfoTicketByPhone(ctx context.Context, info *models.TicketInfoInput) (*models.TicketReturn, error)
	FindTicketByTripAndSeat(ctx context.Context, seatID int32) (*db.SeatTicket, error) // seatID is int32 in db.SeatTicket
	AllTicketsStatus2BySeat(ctx context.Context, seatID int32) (bool, error)           // seatID is int32
	GenerateUniqueTicketID(ctx context.Context) (string, error)
	CreateTicket(ctx context.Context, ticket *db.Ticket, ticketDetail *db.TicketDetail, seatIDs []int32, tripID string) error // << UPDATED signature

	CreateTicketInTransaction(ctx context.Context, params CreateTicketTransactionParams) error

	UpdateTicketPaymentStatus(ctx context.Context, ticketID string, paymentStatus int16, generalStatus int16, tripID string) error
	UpdateSeatTicketsStatus(ctx context.Context, ticketID string, status int16) error // status is int16

	GetTicketFromCache(ctx context.Context, ticketID string) (*db.Ticket, error) // Cache db.Ticket
	CacheTicket(ctx context.Context, ticket *db.Ticket) error                    // Cache db.Ticket
	CacheTemporarySeat(ctx context.Context, seatID int32, tripID string, ticket *db.Ticket) error
	ExtendSeatHoldTime(ctx context.Context, seatID int32, ticketToExtend *db.Ticket, duration time.Duration) error
	ReleaseSeat(ctx context.Context, seatID int32) error // seatID consistent with db
	IsSeatHeld(ctx context.Context, seatID int32) (bool, error)

	AreSeatsAvailable(ctx context.Context, seatIDs []int32) ([]db.AreSeatsAvailableRow, error)

	GetAvailableSeatsByTripID(ctx context.Context, tripID string) ([]models.SeatReturn, error)
	AcquireLock(ctx context.Context, lockKey string, ttl time.Duration) (bool, func(), error)
	CleanupTicketCache(ctx context.Context, ticketID string, seatIDs []int32, tripID string) error

	SetupExpirationHandler(ctx context.Context) error
	UpdateTicketStatus(ctx context.Context, ticketID string, status int16) error // status is int16
	GetSeatIDsByTicketID(ctx context.Context, ticketID string) ([]int32, error)  // Returns []int32 from sqlc

	CacheAvailableSeats(ctx context.Context, tripID string, seats []models.SeatReturn) error
	GetAvailableSeatsFromCache(ctx context.Context, tripID string) ([]models.SeatReturn, error)
	PublishSeatStatusChange(ctx context.Context, message *models.SeatStatusMessage) error
	SubscribeToSeatStatusChanges(ctx context.Context) error
	UpdateCachedAvailableSeats(ctx context.Context, tripID string, seatIDs []int32, action string) error

	GetAllTickets(ctx context.Context, limit int, offset int) ([]*models.TicketReturn, error)
	GetTotalTicketCount(ctx context.Context) (int64, error)
}

type ticketRepositoryImpl struct {
	sqlDB  *sql.DB // For BeginTx and direct execution if no sqlc query exists
	q      *db.Queries
	redis  *redis.Client
	utils  *utils.Utils
	logger utils.Logger
}

func NewTicketRepository(sqlDB *sql.DB, redis *redis.Client, utils *utils.Utils, logger utils.Logger) TicketRepositoryInterface {
	return &ticketRepositoryImpl{
		sqlDB:  sqlDB,
		q:      db.New(sqlDB),
		redis:  redis,
		utils:  utils,
		logger: logger,
	}
}

// AreSeatsAvailable performs a single database query to check the status of multiple seats.
func (r *ticketRepositoryImpl) AreSeatsAvailable(ctx context.Context, seatIDs []int32) ([]db.AreSeatsAvailableRow, error) {
	if len(seatIDs) == 0 {
		return []db.AreSeatsAvailableRow{}, nil
	}
	return r.q.AreSeatsAvailable(ctx, seatIDs)
}

func (r *ticketRepositoryImpl) AcquireLock(ctx context.Context, lockKey string, ttl time.Duration) (bool, func(), error) {
	// ... (Implementation remains the same, but ensure TTL is passed and used)
	lockKey = "lock:" + lockKey
	lockAcquired, err := r.redis.SetNX(ctx, lockKey, "1", ttl).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil, nil // Not an error, just lock not acquired
		}
		r.logger.Error("Redis error on SetNX: %v", err)
		return false, nil, ErrRedisNotAvailable
	}
	if !lockAcquired {
		return false, nil, nil
	}
	unlock := func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := r.redis.Del(unlockCtx, lockKey).Err(); err != nil {
			r.logger.Error("Failed to release lock %s: %v", lockKey, err)
		}
	}
	return true, unlock, nil
}

// CreateTicketInTransaction creates a ticket, its details, seat assignments, and an outbox event within a single robust transaction.
func (r *ticketRepositoryImpl) CreateTicketInTransaction(ctx context.Context, params CreateTicketTransactionParams) error {
	tx, err := r.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// 1. Create Ticket
	createdTicket, err := qtx.CreateTicket(ctx, params.Ticket)
	if err != nil {
		return fmt.Errorf("failed to insert ticket: %w", err)
	}

	// 2. Create Ticket Details
	params.TicketDetail.TicketID = createdTicket.TicketID
	_, err = qtx.CreateTicketDetails(ctx, params.TicketDetail)
	if err != nil {
		return fmt.Errorf("failed to insert ticket details: %w", err)
	}

	// 3. Create Seat Tickets
	for _, seatID := range params.SeatIDsBegin {
		_, err = qtx.CreateSeatTicket(ctx, db.CreateSeatTicketParams{
			SeatID:   seatID,
			TicketID: createdTicket.TicketID,
			Status:   int16(models.SeatStatusPendingPayment),
			TripID:   params.Ticket.TripIDBegin,
		})
		if err != nil {
			return fmt.Errorf("failed to insert seat ticket for seat %d: %w", seatID, err)
		}
	}
	if createdTicket.Type == 1 {
		for _, seatID := range params.SeatIDsEnd {
			_, err = qtx.CreateSeatTicket(ctx, db.CreateSeatTicketParams{
				SeatID:   seatID,
				TicketID: createdTicket.TicketID,
				Status:   int16(models.SeatStatusPendingPayment),
				TripID:   params.Ticket.TripIDEnd.String,
			})
			if err != nil {
				return fmt.Errorf("failed to insert seat ticket for seat %d: %w", seatID, err)
			}
		}
	}

	// 4. Create Outbox Event
	for _, OutboxEvent := range params.OutboxEvents {
		err = qtx.CreateOutboxEvent(ctx, OutboxEvent)
		if err != nil {
			return fmt.Errorf("failed to create outbox event: %w", err)
		}
	}

	return tx.Commit()
}

func (r *ticketRepositoryImpl) CreateTicket(
	ctx context.Context,
	ticket *db.Ticket, // This will now have TripID after sqlc generate
	ticketDetail *db.TicketDetail,
	seatIDs []int32,
	tripID string, // << ADDED tripID as a parameter, or ensure ticket.TripID is set
) error {
	tx, err := r.sqlDB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	// Insert Ticket
	ticketParams := db.CreateTicketParams{
		TicketID:       ticket.TicketID,
		TripIDBegin:    ticket.TripIDBegin, // << ASSIGN TripID here; or ensure ticket.TripID is set correctly before calling
		TripIDEnd:      ticket.TripIDEnd,
		CustomerID:     ticket.CustomerID,
		Phone:          ticket.Phone,
		Type:           ticket.Type,
		Email:          ticket.Email,
		Name:           ticket.Name,
		Price:          ticket.Price,
		Status:         ticket.Status,
		BookingTime:    time.Now(), // sqlc query might use NOW(), check query.sql
		PaymentStatus:  ticket.PaymentStatus,
		BookingChannel: ticket.BookingChannel,
		PolicyID:       ticket.PolicyID,
		BookedBy:       ticket.BookedBy,
	}
	createdTicket, err := qtx.CreateTicket(ctx, ticketParams)
	if err != nil {
		return fmt.Errorf("failed to insert ticket: %w", err)
	}
	ticket.TicketID = createdTicket.TicketID
	ticket.TripIDBegin = createdTicket.TripIDBegin
	ticket.TripIDEnd = createdTicket.TripIDEnd // << Ensure this is updated back
	ticket.CreatedAt = createdTicket.CreatedAt
	ticket.UpdatedAt = createdTicket.UpdatedAt
	ticket.BookingTime = createdTicket.BookingTime

	// Insert TicketDetail
	ticketDetailParams := db.CreateTicketDetailsParams{
		TicketID:             createdTicket.TicketID,
		PickupLocationEnd:    ticketDetail.PickupLocationEnd,
		DropoffLocationEnd:   ticketDetail.DropoffLocationEnd,
		PickupLocationBegin:  ticketDetail.PickupLocationBegin,
		DropoffLocationBegin: ticketDetail.DropoffLocationBegin,
	}
	createdTicketDetail, err := qtx.CreateTicketDetails(ctx, ticketDetailParams)
	if err != nil {
		return fmt.Errorf("failed to insert ticket detail: %w", err)
	}
	ticketDetail.DetailID = createdTicketDetail.DetailID
	ticketDetail.TicketID = createdTicketDetail.TicketID
	ticketDetail.CreatedAt = createdTicketDetail.CreatedAt
	ticketDetail.UpdatedAt = createdTicketDetail.UpdatedAt

	for _, seatIDValue := range seatIDs {
		seatBooked, err := qtx.IsSeatGenerallyBooked(ctx, seatIDValue)
		if err != nil {
			return fmt.Errorf("failed to check seat availability for seat %d: %w", seatIDValue, err)
		}
		if seatBooked {
			return fmt.Errorf("seat %d was booked by another user while processing", seatIDValue)
		}

		seatTicketParams := db.CreateSeatTicketParams{
			SeatID:   seatIDValue,
			TicketID: createdTicket.TicketID,
			Status:   0, // Assuming 0 is 'pending' or initial status
		}
		_, err = qtx.CreateSeatTicket(ctx, seatTicketParams)
		if err != nil {
			return fmt.Errorf("failed to insert seat ticket for seat ID %d: %w", seatIDValue, err)
		}
	}

	return tx.Commit()
}

func (r *ticketRepositoryImpl) UpdateCachedAvailableSeats(ctx context.Context, tripID string, seatIDs []int32, action string) error {
	availableSeatsKey := fmt.Sprintf("available_seats:%s", tripID)

	// Convert seatIDs to a slice of interface{} for Redis commands
	members := make([]interface{}, len(seatIDs))
	for i, id := range seatIDs {
		members[i] = id
	}

	var err error
	if action == "REMOVE" {
		err = r.redis.SRem(ctx, availableSeatsKey, members...).Err()
	} else if action == "ADD" {
		err = r.redis.SAdd(ctx, availableSeatsKey, members...).Err()
	} else {
		return fmt.Errorf("invalid action: %s", action)
	}

	if err != nil {
		r.logger.Error("Failed to update available seats cache for trip %s: %v", tripID, err)
	}
	return err
}

func (r *ticketRepositoryImpl) GetTicketFromCache(ctx context.Context, ticketID string) (*db.Ticket, error) {
	val, err := r.redis.Get(ctx, "ticket:"+ticketID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("redis error: %w", err)
	}
	var ticket db.Ticket
	if err := json.Unmarshal([]byte(val), &ticket); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticket from cache: %w", err)
	}
	return &ticket, nil
}

func (r *ticketRepositoryImpl) CacheTicket(ctx context.Context, ticket *db.Ticket) error {
	// ticket now includes TripID, no change needed here as it marshals the whole struct.
	data, err := json.Marshal(ticket)
	if err != nil {
		return fmt.Errorf("failed to marshal ticket: %w", err)
	}
	return r.redis.Set(ctx, "ticket:"+ticket.TicketID, data, 1*time.Hour).Err()
}

func (r *ticketRepositoryImpl) CacheTemporarySeat(ctx context.Context, seatID int32, tripID string, ticket *db.Ticket) error {
	key := fmt.Sprintf("seat:%d", seatID)
	holdData := map[string]string{
		"ticket_id": ticket.TicketID,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	pipe := r.redis.Pipeline()
	pipe.HMSet(ctx, key, holdData)
	pipe.Expire(ctx, key, 10*time.Minute) // TODO: Make expiration configurable
	pipe.SRem(ctx, fmt.Sprintf("available_seats:%s", tripID), seatID)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *ticketRepositoryImpl) ExtendSeatHoldTime(ctx context.Context, seatID int32, ticketToExtend *db.Ticket, duration time.Duration) error {
	key := fmt.Sprintf("seat:%d", seatID)
	cachedTicketID, err := r.redis.HGet(ctx, key, "ticket_id").Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("seat hold not found or ticket_id field missing")
		}
		return fmt.Errorf("redis error checking seat hold: %w", err)
	}
	if cachedTicketID != ticketToExtend.TicketID {
		return fmt.Errorf("seat is held by a different ticket (expected %s, found %s)", ticketToExtend.TicketID, cachedTicketID)
	}
	if err := r.redis.Expire(ctx, key, duration).Err(); err != nil {
		return fmt.Errorf("failed to extend expiration for seat %d: %w", seatID, err)
	}
	return nil
}

func (r *ticketRepositoryImpl) ReleaseSeat(ctx context.Context, seatID int32) error {
	key := fmt.Sprintf("seat:%d", seatID)
	return r.redis.Del(ctx, key).Err()
}

func (r *ticketRepositoryImpl) IsSeatHeld(ctx context.Context, seatID int32) (bool, error) {
	key := fmt.Sprintf("seat:%d", seatID)
	exists, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, ErrRedisNotAvailable
	}
	return exists > 0, nil
}

// FindTicketByTripAndSeat: The provided sqlc queries do not have a direct equivalent for
// "SELECT * FROM seat_tickets WHERE seat_id = $1" that returns a single SeatTicket.
// GetSeatTicketByID takes the primary key of seat_tickets, not seat_id.
// This method will remain using direct SQL execution.
// Consider adding a specific query to query.sql if this is frequently used and should be managed by sqlc.
func (r *ticketRepositoryImpl) FindTicketByTripAndSeat(ctx context.Context, seatID int32) (*db.SeatTicket, error) {
	// This query is not available in the provided sqlc Querier.
	// It's kept as a direct SQL query.
	// For a full sqlc migration, this query should be added to your SQL files and re-generated.
	const query = "SELECT id, seat_id, ticket_id, status, created_at, updated_at FROM seat_tickets WHERE seat_id = $1 ORDER BY created_at DESC LIMIT 1" // Added order and limit to ensure one record if multiple exist
	row := r.sqlDB.QueryRowContext(ctx, query, seatID)
	var st db.SeatTicket
	err := row.Scan(&st.ID, &st.SeatID, &st.TicketID, &st.Status, &st.CreatedAt, &st.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &st, nil
}

func (r *ticketRepositoryImpl) AllTicketsStatus2BySeat(ctx context.Context, seatID int32) (bool, error) {
	// Uses sqlc.IsSeatGenerallyBooked which checks for status IN (0, 1)
	// The original logic returned !exists. So if IsSeatGenerallyBooked is true, it means it IS booked (status 0 or 1).
	// AllTicketsStatus2BySeat implies it's valid if no tickets are 0 or 1.
	isBooked, err := r.q.IsSeatGenerallyBooked(ctx, seatID)
	if err != nil {
		return false, err
	}
	return !isBooked, nil
}

func (r *ticketRepositoryImpl) GetTicketByCustomer(ctx context.Context, customerID sql.NullInt32) ([]*models.TicketReturn, error) {
	coreTickets, err := r.q.GetTicketsByCustomerIDCore(ctx, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.TicketReturn{}, nil
		}
		return nil, fmt.Errorf("failed to get core tickets by customer ID: %w", err)
	}

	var ticketsReturn []*models.TicketReturn
	for _, coreTicket := range coreTickets {
		t := &models.TicketReturn{
			TicketID:       coreTicket.TicketID,
			TripIDBegin:    coreTicket.TripIDBegin, // << ADDED
			TripIDEnd:      coreTicket.TripIDEnd,
			Type:           coreTicket.Type,
			CustomerID:     coreTicket.CustomerID,
			Phone:          coreTicket.Phone,
			Email:          coreTicket.Email,
			Name:           coreTicket.Name,
			Price:          coreTicket.Price,
			Status:         coreTicket.Status,
			BookingTime:    coreTicket.BookingTime,
			PaymentStatus:  coreTicket.PaymentStatus,
			BookingChannel: coreTicket.BookingChannel,
			CreatedAt:      coreTicket.CreatedAt,
			UpdatedAt:      coreTicket.UpdatedAt,
			PolicyID:       coreTicket.PolicyID,
			BookedBy:       coreTicket.BookedBy,
		}

		details, err := r.q.GetTicketDetailsByTicketID(ctx, coreTicket.TicketID)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to get ticket details for ticket %s: %w", coreTicket.TicketID, err)
		}
		for _, det := range details {
			t.Details = append(t.Details, db.TicketDetail{
				DetailID:             det.DetailID,
				TicketID:             det.TicketID,
				PickupLocationBegin:  det.PickupLocationBegin,
				DropoffLocationBegin: det.DropoffLocationBegin,
				PickupLocationEnd:    det.PickupLocationEnd,
				DropoffLocationEnd:   det.DropoffLocationEnd,
				CreatedAt:            det.CreatedAt,
				UpdatedAt:            det.UpdatedAt,
			})
		}

		seatTicketRows, err := r.q.GetSeatTicketsByTicketID(ctx, coreTicket.TicketID)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to get seat tickets for ticket %s: %w", coreTicket.TicketID, err)
		}
		// 		type GetSeatTicketsByTicketIDRow struct {
		// 	ID         int32          `json:"id"`
		// 	SeatID     int32          `json:"seat_id"`
		// 	TicketID   string         `json:"ticket_id"`
		// 	Status     int16          `json:"status"`
		// 	CreatedAt  sql.NullTime   `json:"created_at"`
		// 	UpdatedAt  sql.NullTime   `json:"updated_at"`
		// 	SeatName   sql.NullString `json:"seat_name"`
		// 	SeatTripID string         `json:"seat_trip_id"`
		// }
		for _, stRow := range seatTicketRows {
			if stRow.TripID == t.TripIDBegin {
				t.SeatTicketsBegin = append(t.SeatTicketsBegin, db.GetSeatTicketsByTicketIDRow{ // Map to models.SeatTicket
					ID:        stRow.ID,
					SeatID:    stRow.SeatID,
					TicketID:  stRow.TicketID,
					Status:    stRow.Status,
					CreatedAt: stRow.CreatedAt,
					UpdatedAt: stRow.UpdatedAt,
					SeatName:  stRow.SeatName,
					TripID:    stRow.TripID,
				})
			} else {
				t.SeatTicketsEnd = append(t.SeatTicketsEnd, db.GetSeatTicketsByTicketIDRow{ // Map to models.SeatTicket
					ID:        stRow.ID,
					SeatID:    stRow.SeatID,
					TicketID:  stRow.TicketID,
					Status:    stRow.Status,
					CreatedAt: stRow.CreatedAt,
					UpdatedAt: stRow.UpdatedAt,
					SeatName:  stRow.SeatName,
					TripID:    stRow.TripID,
				})
			}
		}
		ticketsReturn = append(ticketsReturn, t)
	}
	return ticketsReturn, nil
}

func (r *ticketRepositoryImpl) GetTicketByID(ctx context.Context, ticketID string) (*models.TicketReturn, error) {
	coreTicket, err := r.q.GetTicketCore(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or specific not found error
		}
		return nil, fmt.Errorf("failed to get core ticket by ID %s: %w", ticketID, err)
	}

	t := &models.TicketReturn{
		TicketID:       coreTicket.TicketID,
		TripIDBegin:    coreTicket.TripIDBegin, // << ADDED
		TripIDEnd:      coreTicket.TripIDEnd,
		Type:           coreTicket.Type,
		CustomerID:     coreTicket.CustomerID,
		Phone:          coreTicket.Phone,
		Email:          coreTicket.Email,
		Name:           coreTicket.Name,
		Price:          coreTicket.Price,
		Status:         coreTicket.Status,
		BookingTime:    coreTicket.BookingTime,
		PaymentStatus:  coreTicket.PaymentStatus,
		BookingChannel: coreTicket.BookingChannel,
		CreatedAt:      coreTicket.CreatedAt,
		UpdatedAt:      coreTicket.UpdatedAt,
		PolicyID:       coreTicket.PolicyID,
		BookedBy:       coreTicket.BookedBy,
	}

	details, err := r.q.GetTicketDetailsByTicketID(ctx, ticketID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get ticket details for ticket %s: %w", ticketID, err)
	}
	for _, det := range details {
		t.Details = append(t.Details, db.TicketDetail{
			DetailID:             det.DetailID,
			TicketID:             det.TicketID,
			PickupLocationBegin:  det.PickupLocationBegin,
			DropoffLocationBegin: det.DropoffLocationBegin,
			PickupLocationEnd:    det.PickupLocationEnd,
			DropoffLocationEnd:   det.DropoffLocationEnd,
			CreatedAt:            det.CreatedAt,
			UpdatedAt:            det.UpdatedAt,
		})
	}

	seatTicketRows, err := r.q.GetSeatTicketsByTicketID(ctx, ticketID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get seat tickets for ticket %s: %w", ticketID, err)
	}
	for _, stRow := range seatTicketRows {
		if stRow.TripID == t.TripIDBegin {
			t.SeatTicketsBegin = append(t.SeatTicketsBegin, db.GetSeatTicketsByTicketIDRow{ // Map to models.SeatTicket
				ID:        stRow.ID,
				SeatID:    stRow.SeatID,
				TicketID:  stRow.TicketID,
				Status:    stRow.Status,
				CreatedAt: stRow.CreatedAt,
				UpdatedAt: stRow.UpdatedAt,
				SeatName:  stRow.SeatName,
				TripID:    stRow.TripID,
			})
		} else {
			t.SeatTicketsEnd = append(t.SeatTicketsEnd, db.GetSeatTicketsByTicketIDRow{ // Map to models.SeatTicket
				ID:        stRow.ID,
				SeatID:    stRow.SeatID,
				TicketID:  stRow.TicketID,
				Status:    stRow.Status,
				CreatedAt: stRow.CreatedAt,
				UpdatedAt: stRow.UpdatedAt,
				SeatName:  stRow.SeatName,
				TripID:    stRow.TripID,
			})
		}
	}
	return t, nil
}

func (r *ticketRepositoryImpl) GetInfoTicketByPhone(ctx context.Context, info *models.TicketInfoInput) (*models.TicketReturn, error) {
	params := db.GetTicketByPhoneAndIDCoreParams{
		TicketID: info.TicketID,
		Phone:    sql.NullString{String: info.Phone, Valid: info.Phone != ""},
	}
	coreTicket, err := r.q.GetTicketByPhoneAndIDCore(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or specific not found error
		}
		return nil, fmt.Errorf("failed to get core ticket by ID %s and phone %s: %w", info.TicketID, info.Phone, err)
	}

	t := &models.TicketReturn{
		TicketID:       coreTicket.TicketID,
		TripIDBegin:    coreTicket.TripIDBegin, // << ADDED
		TripIDEnd:      coreTicket.TripIDEnd,
		Type:           coreTicket.Type,
		CustomerID:     coreTicket.CustomerID,
		Phone:          coreTicket.Phone,
		Email:          coreTicket.Email,
		Name:           coreTicket.Name,
		Price:          coreTicket.Price,
		Status:         coreTicket.Status,
		BookingTime:    coreTicket.BookingTime,
		PaymentStatus:  coreTicket.PaymentStatus,
		BookingChannel: coreTicket.BookingChannel,
		CreatedAt:      coreTicket.CreatedAt,
		UpdatedAt:      coreTicket.UpdatedAt,
		PolicyID:       coreTicket.PolicyID,
		BookedBy:       coreTicket.BookedBy,
	}

	details, err := r.q.GetTicketDetailsByTicketID(ctx, info.TicketID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get ticket details for ticket %s: %w", info.TicketID, err)
	}
	for _, det := range details {
		t.Details = append(t.Details, db.TicketDetail{
			DetailID:             det.DetailID,
			TicketID:             det.TicketID,
			PickupLocationBegin:  det.PickupLocationBegin,
			DropoffLocationBegin: det.DropoffLocationBegin,
			PickupLocationEnd:    det.PickupLocationEnd,
			DropoffLocationEnd:   det.DropoffLocationEnd,
			CreatedAt:            det.CreatedAt,
			UpdatedAt:            det.UpdatedAt,
		})
	}

	seatTicketRows, err := r.q.GetSeatTicketsByTicketID(ctx, info.TicketID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get seat tickets for ticket %s: %w", info.TicketID, err)
	}
	for _, stRow := range seatTicketRows {
		if stRow.TripID == t.TripIDBegin {
			t.SeatTicketsBegin = append(t.SeatTicketsBegin, db.GetSeatTicketsByTicketIDRow{ // Map to models.SeatTicket
				ID:        stRow.ID,
				SeatID:    stRow.SeatID,
				TicketID:  stRow.TicketID,
				Status:    stRow.Status,
				CreatedAt: stRow.CreatedAt,
				UpdatedAt: stRow.UpdatedAt,
				SeatName:  stRow.SeatName,
				TripID:    stRow.TripID,
			})
		} else {
			t.SeatTicketsEnd = append(t.SeatTicketsEnd, db.GetSeatTicketsByTicketIDRow{ // Map to models.SeatTicket
				ID:        stRow.ID,
				SeatID:    stRow.SeatID,
				TicketID:  stRow.TicketID,
				Status:    stRow.Status,
				CreatedAt: stRow.CreatedAt,
				UpdatedAt: stRow.UpdatedAt,
				SeatName:  stRow.SeatName,
				TripID:    stRow.TripID,
			})
		}
	}
	return t, nil
}

func (r *ticketRepositoryImpl) ticketIDExists(ctx context.Context, ticketID string) (bool, error) {
	// Use GetTicketCore and check for sql.ErrNoRows
	_, err := r.q.GetTicketCore(ctx, ticketID)
	if err == nil {
		return true, nil // Ticket exists
	}
	if err == sql.ErrNoRows {
		return false, nil // Ticket does not exist
	}
	return false, err // Other error
}

func (r *ticketRepositoryImpl) GenerateUniqueTicketID(ctx context.Context) (string, error) {
	const maxAttempts = 5
	for i := 0; i < maxAttempts; i++ {
		id := r.utils.GenerateRandomID(6) // Ensure this length is okay for Ticket_Id column
		exists, err := r.ticketIDExists(ctx, id)
		if err != nil {
			return "", err
		}
		if !exists {
			return id, nil
		}
	}
	return "", errors.New("failed to generate unique ticket_id after multiple attempts")
}

func (r *ticketRepositoryImpl) GetAvailableSeatsByTripID(ctx context.Context, tripID string) ([]models.SeatReturn, error) {
	// sqlc.ListAvailableSeatsByTripID returns []db.ListAvailableSeatsByTripIDRow
	// db.ListAvailableSeatsByTripIDRow has ID, TripID, SeatName
	// models.SeatReturn has ID, Name
	dbSeats, err := r.q.ListAvailableSeatsByTripID(ctx, tripID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.SeatReturn{}, nil
		}
		return nil, err
	}

	var seatsReturn []models.SeatReturn
	for _, dbSeat := range dbSeats {
		seatsReturn = append(seatsReturn, models.SeatReturn{
			ID:   int(dbSeat.ID),         // models.SeatReturn expects int
			Name: dbSeat.SeatName.String, // Assuming Name in SeatReturn is string
		})
	}
	return seatsReturn, nil
}

// UpdateTicketPaymentStatus updates payment_status and general status of a ticket.
// The sqlc.UpdateTicketPaymentStatus requires both.
func (r *ticketRepositoryImpl) UpdateTicketPaymentStatus(ctx context.Context, ticketID string, paymentStatus int16, generalStatus int16, tripID string) error {
	params := db.UpdateTicketPaymentStatusParams{
		TicketID:      ticketID,
		PaymentStatus: paymentStatus,
		Status:        generalStatus, // The general status of the ticket
	}
	updatedTicket, err := r.q.UpdateTicketPaymentStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update ticket payment status: %w", err)
	}

	// If payment is successful, update cache with the full updated ticket from DB.
	if paymentStatus == int16(models.PaymentStatusPaid) { // Use models constant
		// updatedTicket is db.Ticket. Cache it.
		if errCache := r.CacheTicket(ctx, &updatedTicket); errCache != nil {
			r.logger.Info("Failed to cache ticket after payment status update for ticket %s: %v", ticketID, errCache)
		}
	}
	return nil
}

func (r *ticketRepositoryImpl) UpdateSeatTicketsStatus(ctx context.Context, ticketID string, status int16) error {
	params := db.UpdateSeatTicketStatusByTicketIDParams{
		TicketID: ticketID,
		Status:   status,
	}
	_, err := r.q.UpdateSeatTicketStatusByTicketID(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update seat tickets status: %w", err)
	}
	return nil
}

func (r *ticketRepositoryImpl) SetupExpirationHandler(ctx context.Context) error {
	pubSub := r.redis.Subscribe(ctx, "__keyevent@0__:expired") // Ensure keyspace notifications are enabled in Redis

	go func() { // run in a separate goroutine
		defer pubSub.Close() // Ensure pubSub is closed when goroutine exits
		ch := pubSub.Channel()
		for msg := range ch {
			expiredKey := msg.Payload
			r.logger.Info("Redis key expired: %s", expiredKey)

			if len(expiredKey) > 5 && expiredKey[:5] == "seat:" {
				var seatID int32
				_, err := fmt.Sscanf(expiredKey, "seat:%d", &seatID)
				if err != nil {
					r.logger.Error("Failed to parse seat ID from expired key '%s': %v", expiredKey, err)
					continue
				}

				expCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

				// Find ticket ID related to the seat ID.
				// FindTicketByTripAndSeat might not be the best if a seat can be reused.
				// This needs a robust way to link an expired seat hold to a specific ticket attempt.
				// Assuming FindTicketByTripAndSeat is appropriate here for now.
				seatTicket, err := r.FindTicketByTripAndSeat(expCtx, seatID) // This uses direct SQL
				if err != nil {
					if err == sql.ErrNoRows {
						r.logger.Info("No seat_ticket found for expired seat %d. Already processed or invalid.", seatID)
					} else {
						r.logger.Error("Failed to find ticket for expired seat %d: %v", seatID, err)
					}
					cancel()
					continue
				}
				if seatTicket == nil { // Explicitly check for nil if FindTicketByTripAndSeat returns nil on ErrNoRows
					r.logger.Info("No seat_ticket found for expired seat %d (nil response). Already processed or invalid.", seatID)
					cancel()
					continue
				}

				// Check current status of the seat ticket, only proceed if it's still in a held/pending state
				// This check should ideally be part of an atomic operation or within the transaction if possible
				// For now, let's get its status.
				// Note: GetSeatTicketStatus takes the seat_ticket.id, not seat_id
				currentSeatTicketStatus, err := r.q.GetSeatTicketStatus(expCtx, seatTicket.ID)
				if err != nil {
					r.logger.Error("Failed to get current status for seat_ticket %d: %v", seatTicket.ID, err)
					cancel()
					continue
				}

				// Only cancel if it's in a state that should be cancelled on expiration (e.g., pending payment)
				if currentSeatTicketStatus == int16(models.SeatStatusPendingPayment) { // Assuming 0 is pending
					err = r.UpdateSeatTicketsStatus(expCtx, seatTicket.TicketID, int16(models.SeatStatusCancelled))
					if err != nil {
						r.logger.Error("Failed to update seat status to cancelled for ticket %s: %v", seatTicket.TicketID, err)
						cancel()
						continue
					}

					err = r.UpdateTicketStatus(expCtx, seatTicket.TicketID, int16(models.TicketStatusCancelled))
					if err != nil {
						r.logger.Error("Failed to update ticket status to cancelled for ticket %s: %v", seatTicket.TicketID, err)
						cancel()
						continue
					}
					r.logger.Info("Successfully cancelled ticket %s and seat_ticket %d due to seat hold expiration", seatTicket.TicketID, seatTicket.ID)

					// Publish seat release for cache update
					// Need TripID for handleSeatReleased. FindTicketByTripAndSeat doesn't provide TripID.
					// We might need to fetch Seat details to get TripID.
					seatInfo, errSeat := r.q.GetSeatByID(expCtx, seatTicket.SeatID)
					if errSeat == nil {
						r.PublishSeatStatusChange(expCtx, &models.SeatStatusMessage{
							EventType: "seat_released", // Or specific "seat_hold_expired"
							SeatID:    seatTicket.SeatID,
							TripID:    seatInfo.TripID,
							TicketID:  seatTicket.TicketID,
						})
					} else {
						r.logger.Error("Could not get TripID for seat %d to publish seat release: %v", seatTicket.SeatID, errSeat)
					}

				} else {
					r.logger.Info("Seat hold for seat_ticket %d (ticket %s) expired, but was not in pending status (current: %d). No action taken.", seatTicket.ID, seatTicket.TicketID, currentSeatTicketStatus)
				}
				cancel()
			}
		}
		r.logger.Info("Redis expiration handler goroutine stopped.")
	}()
	r.logger.Info("Subscribed to Redis key expirations.")
	return nil
}

func (r *ticketRepositoryImpl) UpdateTicketStatus(ctx context.Context, ticketID string, status int16) error {
	params := db.UpdateTicketStatusParams{
		TicketID: ticketID,
		Status:   status,
	}
	updatedTicket, err := r.q.UpdateTicketStatus(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}

	// Update cache
	if errCache := r.CacheTicket(ctx, &updatedTicket); errCache != nil {
		r.logger.Info("Failed to cache ticket after status update for ticket %s: %v", ticketID, errCache)
	}
	return nil
}

func (r *ticketRepositoryImpl) GetSeatIDsByTicketID(ctx context.Context, ticketID string) ([]int32, error) {
	ids, err := r.q.GetSeatIDsByTicketID(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []int32{}, nil
		}
		return nil, fmt.Errorf("failed to query seat IDs: %w", err)
	}
	return ids, nil
}

func (r *ticketRepositoryImpl) CleanupTicketCache(ctx context.Context, ticketID string, seatIDs []int32, tripID string) error {
	err := r.redis.Del(ctx, "ticket:"+ticketID).Err()
	if err != nil && err != redis.Nil {
		// Log error but continue, as primary concern is DB.
		r.logger.Error("Failed to delete ticket cache for ticket %s: %v", ticketID, err)
	}

	seatIDs, err = r.GetSeatIDsByTicketID(ctx, ticketID)
	if err != nil {
		// If we can't get seat IDs, we can't clean their cache. Log and return.
		return fmt.Errorf("failed to get seat IDs for ticket %s during cache cleanup: %w", ticketID, err)
	}

	for _, seatID := range seatIDs {
		key := fmt.Sprintf("seat:%d", seatID)
		errDelSeat := r.redis.Del(ctx, key).Err()
		if errDelSeat != nil && errDelSeat != redis.Nil {
			r.logger.Error("Failed to delete seat cache for seat %d (ticket %s): %v", seatID, ticketID, errDelSeat)
			// Continue to try deleting other seat caches
		}
	}
	return nil
}

func (r *ticketRepositoryImpl) CacheAvailableSeats(ctx context.Context, tripID string, seats []models.SeatReturn) error {
	pipe := r.redis.Pipeline()
	availableSeatsKey := fmt.Sprintf("available_seats:%s", tripID)
	seatInfoKey := fmt.Sprintf("seat_info:%s", tripID)

	pipe.Del(ctx, availableSeatsKey)
	pipe.Del(ctx, seatInfoKey) // Also delete seat info key

	if len(seats) > 0 {
		seatInfoMap := make(map[string]interface{})
		for _, seat := range seats {
			pipe.SAdd(ctx, availableSeatsKey, seat.ID) // seat.ID is int, SAdd handles it
			seatInfoMap[fmt.Sprintf("%d", seat.ID)] = seat.Name
		}
		if len(seatInfoMap) > 0 {
			pipe.HMSet(ctx, seatInfoKey, seatInfoMap)
		}
	}

	pipe.Expire(ctx, availableSeatsKey, 30*time.Minute)
	pipe.Expire(ctx, seatInfoKey, 30*time.Minute) // Expire seat info key as well
	_, err := pipe.Exec(ctx)
	return err
}

func (r *ticketRepositoryImpl) GetAvailableSeatsFromCache(ctx context.Context, tripID string) ([]models.SeatReturn, error) {
	availableSeatsKey := fmt.Sprintf("available_seats:%s", tripID)
	seatInfoKey := fmt.Sprintf("seat_info:%s", tripID)

	exists, err := r.redis.Exists(ctx, availableSeatsKey, seatInfoKey).Result()
	if err != nil {
		return nil, fmt.Errorf("redis error checking cache keys for trip %s: %w", tripID, err)
	}
	if exists < 2 { // Both keys must exist
		return nil, nil // Cache miss
	}

	seatIDStrings, err := r.redis.SMembers(ctx, availableSeatsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get available seat IDs from cache for trip %s: %w", tripID, err)
	}
	if len(seatIDStrings) == 0 {
		return []models.SeatReturn{}, nil
	}

	// Fetch seat names using HMGet
	seatNamesResult, err := r.redis.HMGet(ctx, seatInfoKey, seatIDStrings...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get seat info from cache for trip %s: %w", tripID, err)
	}

	var result []models.SeatReturn
	for i, seatIDStr := range seatIDStrings {
		seatID, convErr := strconv.Atoi(seatIDStr)
		if convErr != nil {
			r.logger.Error("Failed to parse seat ID string '%s' from cache for trip %s: %v", seatIDStr, tripID, convErr)
			continue
		}
		if seatNamesResult[i] != nil {
			seatName, ok := seatNamesResult[i].(string)
			if !ok {
				r.logger.Error("Seat name for ID %d is not a string in cache for trip %s", seatID, tripID)
				continue
			}
			result = append(result, models.SeatReturn{ID: seatID, Name: seatName})
		} else {
			// This case (seat ID in set but not in hash) should ideally not happen if cache is consistent
			r.logger.Info("Seat ID %d found in available_seats set but not in seat_info hash for trip %s", seatID, tripID)
		}
	}
	return result, nil
}

func (r *ticketRepositoryImpl) PublishSeatStatusChange(ctx context.Context, message *models.SeatStatusMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal seat status message: %w", err)
	}
	return r.redis.Publish(ctx, "seat_status_changes", data).Err()
}

func (r *ticketRepositoryImpl) SubscribeToSeatStatusChanges(ctx context.Context) error {
	pubsub := r.redis.Subscribe(ctx, "seat_status_changes")
	// Not closing pubsub here, it's closed in the goroutine on exit if needed,
	// or when the application shuts down.
	// For a long-running subscription, the Close would happen on app termination.

	go func() {
		// Need a new context for the goroutine that is not the initial context,
		// or manage its lifecycle carefully. Using context.Background() for simplicity here,
		// but a cancellable context tied to app lifecycle is better.
		goroutineCtx := context.Background() // Or a more managed context
		ch := pubsub.Channel()
		r.logger.Info("Subscribed to seat_status_changes channel on Redis.")

		for msg := range ch {
			var statusMsg models.SeatStatusMessage
			if err := json.Unmarshal([]byte(msg.Payload), &statusMsg); err != nil {
				r.logger.Error("Failed to unmarshal seat status message: %v", err)
				continue
			}

			r.logger.Info("Received seat status change: EventType=%s, SeatID=%d, TripID=%s, TicketID=%s",
				statusMsg.EventType, statusMsg.SeatID, statusMsg.TripID, statusMsg.TicketID)

			switch statusMsg.EventType {
			case "seat_reserved":
				r.handleSeatReserved(goroutineCtx, &statusMsg)
			case "seat_released", "seat_hold_expired": // Handle expired holds as releases
				r.handleSeatReleased(goroutineCtx, &statusMsg)
			case "payment_completed":
				r.handlePaymentCompleted(goroutineCtx, &statusMsg)
			case "payment_failed": // This should also release the seat
				r.handlePaymentFailed(goroutineCtx, &statusMsg)
			default:
				r.logger.Info("Unhandled seat status event type: %s", statusMsg.EventType)
			}
		}
		r.logger.Info("Seat status change subscription channel closed.")
	}()

	return nil // Subscription setup initiated
}

func (r *ticketRepositoryImpl) handleSeatReserved(ctx context.Context, msg *models.SeatStatusMessage) {
	availableSeatsKey := fmt.Sprintf("available_seats:%s", msg.TripID)
	// Seat info key is not directly modified here, only the set of available IDs
	pipe := r.redis.Pipeline()
	pipe.SRem(ctx, availableSeatsKey, msg.SeatID)
	_, err := pipe.Exec(ctx)
	if err != nil {
		r.logger.Error("Failed to update cache for seat reservation (trip %s, seat %d): %v", msg.TripID, msg.SeatID, err)
	} else {
		r.logger.Info("Seat %d reserved for trip %s, removed from available_seats cache.", msg.SeatID, msg.TripID)
	}
}

func (r *ticketRepositoryImpl) handleSeatReleased(ctx context.Context, msg *models.SeatStatusMessage) {
	availableSeatsKey := fmt.Sprintf("available_seats:%s", msg.TripID)
	seatInfoKey := fmt.Sprintf("seat_info:%s", msg.TripID)

	// Fetch seat name from DB to ensure consistency if it was removed or info is stale
	var seatName string
	// Use sqlc GetSeatByID
	dbSeat, err := r.q.GetSeatByID(ctx, int32(msg.SeatID))
	if err != nil {
		r.logger.Error("Failed to get seat details for seat %d (trip %s) during release: %v", msg.SeatID, msg.TripID, err)
		// If we can't get seat details, it's risky to add to cache.
		// Consider invalidating/deleting cache for the trip.
		r.redis.Del(ctx, availableSeatsKey, seatInfoKey)
		return
	}
	seatName = dbSeat.SeatName.String

	pipe := r.redis.Pipeline()
	pipe.SAdd(ctx, availableSeatsKey, msg.SeatID)
	pipe.HSet(ctx, seatInfoKey, fmt.Sprintf("%d", msg.SeatID), seatName)
	_, err = pipe.Exec(ctx)

	if err != nil {
		r.logger.Error("Failed to update cache for seat release (trip %s, seat %d): %v", msg.TripID, msg.SeatID, err)
	} else {
		r.logger.Info("Seat %d released for trip %s, added back to available_seats and seat_info cache.", msg.SeatID, msg.TripID)
	}

	// Also delete any temporary holds on "seat:<id>"
	seatKey := fmt.Sprintf("seat:%d", msg.SeatID)
	r.redis.Del(ctx, seatKey)
}

func (r *ticketRepositoryImpl) handlePaymentCompleted(ctx context.Context, msg *models.SeatStatusMessage) {
	// Seat was already removed from available_seats when reserved.
	// Just clean up any "seat:<id>" temporary hold key.
	seatKey := fmt.Sprintf("seat:%d", msg.SeatID)
	r.redis.Del(ctx, seatKey)
	r.logger.Info("Payment completed for ticket %s (seat %d, trip %s), cleaned up temporary seat hold from cache.", msg.TicketID, msg.SeatID, msg.TripID)
}

func (r *ticketRepositoryImpl) handlePaymentFailed(ctx context.Context, msg *models.SeatStatusMessage) {
	// Payment failed, so the seat should be released.
	// This involves removing the temporary hold and adding it back to available seats.
	r.handleSeatReleased(ctx, msg) // This also deletes the "seat:<id>" key.
	r.logger.Info("Payment failed for ticket %s (seat %d, trip %s), seat returned to available pool in cache.", msg.TicketID, msg.SeatID, msg.TripID)
}

// GetTotalTicketCount retrieves the total number of tickets in the database.
func (r *ticketRepositoryImpl) GetTotalTicketCount(ctx context.Context) (int64, error) {
	return r.q.GetTotalTicketCount(ctx)
}

// GetAllTickets retrieves a paginated and fully detailed list of all tickets.
// This implementation is optimized to prevent N+1 query issues.
func (r *ticketRepositoryImpl) GetAllTickets(ctx context.Context, limit int, offset int) ([]*models.TicketReturn, error) {
	// 1. Fetch the paginated core tickets
	coreTickets, err := r.q.GetAllTickets(ctx, db.GetAllTicketsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.TicketReturn{}, nil
		}
		return nil, fmt.Errorf("failed to get all core tickets: %w", err)
	}

	if len(coreTickets) == 0 {
		return []*models.TicketReturn{}, nil
	}

	// 2. Collect ticket IDs
	ticketIDs := make([]string, len(coreTickets))
	ticketsMap := make(map[string]*models.TicketReturn, len(coreTickets))
	ticketsReturn := make([]*models.TicketReturn, len(coreTickets))

	for i, coreTicket := range coreTickets {
		t := &models.TicketReturn{
			TicketID:       coreTicket.TicketID,
			TripIDBegin:    coreTicket.TripIDBegin,
			TripIDEnd:      coreTicket.TripIDEnd,
			Type:           coreTicket.Type,
			CustomerID:     coreTicket.CustomerID,
			Phone:          coreTicket.Phone,
			Email:          coreTicket.Email,
			Name:           coreTicket.Name,
			Price:          coreTicket.Price,
			Status:         coreTicket.Status,
			BookingTime:    coreTicket.BookingTime,
			PaymentStatus:  coreTicket.PaymentStatus,
			BookingChannel: coreTicket.BookingChannel,
			CreatedAt:      coreTicket.CreatedAt,
			UpdatedAt:      coreTicket.UpdatedAt,
			PolicyID:       coreTicket.PolicyID,
			BookedBy:       coreTicket.BookedBy,
		}
		ticketIDs[i] = t.TicketID
		ticketsMap[t.TicketID] = t
		ticketsReturn[i] = t
	}

	// 3. Fetch all details and seat tickets for the collected ticket IDs in bulk
	details, err := r.q.GetTicketDetailsByTicketIDs(ctx, ticketIDs)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get ticket details for tickets: %w", err)
	}

	seatTicketRows, err := r.q.GetSeatTicketsByTicketIDs(ctx, ticketIDs)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get seat tickets for tickets: %w", err)
	}

	// 4. Map the details back to their respective tickets
	for _, det := range details {
		if ticket, ok := ticketsMap[det.TicketID]; ok {
			ticket.Details = append(ticket.Details, db.TicketDetail{
				DetailID:             det.DetailID,
				TicketID:             det.TicketID,
				PickupLocationBegin:  det.PickupLocationBegin,
				DropoffLocationBegin: det.DropoffLocationBegin,
				PickupLocationEnd:    det.PickupLocationEnd,
				DropoffLocationEnd:   det.DropoffLocationEnd,
			})
		}
	}

	// 5. Map the seat tickets back to their respective tickets
	for _, stRow := range seatTicketRows {
		if ticket, ok := ticketsMap[stRow.TicketID]; ok {
			seatTicket := db.GetSeatTicketsByTicketIDRow{
				ID:        stRow.ID,
				SeatID:    stRow.SeatID,
				TicketID:  stRow.TicketID,
				Status:    stRow.Status,
				TripID:    stRow.TripID,
				CreatedAt: stRow.CreatedAt,
				UpdatedAt: stRow.UpdatedAt,
				SeatName:  stRow.SeatName,
			}
			if stRow.TripID == ticket.TripIDBegin {
				ticket.SeatTicketsBegin = append(ticket.SeatTicketsBegin, seatTicket)
			} else {
				ticket.SeatTicketsEnd = append(ticket.SeatTicketsEnd, seatTicket)
			}
		}
	}

	return ticketsReturn, nil
}
