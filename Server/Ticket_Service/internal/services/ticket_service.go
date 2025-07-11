package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	// "strconv" // No longer needed here if payment logic is removed
	"ticket-service/config"
	"ticket-service/domain/models" // Make sure this path is correct
	"ticket-service/internal/db"
	"ticket-service/internal/repositories"
	"ticket-service/pkg/kafkaclient"
	"ticket-service/pkg/utils"
	"ticket-service/pkg/websocket"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Helper to convert string to sql.NullString (already present in manager_ticket_service.go, ensure accessible or duplicate)
const pendingWsConnectionsKey = "pending_ws_connections"

type ITicketService interface {
	GetTicketByID(ctx context.Context, ticketID string) (*models.TicketReturn, error)
	GetInfoTicketByPhone(ctx context.Context, info *models.TicketInfoInput) (*models.TicketReturn, error)
	GetTicketByCustomer(ctx context.Context, customerID sql.NullInt32) ([]*models.TicketReturn, error)
	// Updated CreateTicket to remove PaymentInitiationInfo
	CreateTicket(ctx context.Context, input *models.TicketInput, customerID sql.NullInt32) (*db.Ticket, error)
	CreateTicketByStaff(ctx context.Context, input *models.TicketInput, staffID string) (*db.Ticket, error)
	GetAvailableSeatsByTripID(ctx context.Context, tripID string) ([]models.SeatReturn, error)
	ExtendSeatHoldTime(ctx context.Context, ticketID string, seatIDs []int32, extendMinutes int) error
	ReleaseHeldSeats(ctx context.Context, ticketID string, seatIDs []int32) error
	UpdateTicketPaymentStatus(ctx context.Context, ticketID string, paymentStatus int16, ticketStatus int16, tripID string) error

	CreateTicketAndNotify(ctx context.Context, input *models.TicketInput, customerID sql.NullInt32, bookingID string)
	CancelTicket(ctx context.Context, ticketID string) error
	QueueNewBooking(ctx context.Context, bookingID string, input *models.TicketInput, customerID sql.NullInt32) error
	GetAllTickets(ctx context.Context, page, limit int) (*models.PaginatedTickets, error)
}

type TicketService struct {
	ticketRepository repositories.TicketRepositoryInterface
	utils            *utils.Utils
	logger           utils.Logger
	cfg              config.Config
	publisher        *kafkaclient.Publisher // << UPDATED
	redisClient      *redis.Client
}

func NewTicketService(ticketRepository repositories.TicketRepositoryInterface, utils *utils.Utils, logger utils.Logger, cfg config.Config, publisher *kafkaclient.Publisher, redisClient *redis.Client) ITicketService {
	return &TicketService{
		ticketRepository: ticketRepository,
		utils:            utils,
		logger:           logger,
		cfg:              cfg,
		publisher:        publisher, // << UPDATED
		redisClient:      redisClient,
	}
}

func (t *TicketService) GetTicketByID(ctx context.Context, ticketID string) (*models.TicketReturn, error) {
	ticket, err := t.ticketRepository.GetTicketByID(ctx, ticketID)
	if err != nil {
		t.logger.Error("Error retrieving ticket %s: %v", ticketID, err)
		return nil, err
	}
	return ticket, nil
}

func (t *TicketService) GetInfoTicketByPhone(ctx context.Context, info *models.TicketInfoInput) (*models.TicketReturn, error) {
	ticket, err := t.ticketRepository.GetInfoTicketByPhone(ctx, info)
	if err != nil {
		t.logger.Error("Error retrieving ticket by phone: %v", err)
		return nil, err
	}
	return ticket, nil
}

func (t *TicketService) GetTicketByCustomer(ctx context.Context, customerID sql.NullInt32) ([]*models.TicketReturn, error) {
	tickets, err := t.ticketRepository.GetTicketByCustomer(ctx, customerID)
	if err != nil {
		logMsg := "Error retrieving tickets"
		if customerID.Valid {
			logMsg = fmt.Sprintf("%s for customer %d", logMsg, customerID.Int32)
		} else {
			logMsg = fmt.Sprintf("%s for nil customer_id", logMsg)
		}
		t.logger.Error("%s: %v", logMsg, err)
		return nil, err
	}

	// // Iterate through tickets and enrich them with trip details
	// for i, ticket := range tickets { // Use index to modify the original slice element if tickets are values
	// 	if ticket.TripID != "" {
	// 		tripDetails, err := GetTripDetails(ticket.TripID) // Calling the new function
	// 		if err != nil {
	// 			// Decide on error handling:
	// 			// 1. Fail the entire request (return nil, err)
	// 			// 2. Log the error and continue (tripDetails will be nil for this ticket)
	// 			// 3. Return partial data with an error indicator
	// 			// For now, let's log and continue, so TripDetails might be nil.
	// 			t.logger.Error("Failed to retrieve details for TripID %s: %v", ticket.TripID, err)
	// 			// tickets[i].TripDetails will remain nil or you can assign an empty struct
	// 		}
	// 		// Ensure you are modifying the actual element in the slice if `tickets` is a slice of values.
	// 		// If `tickets` is a slice of pointers (`[]*models.TicketReturn`), then `ticket.TripDetails = tripDetails` is fine.
	// 		// Based on your original code, `tickets` is `[]*models.TicketReturn`, so direct assignment is fine.
	// 		if tickets[i] != nil { // Defensive check
	// 			tickets[i].TripDetails = tripDetails
	// 		}
	// 	}
	// }

	return tickets, nil
}

// CreateTicket handles ticket creation. Payment initiation is removed.
// func (t *TicketService) CreateTicket(ctx context.Context, input *models.TicketInput, customerID sql.NullInt32) (*db.Ticket, error) {
// 	// Reduced timeout as external payment calls are removed. Adjusted to 15 seconds for DB and Redis operations.
// 	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
// 	defer cancel()

// 	var lockAcquired bool
// 	var unlock func()
// 	var err error

// 	lockKey := fmt.Sprintf("seat-lock:%s", input.TripID)
// 	for i := 0; i < 3; i++ {
// 		lockAcquired, unlock, err = t.ticketRepository.AcquireLock(ctx, lockKey)
// 		if err != nil && !errors.Is(err, repositories.ErrRedisNotAvailable) {
// 			t.logger.Error("Error acquiring lock for trip %s: %v", input.TripID, err)
// 			return nil, fmt.Errorf("could not acquire lock: %w", err)
// 		}
// 		if lockAcquired || errors.Is(err, repositories.ErrRedisNotAvailable) {
// 			break
// 		}
// 		time.Sleep(time.Duration(100*(i+1)) * time.Millisecond)
// 	}
// 	if !lockAcquired && !errors.Is(err, repositories.ErrRedisNotAvailable) {
// 		return nil, errors.New("another booking is in progress for this trip, please try again shortly")
// 	}
// 	if lockAcquired {
// 		defer unlock()
// 	}

// 	for _, seatID := range input.SeatID {
// 		isHeld, err := t.ticketRepository.IsSeatHeld(ctx, seatID)
// 		if err != nil && !errors.Is(err, repositories.ErrRedisNotAvailable) {
// 			t.logger.Error("Error checking if seat %d is held: %v", seatID, err)
// 			return nil, err
// 		}
// 		if isHeld {
// 			t.logger.Info("Seat %d is temporarily held", seatID)
// 			return nil, fmt.Errorf("seat %d is temporarily held by another user", seatID)
// 		}
// 		isBookedInDB, err := t.ticketRepository.AllTicketsStatus2BySeat(ctx, seatID)
// 		if err != nil {
// 			t.logger.Error("Database error checking seat %d: %v", seatID, err)
// 			return nil, err
// 		}
// 		if !isBookedInDB {
// 			t.logger.Info("Seat %d already booked or unavailable in database", seatID)
// 			return nil, fmt.Errorf("seat %d already booked or unavailable", seatID)
// 		}
// 	}

// 	ticketID, err := t.ticketRepository.GenerateUniqueTicketID(ctx)
// 	if err != nil {
// 		t.logger.Error("Error generating unique ticket ID: %v", err)
// 		return nil, err
// 	}

// 	bookedBy := "customer"
// 	if input.BookedBy != "" {
// 		bookedBy = input.BookedBy
// 	}

// 	// Determine initial status based on whether it's an online booking (pending payment)
// 	// or another channel that might imply immediate payment or different handling.
// 	// For now, assume online bookings are pending.
// 	// This might need adjustment based on how payment is handled post-ticket creation.
// 	initialTicketStatus := models.TicketStatusPendingConfirmation
// 	initialPaymentStatus := models.PaymentStatusPending

// 	// If the booking channel suggests immediate payment (e.g. Counter Sale was part of this logic,
// 	// but now CreateTicketByStaff handles that. For generic CreateTicket, assume PENDING)
// 	// we might set different statuses. For this modified function, we assume PENDING for all
// 	// tickets created via this generic method.
// 	// If CreateTicket is *only* for online channels that require a separate payment step, then PENDING is correct.

// 	ticket := &db.Ticket{
// 		TicketID:       ticketID,
// 		CustomerID:     customerID,
// 		Price:          input.Price,
// 		Status:         initialTicketStatus,
// 		PaymentStatus:  initialPaymentStatus,
// 		Name:           utils.ToNullString(input.Name),
// 		BookedBy:       utils.ToNullString(bookedBy),
// 		BookingChannel: input.BookingChannel, // Keep booking channel for context
// 		PolicyID:       input.PolicyID,
// 		Phone:          utils.ToNullString(input.Phone),
// 		Email:          utils.ToNullString(input.Email),
// 	}

// 	ticketDetail := &db.TicketDetail{
// 		TicketID:        ticketID,
// 		PickupLocation:  sql.NullInt32{Int32: input.PickupLocation, Valid: true},
// 		DropoffLocation: sql.NullInt32{Int32: input.DropoffLocation, Valid: true},
// 	}

// 	err = t.ticketRepository.CreateTicket(ctx, ticket, ticketDetail, input.SeatID, input.TripID)
// 	if err != nil {
// 		t.logger.Error("Error creating ticket in database: %v", err)
// 		return nil, err
// 	}

// 	go func() {
// 		eventCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()

// 		seatCount := len(input.SeatID)
// 		eventPayload := kafkaclient.SeatUpdateEvent{
// 			TripID:    input.TripID,
// 			SeatCount: seatCount,
// 		}

// 		err := t.publisher.Publish(eventCtx, t.cfg.Kafka.Topics.SeatsReserved.Topic, []byte(input.TripID), eventPayload)
// 		if err != nil {
// 			t.logger.Error("CRITICAL: Failed to publish seats_reserved event for TripID %s. Error: %v", input.TripID, err)
// 		}
// 	}()

// 	for _, seatID := range input.SeatID {
// 		go func(sID int32) {
// 			holdCtx, holdCancel := context.WithTimeout(context.Background(), 5*time.Second)
// 			defer holdCancel()
// 			if err := t.ticketRepository.CacheTemporarySeat(holdCtx, sID, input.TripID, ticket); err != nil {
// 				t.logger.Error("Warning: Failed to place temporary hold on seat %d for trip %s: %v", sID, input.TripID, err)
// 			}
// 			statusMsg := &models.SeatStatusMessage{EventType: "seat_reserved_pending_payment", TicketID: ticket.TicketID, SeatID: sID, TripID: input.TripID, Timestamp: time.Now()}
// 			if err := t.ticketRepository.PublishSeatStatusChange(holdCtx, statusMsg); err != nil {
// 				t.logger.Error("Warning: Failed to publish seat reservation event: %v", err)
// 			}
// 		}(seatID)
// 	}

// 	// Payment initiation logic has been removed.
// 	// The ticket is created with a PENDING status.
// 	// A separate call will be needed to process payment and update ticket status.

// 	t.logger.Info("Ticket %s created successfully with status PENDING. Awaiting payment processing via channel %d.", ticket.TicketID, input.BookingChannel)
// 	return ticket, nil
// }

func (t *TicketService) CreateTicket(ctx context.Context, input *models.TicketInput, customerID sql.NullInt32) (*db.Ticket, error) {
	// ... (logic from previous refactor is correct and kept here)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if input.TicketType > 1 {
		return nil, errors.New("invalid ticket type")
	}

	if len(input.SeatIDBegin) == 0 {
		return nil, errors.New("at least one seat must be selected")
	}

	if input.TicketType == 1 {
		if len(input.SeatIDEnd) == 0 {
			return nil, errors.New("at least one seat")
		}
	}

	for _, seatID := range input.SeatIDBegin {
		lockKey := fmt.Sprintf("trip-lock:%s", string(seatID))
		lockAcquired, unlock, err := t.ticketRepository.AcquireLock(ctx, lockKey, 5*time.Second)
		if err != nil {
			t.logger.Error("Error acquiring lock for trip %s: %v", seatID, err)
			return nil, errors.New("could not contact booking service, please try again")
		}
		if !lockAcquired {
			return nil, errors.New("booking for this trip is currently busy, please try again shortly")
		}
		defer unlock()
	}
	if input.TicketType == 1 {
		for _, seatID := range input.SeatIDEnd {
			lockKey := fmt.Sprintf("trip-lock:%s", string(seatID))
			lockAcquired, unlock, err := t.ticketRepository.AcquireLock(ctx, lockKey, 5*time.Second)
			if err != nil {
				t.logger.Error("Error acquiring lock for trip %s: %v", seatID, err)
				return nil, errors.New("could not contact booking service, please try again")
			}
			if !lockAcquired {
				return nil, errors.New("booking for this trip is currently busy, please try again shortly")
			}
			defer unlock()
		}
	}
	availableSeatRows, err := t.ticketRepository.AreSeatsAvailable(ctx, input.SeatIDBegin)
	if err != nil {
		t.logger.Error("Database error during batch seat validation for trip %s: %v", input.TripIDBegin, err)
		return nil, errors.New("error checking seat availability")
	}

	var availableSeatRowsEnd []db.AreSeatsAvailableRow

	if input.TicketType == 1 {
		availableSeatRowsEnd, err = t.ticketRepository.AreSeatsAvailable(ctx, input.SeatIDEnd)
		if err != nil {
			t.logger.Error("Database error during batch seat validation for trip %s: %v", input.TripIDEnd, err)
			return nil, errors.New("error checking seat availability")
		}
	}

	if len(availableSeatRows) != len(input.SeatIDBegin) || len(availableSeatRowsEnd) != len(input.SeatIDEnd) {
		return nil, errors.New("one or more selected seats do not exist or belong to a different trip")
	}

	for _, seat := range availableSeatRows {
		if seat.IsBooked {
			return nil, fmt.Errorf("seat %d is already booked or held by another user", seat.ID)
		}
	}

	for _, seat := range availableSeatRowsEnd {
		if seat.IsBooked {
			return nil, fmt.Errorf("seat %d is already booked or held by another user", seat.ID)
		}
	}

	ticketID, err := t.ticketRepository.GenerateUniqueTicketID(ctx)
	if err != nil {
		t.logger.Error("Error generating unique ticket ID: %v", err)
		return nil, err
	}

	// Because this is a customer booking, status is pending confirmation/payment
	ticket := &db.Ticket{
		TicketID:       ticketID,
		CustomerID:     customerID,
		Price:          input.Price,
		TripIDBegin:    input.TripIDBegin,
		TripIDEnd:      sql.NullString{String: input.TripIDEnd, Valid: input.TicketType == 1},
		Status:         models.TicketStatusPendingConfirmation,
		Type:           int16(input.TicketType),
		PaymentStatus:  models.PaymentStatusPending,
		Name:           utils.ToNullString(input.Name),
		BookedBy:       utils.ToNullString("customer"),
		BookingChannel: input.BookingChannel,
		PolicyID:       input.PolicyID,
		Phone:          utils.ToNullString(input.Phone),
		Email:          utils.ToNullString(input.Email),
		BookingTime:    time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	eventPayload := kafkaclient.SeatUpdateEvent{
		TripID:    input.TripIDBegin,
		SeatCount: int(len(input.SeatIDBegin)),
	}

	payloadBytes, _ := json.Marshal(eventPayload)

	var eventPayloadEnd kafkaclient.SeatUpdateEvent
	var payloadBytesEnd []byte

	if input.TicketType == 1 {
		eventPayloadEnd = kafkaclient.SeatUpdateEvent{
			TripID:    input.TripIDEnd,
			SeatCount: int(len(input.SeatIDEnd)),
		}
		payloadBytesEnd, _ = json.Marshal(eventPayloadEnd)
	}

	repoParams := repositories.CreateTicketTransactionParams{
		Ticket: db.CreateTicketParams{
			TicketID:       ticket.TicketID,
			TripIDBegin:    input.TripIDBegin,
			TripIDEnd:      sql.NullString{String: input.TripIDEnd, Valid: input.TicketType == 1},
			Type:           ticket.Type,
			CustomerID:     ticket.CustomerID,
			Phone:          ticket.Phone,
			Email:          ticket.Email,
			Name:           ticket.Name,
			Price:          ticket.Price,
			Status:         ticket.Status,
			BookingTime:    time.Now(),
			PaymentStatus:  ticket.PaymentStatus,
			BookingChannel: ticket.BookingChannel,
			PolicyID:       ticket.PolicyID,
			BookedBy:       ticket.BookedBy,
		},
		TicketDetail: db.CreateTicketDetailsParams{
			TicketID:             ticketID,
			PickupLocationBegin:  sql.NullInt32{Int32: input.PickupLocationBegin, Valid: true},
			DropoffLocationEnd:   sql.NullInt32{Int32: input.DropoffLocationEnd, Valid: true},
			PickupLocationEnd:    sql.NullInt32{Int32: input.PickupLocationEnd, Valid: input.TicketType == 1},
			DropoffLocationBegin: sql.NullInt32{Int32: input.DropoffLocationBegin, Valid: input.TicketType == 1},
		},
		SeatIDsBegin: input.SeatIDBegin,
		SeatIDsEnd:   input.SeatIDEnd,
		OutboxEvents: []db.CreateOutboxEventParams{
			{
				ID:      uuid.New(),
				Topic:   t.cfg.Kafka.Topics.SeatsReserved.Topic,
				Key:     input.TripIDBegin,
				Payload: payloadBytes,
			},
			{
				ID:      uuid.New(),
				Topic:   t.cfg.Kafka.Topics.SeatsReserved.Topic,
				Key:     input.TripIDEnd,
				Payload: payloadBytesEnd,
			},
		},
	}

	err = t.ticketRepository.CreateTicketInTransaction(ctx, repoParams)
	if err != nil {
		t.logger.Error("Error in repository's atomic transaction for ticket %s: %v", ticketID, err)
		return nil, fmt.Errorf("failed to finalize booking, please try again")
	}

	go func(tripID string, seatIDs []int32) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t.ticketRepository.UpdateCachedAvailableSeats(bgCtx, tripID, seatIDs, "REMOVE")
	}(input.TripIDBegin, input.SeatIDBegin)

	if input.TicketType == 1 {
		go func(tripID string, seatIDs []int32) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			t.ticketRepository.UpdateCachedAvailableSeats(bgCtx, tripID, seatIDs, "REMOVE")
		}(input.TripIDEnd, input.SeatIDEnd)
	}

	t.logger.Info("Ticket %s created successfully. Transaction and outbox event committed.", ticket.TicketID)
	return ticket, nil
}

/*
// createVNPayPayment calls the payment_service to create a VNPay payment URL.
// This function is no longer called by CreateTicket. Kept for potential other uses or can be removed.
func (t *TicketService) createVNPayPayment(ctx context.Context, paymentReq map[string]interface{}) (*models.VNPayPaymentResponseData, error) {
	reqBody, err := json.Marshal(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal VNPay payment request: %w", err)
	}

	if t.cfg.URL.Payment == "" {
		return nil, errors.New("VNPay payment URL (PAYMENT_URL) is not configured in ticket service")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", t.cfg.URL.Payment, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create VNPay payment HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to VNPay payment service failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("VNPay payment service returned non-OK status: %d. Body: %s", resp.StatusCode, string(bodyBytes))
	}

	var serviceResp models.VNPayServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode VNPay payment service response: %w", err)
	}

	if serviceResp.Data.PaymentURL == "" {
		return nil, errors.New("VNPay payment service returned empty payment URL")
	}

	return &serviceResp.Data, nil
}

// createStripePaymentIntent calls the payment_service to create a Stripe PaymentIntent.
// This function is no longer called by CreateTicket. Kept for potential other uses or can be removed.
func (t *TicketService) createStripePaymentIntent(ctx context.Context, paymentReq models.InitialStripePaymentRequest) (*models.StripePaymentIntentResponse, error) {
	reqBody, err := json.Marshal(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Stripe payment request: %w", err)
	}

	if t.cfg.URL.StripePayment == "" {
		return nil, errors.New("Stripe payment URL (STRIPE_PAYMENT_URL) is not configured in ticket service")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", t.cfg.URL.StripePayment, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe payment HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	var httpResp *http.Response

	maxRetries := 2
	for attempt := 0; attempt < maxRetries; attempt++ {
		httpResp, err = client.Do(req.Clone(ctx))
		if err == nil && httpResp.StatusCode == http.StatusOK {
			break
		}
		if httpResp != nil {
			httpResp.Body.Close()
		}
		if attempt < maxRetries-1 {
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context cancelled during Stripe payment retry: %w", ctx.Err())
			case <-time.After(time.Duration(500*(attempt+1)) * time.Millisecond):
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("request to Stripe payment service failed after retries: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("Stripe payment service returned non-OK status: %d. Body: %s", httpResp.StatusCode, string(bodyBytes))
	}

	var serviceResp struct {
		Success bool                             `json:"success"`
		Message string                           `json:"message"`
		Data    models.StripePaymentIntentResponse `json:"data"`
		Error   *struct {
			Code    int         `json:"code"`
			Details interface{} `json:"details"`
		} `json:"error"`
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&serviceResp); err != nil {
		return nil, fmt.Errorf("failed to decode Stripe payment service response: %w", err)
	}

	if !serviceResp.Success || serviceResp.Data.ClientSecret == "" {
		errMsg := fmt.Sprintf("Stripe payment service indicated failure: %s", serviceResp.Message)
		if serviceResp.Error != nil && serviceResp.Error.Details != nil {
			errMsg = fmt.Sprintf("%s. Details: %v", errMsg, serviceResp.Error.Details)
		} else if !serviceResp.Success && serviceResp.Data.ClientSecret == "" {
			errMsg = fmt.Sprintf("Stripe payment service returned success=false or empty client_secret. Message: %s", serviceResp.Message)
		}
		return nil, errors.New(errMsg)
	}

	return &serviceResp.Data, nil
}
*/

func (t *TicketService) CreateTicketByStaff(ctx context.Context, input *models.TicketInput, staffID string) (*db.Ticket, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Validate ticket type
	if input.TicketType > 1 {
		return nil, errors.New("invalid ticket type")
	}

	// Validate required seats
	if len(input.SeatIDBegin) == 0 {
		return nil, errors.New("at least one seat must be selected")
	}

	// For round trips, validate end seats
	if input.TicketType == 1 {
		if len(input.SeatIDEnd) == 0 {
			return nil, errors.New("at least one seat must be selected for return trip")
		}
	}

	// 1. Acquire locks for begin trip seats (same locking strategy as CreateTicket)
	for _, seatID := range input.SeatIDBegin {
		lockKey := fmt.Sprintf("trip-lock:%s", string(seatID))
		lockAcquired, unlock, err := t.ticketRepository.AcquireLock(ctx, lockKey, 5*time.Second)
		if err != nil {
			t.logger.Error("CreateTicketByStaff: Error acquiring lock for trip %s: %v", seatID, err)
			return nil, errors.New("could not contact booking service, please try again")
		}
		if !lockAcquired {
			return nil, errors.New("booking for this trip is currently busy, please try again shortly")
		}
		defer unlock()
	}

	// Acquire locks for end trip seats if round trip
	if input.TicketType == 1 {
		for _, seatID := range input.SeatIDEnd {
			lockKey := fmt.Sprintf("trip-lock:%s", string(seatID))
			lockAcquired, unlock, err := t.ticketRepository.AcquireLock(ctx, lockKey, 5*time.Second)
			if err != nil {
				t.logger.Error("CreateTicketByStaff: Error acquiring lock for trip %s: %v", seatID, err)
				return nil, errors.New("could not contact booking service, please try again")
			}
			if !lockAcquired {
				return nil, errors.New("booking for this trip is currently busy, please try again shortly")
			}
			defer unlock()
		}
	}

	// 2. Batch validate seats for begin trip
	availableSeatRows, err := t.ticketRepository.AreSeatsAvailable(ctx, input.SeatIDBegin)
	if err != nil {
		t.logger.Error("CreateTicketByStaff: Database error during batch seat validation for trip %s: %v", input.TripIDBegin, err)
		return nil, errors.New("error checking seat availability")
	}

	var availableSeatRowsEnd []db.AreSeatsAvailableRow

	// Validate seats for end trip if round trip
	if input.TicketType == 1 {
		availableSeatRowsEnd, err = t.ticketRepository.AreSeatsAvailable(ctx, input.SeatIDEnd)
		if err != nil {
			t.logger.Error("CreateTicketByStaff: Database error during batch seat validation for trip %s: %v", input.TripIDEnd, err)
			return nil, errors.New("error checking seat availability")
		}
	}

	// Validate seat counts match
	if len(availableSeatRows) != len(input.SeatIDBegin) || len(availableSeatRowsEnd) != len(input.SeatIDEnd) {
		return nil, errors.New("one or more selected seats do not exist or belong to a different trip")
	}

	// Check if seats are available for begin trip
	for _, seat := range availableSeatRows {
		if seat.IsBooked {
			return nil, fmt.Errorf("seat %d is already booked or held by another user", seat.ID)
		}
	}

	// Check if seats are available for end trip
	for _, seat := range availableSeatRowsEnd {
		if seat.IsBooked {
			return nil, fmt.Errorf("seat %d is already booked or held by another user", seat.ID)
		}
	}

	// 3. Generate unique ticket ID
	ticketID, err := t.ticketRepository.GenerateUniqueTicketID(ctx)
	if err != nil {
		t.logger.Error("CreateTicketByStaff: Error generating unique ticket ID: %v", err)
		return nil, err
	}

	// 4. Prepare ticket data - For staff bookings, status is immediately confirmed/paid
	ticket := &db.Ticket{
		TicketID:       ticketID,
		CustomerID:     sql.NullInt32{Int32: 0, Valid: false},
		TripIDBegin:    input.TripIDBegin,
		TripIDEnd:      sql.NullString{String: input.TripIDEnd, Valid: input.TicketType == 1},
		Price:          input.Price,
		Status:         models.TicketStatusConfirmed, // Confirmed directly
		PaymentStatus:  models.PaymentStatusPaid,     // Paid directly
		Type:           int16(input.TicketType),
		Name:           utils.ToNullString(input.Name),
		BookedBy:       utils.ToNullString(staffID),
		BookingChannel: models.BookingChannelCounter, // Staff-specific channel
		PolicyID:       input.PolicyID,
		Phone:          utils.ToNullString(input.Phone),
		Email:          utils.ToNullString(input.Email),
		BookingTime:    time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 5. Prepare QR generation event for begin trip
	var ticketDetailsQR []kafkaclient.TicketDetailForQR
	for _, seatID := range input.SeatIDBegin {
		ticketDetailsQR = append(ticketDetailsQR, kafkaclient.TicketDetailForQR{
			SeatID:    seatID,
			QRContent: fmt.Sprintf("TICKET:%s-SEAT:%d", ticketID, seatID),
		})
	}

	// Add QR details for end trip seats if round trip
	if input.TicketType == 1 {
		for _, seatID := range input.SeatIDEnd {
			ticketDetailsQR = append(ticketDetailsQR, kafkaclient.TicketDetailForQR{
				SeatID:    seatID,
				QRContent: fmt.Sprintf("TICKET:%s-SEAT:%d", ticketID, seatID),
			})
		}
	}

	eventPayload := kafkaclient.OrderQRGenerationRequestEvent{
		OrderID:       ticketID,
		CustomerEmail: ticket.Email.String,
		CustomerName:  ticket.Name.String,
		TotalPrice:    ticket.Price,
		Tickets:       ticketDetailsQR,
	}
	payloadBytes, _ := json.Marshal(eventPayload)

	// 6. Prepare repository transaction parameters (matching CreateTicket structure)
	repoParams := repositories.CreateTicketTransactionParams{
		Ticket: db.CreateTicketParams{
			TicketID:       ticket.TicketID,
			TripIDBegin:    input.TripIDBegin,
			TripIDEnd:      sql.NullString{String: input.TripIDEnd, Valid: input.TicketType == 1},
			Type:           ticket.Type,
			CustomerID:     ticket.CustomerID,
			Phone:          ticket.Phone,
			Email:          ticket.Email,
			Name:           ticket.Name,
			Price:          ticket.Price,
			Status:         ticket.Status,
			BookingTime:    time.Now(),
			PaymentStatus:  ticket.PaymentStatus,
			BookingChannel: ticket.BookingChannel,
			PolicyID:       ticket.PolicyID,
			BookedBy:       ticket.BookedBy,
		},
		TicketDetail: db.CreateTicketDetailsParams{
			TicketID:             ticketID,
			PickupLocationBegin:  sql.NullInt32{Int32: input.PickupLocationBegin, Valid: true},
			DropoffLocationEnd:   sql.NullInt32{Int32: input.DropoffLocationEnd, Valid: true},
			PickupLocationEnd:    sql.NullInt32{Int32: input.PickupLocationEnd, Valid: input.TicketType == 1},
			DropoffLocationBegin: sql.NullInt32{Int32: input.DropoffLocationBegin, Valid: input.TicketType == 1},
		},
		SeatIDsBegin: input.SeatIDBegin,
		SeatIDsEnd:   input.SeatIDEnd,
		OutboxEvents: []db.CreateOutboxEventParams{
			{
				ID:      uuid.New(),
				Topic:   t.cfg.Kafka.Topics.OrderQRRequests.Topic,
				Key:     ticketID,
				Payload: payloadBytes,
			},
		},
	}

	// 7. Execute atomic transaction
	err = t.ticketRepository.CreateTicketInTransaction(ctx, repoParams)
	if err != nil {
		t.logger.Error("CreateTicketByStaff: Error in repository's atomic transaction for ticket %s: %v", ticketID, err)
		return nil, fmt.Errorf("failed to finalize staff booking, please try again")
	}

	// 8. Asynchronously update cache for begin trip
	go func(tripID string, seatIDs []int32) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		t.ticketRepository.UpdateCachedAvailableSeats(bgCtx, tripID, seatIDs, "REMOVE")
	}(input.TripIDBegin, input.SeatIDBegin)

	// Update cache for end trip if round trip
	if input.TicketType == 1 {
		go func(tripID string, seatIDs []int32) {
			bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			t.ticketRepository.UpdateCachedAvailableSeats(bgCtx, tripID, seatIDs, "REMOVE")
		}(input.TripIDEnd, input.SeatIDEnd)
	}

	t.logger.Info("Ticket %s created successfully by staff %s. QR generation event placed in outbox.", ticket.TicketID, staffID)
	return ticket, nil
}

func (t *TicketService) GetAvailableSeatsByTripID(ctx context.Context, tripID string) ([]models.SeatReturn, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	seats, err := t.ticketRepository.GetAvailableSeatsFromCache(ctx, tripID)
	if err == nil && len(seats) > 0 {
		t.logger.Info("Available seats for trip %s retrieved from cache", tripID)
		return seats, nil
	}
	if err != nil && !errors.Is(err, repositories.ErrRedisNotAvailable) && !errors.Is(err, repositories.ErrCacheMiss) {
		t.logger.Error("Error getting available seats from cache for trip %s: %v", tripID, err)
	}

	seats, err = t.ticketRepository.GetAvailableSeatsByTripID(ctx, tripID)
	if err != nil {
		t.logger.Error("Error getting available seats from DB for trip %s: %v", tripID, err)
		return nil, err
	}

	go func() {
		cacheCtx, cacheCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cacheCancel()
		if err := t.ticketRepository.CacheAvailableSeats(cacheCtx, tripID, seats); err != nil {
			t.logger.Error("Warning: Failed to cache available seats for trip %s: %v", tripID, err)
		}
	}()

	return seats, nil
}

func (t *TicketService) ExtendSeatHoldTime(ctx context.Context, ticketID string, seatIDs []int32, extendMinutes int) error {
	if extendMinutes <= 0 || extendMinutes > 30 {
		return errors.New("invalid extension time (must be between 1-30 minutes)")
	}

	ticketReturn, err := t.GetTicketByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket %s not found for extending hold: %w", ticketID, err)
	}

	ticketForCache := &db.Ticket{
		TicketID: ticketReturn.TicketID,
	}

	if ticketReturn.PaymentStatus != models.PaymentStatusPending {
		return fmt.Errorf("cannot extend hold time for ticket %s with payment status %d (expected PENDING)", ticketID, ticketReturn.PaymentStatus)
	}

	newExpiry := time.Duration(extendMinutes) * time.Minute
	for _, seatID := range seatIDs {
		err := t.ticketRepository.ExtendSeatHoldTime(ctx, seatID, ticketForCache, newExpiry)
		if err != nil {
			t.logger.Error("Error extending seat hold time for seat %d on ticket %s: %v", seatID, ticketID, err)
			return err
		}
	}

	t.logger.Info("Successfully extended hold time for ticket %s seats by %d minutes", ticketID, extendMinutes)
	return nil
}

func (t *TicketService) ReleaseHeldSeats(ctx context.Context, ticketID string, seatIDs []int32) error {
	for _, seatID := range seatIDs {
		err := t.ticketRepository.ReleaseSeat(ctx, seatID)
		if err != nil {
			t.logger.Error("Warning: Failed to release seat %d from Redis for ticket %s: %v", seatID, ticketID, err)
		}
	}

	err := t.ticketRepository.UpdateSeatTicketsStatus(ctx, ticketID, models.SeatStatusAvailable)
	if err != nil {
		t.logger.Error("Error updating seat ticket status to available/cancelled in DB for ticket %s: %v", ticketID, err)
		return err
	}

	t.logger.Info("Successfully released held seats for ticket %s", ticketID)
	return nil
}

func (t *TicketService) UpdateTicketPaymentStatus(ctx context.Context, ticketID string, paymentStatus int16, ticketStatus int16, tripID string) error {
	err := t.ticketRepository.UpdateTicketPaymentStatus(ctx, ticketID, paymentStatus, ticketStatus, tripID)
	if err != nil {
		t.logger.Error("Error updating payment/ticket status for ticket %s: %v", ticketID, err)
		return err
	}

	// Event publishing logic (original)
	seatIDs, seatErr := t.ticketRepository.GetSeatIDsByTicketID(ctx, ticketID) // These are seat PKs
	if seatErr != nil {
		t.logger.Error("Error retrieving seat IDs for ticket %s for event publishing: %v", ticketID, seatErr)
	} else {
		eventType := "payment_status_updated"
		if paymentStatus == models.PaymentStatusPaid {
			eventType = "payment_completed" // This event should trigger ManagerTicketService to send email with QR
		} else if paymentStatus == models.PaymentStatusFailed {
			eventType = "payment_failed"
		}
		// ... (rest of event publishing for each seatID) ...
		for _, seatID := range seatIDs {
			go func(sID int32, currentEventType string) {
				eventCtx, eventCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer eventCancel()
				statusMsg := &models.SeatStatusMessage{
					EventType: currentEventType, TicketID: ticketID, SeatID: sID, TripID: tripID, Timestamp: time.Now(),
				}
				if errPub := t.ticketRepository.PublishSeatStatusChange(eventCtx, statusMsg); errPub != nil {
					t.logger.Error("Warning: Failed to publish payment status event for ticket %s, seat %d: %v", ticketID, sID, errPub)
				}
			}(seatID, eventType)
		}
	}

	// Cache update logic (original)
	if paymentStatus == models.PaymentStatusPaid {
		// ... (original logic for caching ticket) ...
		ticketReturn, errGet := t.GetTicketByID(ctx, ticketID) // Use existing method to get full ticket data
		if errGet != nil {
			t.logger.Error("Warning: Couldn't retrieve ticket %s for cache update after payment: %v", ticketID, errGet)
		} else if ticketReturn != nil {
			// Map models.TicketReturn to db.Ticket for caching if necessary, or cache TicketReturn directly
			// The original code maps to db.Ticket, ensure all fields used by caching are present
			updatedTicketForCache := &db.Ticket{ /* map fields from ticketReturn to db.Ticket */ }
			// ... (mapping logic as in original code)
			updatedTicketForCache.TicketID = ticketReturn.TicketID
			// ... (fill other fields)
			updatedTicketForCache.PaymentStatus = ticketReturn.PaymentStatus
			updatedTicketForCache.Status = ticketReturn.Status

			if errCache := t.ticketRepository.CacheTicket(ctx, updatedTicketForCache); errCache != nil {
				t.logger.Error("Warning: Failed to update cached ticket %s after payment: %v", ticketID, errCache)
			}
		}
	}
	// NOTE: The actual QR generation and email sending for regular payment success
	// will be handled by ManagerTicketService reacting to the "payment_completed" event
	// or by its UpdateStatusByTicketID method if called directly.

	t.logger.Info("Successfully updated payment status to %d and ticket status to %d for ticket %s", paymentStatus, ticketStatus, ticketID)
	return nil
}

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
	// This URL should come from configuration in a real application
	tripServiceBaseURL = "http://localhost:8082" // Assuming trip-service runs on port 8081
)

// GetTripDetails fetches trip details from the trip-service.
func GetTripDetails(tripID string) (*models.TripInfo, error) {
	if tripID == "" {
		return nil, nil // Or an error: fmt.Errorf("tripID cannot be empty")
	}

	requestURL := fmt.Sprintf("%s/api/v1/trips/%s", tripServiceBaseURL, tripID)

	resp, err := httpClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to request trip details from trip-service for tripID %s: %w", tripID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Trip not found in the other service, not necessarily an error for the whole process
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("trip-service returned status %d for tripID %s: %s", resp.StatusCode, tripID, string(bodyBytes))
	}

	// The Java service wraps its response: {"code": ..., "message": ..., "data": TripObject}
	var serviceResponse struct {
		Code    int              `json:"code"`
		Message string           `json:"message"`
		Data    *models.TripInfo `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&serviceResponse); err != nil {
		return nil, fmt.Errorf("failed to decode trip details response for tripID %s: %w", tripID, err)
	}

	// Check the business logic code from the service if necessary
	if serviceResponse.Code != http.StatusOK {
		// This could mean the trip was not found or another business error occurred
		if serviceResponse.Code == http.StatusNotFound {
			return nil, nil // Consistent with direct 404 handling
		}
		return nil, fmt.Errorf("trip-service reported an issue for tripID %s (message: %s, code: %d)", tripID, serviceResponse.Message, serviceResponse.Code)
	}

	return serviceResponse.Data, nil
}

// CreateTicketAndNotify xử lý việc tạo vé và publish kết quả vào Redis
func (t *TicketService) CreateTicketAndNotify(ctx context.Context, input *models.TicketInput, customerID sql.NullInt32, bookingID string) {
	// 1. Gọi logic tạo vé cốt lõi
	ticket, err := t.CreateTicket(ctx, input, customerID)

	var messageToPublish websocket.Message

	// 2. Chuẩn bị message kết quả
	if err != nil {
		t.logger.Error("[Notify] Failed to create ticket for bookingId %s: %v", bookingID, err)
		errorPayload := gin.H{"error": err.Error()}
		messageToPublish = websocket.Message{Type: "error", Payload: errorPayload}
	} else {
		t.logger.Info("[Notify] Ticket %s created successfully for bookingId %s.", ticket.TicketID, bookingID)
		// Lấy thông tin đầy đủ của vé để trả về cho client
		ticketReturn, repoErr := t.ticketRepository.GetTicketByID(ctx, ticket.TicketID)
		if repoErr != nil {
			t.logger.Error("[Notify] Failed to retrieve ticket details for bookingId %s: %v", bookingID, repoErr)
			errorPayload := gin.H{"error": "Không thể lấy chi tiết vé sau khi tạo."}
			messageToPublish = websocket.Message{Type: "error", Payload: errorPayload}
		} else {
			messageToPublish = websocket.Message{Type: "result", Payload: ticketReturn}
		}
	}

	// 3. Marshal message thành JSON
	payloadBytes, jsonErr := json.Marshal(messageToPublish)
	if jsonErr != nil {
		t.logger.Error("[Notify] Failed to marshal result message for bookingId %s: %v", bookingID, jsonErr)
		return
	}

	// 4. Publish message vào kênh Redis
	channel := fmt.Sprintf("booking-result:%s", bookingID)
	pubErr := t.redisClient.Publish(ctx, channel, payloadBytes).Err()
	if pubErr != nil {
		t.logger.Error("[Notify] Failed to publish result to Redis for bookingId %s: %v", bookingID, pubErr)
	} else {
		t.logger.Info("[Notify] Successfully published result to Redis channel %s", channel)
	}
}

// CancelTicket xử lý việc hủy vé khi ACK timeout
func (t *TicketService) CancelTicket(ctx context.Context, ticketID string) error {
	t.logger.Info("[Cancel] Initiating cancellation for ticketID: %s due to ACK timeout.", ticketID)
	ticket, err := t.ticketRepository.GetTicketByID(ctx, ticketID)
	if err != nil || ticket == nil {
		t.logger.Error("Failed to get ticket info for ticket %s: %v", ticketID, err)
		return fmt.Errorf("ticket %s not found", ticketID)
	}
	var params repositories.UpdateStatusTransactionParams
	params.TicketID = ticketID
	params.OutboxEvents = []db.CreateOutboxEventParams{}
	params.PaymentStatus = models.TicketStatusCancelled
	params.GeneralTicketStatus = models.TicketStatusCancelled
	params.SeatTicketStatus = models.SeatStatusCancelled
	var eventPayload kafkaclient.SeatUpdateEvent
	seatCount := len(ticket.SeatTicketsBegin)
	if seatCount > 0 {
		eventPayload = kafkaclient.SeatUpdateEvent{
			TripID:    ticket.TripIDBegin,
			SeatCount: seatCount,
		}
	}
	releasePayloadBytes, _ := json.Marshal(eventPayload)
	params.OutboxEvents = append(params.OutboxEvents, db.CreateOutboxEventParams{
		ID: uuid.New(), Topic: t.cfg.Kafka.Topics.SeatsReleased.Topic, Key: ticket.TripIDBegin, Payload: releasePayloadBytes,
	})

	if ticket.Type == 1 {
		var eventPayloadEnd kafkaclient.SeatUpdateEvent
		seatCount := len(ticket.SeatTicketsEnd)
		if seatCount > 0 {
			eventPayloadEnd = kafkaclient.SeatUpdateEvent{
				TripID:    ticket.TripIDEnd.String,
				SeatCount: seatCount,
			}
		}
		releasePayloadBytes, _ := json.Marshal(eventPayloadEnd)
		params.OutboxEvents = append(params.OutboxEvents, db.CreateOutboxEventParams{
			ID: uuid.New(), Topic: t.cfg.Kafka.Topics.SeatsReleased.Topic, Key: ticket.TripIDBegin, Payload: releasePayloadBytes,
		})
	}
	var seatIDs []int32
	for _, seatTicket := range ticket.SeatTicketsBegin {
		seatIDs = append(seatIDs, seatTicket.SeatID)
	}

	var seatIDsEnd []int32
	if ticket.Type == 1 {
		for _, seatTicket := range ticket.SeatTicketsEnd {
			seatIDsEnd = append(seatIDsEnd, seatTicket.SeatID)
		}
	}
	if err := t.ticketRepository.UpdateCachedAvailableSeats(ctx, ticket.TripIDBegin, seatIDs, "ADD"); err != nil {
		t.logger.Error("Failed to update cached available seats for trip %s: %v", ticket.TripIDBegin, err)
	}

	go t.ticketRepository.CleanupTicketCache(context.Background(), ticketID, nil, ticket.TripIDBegin)

	if ticket.Type == 1 {
		if err := t.ticketRepository.UpdateCachedAvailableSeats(ctx, ticket.TripIDEnd.String, seatIDsEnd, "ADD"); err != nil {
			t.logger.Error("Failed to update cached available seats for trip %s: %v", ticket.TripIDEnd, err)
		}

		go t.ticketRepository.CleanupTicketCache(context.Background(), ticketID, nil, ticket.TripIDEnd.String)
	}
	// TODO: Thêm logic để giải phóng ghế trong bảng seat_tickets và cập nhật lại cache Redis
	// Ví dụ: t.ticketRepository.ReleaseSeatsForTicket(ctx, ticketID)

	t.logger.Info("[Cancel] Successfully cancelled ticket %s.", ticketID)
	return nil
}

func (t *TicketService) QueueNewBooking(ctx context.Context, bookingID string, input *models.TicketInput, customerID sql.NullInt32) error {
	t.logger.Info("Service: Queueing new booking request for bookingId: %s", bookingID)

	// 1. Tạo trạng thái ban đầu trong Redis Hash
	redisStateKey := fmt.Sprintf("booking:state:%s", bookingID)
	inputBytes, _ := json.Marshal(input)
	stateData := map[string]interface{}{
		"status":       "QUEUED",
		"request_data": string(inputBytes),
		"submitted_at": time.Now().UTC().Format(time.RFC3339),
	}
	if err := t.redisClient.HSet(ctx, redisStateKey, stateData).Err(); err != nil {
		t.logger.Error("Failed to set initial state in Redis for %s: %v", bookingID, err)
		return fmt.Errorf("failed to create booking session: %w", err)
	}
	// Đặt TTL dài (10 phút) để dọn dẹp nếu có lỗi không xử lý được
	t.redisClient.Expire(ctx, redisStateKey, 10*time.Minute)

	// --- THAY ĐỔI MỚI: Thêm vào Sorted Set để theo dõi timeout kết nối WebSocket ---
	connectTimeout := 30 * time.Second // Timeout 30 giây
	expirationTimestamp := time.Now().Add(connectTimeout).UnixMilli()
	if err := t.redisClient.ZAdd(ctx, pendingWsConnectionsKey, redis.Z{Score: float64(expirationTimestamp), Member: bookingID}).Err(); err != nil {
		t.logger.Error("Failed to add bookingId %s to timeout sorted set: %v", bookingID, err)
		t.redisClient.Del(context.Background(), redisStateKey) // Rollback
		return fmt.Errorf("failed to schedule booking timeout: %w", err)
	}

	// 2. Tạo và Publish sự kiện vào Kafka
	event := kafkaclient.BookingRequestEvent{
		BookingID:  bookingID,
		Input:      *input,
		CustomerID: customerID,
	}

	publishCtx := context.Background()

	err := t.publisher.Publish(publishCtx, t.cfg.Kafka.Topics.BookingRequests.Topic, []byte(bookingID), event)
	if err != nil {
		t.logger.Error("Failed to publish booking request event to Kafka for bookingId %s: %v", bookingID, err)
		// Rollback Redis
		t.redisClient.Del(context.Background(), redisStateKey)
		t.redisClient.ZRem(context.Background(), pendingWsConnectionsKey, bookingID)
		return fmt.Errorf("failed to queue booking request: %w", err)
	}

	t.logger.Info("Successfully queued booking request for bookingId: %s", bookingID)
	return nil
}

func (t *TicketService) GetAllTickets(ctx context.Context, page, limit int) (*models.PaginatedTickets, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 10 // Default limit
	}
	offset := (page - 1) * limit

	total, err := t.ticketRepository.GetTotalTicketCount(ctx)
	if err != nil {
		t.logger.Error("Error getting total ticket count: %v", err)
		return nil, errors.New("failed to retrieve ticket count")
	}

	tickets, err := t.ticketRepository.GetAllTickets(ctx, limit, offset)
	if err != nil {
		t.logger.Error("Error retrieving all tickets: %v", err)
		return nil, errors.New("failed to retrieve tickets")
	}

	return &models.PaginatedTickets{
		Tickets: tickets,
		Total:   total,
		Page:    page,
		Limit:   limit,
	}, nil
}
