package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"payment_service/domain/model"        // Your domain models
	"payment_service/internal/db"         // SQLC generated
	"payment_service/internal/repository" // Your repository interface
	// "payment_service/pkg/utils" // Your utility functions
)

// BankServiceInterface defines the methods for bank payment processing.
type BankServiceInterface interface {
	CreateBankPaymentRequest(ctx context.Context, req model.InitialBankPaymentRequest) (*model.BankPaymentDetailsResponse, error)
	ConfirmBankPayment(ctx context.Context, req model.BankPaymentConfirmationRequest) (db.Invoice, error)
	HandleBankPaymentFailed(ctx context.Context, invoiceID uuid.UUID, reason string) (db.Invoice, error)
	// Add RefundBankPayment if needed
}

// BankService handles bank payment logic.
type BankService struct {
	invoiceRepo           repository.InvoiceRepositoryInterface
	invoiceService        InvoiceServiceInterface // To use common invoice logic if any
	accountServiceBaseURL string                  // Base URL for the external Account Service
	httpClient            *http.Client            // HTTP client for external calls
	// Add any other dependencies, like config for bank details to display
}

// NewBankService creates a new BankService.
func NewBankService(
	invoiceRepo repository.InvoiceRepositoryInterface,
	invoiceService InvoiceServiceInterface,
	accountServiceBaseURL string,
	httpClient *http.Client,
) BankServiceInterface {
	return &BankService{
		invoiceRepo:           invoiceRepo,
		invoiceService:        invoiceService,
		accountServiceBaseURL: accountServiceBaseURL,
		httpClient:            httpClient,
	}
}

// AccountServiceAccountResponse mirrors the structure of models.AccountResponse from the bank service.
// It's defined here because bank_service cannot directly import internal models of another service.
type AccountServiceAccountResponse struct {
	ID        int64     `json:"id"`
	OwnerName string    `json:"owner_name"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// convertToSmallestUnit converts a float amount in a major currency unit
// to an int64 amount in the smallest currency unit (e.g., cents).
func convertToSmallestUnit(amount float64, currency string) (int64, error) {
	var multiplier float64
	// This should be expanded based on all supported currencies and their decimal places.
	// The currency codes should ideally come from a shared constant or configuration.
	switch currency {
	case "USD", "EUR": // Currencies with 2 decimal places
		multiplier = 100.0
	case "VND", "JPY": // Currencies with 0 decimal places
		multiplier = 1.0
	default:
		return 0, fmt.Errorf("currency %s is not configured for smallest unit conversion or is unsupported", currency)
	}

	if amount < 0 {
		return 0, fmt.Errorf("amount cannot be negative")
	}

	// Basic check for potential overflow, though a more robust solution might be needed for very large amounts.
	// Max float64 can be much larger than max int64 after multiplication.
	// This is a simplified check.
	if (amount * multiplier) > float64(^uint64(0)>>1) { // Approx check against MaxInt64
		return 0, fmt.Errorf("amount %f %s is too large to convert to smallest unit", amount, currency)
	}

	return int64(amount * multiplier), nil
}

// getAccountDetails fetches account details from the external Account Service.
func (s *BankService) getAccountDetails(ctx context.Context, accountID int64) (*AccountServiceAccountResponse, error) {
	if s.httpClient == nil {
		return nil, fmt.Errorf("HTTP client not configured for AccountService communication")
	}
	if s.accountServiceBaseURL == "" {
		return nil, fmt.Errorf("AccountService base URL not configured")
	}

	// === THAY ĐỔI 1: Xây dựng URL mới ===
	// URL bây giờ trỏ đến endpoint "/me" thay vì "/{id}"
	url := fmt.Sprintf("%s/api/v1/accounts/me", s.accountServiceBaseURL)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to AccountService: %w", err)
	}

	// === THAY ĐỔI 2: Thêm header X-User-ID vào request ===
	httpReq.Header.Set("X-User-ID", strconv.FormatInt(accountID, 10))
	httpReq.Header.Set("Accept", "application/json")
	// httpReq.Header.Set("Authorization", "Bearer <token>") // Nếu cần xác thực giữa các service

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("Error making request to AccountService at %s: %v", url, err)
		return nil, fmt.Errorf("request to AccountService failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Printf("Account with ID %d not found in AccountService (URL: %s)", accountID, url)
		return nil, fmt.Errorf("payer account (ID: %d) not found", accountID)
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("AccountService at %s returned error: %s (status %d), Header(X-User-ID: %d), Body: %s", url, resp.Status, resp.StatusCode, accountID, string(bodyBytes))
		return nil, fmt.Errorf("AccountService error: %s (status %d)", resp.Status, resp.StatusCode)
	}

	var accDetails AccountServiceAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&accDetails); err != nil {
		log.Printf("Failed to decode AccountService response from %s: %v", url, err)
		return nil, fmt.Errorf("failed to decode AccountService response: %w", err)
	}

	return &accDetails, nil
}

// CreateBankPaymentRequest prepares an invoice for bank payment and returns details for the user.
// It now includes a check with an external Account Service for balance and account status.
func (s *BankService) CreateBankPaymentRequest(ctx context.Context, req model.InitialBankPaymentRequest) (*model.BankPaymentDetailsResponse, error) {
	// 0. Validate CustomerID format and parse to int64 for AccountService.
	//    This assumes req.CustomerID is a string representation of the account's int64 ID.
	if req.CustomerID == "" {
		return nil, fmt.Errorf("customer ID is required")
	}
	accountID, err := strconv.ParseInt(req.CustomerID, 10, 64)
	if err != nil {
		log.Printf("Invalid CustomerID format: '%s', error: %v", req.CustomerID, err)
		return nil, fmt.Errorf("invalid customer ID format: must be an integer string")
	}

	// 1. Call Account Service to get account details and perform pre-checks.
	log.Printf("Fetching account details for AccountID: %d for bank payment request.", accountID)
	accountDetails, err := s.getAccountDetails(ctx, accountID)
	if err != nil {
		// Error already logged in getAccountDetails or if it's a specific known error like "account not found"
		return nil, fmt.Errorf("failed to verify payer account: %w", err) // Wrap error for context
	}

	// 2. Validate account status, currency, and balance.
	// Account status check (e.g., "active" as per bank/internal/models/account.go constants like StatusActive)
	// Ideally, "active" should be a shared constant if possible.
	const expectedAccountStatus = "active"
	if accountDetails.Status != expectedAccountStatus {
		log.Printf("Payer account %d is not active. Status: %s", accountID, accountDetails.Status)
		return nil, fmt.Errorf("payer account is not active (status: %s)", accountDetails.Status)
	}

	// Currency check
	if accountDetails.Currency != req.Currency {
		log.Printf("Payer account %d currency mismatch. Account: %s, Request: %s", accountID, accountDetails.Currency, req.Currency)
		return nil, fmt.Errorf("payer account currency (%s) does not match payment currency (%s)", accountDetails.Currency, req.Currency)
	}

	// Balance check
	// Convert requested amount (float64 in major unit) to smallest unit (int64) for comparison.
	requiredAmountSmallestUnit, err := convertToSmallestUnit(req.Amount, req.Currency)
	if err != nil {
		log.Printf("Currency conversion error for account %d, amount %.2f %s: %v", accountID, req.Amount, req.Currency, err)
		return nil, fmt.Errorf("payment currency processing error: %w", err)
	}

	if accountDetails.Balance < requiredAmountSmallestUnit {
		log.Printf("Insufficient balance for payer account %d. Available: %d %s, Required: %d %s",
			accountID, accountDetails.Balance, accountDetails.Currency, requiredAmountSmallestUnit, req.Currency)
		return nil, fmt.Errorf("insufficient funds in payer account (available: %d %s, required: %d %s in smallest unit)",
			accountDetails.Balance, accountDetails.Currency, requiredAmountSmallestUnit, req.Currency)
	}

	log.Printf("Account %d validated successfully for payment. Balance: %d %s, Requested: %.2f %s",
		accountID, accountDetails.Balance, accountDetails.Currency, req.Amount, req.Currency)

	// 3. If all checks passed, proceed to create the invoice.
	invoiceID := uuid.New()
	invoiceNumber := fmt.Sprintf("INV-BANK-%d", time.Now().UnixNano())
	bankTransferCode := fmt.Sprintf("BK-%s", uuid.New().String()[:8])

	createParams := db.CreateInvoiceParams{
		InvoiceID:        invoiceID,
		InvoiceNumber:    invoiceNumber,
		InvoiceType:      sql.NullString{String: req.InvoiceType, Valid: req.InvoiceType != ""},
		CustomerID:       req.CustomerID, // Store original CustomerID
		TicketID:         req.TicketID,
		TotalAmount:      req.Amount,
		FinalAmount:      req.Amount,
		Currency:         sql.NullString{String: req.Currency, Valid: req.Currency != ""},
		PaymentMethod:    sql.NullString{String: string(model.PaymentMethodBank), Valid: true},
		PaymentStatus:    sql.NullString{String: string(model.PaymentStatusAwaitingConfirmation), Valid: true},
		IssueDate:        sql.NullTime{Time: time.Now(), Valid: true},
		Notes:            "Please use the provided transfer code in your bank transfer reference. Account balance verified.",
		BankTransferCode: sql.NullString{String: bankTransferCode, Valid: true},
		// Assuming PayerAccountID might be useful to store on the invoice for reconciliation.
		// This field would need to be added to your db.CreateInvoiceParams and schema if desired.
		// PayerAccountID: sql.NullInt64{Int64: accountID, Valid: true},
	}

	dbInvoice, err := s.invoiceRepo.CreateInvoice(ctx, createParams)
	if err != nil {
		log.Printf("Error creating invoice for bank payment after account check: %v", err)
		return nil, fmt.Errorf("failed to create invoice for bank payment: %w", err)
	}

	paymentInstructions := model.BankPaymentDetailsResponse{
		InvoiceID:                    dbInvoice.InvoiceID,
		BankTransferCode:             bankTransferCode,
		PayableAmount:                dbInvoice.FinalAmount,
		Currency:                     dbInvoice.Currency.String,
		OurBankAccountName:           "YOUR COMPANY NAME",        // From config
		OurBankAccountNumber:         "YOUR BANK ACCOUNT NUMBER", // From config
		OurBankName:                  "YOUR BANK NAME",           // From config
		PaymentReferenceInstructions: fmt.Sprintf("Please include this code in your transfer reference: %s", bankTransferCode),
		DueDate:                      time.Now().Add(48 * time.Hour),
	}

	log.Printf("Bank payment request created for invoice %s, transfer code %s, after successful account check for %d", dbInvoice.InvoiceID, bankTransferCode, accountID)
	return &paymentInstructions, nil
}

// ConfirmBankPayment is called when your system verifies a bank transfer has been received.
func (s *BankService) ConfirmBankPayment(ctx context.Context, req model.BankPaymentConfirmationRequest) (db.Invoice, error) {
	// 1. Fetch the invoice - This logic remains the same.
	var existingInvoice db.Invoice
	var err error

	if req.InvoiceID != uuid.Nil {
		existingInvoice, err = s.invoiceRepo.GetInvoiceByID(ctx, req.InvoiceID)
	} else if req.BankTransferCode != "" {
		existingInvoice, err = s.invoiceRepo.GetInvoiceByBankTransferCode(ctx, sql.NullString{String: req.BankTransferCode, Valid: true})
	} else {
		return db.Invoice{}, fmt.Errorf("invoiceID or bankTransferCode is required to confirm payment")
	}

	if err != nil {
		log.Printf("Error fetching invoice for bank payment confirmation (%s/%s): %v", req.InvoiceID, req.BankTransferCode, err)
		if err == sql.ErrNoRows {
			return db.Invoice{}, fmt.Errorf("invoice not found")
		}
		return db.Invoice{}, fmt.Errorf("invoice not found for confirmation: %w", err)
	}

	// 2. Validate invoice status - This logic remains the same.
	if existingInvoice.PaymentStatus.String != string(model.PaymentStatusAwaitingConfirmation) && existingInvoice.PaymentStatus.String != string(model.PaymentStatusPending) {
		log.Printf("Invoice %s cannot be confirmed. Status: %s", existingInvoice.InvoiceID, existingInvoice.PaymentStatus.String)
		return db.Invoice{}, fmt.Errorf("invoice %s is not awaiting confirmation (status: %s)", existingInvoice.InvoiceID, existingInvoice.PaymentStatus.String)
	}
	// 3b. Convert invoice amount to the smallest unit for the payment service.
	amountToDebit, err := convertToSmallestUnit(existingInvoice.FinalAmount, existingInvoice.Currency.String)
	if err != nil {
		log.Printf("Currency conversion error for invoice %s. Amount: %.2f %s. Error: %v. Marking as failed.",
			existingInvoice.InvoiceID, existingInvoice.FinalAmount, existingInvoice.Currency.String, err)
		return s.HandleBankPaymentFailed(ctx, existingInvoice.InvoiceID, "Internal error during currency conversion.")
	}

	// 3c. Call the Account Service to make the payment.
	log.Printf("Attempting to process payment for invoice %s via AccountService for account %d. Amount: %d %s",
		existingInvoice.InvoiceID, existingInvoice.CustomerID, amountToDebit, existingInvoice.Currency.String)

	err = s.makePaymentOnAccount(ctx, existingInvoice.CustomerID, amountToDebit, existingInvoice.Currency.String)
	if err != nil {
		// The payment failed (e.g., insufficient funds). Mark the invoice as FAILED.
		log.Printf("Payment via AccountService failed for invoice %s. Reason: %v. Marking invoice as failed.", existingInvoice.InvoiceID, err)
		failureReason := fmt.Sprintf("Payment processing failed: %v", err)
		return s.HandleBankPaymentFailed(ctx, existingInvoice.InvoiceID, failureReason)
	}
	log.Printf("Payment processed successfully via AccountService for invoice %s.", existingInvoice.InvoiceID)
	// === END OF NEW LOGIC ===

	// 4. Update local invoice status to COMPLETED - This proceeds only if the external payment was successful.
	updateParams := db.UpdateInvoiceBankPaymentConfirmationParams{
		InvoiceID:          existingInvoice.InvoiceID,
		PaymentStatus:      sql.NullString{String: string(model.PaymentStatusCompleted), Valid: true},
		BankAccountName:    sql.NullString{String: req.PayerAccountName, Valid: req.PayerAccountName != ""},
		BankAccountNumber:  sql.NullString{String: req.PayerAccountNumber, Valid: req.PayerAccountNumber != ""},
		BankName:           sql.NullString{String: req.PayerBankName, Valid: req.PayerBankName != ""},
		BankTransactionID:  sql.NullString{String: req.BankTransactionID, Valid: req.BankTransactionID != ""},
		BankPaymentDetails: req.ConfirmationDetails,
		Notes:              "Payment confirmed and processed via Account Service.", // Updated note
	}

	updatedInvoice, err := s.invoiceRepo.UpdateInvoiceBankPaymentConfirmation(ctx, updateParams)
	if err != nil {
		// CRITICAL: The payment was made, but the local status update failed.
		// This state requires manual intervention or an automated reconciliation process.
		log.Printf("CRITICAL ERROR: Error confirming bank payment for invoice %s after successful AccountService transaction: %v", existingInvoice.InvoiceID, err)
		return db.Invoice{}, fmt.Errorf("CRITICAL: failed to update invoice %s after payment was processed: %w", existingInvoice.InvoiceID, err)
	}

	log.Printf("Bank payment confirmed for invoice %s. Status: COMPLETED", updatedInvoice.InvoiceID)

	// 5. Trigger post-payment actions - This logic remains the same.
	if s.invoiceService != nil {
		s.invoiceService.UpdateTicketStatus(ctx, existingInvoice.TicketID, model.TicketStatusPaid, existingInvoice.InvoiceID)
	}

	return updatedInvoice, nil
}

// HandleBankPaymentFailed marks a bank payment as failed.
func (s *BankService) HandleBankPaymentFailed(ctx context.Context, invoiceID uuid.UUID, reason string) (db.Invoice, error) {
	params := db.UpdateInvoicePaymentFailedParams{
		InvoiceID:     invoiceID,
		PaymentStatus: sql.NullString{String: string(model.PaymentStatusFailed), Valid: true},
		Notes:         reason,
	}
	failedInvoice, err := s.invoiceRepo.UpdateInvoicePaymentFailed(ctx, params)
	if err != nil {
		log.Printf("Error marking bank payment as failed for invoice %s: %v", invoiceID, err)
		return db.Invoice{}, fmt.Errorf("failed to mark bank payment as failed for invoice %s: %w", invoiceID, err)
	}
	log.Printf("Bank payment marked as FAILED for invoice %s. Reason: %s", invoiceID, reason)
	return failedInvoice, nil
}

// Implement RefundBankPayment if necessary

// makePaymentOnAccount calls the external Account Service to perform a payment.
func (s *BankService) makePaymentOnAccount(ctx context.Context, accountID string, amount int64, currency string) error {
	if s.httpClient == nil {
		return fmt.Errorf("HTTP client not configured for AccountService communication")
	}
	if s.accountServiceBaseURL == "" {
		return fmt.Errorf("AccountService base URL not configured")
	}

	// This endpoint is derived from the provided controller code for making a payment.
	url := fmt.Sprintf("%s/api/v1/accounts/payment", s.accountServiceBaseURL)

	// Prepare the request body based on the Account Service's models.PaymentRequest.
	paymentReqBody := map[string]interface{}{
		"amount":   amount,
		"currency": currency,
	}
	bodyBytes, err := json.Marshal(paymentReqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal payment request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create payment request to AccountService: %w", err)
	}

	// Set headers, mirroring the authentication pattern in getAccountDetails.
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", accountID)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		log.Printf("Error making payment request to AccountService at %s: %v", url, err)
		return fmt.Errorf("payment request to AccountService failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-OK responses, which indicate a failure (e.g., insufficient funds).
	if resp.StatusCode != http.StatusOK {
		respBodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("AccountService payment endpoint at %s returned error: %s (status %d) for AccountID %d, Body: %s", url, resp.Status, resp.StatusCode, accountID, string(respBodyBytes))

		// Return a specific error message that can be used upstream.
		return fmt.Errorf("accountService payment failed with status %d: %s", resp.StatusCode, string(respBodyBytes))
	}

	log.Printf("Successfully processed payment of %d %s for account %d via AccountService", amount, currency, accountID)
	return nil
}
