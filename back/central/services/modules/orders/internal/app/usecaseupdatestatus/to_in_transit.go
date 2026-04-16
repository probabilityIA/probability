package usecaseupdatestatus

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseUpdateStatus) toInTransit(order *entities.ProbabilityOrder, req *dtos.ChangeStatusRequest) {
	order.Status = string(entities.OrderStatusInTransit)

	if req.Metadata != nil {
		if trackingNumber, ok := req.Metadata["tracking_number"].(string); ok {
			order.TrackingNumber = &trackingNumber
		}
		if trackingLink, ok := req.Metadata["tracking_link"].(string); ok {
			order.TrackingLink = &trackingLink
		}
	}
}
