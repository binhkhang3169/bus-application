package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"shipment-service/internal/db"
	"shipment-service/internal/model"
	"shipment-service/internal/repository"
	"time"
)

type shipmentService struct {
	shipmentRepo repository.ShipmentRepository
	invoiceRepo  repository.InvoiceRepository
	db           *sql.DB
}

func NewShipmentService(
	shipmentRepo repository.ShipmentRepository,
	invoiceRepo repository.InvoiceRepository,
	dbPool *sql.DB,
) ShipmentService {
	return &shipmentService{
		shipmentRepo: shipmentRepo,
		invoiceRepo:  invoiceRepo,
		db:           dbPool,
	}
}

func (s *shipmentService) executeInTransaction(ctx context.Context, fn func(queries db.Querier) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("transaction rollback failed: %v", rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()
	qtx := db.New(tx)
	err = fn(qtx)
	return err
}

func (s *shipmentService) CreateShipment(ctx context.Context, req model.CreateShipmentRequest) (db.Shipment, error) {
	var createdShipment db.Shipment

	if req.PayerType != "sender" && req.PayerType != "receiver" {
		return db.Shipment{}, fmt.Errorf("invalid payer_type: must be 'sender' or 'receiver'")
	}

	err := s.executeInTransaction(ctx, func(qtx db.Querier) error {
		physicalVolumeCm3 := req.Dimensions.Length * req.Dimensions.Width * req.Dimensions.Height

		shipmentParams := db.CreateShipmentParams{
			TripID:       req.TripID,
			SenderName:   req.SenderName,
			ReceiverName: req.ReceiverName,
			ItemName:     req.ItemName,
			ItemType:     req.ItemType,
			Weight:       req.Weight,
			Length:       req.Dimensions.Length,
			Width:        req.Dimensions.Width,
			Height:       req.Dimensions.Height,
			Volume:       physicalVolumeCm3,
			Price:        req.Price,     // Price from frontend
			PayerType:    req.PayerType, // Payer from frontend
			Note:         req.Note,
		}

		var txErr error
		createdShipment, txErr = s.shipmentRepo.CreateShipment(ctx, qtx, shipmentParams)
		if txErr != nil {
			return fmt.Errorf("repository: failed to create shipment: %w", txErr)
		}

		invoiceParams := db.CreateInvoiceParams{
			ShipmentID: createdShipment.ID,
			Amount:     createdShipment.Price,
			IssuedAt:   time.Now(),
		}

		_, txErr = s.invoiceRepo.CreateInvoice(ctx, qtx, invoiceParams)
		if txErr != nil {
			return fmt.Errorf("repository: failed to create invoice: %w", txErr)
		}
		return nil
	})

	if err != nil {
		return db.Shipment{}, err
	}

	return createdShipment, nil
}

func (s *shipmentService) GetShipmentByID(ctx context.Context, shipmentID int32) (db.Shipment, error) {
	q := db.New(s.db)
	shipment, err := s.shipmentRepo.GetShipmentByID(ctx, q, shipmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Shipment{}, fmt.Errorf("shipment with ID %d not found", shipmentID)
		}
		return db.Shipment{}, fmt.Errorf("repository: failed to get shipment: %w", err)
	}
	return shipment, nil
}

func (s *shipmentService) ListShipments(ctx context.Context, params db.ListShipmentsParams) ([]db.Shipment, error) {
	q := db.New(s.db)
	shipments, err := s.shipmentRepo.ListShipments(ctx, q, params)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list shipments: %w", err)
	}
	return shipments, nil
}

func (s *shipmentService) GetShipmentsByTripID(ctx context.Context, tripID int32) ([]db.Shipment, error) {
	q := db.New(s.db)
	shipments, err := s.shipmentRepo.ListShipmentsByTripID(ctx, q, tripID)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list shipments by trip ID %d: %w", tripID, err)
	}
	return shipments, nil
}

func (s *shipmentService) GetInvoiceByID(ctx context.Context, invoiceID int32) (db.Invoice, error) {
	q := db.New(s.db)
	invoice, err := s.invoiceRepo.GetInvoiceByID(ctx, q, invoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Invoice{}, fmt.Errorf("invoice with ID %d not found", invoiceID)
		}
		return db.Invoice{}, fmt.Errorf("repository: failed to get invoice: %w", err)
	}
	return invoice, nil
}

func (s *shipmentService) GetInvoiceForShipment(ctx context.Context, shipmentID int32) (db.Invoice, error) {
	q := db.New(s.db)
	// First, check if shipment exists to give a clear error message
	_, err := s.shipmentRepo.GetShipmentByID(ctx, q, shipmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Invoice{}, fmt.Errorf("shipment with ID %d not found", shipmentID)
		}
		return db.Invoice{}, fmt.Errorf("repository: failed to verify shipment %d: %w", shipmentID, err)
	}

	invoice, err := s.invoiceRepo.GetInvoiceByShipmentID(ctx, q, shipmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return db.Invoice{}, fmt.Errorf("invoice for shipment ID %d not found", shipmentID)
		}
		return db.Invoice{}, fmt.Errorf("repository: failed to get invoice by shipment ID: %w", err)
	}
	return invoice, nil
}

func (s *shipmentService) ListInvoices(ctx context.Context, params db.ListInvoicesParams) ([]db.Invoice, error) {
	q := db.New(s.db)
	invoices, err := s.invoiceRepo.ListInvoices(ctx, q, params)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list invoices: %w", err)
	}
	return invoices, nil
}

func (s *shipmentService) ListInvoicesByTripID(ctx context.Context, tripID int32) ([]db.Invoice, error) {
	q := db.New(s.db)
	invoices, err := s.invoiceRepo.ListInvoicesByTripID(ctx, q, tripID)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to list invoices by trip ID %d: %w", tripID, err)
	}
	return invoices, nil
}
