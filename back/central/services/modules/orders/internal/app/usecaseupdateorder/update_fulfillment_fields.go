package usecaseupdateorder

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

// updateFulfillmentFields actualiza los campos relacionados con fulfillment desde Shipments
func (uc *UseCaseUpdateOrder) updateFulfillmentFields(order *entities.ProbabilityOrder, dto *dtos.ProbabilityOrderDTO) bool {
	if len(dto.Shipments) == 0 {
		return false
	}

	changed := false
	shipment := dto.Shipments[0]

	if shipment.WarehouseID != nil && (order.WarehouseID == nil || *order.WarehouseID != *shipment.WarehouseID) {
		order.WarehouseID = shipment.WarehouseID
		changed = true
	}

	if shipment.WarehouseName != "" && order.WarehouseName != shipment.WarehouseName {
		order.WarehouseName = shipment.WarehouseName
		changed = true
	}

	if shipment.DriverID != nil && (order.DriverID == nil || *order.DriverID != *shipment.DriverID) {
		order.DriverID = shipment.DriverID
		changed = true
	}

	if shipment.DriverName != "" && order.DriverName != shipment.DriverName {
		order.DriverName = shipment.DriverName
		changed = true
	}

	if order.IsLastMile != shipment.IsLastMile {
		order.IsLastMile = shipment.IsLastMile
		changed = true
	}

	return changed
}
