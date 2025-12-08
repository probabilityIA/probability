package usecaseshipment

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

// ───────────────────────────────────────────
//
//	CREATE SHIPMENT
//
// ───────────────────────────────────────────

// CreateShipment crea un nuevo envío
func (uc *UseCaseShipment) CreateShipment(ctx context.Context, req *domain.CreateShipmentRequest) (*domain.ShipmentResponse, error) {
	// Validar order_id
	if req.OrderID == "" {
		return nil, domain.ErrOrderIDRequired
	}

	// Si hay tracking number, validar que no exista otro envío con el mismo tracking para la misma orden
	if req.TrackingNumber != nil && *req.TrackingNumber != "" {
		exists, err := uc.repo.ShipmentExists(ctx, req.OrderID, *req.TrackingNumber)
		if err != nil {
			return nil, fmt.Errorf("error checking if shipment exists: %w", err)
		}
		if exists {
			return nil, domain.ErrShipmentAlreadyExists
		}
	}

	// Establecer status por defecto si no se proporciona
	status := req.Status
	if status == "" {
		status = "pending"
	}

	// Crear el modelo de envío
	shipment := &domain.Shipment{
		OrderID: req.OrderID,

		TrackingNumber: req.TrackingNumber,
		TrackingURL:    req.TrackingURL,
		Carrier:        req.Carrier,
		CarrierCode:    req.CarrierCode,

		GuideID:  req.GuideID,
		GuideURL: req.GuideURL,

		Status:      status,
		ShippedAt:   req.ShippedAt,
		DeliveredAt: req.DeliveredAt,

		ShippingAddressID: req.ShippingAddressID,

		ShippingCost:  req.ShippingCost,
		InsuranceCost: req.InsuranceCost,
		TotalCost:     req.TotalCost,

		Weight: req.Weight,
		Height: req.Height,
		Width:  req.Width,
		Length: req.Length,

		WarehouseID:   req.WarehouseID,
		WarehouseName: req.WarehouseName,
		DriverID:      req.DriverID,
		DriverName:    req.DriverName,
		IsLastMile:    req.IsLastMile,

		EstimatedDelivery: req.EstimatedDelivery,
		DeliveryNotes:     req.DeliveryNotes,
		Metadata:          req.Metadata,
	}

	// Guardar en la base de datos
	if err := uc.repo.CreateShipment(ctx, shipment); err != nil {
		return nil, fmt.Errorf("error creating shipment: %w", err)
	}

	// Retornar la respuesta
	return mapShipmentToResponse(shipment), nil
}

// ───────────────────────────────────────────
//
//	GET SHIPMENT BY ID
//
// ───────────────────────────────────────────

// GetShipmentByID obtiene un envío por su ID
func (uc *UseCaseShipment) GetShipmentByID(ctx context.Context, id uint) (*domain.ShipmentResponse, error) {
	if id == 0 {
		return nil, errors.New("shipment ID is required")
	}

	shipment, err := uc.repo.GetShipmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting shipment: %w", err)
	}

	if shipment == nil {
		return nil, domain.ErrShipmentNotFound
	}

	return mapShipmentToResponse(shipment), nil
}

// ───────────────────────────────────────────
//
//	LIST SHIPMENTS
//
// ───────────────────────────────────────────

// ListShipments obtiene una lista paginada de envíos con filtros
func (uc *UseCaseShipment) ListShipments(ctx context.Context, page, pageSize int, filters map[string]interface{}) (*domain.ShipmentsListResponse, error) {
	// Validar paginación
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Obtener envíos del repositorio
	shipments, total, err := uc.repo.ListShipments(ctx, page, pageSize, filters)
	if err != nil {
		return nil, fmt.Errorf("error listing shipments: %w", err)
	}

	// Mapear a respuestas
	shipmentResponses := make([]domain.ShipmentResponse, len(shipments))
	for i, shipment := range shipments {
		shipmentResponses[i] = *mapShipmentToResponse(&shipment)
	}

	// Calcular total de páginas
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &domain.ShipmentsListResponse{
		Data:       shipmentResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ───────────────────────────────────────────
//
//	UPDATE SHIPMENT
//
// ───────────────────────────────────────────

// UpdateShipment actualiza un envío existente
func (uc *UseCaseShipment) UpdateShipment(ctx context.Context, id uint, req *domain.UpdateShipmentRequest) (*domain.ShipmentResponse, error) {
	if id == 0 {
		return nil, errors.New("shipment ID is required")
	}

	// Obtener el envío existente
	shipment, err := uc.repo.GetShipmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting shipment: %w", err)
	}

	if shipment == nil {
		return nil, domain.ErrShipmentNotFound
	}

	// Actualizar solo los campos proporcionados
	if req.TrackingNumber != nil {
		// Si se cambia el tracking number, verificar que no exista otro envío con ese tracking para la misma orden
		if *req.TrackingNumber != "" && (shipment.TrackingNumber == nil || *req.TrackingNumber != *shipment.TrackingNumber) {
			exists, err := uc.repo.ShipmentExists(ctx, shipment.OrderID, *req.TrackingNumber)
			if err != nil {
				return nil, fmt.Errorf("error checking if shipment exists: %w", err)
			}
			if exists {
				return nil, domain.ErrShipmentAlreadyExists
			}
		}
		shipment.TrackingNumber = req.TrackingNumber
	}
	if req.TrackingURL != nil {
		shipment.TrackingURL = req.TrackingURL
	}
	if req.Carrier != nil {
		shipment.Carrier = req.Carrier
	}
	if req.CarrierCode != nil {
		shipment.CarrierCode = req.CarrierCode
	}
	if req.GuideID != nil {
		shipment.GuideID = req.GuideID
	}
	if req.GuideURL != nil {
		shipment.GuideURL = req.GuideURL
	}
	if req.Status != nil {
		shipment.Status = *req.Status
	}
	if req.ShippedAt != nil {
		shipment.ShippedAt = req.ShippedAt
	}
	if req.DeliveredAt != nil {
		shipment.DeliveredAt = req.DeliveredAt
	}
	if req.ShippingAddressID != nil {
		shipment.ShippingAddressID = req.ShippingAddressID
	}
	if req.ShippingCost != nil {
		shipment.ShippingCost = req.ShippingCost
	}
	if req.InsuranceCost != nil {
		shipment.InsuranceCost = req.InsuranceCost
	}
	if req.TotalCost != nil {
		shipment.TotalCost = req.TotalCost
	}
	if req.Weight != nil {
		shipment.Weight = req.Weight
	}
	if req.Height != nil {
		shipment.Height = req.Height
	}
	if req.Width != nil {
		shipment.Width = req.Width
	}
	if req.Length != nil {
		shipment.Length = req.Length
	}
	if req.WarehouseID != nil {
		shipment.WarehouseID = req.WarehouseID
	}
	if req.WarehouseName != nil {
		shipment.WarehouseName = *req.WarehouseName
	}
	if req.DriverID != nil {
		shipment.DriverID = req.DriverID
	}
	if req.DriverName != nil {
		shipment.DriverName = *req.DriverName
	}
	if req.IsLastMile != nil {
		shipment.IsLastMile = *req.IsLastMile
	}
	if req.EstimatedDelivery != nil {
		shipment.EstimatedDelivery = req.EstimatedDelivery
	}
	if req.DeliveryNotes != nil {
		shipment.DeliveryNotes = req.DeliveryNotes
	}
	if req.Metadata != nil {
		shipment.Metadata = req.Metadata
	}

	// Guardar cambios
	if err := uc.repo.UpdateShipment(ctx, shipment); err != nil {
		return nil, fmt.Errorf("error updating shipment: %w", err)
	}

	return mapShipmentToResponse(shipment), nil
}

// ───────────────────────────────────────────
//
//	DELETE SHIPMENT
//
// ───────────────────────────────────────────

// DeleteShipment elimina (soft delete) un envío
func (uc *UseCaseShipment) DeleteShipment(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("shipment ID is required")
	}

	// Verificar que el envío existe
	shipment, err := uc.repo.GetShipmentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error getting shipment: %w", err)
	}

	if shipment == nil {
		return domain.ErrShipmentNotFound
	}

	// Eliminar el envío
	if err := uc.repo.DeleteShipment(ctx, id); err != nil {
		return fmt.Errorf("error deleting shipment: %w", err)
	}

	return nil
}

// ───────────────────────────────────────────
//
//	GET SHIPMENTS BY ORDER ID
//
// ───────────────────────────────────────────

// GetShipmentsByOrderID obtiene todos los envíos de una orden
func (uc *UseCaseShipment) GetShipmentsByOrderID(ctx context.Context, orderID string) ([]domain.Shipment, error) {
	if orderID == "" {
		return nil, domain.ErrOrderIDRequired
	}

	shipments, err := uc.repo.GetShipmentsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error getting shipments by order ID: %w", err)
	}

	return shipments, nil
}

// ───────────────────────────────────────────
//
//	GET SHIPMENT BY TRACKING NUMBER
//
// ───────────────────────────────────────────

// GetShipmentByTrackingNumber obtiene un envío por su número de tracking
func (uc *UseCaseShipment) GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.Shipment, error) {
	if trackingNumber == "" {
		return nil, errors.New("tracking number is required")
	}

	shipment, err := uc.repo.GetShipmentByTrackingNumber(ctx, trackingNumber)
	if err != nil {
		return nil, fmt.Errorf("error getting shipment by tracking number: %w", err)
	}

	if shipment == nil {
		return nil, domain.ErrShipmentNotFound
	}

	return shipment, nil
}

// ───────────────────────────────────────────
//
//	HELPER FUNCTIONS
//
// ───────────────────────────────────────────

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

