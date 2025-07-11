package service

import (
	"context"
	"shipment-service/internal/db"
	"shipment-service/internal/model"
)

// ShipmentService defines the business logic for managing shipments and invoices.
type ShipmentService interface {
	CreateShipment(ctx context.Context, req model.CreateShipmentRequest) (db.Shipment, error)
	GetShipmentByID(ctx context.Context, shipmentID int32) (db.Shipment, error)
	ListShipments(ctx context.Context, params db.ListShipmentsParams) ([]db.Shipment, error)
	GetShipmentsByTripID(ctx context.Context, tripID int32) ([]db.Shipment, error)

	GetInvoiceByID(ctx context.Context, invoiceID int32) (db.Invoice, error)
	GetInvoiceForShipment(ctx context.Context, shipmentID int32) (db.Invoice, error)
	ListInvoices(ctx context.Context, params db.ListInvoicesParams) ([]db.Invoice, error)
	ListInvoicesByTripID(ctx context.Context, tripID int32) ([]db.Invoice, error)
}
