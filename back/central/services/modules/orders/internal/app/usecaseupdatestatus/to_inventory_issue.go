package usecaseupdatestatus

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseUpdateStatus) toInventoryIssue(order *entities.ProbabilityOrder, req *dtos.ChangeStatusRequest) {
	order.Status = string(entities.OrderStatusInventoryIssue)

	if req.Metadata != nil {
		if notes, ok := req.Metadata["notes"].(string); ok {
			order.Novelty = &notes
		}
	}
}
