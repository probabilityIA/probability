package usecaseupdatestatus

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseUpdateStatus) toRejected(order *entities.ProbabilityOrder, req *dtos.ChangeStatusRequest) {
	order.Status = string(entities.OrderStatusRejected)

	if req.Metadata != nil {
		if reason, ok := req.Metadata["reason"].(string); ok {
			order.Notes = &reason
		}
	}
}
