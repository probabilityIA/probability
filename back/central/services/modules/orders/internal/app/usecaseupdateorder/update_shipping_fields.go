package usecaseupdateorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// updateShippingFields actualiza los campos relacionados con el envío
func (uc *UseCaseUpdateOrder) updateShippingFields(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	changed := false

	// Actualizar información de tracking desde Shipments
	if uc.updateTrackingFields(order, dto) {
		changed = true
	}

	// Actualizar dirección de envío desde Addresses
	if uc.updateShippingAddress(order, dto) {
		changed = true
	}

	return changed
}

// updateTrackingFields actualiza los campos de tracking desde Shipments
func (uc *UseCaseUpdateOrder) updateTrackingFields(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	if len(dto.Shipments) == 0 {
		return false
	}

	changed := false
	shipment := dto.Shipments[0]

	if shipment.TrackingNumber != nil && (order.TrackingNumber == nil || *order.TrackingNumber != *shipment.TrackingNumber) {
		order.TrackingNumber = shipment.TrackingNumber
		changed = true
	}

	if shipment.TrackingURL != nil && (order.TrackingLink == nil || *order.TrackingLink != *shipment.TrackingURL) {
		order.TrackingLink = shipment.TrackingURL
		changed = true
	}

	if shipment.GuideID != nil && (order.GuideID == nil || *order.GuideID != *shipment.GuideID) {
		order.GuideID = shipment.GuideID
		changed = true
	}

	if shipment.GuideURL != nil && (order.GuideLink == nil || *order.GuideLink != *shipment.GuideURL) {
		order.GuideLink = shipment.GuideURL
		changed = true
	}

	if shipment.DeliveredAt != nil && (order.DeliveredAt == nil || !order.DeliveredAt.Equal(*shipment.DeliveredAt)) {
		order.DeliveredAt = shipment.DeliveredAt
		changed = true
	}

	if shipment.ShippedAt != nil && (order.DeliveryDate == nil || !order.DeliveryDate.Equal(*shipment.ShippedAt)) {
		order.DeliveryDate = shipment.ShippedAt
		changed = true
	}

	return changed
}

// updateShippingAddress actualiza la dirección de envío desde Addresses
func (uc *UseCaseUpdateOrder) updateShippingAddress(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	if len(dto.Addresses) == 0 {
		return false
	}

	changed := false

	for _, addr := range dto.Addresses {
		if addr.Type == "shipping" {
			if addr.Street != "" && order.ShippingStreet != addr.Street {
				order.ShippingStreet = addr.Street
				changed = true
			}

			if addr.Street2 != "" && order.Address2 != addr.Street2 {
				order.Address2 = addr.Street2
				changed = true
			}

			if addr.City != "" && order.ShippingCity != addr.City {
				order.ShippingCity = addr.City
				changed = true
			}

			if addr.State != "" && order.ShippingState != addr.State {
				order.ShippingState = addr.State
				changed = true
			}

			if addr.Country != "" && order.ShippingCountry != addr.Country {
				order.ShippingCountry = addr.Country
				changed = true
			}

			if addr.PostalCode != "" && order.ShippingPostalCode != addr.PostalCode {
				order.ShippingPostalCode = addr.PostalCode
				changed = true
			}

			if addr.Latitude != nil && (order.ShippingLat == nil || *order.ShippingLat != *addr.Latitude) {
				order.ShippingLat = addr.Latitude
				changed = true
			}

			if addr.Longitude != nil && (order.ShippingLng == nil || *order.ShippingLng != *addr.Longitude) {
				order.ShippingLng = addr.Longitude
				changed = true
			}

			break
		}
	}

	return changed
}
