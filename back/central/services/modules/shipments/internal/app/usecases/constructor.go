package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecaseshipment"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// UseCases contiene todos los casos de uso del módulo shipments
type UseCases struct {
	repo domain.IRepository

	// Casos de uso modulares
	ShipmentCRUD *usecaseshipment.UseCaseShipment
}

// New crea una nueva instancia de UseCases
func New(repo domain.IRepository) *UseCases {
	return &UseCases{
		repo:         repo,
		ShipmentCRUD: usecaseshipment.New(repo),
	}
}

// ───────────────────────────────────────────
// MÉTODOS DE COMPATIBILIDAD - Delegar al CRUD
// ───────────────────────────────────────────

// CreateShipment delega al caso de uso CRUD
func (uc *UseCases) CreateShipment(ctx context.Context, req *domain.CreateShipmentRequest) (*domain.ShipmentResponse, error) {
	return uc.ShipmentCRUD.CreateShipment(ctx, req)
}

// GetShipmentByID delega al caso de uso CRUD
func (uc *UseCases) GetShipmentByID(ctx context.Context, id uint) (*domain.ShipmentResponse, error) {
	return uc.ShipmentCRUD.GetShipmentByID(ctx, id)
}

// ListShipments delega al caso de uso CRUD
func (uc *UseCases) ListShipments(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*domain.ShipmentsListResponse, error) {
	return uc.ShipmentCRUD.ListShipments(ctx, page, pageSize, filters)
}

// UpdateShipment delega al caso de uso CRUD
func (uc *UseCases) UpdateShipment(ctx context.Context, id uint, req *domain.UpdateShipmentRequest) (*domain.ShipmentResponse, error) {
	return uc.ShipmentCRUD.UpdateShipment(ctx, id, req)
}

// DeleteShipment delega al caso de uso CRUD
func (uc *UseCases) DeleteShipment(ctx context.Context, id uint) error {
	return uc.ShipmentCRUD.DeleteShipment(ctx, id)
}

// GetShipmentsByOrderID delega al caso de uso CRUD
func (uc *UseCases) GetShipmentsByOrderID(ctx context.Context, orderID string) ([]domain.Shipment, error) {
	return uc.ShipmentCRUD.GetShipmentsByOrderID(ctx, orderID)
}

// GetShipmentByTrackingNumber delega al caso de uso CRUD
func (uc *UseCases) GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.ShipmentResponse, error) {
	shipment, err := uc.ShipmentCRUD.GetShipmentByTrackingNumber(ctx, trackingNumber)
	if err != nil {
		return nil, err
	}
	return mapShipmentToResponse(shipment), nil
}

// mapShipmentToResponse convierte un modelo Shipment a ShipmentResponse
func mapShipmentToResponse(shipment *domain.Shipment) *domain.ShipmentResponse {
	return &domain.ShipmentResponse{
		ID:         shipment.ID,
		CreatedAt:  shipment.CreatedAt,
		UpdatedAt:  shipment.UpdatedAt,
		DeletedAt:  shipment.DeletedAt,
		OrderID:    shipment.OrderID,
		TrackingNumber: shipment.TrackingNumber,
		TrackingURL:    shipment.TrackingURL,
		Carrier:        shipment.Carrier,
		CarrierCode:    shipment.CarrierCode,
		GuideID:        shipment.GuideID,
		GuideURL:       shipment.GuideURL,
		Status:         shipment.Status,
		ShippedAt:      shipment.ShippedAt,
		DeliveredAt:    shipment.DeliveredAt,
		ShippingAddressID: shipment.ShippingAddressID,
		ShippingCost:   shipment.ShippingCost,
		InsuranceCost:  shipment.InsuranceCost,
		TotalCost:      shipment.TotalCost,
		Weight:         shipment.Weight,
		Height:         shipment.Height,
		Width:          shipment.Width,
		Length:         shipment.Length,
		WarehouseID:    shipment.WarehouseID,
		WarehouseName:  shipment.WarehouseName,
		DriverID:       shipment.DriverID,
		DriverName:     shipment.DriverName,
		IsLastMile:    shipment.IsLastMile,
		EstimatedDelivery: shipment.EstimatedDelivery,
		DeliveryNotes:     shipment.DeliveryNotes,
		Metadata:          shipment.Metadata,
	}
}

