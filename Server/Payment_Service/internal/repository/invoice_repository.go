package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"payment_service/internal/db" // Your SQLC generated package
)

// InvoiceRepositoryInterface defines the methods for invoice repository
type InvoiceRepositoryInterface interface {
	CreateInvoice(ctx context.Context, arg db.CreateInvoiceParams) (db.Invoice, error)
	GetInvoiceByID(ctx context.Context, id uuid.UUID) (db.Invoice, error)
	GetInvoiceByVNPayTxnRef(ctx context.Context, vnpayTxnRef sql.NullString) (db.Invoice, error)
	GetInvoiceByStripePaymentIntentID(ctx context.Context, stripePaymentIntentID sql.NullString) (db.Invoice, error)
	GetInvoiceByBankTransferCode(ctx context.Context, bankTransferCode sql.NullString) (db.Invoice, error) // New method
	GetLatestCompletedInvoiceByTicketID(ctx context.Context, ticketID string) (db.Invoice, error)
	ListInvoicesByCustomerID(ctx context.Context, customerID string) ([]db.Invoice, error)
	UpdateInvoiceVNPayStatus(ctx context.Context, arg db.UpdateInvoiceVNPayStatusParams) (db.Invoice, error)
	UpdateInvoiceStripePaymentSuccess(ctx context.Context, arg db.UpdateInvoiceStripePaymentSuccessParams) (db.Invoice, error)
	UpdateInvoiceStripePaymentIntent(ctx context.Context, arg db.UpdateInvoiceStripePaymentIntentParams) (db.Invoice, error)
	UpdateInvoiceBankPaymentRequest(ctx context.Context, arg db.UpdateInvoiceBankPaymentRequestParams) (db.Invoice, error)           // New method
	UpdateInvoiceBankPaymentConfirmation(ctx context.Context, arg db.UpdateInvoiceBankPaymentConfirmationParams) (db.Invoice, error) // New method
	UpdateInvoicePaymentFailed(ctx context.Context, arg db.UpdateInvoicePaymentFailedParams) (db.Invoice, error)
	UpdateInvoiceStatusGeneral(ctx context.Context, arg db.UpdateInvoiceStatusGeneralParams) (db.Invoice, error)
	GetDB() *sql.DB
}

// InvoiceRepository handles database operations for invoices
type InvoiceRepository struct {
	dbConn *sql.DB
	*db.Queries
}

// NewInvoiceRepository creates a new InvoiceRepository
func NewInvoiceRepository(dbConn *sql.DB) InvoiceRepositoryInterface {
	return &InvoiceRepository{
		dbConn:  dbConn,
		Queries: db.New(dbConn), // db.New comes from SQLC generated code
	}
}

// GetDB returns the underlying sql.DB instance
func (r *InvoiceRepository) GetDB() *sql.DB {
	return r.dbConn
}

// CreateInvoice creates a new invoice in the database
func (r *InvoiceRepository) CreateInvoice(ctx context.Context, arg db.CreateInvoiceParams) (db.Invoice, error) {
	if arg.InvoiceID == uuid.Nil {
		arg.InvoiceID = uuid.New()
	}
	invoice, err := r.Queries.CreateInvoice(ctx, arg)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("repository: CreateInvoice failed: %w", err)
	}
	return invoice, nil
}

// GetInvoiceByID retrieves an invoice by its ID
func (r *InvoiceRepository) GetInvoiceByID(ctx context.Context, id uuid.UUID) (db.Invoice, error) {
	invoice, err := r.Queries.GetInvoiceByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Invoice{}, fmt.Errorf("repository: GetInvoiceByID - invoice with ID %s not found: %w", id, err)
		}
		return db.Invoice{}, fmt.Errorf("repository: GetInvoiceByID failed for ID %s: %w", id, err)
	}
	return invoice, nil
}

// GetInvoiceByVNPayTxnRef retrieves an invoice by its VNPay transaction reference
func (r *InvoiceRepository) GetInvoiceByVNPayTxnRef(ctx context.Context, vnpayTxnRef sql.NullString) (db.Invoice, error) {
	invoice, err := r.Queries.GetInvoiceByVNPayTxnRef(ctx, vnpayTxnRef)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Invoice{}, fmt.Errorf("repository: GetInvoiceByVNPayTxnRef - invoice with VNPay TxnRef %s not found: %w", vnpayTxnRef.String, err)
		}
		return db.Invoice{}, fmt.Errorf("repository: GetInvoiceByVNPayTxnRef failed for TxnRef %s: %w", vnpayTxnRef.String, err)
	}
	return invoice, nil
}

// GetInvoiceByStripePaymentIntentID retrieves an invoice by its Stripe Payment Intent ID
func (r *InvoiceRepository) GetInvoiceByStripePaymentIntentID(ctx context.Context, stripePaymentIntentID sql.NullString) (db.Invoice, error) {
	invoice, err := r.Queries.GetInvoiceByStripePaymentIntentID(ctx, stripePaymentIntentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Invoice{}, fmt.Errorf("repository: GetInvoiceByStripePaymentIntentID - invoice with Stripe PI ID %s not found: %w", stripePaymentIntentID.String, err)
		}
		return db.Invoice{}, fmt.Errorf("repository: GetInvoiceByStripePaymentIntentID failed for PI ID %s: %w", stripePaymentIntentID.String, err)
	}
	return invoice, nil
}

// GetInvoiceByBankTransferCode retrieves an invoice by its bank transfer code
func (r *InvoiceRepository) GetInvoiceByBankTransferCode(ctx context.Context, bankTransferCode sql.NullString) (db.Invoice, error) {
	invoice, err := r.Queries.GetInvoiceByBankTransferCode(ctx, bankTransferCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Invoice{}, fmt.Errorf("repository: GetInvoiceByBankTransferCode - invoice with Bank Transfer Code %s not found: %w", bankTransferCode.String, err)
		}
		return db.Invoice{}, fmt.Errorf("repository: GetInvoiceByBankTransferCode failed for Code %s: %w", bankTransferCode.String, err)
	}
	return invoice, nil
}

// GetLatestCompletedInvoiceByTicketID retrieves the latest completed invoice for a given ticket ID
func (r *InvoiceRepository) GetLatestCompletedInvoiceByTicketID(ctx context.Context, ticketID string) (db.Invoice, error) {
	invoice, err := r.Queries.GetLatestCompletedInvoiceByTicketID(ctx, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Invoice{}, fmt.Errorf("repository: GetLatestCompletedInvoiceByTicketID - no completed invoice found for ticket ID %s: %w", ticketID, err)
		}
		return db.Invoice{}, fmt.Errorf("repository: GetLatestCompletedInvoiceByTicketID failed for ticket ID %s: %w", ticketID, err)
	}
	return invoice, nil
}

// ListInvoicesByCustomerID retrieves all invoices for a given customer ID
func (r *InvoiceRepository) ListInvoicesByCustomerID(ctx context.Context, customerID string) ([]db.Invoice, error) {
	invoices, err := r.Queries.ListInvoicesByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("repository: ListInvoicesByCustomerID failed for CustomerID %s: %w", customerID, err)
	}
	if invoices == nil { // SQLC might return nil slice if no rows, ensure consistent empty slice
		return []db.Invoice{}, nil
	}
	return invoices, nil
}

// UpdateInvoiceVNPayStatus updates the status of an invoice for a VNPay transaction
func (r *InvoiceRepository) UpdateInvoiceVNPayStatus(ctx context.Context, arg db.UpdateInvoiceVNPayStatusParams) (db.Invoice, error) {
	invoice, err := r.Queries.UpdateInvoiceVNPayStatus(ctx, arg)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("repository: UpdateInvoiceVNPayStatus failed: %w", err)
	}
	return invoice, nil
}

// UpdateInvoiceStripePaymentSuccess updates the status of an invoice for a successful Stripe payment
func (r *InvoiceRepository) UpdateInvoiceStripePaymentSuccess(ctx context.Context, arg db.UpdateInvoiceStripePaymentSuccessParams) (db.Invoice, error) {
	invoice, err := r.Queries.UpdateInvoiceStripePaymentSuccess(ctx, arg)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("repository: UpdateInvoiceStripePaymentSuccess failed: %w", err)
	}
	return invoice, nil
}

// UpdateInvoiceStripePaymentIntent updates an invoice with Stripe Payment Intent details
func (r *InvoiceRepository) UpdateInvoiceStripePaymentIntent(ctx context.Context, arg db.UpdateInvoiceStripePaymentIntentParams) (db.Invoice, error) {
	invoice, err := r.Queries.UpdateInvoiceStripePaymentIntent(ctx, arg)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("repository: UpdateInvoiceStripePaymentIntent failed: %w", err)
	}
	return invoice, nil
}

// UpdateInvoiceBankPaymentRequest updates an invoice when a bank payment is initiated
func (r *InvoiceRepository) UpdateInvoiceBankPaymentRequest(ctx context.Context, arg db.UpdateInvoiceBankPaymentRequestParams) (db.Invoice, error) {
	invoice, err := r.Queries.UpdateInvoiceBankPaymentRequest(ctx, arg)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("repository: UpdateInvoiceBankPaymentRequest failed: %w", err)
	}
	return invoice, nil
}

// UpdateInvoiceBankPaymentConfirmation updates an invoice when a bank payment is confirmed
func (r *InvoiceRepository) UpdateInvoiceBankPaymentConfirmation(ctx context.Context, arg db.UpdateInvoiceBankPaymentConfirmationParams) (db.Invoice, error) {
	invoice, err := r.Queries.UpdateInvoiceBankPaymentConfirmation(ctx, arg)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("repository: UpdateInvoiceBankPaymentConfirmation failed: %w", err)
	}
	return invoice, nil
}

// UpdateInvoicePaymentFailed updates the status of an invoice to 'FAILED'
func (r *InvoiceRepository) UpdateInvoicePaymentFailed(ctx context.Context, arg db.UpdateInvoicePaymentFailedParams) (db.Invoice, error) {
	invoice, err := r.Queries.UpdateInvoicePaymentFailed(ctx, arg)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("repository: UpdateInvoicePaymentFailed failed: %w", err)
	}
	return invoice, nil
}

// UpdateInvoiceStatusGeneral updates the status of an invoice for general purposes (e.g., REFUNDED)
func (r *InvoiceRepository) UpdateInvoiceStatusGeneral(ctx context.Context, arg db.UpdateInvoiceStatusGeneralParams) (db.Invoice, error) {
	invoice, err := r.Queries.UpdateInvoiceStatusGeneral(ctx, arg)
	if err != nil {
		return db.Invoice{}, fmt.Errorf("repository: UpdateInvoiceStatusGeneral failed: %w", err)
	}
	return invoice, nil
}
