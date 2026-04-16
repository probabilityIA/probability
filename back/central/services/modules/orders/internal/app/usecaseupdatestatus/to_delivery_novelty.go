package usecaseupdatestatus

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseUpdateStatus) toDeliveryNovelty(order *entities.ProbabilityOrder, req *dtos.ChangeStatusRequest) {
	order.Status = string(entities.OrderStatusDeliveryNovelty)

	if req.Metadata != nil {
		if reason, ok := req.Metadata["reason"].(string); ok {
			order.Novelty = &reason
		}
	}
}
