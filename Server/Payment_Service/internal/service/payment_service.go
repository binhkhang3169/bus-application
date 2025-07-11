package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"payment_service/internal/repository"
	"payment_service/pkg/kafkaclient"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	// Đường dẫn đến package config của bạn
	"payment_service/domain/model" // Các model cho API request/response
	"payment_service/internal/db"  // Thư mục chứa code sqlc đã tạo

	"github.com/redis/go-redis/v9"
)

// InvoiceServiceInterface định nghĩa các phương thức cho invoice service
type InvoiceServiceInterface interface {
	CreateInvoiceForVNPay(ctx context.Context, req model.VNPayPaymentRequest, txnRef string) (db.Invoice, error)
	CreateInvoiceForStripe(ctx context.Context, req model.InitialStripePaymentRequest, paymentIntentID string) (db.Invoice, error)
	GetInvoiceByID(ctx context.Context, id uuid.UUID) (db.Invoice, error)
	GetInvoiceByVNPayTxnRef(ctx context.Context, txnRef string) (db.Invoice, error)
	GetInvoiceByStripePaymentIntentID(ctx context.Context, paymentIntentID string) (db.Invoice, error)
	GetLatestCompletedInvoiceByTicketID(ctx context.Context, ticketID string) (db.Invoice, error) // New
	UpdateInvoiceStatusForVNPaySuccess(ctx context.Context, txnRef, vnpBankCode, vnpTxnNo, vnpPayDate string) (db.Invoice, error)
	UpdateInvoiceStatusForStripeSuccess(ctx context.Context, paymentIntentID, chargeID, paymentMethodDetailsJSON string) (db.Invoice, error)
	UpdateInvoiceStatusForPaymentFailure(ctx context.Context, identifier string, method model.PaymentMethod, reason string) (db.Invoice, error)
	UpdateInvoiceStatusForRefund(ctx context.Context, invoiceID uuid.UUID, reason string, refundSpecificIdentifier string) (db.Invoice, error) // Added refundSpecificIdentifier
	UpdateInvoiceStatusForPaymentFailureForUUID(ctx context.Context, identifier uuid.UUID, method model.PaymentMethod, reason string) (db.Invoice, error)
	GetInvoicesByCustomerID(ctx context.Context, customerID string) ([]db.Invoice, error)
	MapDbInvoiceToAPIResponse(invoice db.Invoice) model.GetInvoiceResponse
	MapDbInvoicesToAPIResponses(invoices []db.Invoice) []model.GetInvoiceResponse
	GetAmountInSmallestUnit(amount float64, currencyCode string) (int64, error)
	ConvertSmallestUnitToFloat(amount int64, currencyCode string) (float64, error)
	UpdateTicketStatus(ctx context.Context, ticketID string, statusCode string, invoiceID uuid.UUID) error // Expose for direct use by other services if needed

	CreateInvoiceForBankPayment(ctx context.Context, req model.InitialBankPaymentRequest, bankTransferCode string, dueDate time.Time) (db.Invoice, error) // Assuming this exists or will be added for bank_controller
	ConfirmInvoiceForBankPayment(ctx context.Context, req model.BankPaymentConfirmationRequest) (db.Invoice, error)                                       // Assuming this exists
	MarkInvoiceAsFailedForBankPayment(ctx context.Context, invoiceID uuid.UUID, reason string) (db.Invoice, error)                                        // Assuming this exists
	ProcessStaffDirectPayment(ctx context.Context, req model.StaffDirectPaymentRequest) (db.Invoice, error)                                               // New Method
}

// InvoiceService xử lý logic nghiệp vụ liên quan đến hóa đơn
type InvoiceService struct {
	repo        repository.InvoiceRepositoryInterface // Sử dụng interface
	publisher   *kafkaclient.Publisher                // << THAY ĐỔI: Thay thế URL và http client bằng producer
	redisClient *redis.Client
}

// NewInvoiceService tạo một invoice service mới
func NewInvoiceService(repo repository.InvoiceRepositoryInterface, publisher *kafkaclient.Publisher, redisClient *redis.Client) InvoiceServiceInterface {
	return &InvoiceService{
		repo:        repo,
		publisher:   publisher,
		redisClient: redisClient,
	}
}

func (s *InvoiceService) setInvoiceExpiration(ctx context.Context, invoiceID string, expiration time.Duration) error {
	key := fmt.Sprintf("invoice_expiry:%s", invoiceID)
	return s.redisClient.Set(ctx, key, invoiceID, expiration).Err()
}

func (s *InvoiceService) clearInvoiceExpiration(ctx context.Context, invoiceID string) error {
	key := fmt.Sprintf("invoice_expiry:%s", invoiceID)
	return s.redisClient.Del(ctx, key).Err()
}

// generateInvoiceNumber tạo một số hóa đơn duy nhất
func generateInvoiceNumber(prefix string) string {
	return fmt.Sprintf("%s-%s-%s", prefix, time.Now().Format("20060102"), uuid.New().String()[:8])
}

func floatToDecimalString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func int64ToDecimalString(amountSmallestUnit int64, currency string) string {
	isVND := strings.ToLower(currency) == "vnd"
	if isVND {
		return fmt.Sprintf("%d.00", amountSmallestUnit)
	}
	dollars := amountSmallestUnit / 100
	cents := amountSmallestUnit % 100
	return fmt.Sprintf("%d.%02d", dollars, cents)
}

func int64ToDecimalFloat(amountSmallestUnit int64, currency string) float64 {
	isVND := strings.ToLower(currency) == "vnd"
	if isVND {
		return float64(amountSmallestUnit)
	}
	return float64(amountSmallestUnit) / 100.0
}

// GetAmountInSmallestUnit converts a float amount to its smallest currency unit (e.g., cents).
func (s *InvoiceService) GetAmountInSmallestUnit(amount float64, currencyCode string) (int64, error) {
	currency := strings.ToLower(currencyCode)
	// For currencies like JPY, VND that don't have subunits, or where Stripe expects the amount as is.
	// Stripe treats VND as a zero-decimal currency.
	if currency == "vnd" || currency == "jpy" { // Add other zero-decimal currencies if needed
		return int64(math.Round(amount)), nil
	}
	// For currencies with 2 decimal places like USD, EUR
	return int64(math.Round(amount * 100)), nil
}

// ConvertSmallestUnitToFloat converts an amount from smallest unit to float.
func (s *InvoiceService) ConvertSmallestUnitToFloat(amount int64, currencyCode string) (float64, error) {
	currency := strings.ToLower(currencyCode)
	if currency == "vnd" || currency == "jpy" {
		return float64(amount), nil
	}
	return float64(amount) / 100.0, nil
}

// CreateInvoiceForVNPay tạo hóa đơn mới cho giao dịch VNPay
func (s *InvoiceService) CreateInvoiceForVNPay(ctx context.Context, req model.VNPayPaymentRequest, txnRef string) (db.Invoice, error) {
	finalAmount := req.Amount - req.DiscountAmount + req.TaxAmount
	invoiceID := uuid.New()

	params := db.CreateInvoiceParams{
		InvoiceID:      invoiceID,
		InvoiceNumber:  generateInvoiceNumber("VNP"),
		InvoiceType:    sql.NullString{String: req.InvoiceType, Valid: req.InvoiceType != ""},
		CustomerID:     req.CustomerID,
		TicketID:       req.TicketID,
		TotalAmount:    req.Amount,
		DiscountAmount: sql.NullString{String: floatToDecimalString(req.DiscountAmount), Valid: req.DiscountAmount != 0},
		TaxAmount:      sql.NullString{String: floatToDecimalString(req.TaxAmount), Valid: req.TaxAmount != 0},
		FinalAmount:    finalAmount,
		Currency:       sql.NullString{String: "vnd", Valid: true}, // VNPay is typically VND
		PaymentStatus:  sql.NullString{String: string(model.PaymentStatusPending), Valid: true},
		PaymentMethod:  sql.NullString{String: string(model.PaymentMethodVNPay), Valid: true},
		IssueDate:      sql.NullTime{Time: time.Now(), Valid: true},
		Notes:          req.Notes,
		VnpayTxnRef:    sql.NullString{String: txnRef, Valid: txnRef != ""},
	}

	createdInvoice, err := s.repo.CreateInvoice(ctx, params)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to create VNPay invoice: %w", err)
	}

	expiration := 15 * time.Minute
	if err := s.setInvoiceExpiration(ctx, createdInvoice.InvoiceID.String(), expiration); err != nil {
		// Nếu không set được key Redis, có thể log lỗi nhưng vẫn tiếp tục
		// hoặc coi đây là lỗi nghiêm trọng và rollback.
		log.Printf("CẢNH BÁO: Không thể set key hết hạn cho hóa đơn %s: %v", createdInvoice.InvoiceID, err)
	}

	return createdInvoice, nil
}

// CreateInvoiceForStripe chuẩn bị hóa đơn cho Stripe PaymentIntent
func (s *InvoiceService) CreateInvoiceForStripe(ctx context.Context, req model.InitialStripePaymentRequest, paymentIntentID string) (db.Invoice, error) {
	totalAmountFloat, _ := s.ConvertSmallestUnitToFloat(req.Amount, req.Currency)
	discountAmountFloat, _ := s.ConvertSmallestUnitToFloat(req.DiscountAmount, req.Currency)
	taxAmountFloat, _ := s.ConvertSmallestUnitToFloat(req.TaxAmount, req.Currency)

	finalAmountFloat := totalAmountFloat - discountAmountFloat + taxAmountFloat
	invoiceID := uuid.New()

	params := db.CreateInvoiceParams{
		InvoiceID:             invoiceID,
		InvoiceNumber:         generateInvoiceNumber("STR"),
		InvoiceType:           sql.NullString{String: req.InvoiceType, Valid: req.InvoiceType != ""},
		CustomerID:            req.CustomerID,
		TicketID:              req.TicketID,
		TotalAmount:           totalAmountFloat,
		DiscountAmount:        sql.NullString{String: floatToDecimalString(discountAmountFloat), Valid: req.DiscountAmount != 0},
		TaxAmount:             sql.NullString{String: floatToDecimalString(taxAmountFloat), Valid: req.TaxAmount != 0},
		FinalAmount:           finalAmountFloat,
		Currency:              sql.NullString{String: strings.ToLower(req.Currency), Valid: req.Currency != ""},
		PaymentStatus:         sql.NullString{String: string(model.PaymentStatusPending), Valid: true},
		PaymentMethod:         sql.NullString{String: string(model.PaymentMethodStripe), Valid: true},
		IssueDate:             sql.NullTime{Time: time.Now(), Valid: true},
		Notes:                 req.Notes,
		StripePaymentIntentID: sql.NullString{String: paymentIntentID, Valid: paymentIntentID != ""},
	}

	createdInvoice, err := s.repo.CreateInvoice(ctx, params)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to create Stripe invoice: %w", err)
	}

	expiration := 15 * time.Minute
	if err := s.setInvoiceExpiration(ctx, createdInvoice.InvoiceID.String(), expiration); err != nil {
		log.Printf("CẢNH BÁO: Không thể set key hết hạn cho hóa đơn %s: %v", createdInvoice.InvoiceID, err)
	}

	return createdInvoice, nil
}

// GetInvoiceByID lấy hóa đơn theo ID
func (s *InvoiceService) GetInvoiceByID(ctx context.Context, id uuid.UUID) (db.Invoice, error) {
	invoice, err := s.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: GetInvoiceByID failed: %w", err)
	}
	return invoice, nil
}

// GetInvoiceByVNPayTxnRef lấy hóa đơn theo VNPay TxnRef
func (s *InvoiceService) GetInvoiceByVNPayTxnRef(ctx context.Context, txnRef string) (db.Invoice, error) {
	nsTxnRef := sql.NullString{String: txnRef, Valid: txnRef != ""}
	invoice, err := s.repo.GetInvoiceByVNPayTxnRef(ctx, nsTxnRef)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: GetInvoiceByVNPayTxnRef failed for %s: %w", txnRef, err)
	}
	return invoice, nil
}

// GetInvoiceByStripePaymentIntentID lấy hóa đơn theo Stripe PaymentIntentID
func (s *InvoiceService) GetInvoiceByStripePaymentIntentID(ctx context.Context, paymentIntentID string) (db.Invoice, error) {
	nsPIID := sql.NullString{String: paymentIntentID, Valid: paymentIntentID != ""}
	invoice, err := s.repo.GetInvoiceByStripePaymentIntentID(ctx, nsPIID)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: GetInvoiceByStripePaymentIntentID failed for %s: %w", paymentIntentID, err)
	}
	return invoice, nil
}

// GetInvoiceByBankTransferCode lấy hoá đơn theo Bank bankTransferCode
func (s *InvoiceService) GetInvoiceByBankTransferCode(ctx context.Context, bankTransferCode string) (db.Invoice, error) {
	nsBank := sql.NullString{String: bankTransferCode, Valid: bankTransferCode != ""}
	invoice, err := s.repo.GetInvoiceByBankTransferCode(ctx, nsBank)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: GetInvoiceByBankTransferCode failed for bank")
	}
	return invoice, nil
}

// GetLatestCompletedInvoiceByTicketID retrieves the latest completed invoice for a ticket ID
func (s *InvoiceService) GetLatestCompletedInvoiceByTicketID(ctx context.Context, ticketID string) (db.Invoice, error) {
	invoice, err := s.repo.GetLatestCompletedInvoiceByTicketID(ctx, ticketID)
	if err != nil {
		// Error is already specific from repository
		return db.Invoice{}, fmt.Errorf("service: %w", err)
	}
	return invoice, nil
}

// UpdateInvoiceStatusForVNPaySuccess cập nhật trạng thái hóa đơn khi VNPay thành công
func (s *InvoiceService) UpdateInvoiceStatusForVNPaySuccess(ctx context.Context, txnRef, vnpBankCode, vnpTxnNo, vnpPayDate string) (db.Invoice, error) {
	params := db.UpdateInvoiceVNPayStatusParams{
		VnpayTxnRef:   sql.NullString{String: txnRef, Valid: txnRef != ""},
		PaymentStatus: sql.NullString{String: string(model.PaymentStatusCompleted), Valid: true},
		VnpayBankCode: sql.NullString{String: vnpBankCode, Valid: vnpBankCode != ""},
		VnpayTxnNo:    sql.NullString{String: vnpTxnNo, Valid: vnpTxnNo != ""},
		VnpayPayDate:  sql.NullString{String: vnpPayDate, Valid: vnpPayDate != ""},
	}
	invoice, err := s.repo.UpdateInvoiceVNPayStatus(ctx, params)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to update invoice for VNPay success (TxnRef: %s): %w", txnRef, err)
	}

	if err := s.UpdateTicketStatus(ctx, invoice.TicketID, model.TicketStatusPaid, invoice.InvoiceID); err != nil {
		log.Printf("Warning: service: failed to update ticket status for invoice %s (ticket %s) after VNPay success: %v", invoice.InvoiceID, invoice.TicketID, err)
	}
	s.publishSuccessNotification(invoice)

	return invoice, nil
}

// UpdateInvoiceStatusForStripeSuccess cập nhật trạng thái hóa đơn khi Stripe thành công
func (s *InvoiceService) UpdateInvoiceStatusForStripeSuccess(ctx context.Context, paymentIntentID, chargeID, paymentMethodDetailsJSON string) (db.Invoice, error) {
	params := db.UpdateInvoiceStripePaymentSuccessParams{
		StripePaymentIntentID:      sql.NullString{String: paymentIntentID, Valid: paymentIntentID != ""},
		PaymentStatus:              sql.NullString{String: string(model.PaymentStatusCompleted), Valid: true},
		StripeChargeID:             sql.NullString{String: chargeID, Valid: chargeID != ""},
		StripePaymentMethodDetails: paymentMethodDetailsJSON,
	}
	invoice, err := s.repo.UpdateInvoiceStripePaymentSuccess(ctx, params)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to update invoice for Stripe success (PI_ID: %s): %w", paymentIntentID, err)
	}

	if err := s.UpdateTicketStatus(ctx, invoice.TicketID, model.TicketStatusPaid, invoice.InvoiceID); err != nil {
		log.Printf("Warning: service: failed to update ticket status for invoice %s (ticket %s) after Stripe success: %v", invoice.InvoiceID, invoice.TicketID, err)
	}
	s.publishSuccessNotification(invoice)
	return invoice, nil
}

func (s *InvoiceService) publishSuccessNotification(invoice db.Invoice) {
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		notificationTopic := "notifications_topic" // Nên lấy từ config

		userID := invoice.CustomerID // Lấy CustomerID

		event := kafkaclient.NotificationEvent{
			UserID:  &userID,
			Type:    "PAYMENT_SUCCESS",
			Title:   "Thanh toán thành công",
			Message: fmt.Sprintf("Thanh toán cho hóa đơn %s trị giá %.2f %s đã được xác nhận thành công.", invoice.InvoiceNumber, invoice.FinalAmount, invoice.Currency.String),
		}

		if err := s.publisher.Publish(bgCtx, notificationTopic, []byte(userID), event); err != nil {
			log.Printf("CRITICAL: Failed to publish payment success notification for invoice %s: %v", invoice.InvoiceID, err)
		}
	}()
}

func (s *InvoiceService) publishFailedNotification(invoice db.Invoice) {
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		notificationTopic := "notifications_topic" // Nên lấy từ config

		userID := invoice.CustomerID // Lấy CustomerID

		event := kafkaclient.NotificationEvent{
			UserID:  &userID,
			Type:    "PAYMENT_FAILED",
			Title:   "Thanh toán thất bại",
			Message: fmt.Sprintf("Thanh toán cho hóa đơn %s trị giá %.2f %s đã được xác nhận thất bại.", invoice.InvoiceNumber, invoice.FinalAmount, invoice.Currency.String),
		}

		if err := s.publisher.Publish(bgCtx, notificationTopic, []byte(userID), event); err != nil {
			log.Printf("CRITICAL: Failed to publish payment success notification for invoice %s: %v", invoice.InvoiceID, err)
		}
	}()
}

// UpdateInvoiceStatusForPaymentFailure cập nhật trạng thái hóa đơn khi thanh toán thất bại
func (s *InvoiceService) UpdateInvoiceStatusForPaymentFailureForUUID(ctx context.Context, identifier uuid.UUID, method model.PaymentMethod, reason string) (db.Invoice, error) {

	invoice, err := s.GetInvoiceByID(ctx, identifier)

	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: could not find invoice by identifier")
	}

	err = s.clearInvoiceExpiration(ctx, invoice.InvoiceID.String())
	if err != nil {
		log.Printf("Info: %w", err)
	}

	if invoice.PaymentStatus.String == string(model.PaymentStatusCompleted) ||
		invoice.PaymentStatus.String == string(model.PaymentStatusRefunded) ||
		invoice.PaymentStatus.String == string(model.PaymentStatusFailed) {
		log.Printf("Info: service: invoice %s already in final state %s. Skipping failure update.", invoice.InvoiceID, invoice.PaymentStatus.String)
		return invoice, nil
	}

	updateParams := db.UpdateInvoicePaymentFailedParams{
		InvoiceID:     invoice.InvoiceID,
		PaymentStatus: sql.NullString{String: string(model.PaymentStatusFailed), Valid: true},
		Notes:         fmt.Sprintf("Payment failed via %s. Reason: %s", method, reason),
	}

	updatedInvoice, err := s.repo.UpdateInvoicePaymentFailed(ctx, updateParams)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to update invoice %s for payment failure (%s): %w", invoice.InvoiceID, method, err)
	}

	if err := s.UpdateTicketStatus(ctx, updatedInvoice.TicketID, model.TicketStatusFailed, updatedInvoice.InvoiceID); err != nil {
		log.Printf("Warning: service: failed to update ticket status for invoice %s (ticket %s) after payment failure: %v", updatedInvoice.InvoiceID, updatedInvoice.TicketID, err)
	}
	s.publishFailedNotification(updatedInvoice)
	return updatedInvoice, nil
}

// UpdateInvoiceStatusForPaymentFailure cập nhật trạng thái hóa đơn khi thanh toán thất bại
func (s *InvoiceService) UpdateInvoiceStatusForPaymentFailure(ctx context.Context, identifier string, method model.PaymentMethod, reason string) (db.Invoice, error) {
	var invoice db.Invoice
	var err error

	switch method {
	case model.PaymentMethodStripe:
		invoice, err = s.GetInvoiceByStripePaymentIntentID(ctx, identifier)
	case model.PaymentMethodVNPay:
		invoice, err = s.GetInvoiceByVNPayTxnRef(ctx, identifier)
	case model.PaymentMethodBank:
		invoice, err = s.GetInvoiceByBankTransferCode(ctx, identifier)
	default:
		return db.Invoice{}, fmt.Errorf("service: unsupported payment method for failure update: %s", method)
	}

	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: could not find invoice by identifier '%s' for payment method '%s' to mark as failed: %w", identifier, method, err)
	}

	if invoice.PaymentStatus.String == string(model.PaymentStatusCompleted) ||
		invoice.PaymentStatus.String == string(model.PaymentStatusRefunded) ||
		invoice.PaymentStatus.String == string(model.PaymentStatusFailed) {
		log.Printf("Info: service: invoice %s already in final state %s. Skipping failure update.", invoice.InvoiceID, invoice.PaymentStatus.String)
		return invoice, nil
	}

	updateParams := db.UpdateInvoicePaymentFailedParams{
		InvoiceID:     invoice.InvoiceID,
		PaymentStatus: sql.NullString{String: string(model.PaymentStatusFailed), Valid: true},
		Notes:         fmt.Sprintf("Payment failed via %s. Reason: %s", method, reason),
	}

	updatedInvoice, err := s.repo.UpdateInvoicePaymentFailed(ctx, updateParams)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to update invoice %s for payment failure (%s): %w", invoice.InvoiceID, method, err)
	}

	if err := s.UpdateTicketStatus(ctx, updatedInvoice.TicketID, model.TicketStatusFailed, updatedInvoice.InvoiceID); err != nil {
		log.Printf("Warning: service: failed to update ticket status for invoice %s (ticket %s) after payment failure: %v", updatedInvoice.InvoiceID, updatedInvoice.TicketID, err)
	}
	s.publishFailedNotification(updatedInvoice)
	return updatedInvoice, nil
}

// UpdateInvoiceStatusForRefund cập nhật trạng thái hóa đơn khi hoàn tiền
func (s *InvoiceService) UpdateInvoiceStatusForRefund(ctx context.Context, invoiceID uuid.UUID, reason string, refundSpecificIdentifier string) (db.Invoice, error) {
	invoice, err := s.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: could not find invoice %s to mark as refunded: %w", invoiceID, err)
	}

	if invoice.PaymentStatus.String != string(model.PaymentStatusCompleted) {
		return db.Invoice{}, fmt.Errorf("service: invoice %s is not completed, cannot refund. Current status: %s", invoiceID, invoice.PaymentStatus.String)
	}

	fullReason := reason
	if refundSpecificIdentifier != "" {
		fullReason = fmt.Sprintf("Refund ID: %s. Reason: %s", refundSpecificIdentifier, reason)
	}

	params := db.UpdateInvoiceStatusGeneralParams{
		InvoiceID:     invoiceID,
		PaymentStatus: sql.NullString{String: string(model.PaymentStatusRefunded), Valid: true},
		Notes:         fullReason,
	}

	updatedInvoice, err := s.repo.UpdateInvoiceStatusGeneral(ctx, params)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to update invoice %s to refunded: %w", invoiceID, err)
	}

	if err := s.UpdateTicketStatus(ctx, updatedInvoice.TicketID, model.TicketStatusRefunded, updatedInvoice.InvoiceID); err != nil {
		log.Printf("Warning: service: failed to update ticket status for invoice %s (ticket %s) after refund: %v", updatedInvoice.InvoiceID, updatedInvoice.TicketID, err)
	}
	return updatedInvoice, nil
}

// GetInvoicesByCustomerID lấy tất cả hóa đơn của một khách hàng
func (s *InvoiceService) GetInvoicesByCustomerID(ctx context.Context, customerID string) ([]db.Invoice, error) {
	invoices, err := s.repo.ListInvoicesByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("service: GetInvoicesByCustomerID failed for customer %s: %w", customerID, err)
	}
	return invoices, nil
}

// UpdateTicketStatus gọi ticket service bên ngoài để cập nhật trạng thái vé
func (s *InvoiceService) UpdateTicketStatus(ctx context.Context, ticketID string, statusCode string, invoiceID uuid.UUID) error {
	if s.publisher == nil {
		log.Println("Warning: Kafka publisher is not configured. Skipping ticket status update.")
		return nil
	}
	if ticketID == "" {
		return nil
	}

	err := s.clearInvoiceExpiration(ctx, invoiceID.String())
	if err != nil {
		log.Printf("Info: %w", err)
	}
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		ticketTopic := "ticket_status_updates" // Nên lấy từ config
		event := kafkaclient.TicketStatusUpdateEvent{
			TicketID:   ticketID,
			StatusCode: statusCode,
		}

		if err := s.publisher.Publish(bgCtx, ticketTopic, []byte(ticketID), event); err != nil {
			log.Printf("CRITICAL: Failed to publish ticket status update event for ticket %s: %v", ticketID, err)
		}
	}()

	return nil
}

// MapDbInvoiceToAPIResponse chuyển đổi db.Invoice sang model.GetInvoiceResponse
func (s *InvoiceService) MapDbInvoiceToAPIResponse(invoice db.Invoice) model.GetInvoiceResponse {
	totalAmount := invoice.TotalAmount
	var discountAmount float64
	if invoice.DiscountAmount.Valid {
		parsed, err := strconv.ParseFloat(invoice.DiscountAmount.String, 64)
		if err == nil {
			discountAmount = parsed
		} else {
			log.Printf("Warning: MapDbInvoiceToAPIResponse: Failed to parse DiscountAmount '%s': %v", invoice.DiscountAmount.String, err)
		}
	}

	var taxAmount float64
	if invoice.TaxAmount.Valid {
		parsed, err := strconv.ParseFloat(invoice.TaxAmount.String, 64)
		if err == nil {
			taxAmount = parsed
		} else {
			log.Printf("Warning: MapDbInvoiceToAPIResponse: Failed to parse TaxAmount '%s': %v", invoice.TaxAmount.String, err)
		}
	}
	finalAmount := invoice.FinalAmount

	var issueDateStr string
	if invoice.IssueDate.Valid {
		issueDateStr = invoice.IssueDate.Time.Format("2006-01-02 15:04:05")
	}

	return model.GetInvoiceResponse{
		InvoiceID:                  invoice.InvoiceID,
		InvoiceNumber:              invoice.InvoiceNumber,
		InvoiceType:                invoice.InvoiceType.String,
		CustomerID:                 invoice.CustomerID,
		TicketID:                   invoice.TicketID,
		TotalAmount:                totalAmount,
		DiscountAmount:             discountAmount,
		TaxAmount:                  taxAmount,
		FinalAmount:                finalAmount,
		Currency:                   invoice.Currency.String, // Added currency
		PaymentStatus:              invoice.PaymentStatus.String,
		PaymentMethod:              invoice.PaymentMethod.String,
		IssueDate:                  issueDateStr,
		Notes:                      invoice.Notes, // db.Invoice.Notes is sql.NullString
		CreatedAt:                  invoice.CreatedAt.Time.Format("2006-01-02 15:04:05"),
		UpdatedAt:                  invoice.UpdatedAt.Time.Format("2006-01-02 15:04:05"),
		VNPayTxnRef:                invoice.VnpayTxnRef.String,
		VNPayBankCode:              invoice.VnpayBankCode.String,
		VNPayTxnNo:                 invoice.VnpayTxnNo.String,
		VNPayPayDate:               invoice.VnpayPayDate.String,
		StripePaymentIntentID:      invoice.StripePaymentIntentID.String,
		StripeChargeID:             invoice.StripeChargeID.String,
		StripeCustomerID:           invoice.StripeCustomerID.String,
		StripePaymentMethodDetails: invoice.StripePaymentMethodDetails, // db.Invoice.StripePaymentMethodDetails is sql.NullString
	}
}

// MapDbInvoicesToAPIResponses chuyển đổi một slice db.Invoice sang slice model.GetInvoiceResponse
func (s *InvoiceService) MapDbInvoicesToAPIResponses(invoices []db.Invoice) []model.GetInvoiceResponse {
	apiInvoices := make([]model.GetInvoiceResponse, len(invoices))
	for i, inv := range invoices {
		apiInvoices[i] = s.MapDbInvoiceToAPIResponse(inv)
	}
	return apiInvoices
}

// CreateInvoiceForBankPayment creates an invoice for a bank payment request.
func (s *InvoiceService) CreateInvoiceForBankPayment(ctx context.Context, req model.InitialBankPaymentRequest, bankTransferCode string, dueDate time.Time) (db.Invoice, error) {
	invoiceID := uuid.New()
	// Assuming req.Amount is the final amount for bank transfers, as InitialBankPaymentRequest doesn't have discount/tax.
	// If discount/tax were applicable, they'd need to be part of the request model.
	finalAmount := req.Amount

	notes := fmt.Sprintf("Bank transfer payment. Please use reference code: %s. Payment due by: %s.",
		bankTransferCode, dueDate.Format("2006-01-02"))

	params := db.CreateInvoiceParams{
		InvoiceID:          invoiceID,
		InvoiceNumber:      generateInvoiceNumber("BANK"),
		InvoiceType:        sql.NullString{String: req.InvoiceType, Valid: req.InvoiceType != ""},
		CustomerID:         req.CustomerID,                                                                       //
		TicketID:           req.TicketID,                                                                         //
		TotalAmount:        req.Amount,                                                                           // Assuming TotalAmount is the same as FinalAmount for this request type
		DiscountAmount:     sql.NullString{String: "0.00", Valid: true},                                          // No discount in InitialBankPaymentRequest
		TaxAmount:          sql.NullString{String: "0.00", Valid: true},                                          // No tax in InitialBankPaymentRequest
		FinalAmount:        finalAmount,                                                                          //
		Currency:           sql.NullString{String: strings.ToLower(req.Currency), Valid: req.Currency != ""},     //
		PaymentStatus:      sql.NullString{String: string(model.PaymentStatusAwaitingConfirmation), Valid: true}, //
		PaymentMethod:      sql.NullString{String: string(model.PaymentMethodBank), Valid: true},                 //
		IssueDate:          sql.NullTime{Time: time.Now(), Valid: true},
		Notes:              notes,
		BankTransferCode:   sql.NullString{String: bankTransferCode, Valid: bankTransferCode != ""}, //
		BankPaymentDetails: fmt.Sprintf("Due Date: %s", dueDate.Format("2006-01-02")),
	}

	createdInvoice, err := s.repo.CreateInvoice(ctx, params) //
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to create bank payment invoice: %w", err)
	}

	expiration := 15 * time.Minute
	if err := s.setInvoiceExpiration(ctx, createdInvoice.InvoiceID.String(), expiration); err != nil {
		log.Printf("CẢNH BÁO: Không thể set key hết hạn cho hóa đơn %s: %v", createdInvoice.InvoiceID, err)
	}

	return createdInvoice, nil
}

// ConfirmInvoiceForBankPayment confirms a bank payment after funds are verified.
func (s *InvoiceService) ConfirmInvoiceForBankPayment(ctx context.Context, req model.BankPaymentConfirmationRequest) (db.Invoice, error) {
	var invoice db.Invoice
	var err error

	if req.InvoiceID != uuid.Nil {
		invoice, err = s.repo.GetInvoiceByID(ctx, req.InvoiceID) //
	} else if req.BankTransferCode != "" {
		invoice, err = s.repo.GetInvoiceByBankTransferCode(ctx, sql.NullString{String: req.BankTransferCode, Valid: true}) //
	} else {
		return db.Invoice{}, fmt.Errorf("service: either InvoiceID or BankTransferCode is required to confirm bank payment")
	}

	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to retrieve invoice for bank payment confirmation: %w", err)
	}

	// Check current status
	if invoice.PaymentStatus.String != string(model.PaymentStatusAwaitingConfirmation) && invoice.PaymentStatus.String != string(model.PaymentStatusPending) {
		// Allow confirmation even if it was PENDING if somehow AWAITING_CONFIRMATION was skipped
		// but ideally it should be AWAITING_CONFIRMATION
		if invoice.PaymentStatus.String == string(model.PaymentStatusCompleted) {
			log.Printf("Info: service: bank payment for invoice %s already confirmed.", invoice.InvoiceID)
			return invoice, nil // Already completed
		}
		return db.Invoice{}, fmt.Errorf("service: invoice %s is not in a confirmable state for bank payment. Current status: %s", invoice.InvoiceID, invoice.PaymentStatus.String)
	}

	// Optional: Validate amount received
	if req.AmountReceived < invoice.FinalAmount {
		// Handle partial payment or underpayment scenario based on business rules.
		// For now, log a warning. Confirmation might proceed, or you might mark as FAILED or requires manual intervention.
		log.Printf("Warning: service: Amount received (%.2f %s) for invoice %s is less than final amount (%.2f %s).",
			req.AmountReceived, req.CurrencyReceived, invoice.InvoiceID, invoice.FinalAmount, invoice.Currency.String)
		// Or return an error:
		// return db.Invoice{}, fmt.Errorf("service: amount received %.2f is less than invoice amount %.2f for invoice %s", req.AmountReceived, invoice.FinalAmount, invoice.InvoiceID)
	}
	if strings.ToLower(req.CurrencyReceived) != strings.ToLower(invoice.Currency.String) {
		log.Printf("Warning: service: Currency received (%s) for invoice %s does not match invoice currency (%s).",
			req.CurrencyReceived, invoice.InvoiceID, invoice.Currency.String)
		// Or return an error
	}

	confirmationNotes := fmt.Sprintf("Payment confirmed by %s on %s.",
		req.ConfirmedBy, req.ConfirmationTimestamp.Format("2006-01-02 15:04:05"))
	if req.ConfirmationDetails != "" {
		confirmationNotes += " Details: " + req.ConfirmationDetails
	}
	if invoice.Notes != "" {
		confirmationNotes = invoice.Notes + " | " + confirmationNotes
	}

	params := db.UpdateInvoiceBankPaymentConfirmationParams{ //
		InvoiceID:          invoice.InvoiceID,
		PaymentStatus:      sql.NullString{String: string(model.PaymentStatusCompleted), Valid: true},           //
		BankAccountName:    sql.NullString{String: req.PayerAccountName, Valid: req.PayerAccountName != ""},     //
		BankAccountNumber:  sql.NullString{String: req.PayerAccountNumber, Valid: req.PayerAccountNumber != ""}, //
		BankName:           sql.NullString{String: req.PayerBankName, Valid: req.PayerBankName != ""},           //
		BankTransactionID:  sql.NullString{String: req.BankTransactionID, Valid: req.BankTransactionID != ""},   //
		BankPaymentDetails: fmt.Sprintf("Confirmed Amount: %.2f %s. Timestamp: %s. Ref: %s", req.AmountReceived, req.CurrencyReceived, req.ConfirmationTimestamp.Format(time.RFC3339), req.BankTransactionID),
		Notes:              confirmationNotes,
	}

	updatedInvoice, err := s.repo.UpdateInvoiceBankPaymentConfirmation(ctx, params) //
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to update invoice for bank payment confirmation (InvoiceID: %s): %w", invoice.InvoiceID, err)
	}

	if err := s.UpdateTicketStatus(ctx, updatedInvoice.TicketID, model.TicketStatusPaid, updatedInvoice.InvoiceID); err != nil { //
		log.Printf("Warning: service: failed to update ticket status for invoice %s (ticket %s) after bank payment confirmation: %v", updatedInvoice.InvoiceID, updatedInvoice.TicketID, err)
	}
	return updatedInvoice, nil
}

// MarkInvoiceAsFailedForBankPayment marks a bank payment invoice as failed.
func (s *InvoiceService) MarkInvoiceAsFailedForBankPayment(ctx context.Context, invoiceID uuid.UUID, reason string) (db.Invoice, error) {
	invoice, err := s.repo.GetInvoiceByID(ctx, invoiceID) //
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to retrieve invoice %s to mark as failed for bank payment: %w", invoiceID, err)
	}

	// Check current status to avoid marking already finalized (completed/refunded) or already failed invoices again
	currentStatus := model.PaymentStatus(invoice.PaymentStatus.String)
	if currentStatus == model.PaymentStatusCompleted || currentStatus == model.PaymentStatusRefunded || currentStatus == model.PaymentStatusFailed {
		log.Printf("Info: service: invoice %s already in a final state (%s). Cannot mark as failed for bank payment.", invoiceID, currentStatus)
		return invoice, fmt.Errorf("service: invoice %s already in a final state (%s)", invoiceID, currentStatus)
	}

	failureNotes := fmt.Sprintf("Bank payment failed. Reason: %s.", reason)
	if invoice.Notes != "" {
		failureNotes = invoice.Notes + " | " + failureNotes
	}

	params := db.UpdateInvoicePaymentFailedParams{ //
		InvoiceID:     invoiceID,
		PaymentStatus: sql.NullString{String: string(model.PaymentStatusFailed), Valid: true}, //
		Notes:         failureNotes,
	}

	updatedInvoice, err := s.repo.UpdateInvoicePaymentFailed(ctx, params) //
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to update invoice %s to failed for bank payment: %w", invoiceID, err)
	}

	if err := s.UpdateTicketStatus(ctx, updatedInvoice.TicketID, model.TicketStatusFailed, updatedInvoice.InvoiceID); err != nil { //
		log.Printf("Warning: service: failed to update ticket status for invoice %s (ticket %s) after bank payment failure: %v", updatedInvoice.InvoiceID, updatedInvoice.TicketID, err)
	}
	return updatedInvoice, nil
}

// ProcessStaffDirectPayment handles a direct payment processed by staff
func (s *InvoiceService) ProcessStaffDirectPayment(ctx context.Context, req model.StaffDirectPaymentRequest) (db.Invoice, error) {
	finalAmount := req.Amount - req.DiscountAmount + req.TaxAmount
	if finalAmount <= 0 && req.Amount > 0 { // Ensure final amount is positive if initial amount was positive
		finalAmount = req.Amount // Or handle as an error if preferred
	}
	invoiceID := uuid.New()

	// Combine notes
	fullNotes := fmt.Sprintf("Processed by Staff ID: %s.", req.StaffID)
	if req.TransactionReference != "" {
		fullNotes += fmt.Sprintf(" Reference: %s.", req.TransactionReference)
	}
	if req.Notes != "" {
		fullNotes += fmt.Sprintf(" Notes: %s.", req.Notes)
	}

	// Determine currency and smallest unit conversion if necessary.
	// For simplicity, assuming req.Amount is already in the major unit for display and storage.
	// If amounts were to be stored in smallest units for these direct payments, conversion logic would be needed here.

	params := db.CreateInvoiceParams{
		InvoiceID:      invoiceID,
		InvoiceNumber:  generateInvoiceNumber("STAFF"),
		InvoiceType:    sql.NullString{String: req.InvoiceType, Valid: req.InvoiceType != ""},
		CustomerID:     req.CustomerID,
		TicketID:       req.TicketID,
		TotalAmount:    req.Amount, // Assuming req.Amount is the gross amount before discounts
		DiscountAmount: sql.NullString{String: floatToDecimalString(req.DiscountAmount), Valid: req.DiscountAmount > 0},
		TaxAmount:      sql.NullString{String: floatToDecimalString(req.TaxAmount), Valid: req.TaxAmount > 0},
		FinalAmount:    finalAmount,
		Currency:       sql.NullString{String: strings.ToLower(req.Currency), Valid: req.Currency != ""},
		PaymentStatus:  sql.NullString{String: string(model.PaymentStatusCompleted), Valid: true}, // Direct payment is completed
		PaymentMethod:  sql.NullString{String: req.PaymentMethod, Valid: req.PaymentMethod != ""},
		IssueDate:      sql.NullTime{Time: time.Now(), Valid: true},
		Notes:          fullNotes,
		// VNPay, Stripe, and standard Bank Transfer specific fields will be null/empty by default
	}

	createdInvoice, err := s.repo.CreateInvoice(ctx, params)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("service: failed to create staff direct payment invoice: %w", err)
	}

	// Update ticket status
	if err := s.UpdateTicketStatus(ctx, createdInvoice.TicketID, model.TicketStatusPaid, createdInvoice.InvoiceID); err != nil {
		log.Printf("Warning: service: failed to update ticket status for invoice %s (ticket %s) after staff direct payment: %v", createdInvoice.InvoiceID, createdInvoice.TicketID, err)
		// Depending on policy, you might want to decide if the invoice creation should be rolled back or if logging is sufficient.
	}

	return createdInvoice, nil
}
