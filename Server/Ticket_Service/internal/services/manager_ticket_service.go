package services

import (
	"context"
	"encoding/json"
	"fmt"
	"ticket-service/config"
	"ticket-service/domain/models"
	"ticket-service/internal/db"
	"ticket-service/internal/repositories"
	"ticket-service/pkg/emailclient"
	"ticket-service/pkg/kafkaclient"
	"ticket-service/pkg/utils"

	"github.com/google/uuid"
)

type IManagerTicketService interface {
	GenerateSeats() []string
	CreateManagerTicketsForTrip(ctx context.Context, tripID string) ([]db.Seat, error)
	UpdateStatusByTicketID(ctx context.Context, ticketID string, statusCode string) error
	InitializeExpirationHandlers(ctx context.Context) error
}

type ManagerTicketService struct {
	managerTicketRepository repositories.ManagerTicketInterface
	ticketRepository        repositories.TicketRepositoryInterface
	logger                  utils.Logger
	cfg                     config.Config
	publisher               *kafkaclient.Publisher // << UPDATED
	KafkaEmailPublisher     *emailclient.EmailClient
}

func NewManagerTicketService(
	managerTicketRepository repositories.ManagerTicketInterface,
	ticketRepository repositories.TicketRepositoryInterface,
	logger utils.Logger,
	cfg config.Config,
	publisher *kafkaclient.Publisher, // << UPDATED,
	emailclient *emailclient.EmailClient,

) IManagerTicketService {
	return &ManagerTicketService{
		managerTicketRepository: managerTicketRepository,
		ticketRepository:        ticketRepository,
		logger:                  logger,
		cfg:                     cfg,
		publisher:               publisher, // << UPDATED
		KafkaEmailPublisher:     emailclient,
	}
}

func (m *ManagerTicketService) GenerateSeats() []string {
	seats := []string{}
	for i := 1; i <= 15; i++ {
		seatNumber := fmt.Sprintf("%d", i) // Convert i to string
		seats = append(seats, "A"+seatNumber)
		seats = append(seats, "B"+seatNumber)
	}
	return seats
}

func (s *ManagerTicketService) CreateManagerTicketsForTrip(ctx context.Context, tripID string) ([]db.Seat, error) {
	seats := s.GenerateSeats() // giả sử trả về danh sách seatID như ["A1", "A2", ...]
	var created []db.Seat

	for _, seatID := range seats {
		seat := db.Seat{
			TripID:   tripID,
			SeatName: utils.ToNullString(seatID),
		}
		err := s.managerTicketRepository.CreateManagerTicket(ctx, &seat)
		if err != nil {
			return nil, err
		}
		created = append(created, seat)
	}
	return created, nil
}

// Khởi tạo handler để lắng nghe sự kiện hết hạn của Redis
func (s *ManagerTicketService) InitializeExpirationHandlers(ctx context.Context) error {
	s.logger.Info("Initializing Redis expiration handlers...")
	err := s.ticketRepository.SetupExpirationHandler(ctx)
	if err != nil {
		s.logger.Error("Failed to setup Redis expiration handler: %v", err)
		return err
	}
	s.logger.Info("Redis expiration handlers initialized successfully")
	return nil
}

// // Cập nhật trạng thái vé theo kết quả thanh toán
// func (s *ManagerTicketService) UpdateStatusByTicketID(ctx context.Context, ticketID string, statusCode string) error {
// 	s.logger.Info("Processing update for ticket %s with status code %s", ticketID, statusCode)

// 	// 1. Fetch ticket information ONCE.
// 	ticket, err := s.ticketRepository.GetTicketByID(ctx, ticketID)
// 	if err != nil || ticket == nil {
// 		s.logger.Error("Failed to get ticket information for ticket %s: %v", ticketID, err)
// 		return fmt.Errorf("ticket %s not found", ticketID)
// 	}

// 	// 2. Call repository to update status in DB ONCE.
// 	err = s.managerTicketRepository.UpdateStatusByTicketID(ctx, ticketID, statusCode)
// 	if err != nil {
// 		s.logger.Error("Failed to update ticket status in DB for ticket %s: %v", ticketID, err)
// 		return err
// 	}

// 	// 3. Handle specific logic based on the outcome
// 	switch statusCode {
// 	case "1": // Payment successful
// 		s.logger.Info("Payment successful for ticket %s. Publishing QR generation requests.", ticketID)
// 		// << REPLACED: Thay vì tạo QR và gửi mail, hãy gửi sự kiện cho mỗi ghế
// 		var ticketDetails []kafkaclient.TicketDetailForQR
// 		for _, seatTicket := range ticket.SeatTickets {
// 			ticketDetails = append(ticketDetails, kafkaclient.TicketDetailForQR{
// 				SeatID:    seatTicket.SeatID,
// 				QRContent: fmt.Sprintf("TICKET:%s-SEAT:%d", ticket.TicketID, seatTicket.SeatID),
// 			})
// 		}

// 		// Tạo sự kiện tổng cho cả đơn hàng
// 		eventPayload := kafkaclient.OrderQRGenerationRequestEvent{
// 			OrderID:       ticket.TicketID,
// 			CustomerEmail: ticket.Email.String,
// 			CustomerName:  ticket.Name.String,
// 			TotalPrice:    ticket.Price,
// 			Tickets:       ticketDetails,
// 		}

// 		// Gửi đi một sự kiện duy nhất
// 		eventCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()

// 		// Sử dụng một topic mới cho các yêu cầu theo lô
// 		err := s.publisher.Publish(eventCtx, s.cfg.Kafka.Topics.OrderQRRequests.Topic, []byte(ticket.TicketID), eventPayload)
// 		if err != nil {
// 			s.logger.Error("CRITICAL: Failed to publish order_qr_request event for Order %s. Error: %v", ticket.TicketID, err)
// 		} else {
// 			s.logger.Info("Successfully published QR generation request for order %s", ticket.TicketID)
// 		}

// 	case "2": // Payment failed
// 		// Publish an event to release the held seats
// 		go func() {
// 			eventCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 			defer cancel()
// 			seatCount := len(ticket.SeatTickets)
// 			if seatCount > 0 {
// 				eventPayload := kafkaclient.SeatUpdateEvent{
// 					TripID:    ticket.TripID,
// 					SeatCount: seatCount,
// 				}
// 				err := s.publisher.Publish(eventCtx, s.cfg.Kafka.Topics.SeatsReleased.Topic, []byte(ticket.TripID), eventPayload)
// 				if err != nil {
// 					s.logger.Error("CRITICAL: Failed to publish seats_released event for TripID %s. Error: %v", ticket.TripID, err)
// 				}
// 			}
// 		}()
// 	}

// 	s.logger.Info("Successfully processed status update for ticket %s to %s and handled notifications.", ticketID, statusCode)
// 	return nil
// }

func (s *ManagerTicketService) UpdateStatusByTicketID(ctx context.Context, ticketID string, statusCode string) error {
	s.logger.Info("Processing status update for ticket %s with status code %s", ticketID, statusCode)

	ticket, err := s.ticketRepository.GetTicketByID(ctx, ticketID)
	if err != nil || ticket == nil {
		s.logger.Error("Failed to get ticket info for ticket %s: %v", ticketID, err)
		return fmt.Errorf("ticket %s not found", ticketID)
	}

	isSuccess := statusCode == "1"
	expectedStatus := models.TicketStatusConfirmed
	if !isSuccess {
		expectedStatus = int(models.TicketStatusCancelled)
	}
	if int(ticket.Status) == expectedStatus {
		s.logger.Info("Ticket %s is already in the target status %d. Skipping update.", ticketID, expectedStatus)
		return nil
	}

	var params repositories.UpdateStatusTransactionParams
	params.TicketID = ticketID
	params.OutboxEvents = []db.CreateOutboxEventParams{}

	switch statusCode {
	case "1":
		params.PaymentStatus = models.PaymentStatusPaid
		params.GeneralTicketStatus = models.TicketStatusConfirmed
		params.SeatTicketStatus = models.SeatStatusConfirmed

		var ticketDetails []kafkaclient.TicketDetailForQR
		for _, seatTicket := range append(ticket.SeatTicketsBegin, ticket.SeatTicketsEnd...) {
			ticketDetails = append(ticketDetails, kafkaclient.TicketDetailForQR{
				SeatID:    seatTicket.SeatID,
				QRContent: fmt.Sprintf("TICKET:%s-SEAT:%d", ticket.TicketID, seatTicket.SeatID),
			})
		}

		// Tạo sự kiện tổng cho cả đơn hàng
		eventPayload := kafkaclient.OrderQRGenerationRequestEvent{
			OrderID:       ticket.TicketID,
			CustomerEmail: ticket.Email.String,
			CustomerName:  ticket.Name.String,
			TotalPrice:    ticket.Price,
			Tickets:       ticketDetails,
		}

		// Gửi đi một sự kiện duy nhất
		// eventCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()

		// // // Sử dụng một topic mới cho các yêu cầu theo lô
		// // err := s.publisher.Publish(eventCtx, s.cfg.Kafka.Topics.OrderQRRequests.Topic, []byte(ticket.TicketID), eventPayload)
		// // if err != nil {
		// // 	s.logger.Error("CRITICAL: Failed to publish order_qr_request event for Order %s. Error: %v", ticket.TicketID, err)
		// // } else {
		// // 	s.logger.Info("Successfully published QR generation request for order %s", ticket.TicketID)
		// // }

		payloadBytes, _ := json.Marshal(eventPayload)
		params.OutboxEvents = append(params.OutboxEvents, db.CreateOutboxEventParams{
			ID: uuid.New(), Topic: s.cfg.Kafka.Topics.OrderQRRequests.Topic, Key: ticketID, Payload: payloadBytes,
		})

	case "3":
		params.PaymentStatus = models.PaymentStatusFailed
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
		// go func() {
		// 	// eventCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// 	// defer cancel()

		// 		// err := s.publisher.Publish(eventCtx, s.cfg.Kafka.Topics.SeatsReleased.Topic, []byte(ticket.TripIDBegin), eventPayload)
		// 		// if err != nil {
		// 		// 	s.logger.Error("CRITICAL: Failed to publish seats_released event for TripID %s. Error: %v", ticket.TripIDBegin, err)
		// 		// }
		// 	}
		// }()
		releasePayloadBytes, _ := json.Marshal(eventPayload)
		params.OutboxEvents = append(params.OutboxEvents, db.CreateOutboxEventParams{
			ID: uuid.New(), Topic: s.cfg.Kafka.Topics.SeatsReleased.Topic, Key: ticket.TripIDBegin, Payload: releasePayloadBytes,
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
			// go func() {
			// 	eventCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			// 	defer cancel()

			// 		err := s.publisher.Publish(eventCtx, s.cfg.Kafka.Topics.SeatsReleased.Topic, []byte(ticket.TripIDEnd.String), eventPayload)
			// 		if err != nil {
			// 			s.logger.Error("CRITICAL: Failed to publish seats_released event for TripID %s. Error: %v", ticket.TripIDEnd.String, err)
			// 		}
			// 	}
			// }()
			releasePayloadBytes, _ := json.Marshal(eventPayloadEnd)
			params.OutboxEvents = append(params.OutboxEvents, db.CreateOutboxEventParams{
				ID: uuid.New(), Topic: s.cfg.Kafka.Topics.SeatsReleased.Topic, Key: ticket.TripIDBegin, Payload: releasePayloadBytes,
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
		if err := s.ticketRepository.UpdateCachedAvailableSeats(ctx, ticket.TripIDBegin, seatIDs, "ADD"); err != nil {
			s.logger.Error("Failed to update cached available seats for trip %s: %v", ticket.TripIDBegin, err)
		}

		go s.ticketRepository.CleanupTicketCache(context.Background(), ticketID, nil, ticket.TripIDBegin)

		if ticket.Type == 1 {
			if err := s.ticketRepository.UpdateCachedAvailableSeats(ctx, ticket.TripIDEnd.String, seatIDsEnd, "ADD"); err != nil {
				s.logger.Error("Failed to update cached available seats for trip %s: %v", ticket.TripIDEnd, err)
			}

			go s.ticketRepository.CleanupTicketCache(context.Background(), ticketID, nil, ticket.TripIDEnd.String)
		}

	default:
		return fmt.Errorf("invalid status code: %s", statusCode)
	}

	if err := s.managerTicketRepository.UpdateStatusInTransaction(ctx, params); err != nil {
		s.logger.Error("Failed to execute status update transaction for ticket %s: %v", ticketID, err)
		return err
	}

	s.logger.Info("Successfully committed status update for ticket %s to %s.", ticketID, statusCode)
	return nil
}
