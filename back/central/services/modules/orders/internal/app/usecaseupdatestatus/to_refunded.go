package usecaseupdatestatus

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseUpdateStatus) toRefunded(order *entities.ProbabilityOrder, _ *dtos.ChangeStatusRequest) {
	order.Status = string(entities.OrderStatusRefunded)
}
