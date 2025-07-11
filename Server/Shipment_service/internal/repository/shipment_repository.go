package repository

import (
	"context"
	"shipment-service/internal/db"
)

type shipmentRepository struct{}

func NewShipmentRepository() ShipmentRepository {
	return &shipmentRepository{}
}

func (r *shipmentRepository) CreateShipment(ctx context.Context, querier db.Querier, arg db.CreateShipmentParams) (db.Shipment, error) {
	return querier.CreateShipment(ctx, arg)
}

func (r *shipmentRepository) GetShipmentByID(ctx context.Context, querier db.Querier, id int32) (db.Shipment, error) {
	return querier.GetShipmentByID(ctx, id)
}

func (r *shipmentRepository) ListShipments(ctx context.Context, querier db.Querier, arg db.ListShipmentsParams) ([]db.Shipment, error) {
	return querier.ListShipments(ctx, arg)
}

func (r *shipmentRepository) ListShipmentsByTripID(ctx context.Context, querier db.Querier, tripID int32) ([]db.Shipment, error) {
	return querier.ListShipmentsByTripID(ctx, tripID)
}
