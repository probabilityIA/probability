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
		OrderID:            s.OrderID,
		ClientName:         s.ClientName,
		DestinationAddress: s.DestinationAddress,
		TrackingNumber:     s.TrackingNumber,
		TrackingURL:        s.TrackingURL,
		Carrier:            s.Carrier,
		CarrierCode:        s.CarrierCode,
		GuideID:            s.GuideID,
		GuideURL:           s.GuideURL,
		Status:             s.Status,
		ShippedAt:          s.ShippedAt,
		DeliveredAt:        s.DeliveredAt,
		ShippingAddressID:  s.ShippingAddressID,
		ShippingCost:       s.ShippingCost,
		InsuranceCost:      s.InsuranceCost,
		TotalCost:          s.TotalCost,
		Weight:             s.Weight,
		Height:             s.Height,
		Width:              s.Width,
		Length:             s.Length,
		WarehouseID:        s.WarehouseID,
		WarehouseName:      s.WarehouseName,
		DriverID:           s.DriverID,
		DriverName:         s.DriverName,
		IsLastMile:         s.IsLastMile,
		IsTest:             s.IsTest,
		EstimatedDelivery:  s.EstimatedDelivery,
		DeliveryNotes:      s.DeliveryNotes,
		Metadata:           s.Metadata,
	}
	if s.DeletedAt != nil {
		dbShipment.DeletedAt = gorm.DeletedAt{Time: *s.DeletedAt, Valid: true}
	}
	return dbShipment
}

// ToDBShipmentWithoutCustomerData es para guardar, sin incluir datos de cliente
func ToDBShipmentWithoutCustomerData(s *domain.Shipment) *models.Shipment {
	return ToDBShipment(s)
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

	shipment := &domain.Shipment{
		ID:                 s.ID,
		CreatedAt:          s.CreatedAt,
		UpdatedAt:          s.UpdatedAt,
		DeletedAt:          deletedAt,
		OrderID:            s.OrderID,
		ClientName:         s.ClientName,
		DestinationAddress: s.DestinationAddress,
		TrackingNumber:     s.TrackingNumber,
		TrackingURL:        s.TrackingURL,
		Carrier:            s.Carrier,
		CarrierCode:        s.CarrierCode,
		GuideID:            s.GuideID,
		GuideURL:           s.GuideURL,
		Status:             s.Status,
		ShippedAt:          s.ShippedAt,
		DeliveredAt:        s.DeliveredAt,
		ShippingAddressID:  s.ShippingAddressID,
		ShippingCost:       s.ShippingCost,
		InsuranceCost:      s.InsuranceCost,
		TotalCost:          s.TotalCost,
		Weight:             s.Weight,
		Height:             s.Height,
		Width:              s.Width,
		Length:             s.Length,
		WarehouseID:        s.WarehouseID,
		WarehouseName:      s.WarehouseName,
		DriverID:           s.DriverID,
		DriverName:         s.DriverName,
		IsLastMile:         s.IsLastMile,
		IsTest:             s.IsTest,
		EstimatedDelivery:  s.EstimatedDelivery,
		DeliveryNotes:      s.DeliveryNotes,
		Metadata:           s.Metadata,
	}

	// Incluir datos del cliente desde la orden si existe
	if s.Order != nil {
		shipment.CustomerName = s.Order.CustomerName
		shipment.CustomerEmail = s.Order.CustomerEmail
		shipment.CustomerPhone = s.Order.CustomerPhone
		shipment.CustomerDNI = s.Order.CustomerDNI
		shipment.OrderNumber = s.Order.OrderNumber
	}

	return shipment
}
