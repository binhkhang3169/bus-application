package repository

import (
	"context"
	"shipment-service/internal/db"
)

type invoiceRepository struct{}

func NewInvoiceRepository() InvoiceRepository {
	return &invoiceRepository{}
}

func (r *invoiceRepository) CreateInvoice(ctx context.Context, querier db.Querier, arg db.CreateInvoiceParams) (db.Invoice, error) {
	return querier.CreateInvoice(ctx, arg)
}

func (r *invoiceRepository) GetInvoiceByShipmentID(ctx context.Context, querier db.Querier, shipmentID int32) (db.Invoice, error) {
	return querier.GetInvoiceByShipmentID(ctx, shipmentID)
}

func (r *invoiceRepository) GetInvoiceByID(ctx context.Context, querier db.Querier, id int32) (db.Invoice, error) {
	return querier.GetInvoiceByID(ctx, id)
}

func (r *invoiceRepository) ListInvoices(ctx context.Context, querier db.Querier, arg db.ListInvoicesParams) ([]db.Invoice, error) {
	return querier.ListInvoices(ctx, arg)
}

func (r *invoiceRepository) ListInvoicesByTripID(ctx context.Context, querier db.Querier, tripID int32) ([]db.Invoice, error) {
	return querier.ListInvoicesByTripID(ctx, tripID)
}
