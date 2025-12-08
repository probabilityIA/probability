package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ToDBShipment convierte un envío de dominio a modelo de base de datos
func ToDBShipment(s *domain.Shipment) *models.Shipment {
	if s == nil {
		return nil
	}
	dbShipment := &models.Shipment{
		Model: gorm.Model{
			ID:        s.ID,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		},
		OrderID:         s.OrderID,
		TrackingNumber:  s.TrackingNumber,
		TrackingURL:     s.TrackingURL,
		Carrier:          s.Carrier,
		CarrierCode:      s.CarrierCode,
		GuideID:          s.GuideID,
		GuideURL:         s.GuideURL,
		Status:           s.Status,
		ShippedAt:        s.ShippedAt,
		DeliveredAt:      s.DeliveredAt,
		ShippingAddressID: s.ShippingAddressID,
		ShippingCost:     s.ShippingCost,
		InsuranceCost:    s.InsuranceCost,
		TotalCost:        s.TotalCost,
		Weight:           s.Weight,
		Height:           s.Height,
		Width:            s.Width,
		Length:           s.Length,
		WarehouseID:      s.WarehouseID,
		WarehouseName:    s.WarehouseName,
		DriverID:         s.DriverID,
		DriverName:       s.DriverName,
		IsLastMile:       s.IsLastMile,
		EstimatedDelivery: s.EstimatedDelivery,
		DeliveryNotes:     s.DeliveryNotes,
		Metadata:          s.Metadata,
	}
	if s.DeletedAt != nil {
		dbShipment.DeletedAt = gorm.DeletedAt{Time: *s.DeletedAt, Valid: true}
	}
	return dbShipment
}

// ToDomainShipment convierte un envío de base de datos a dominio
func ToDomainShipment(s *models.Shipment) *domain.Shipment {
	if s == nil {
		return nil
	}
	var deletedAt *time.Time
	if s.DeletedAt.Valid {
		deletedAt = &s.DeletedAt.Time
	}
	return &domain.Shipment{
		ID:         s.ID,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
		DeletedAt:  deletedAt,
		OrderID:    s.OrderID,
		TrackingNumber: s.TrackingNumber,
		TrackingURL:    s.TrackingURL,
		Carrier:        s.Carrier,
		CarrierCode:    s.CarrierCode,
		GuideID:        s.GuideID,
		GuideURL:       s.GuideURL,
		Status:         s.Status,
		ShippedAt:      s.ShippedAt,
		DeliveredAt:    s.DeliveredAt,
		ShippingAddressID: s.ShippingAddressID,
		ShippingCost:   s.ShippingCost,
		InsuranceCost:  s.InsuranceCost,
		TotalCost:      s.TotalCost,
		Weight:         s.Weight,
		Height:         s.Height,
		Width:          s.Width,
		Length:         s.Length,
		WarehouseID:    s.WarehouseID,
		WarehouseName:  s.WarehouseName,
		DriverID:       s.DriverID,
		DriverName:     s.DriverName,
		IsLastMile:    s.IsLastMile,
		EstimatedDelivery: s.EstimatedDelivery,
		DeliveryNotes:     s.DeliveryNotes,
		Metadata:          s.Metadata,
	}
}

