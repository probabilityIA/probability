package usecaseupdatestatus

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseUpdateStatus) toAssignedToDriver(order *entities.ProbabilityOrder, req *dtos.ChangeStatusRequest) {
	order.Status = string(entities.OrderStatusAssignedToDriver)

	if req.Metadata != nil {
		if driverID, ok := req.Metadata["driver_id"].(float64); ok {
			id := uint(driverID)
			order.DriverID = &id
		}
		if driverName, ok := req.Metadata["driver_name"].(string); ok {
			order.DriverName = driverName
		}
	}
}
