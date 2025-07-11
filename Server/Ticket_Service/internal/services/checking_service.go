package services

import (
	"context"
	"fmt"
	"ticket-service/domain/models"
	"ticket-service/internal/repositories"
	"ticket-service/pkg/utils"
)

type ICheckinService interface {
	ProcessCheckin(ctx context.Context, ticketID string, seatTicketID int32, tripID string, note string) (*models.CheckinResponse, error)
	// Mới: Lấy danh sách check-in cho một chuyến đi
	GetCheckinsForTrip(ctx context.Context, tripID string) ([]models.CheckinResponse, error)
}

type CheckinService struct {
	checkinRepo repositories.CheckinRepositoryInterface
	logger      utils.Logger
	// You might need ticketRepository if you need more ticket details or validation not covered by checkinRepo
}

func NewCheckinService(checkinRepo repositories.CheckinRepositoryInterface, logger utils.Logger) ICheckinService {
	return &CheckinService{
		checkinRepo: checkinRepo,
		logger:      logger,
	}
}

// Sửa đổi: Logic xử lý check-in
func (s *CheckinService) ProcessCheckin(ctx context.Context, ticketID string, seatTicketID int32, tripID string, note string) (*models.CheckinResponse, error) {
	s.logger.Info("Processing check-in for SeatTicketID: %d, TicketID: %s, TripID: %s", seatTicketID, ticketID, tripID)

	// Lấy thông tin chi tiết của vé ghế bằng ID của nó
	seatTicketDetails, err := s.checkinRepo.GetSeatTicketDetails(ctx, seatTicketID, ticketID)
	if err != nil {
		s.logger.Error("Failed to get details for seat_ticket_id %d: %v", seatTicketID, err)
		return nil, fmt.Errorf("invalid or non-existent seat ticket ID: %d", seatTicketID)
	}

	// Xác thực: Kiểm tra xem ticket_id và trip_id từ QR có khớp với dữ liệu trong DB không
	if seatTicketDetails.TicketID != ticketID {
		s.logger.Info("TicketID mismatch for seat_ticket_id %d. Expected %s, got %s", seatTicketID, seatTicketDetails.TicketID, ticketID)
		return nil, fmt.Errorf("check-in failed: ticket information mismatch")
	}
	if seatTicketDetails.TripID != tripID {
		s.logger.Info("TripID mismatch for seat_ticket_id %d. Expected %s, got %s", seatTicketID, seatTicketDetails.TripID, tripID)
		return nil, fmt.Errorf("check-in failed: this ticket is not for this trip")
	}

	// Xác thực trạng thái vé
	if seatTicketDetails.Status != models.SeatStatusConfirmed {
		s.logger.Info("SeatTicket %d is not in a checkable state. Current status: %d", seatTicketID, seatTicketDetails.Status)
		var statusMsg string
		switch seatTicketDetails.Status {
		case models.SeatStatusPendingPayment:
			statusMsg = "pending payment"
		case models.SeatStatusCancelled:
			statusMsg = "cancelled"
		case models.SeatStatusCheckedIn:
			statusMsg = "already checked in"
		default:
			statusMsg = "in an invalid state for check-in"
		}
		return nil, fmt.Errorf("ticket %s is %s", ticketID, statusMsg)
	}

	// Thực hiện check-in
	checkedInRecord, err := s.checkinRepo.PerformCheckin(
		ctx,
		seatTicketID,
		ticketID,
		tripID,
		seatTicketDetails.SeatName,
		note,
		models.SeatStatusCheckedIn,
		models.TicketStatusUsed,
	)
	if err != nil {
		s.logger.Error("Failed to perform check-in for seat_ticket_id %d: %v", seatTicketID, err)
		return nil, fmt.Errorf("check-in failed: %w", err)
	}

	s.logger.Info("Successfully checked in Ticket %s, SeatTicketID %d, Seat: %s, Trip: %s", ticketID, seatTicketID, seatTicketDetails.SeatName.String, tripID)

	return &models.CheckinResponse{
		TicketID:    checkedInRecord.TicketID,
		TripID:      checkedInRecord.TripID,
		SeatName:    checkedInRecord.SeatName,
		CheckedInAt: checkedInRecord.CheckedInAt,
		CheckinNote: checkedInRecord.Note,
		Message:     "Check-in successful",
	}, nil
}

// Mới: Implement logic lấy danh sách check-in
func (s *CheckinService) GetCheckinsForTrip(ctx context.Context, tripID string) ([]models.CheckinResponse, error) {
	s.logger.Info("Fetching all check-in records for TripID: %s", tripID)

	checkinRecords, err := s.checkinRepo.GetAllCheckinsByTripID(ctx, tripID)
	if err != nil {
		s.logger.Error("Failed to get check-in records for TripID %s: %v", tripID, err)
		return nil, err
	}

	// Map từ db model sang response model
	var response []models.CheckinResponse
	for _, record := range checkinRecords {
		response = append(response, models.CheckinResponse{
			TicketID:    record.TicketID,
			TripID:      record.TripID,
			SeatName:    record.SeatName,
			CheckedInAt: record.CheckedInAt,
			CheckinNote: record.Note,
			Message:     "Retrieved check-in record",
		})
	}

	return response, nil
}
