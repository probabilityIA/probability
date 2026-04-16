package usecaseupdatestatus

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

func (uc *UseCaseUpdateStatus) toDelivered(order *entities.ProbabilityOrder, _ *dtos.ChangeStatusRequest) {
	order.Status = string(entities.OrderStatusDelivered)

	now := time.Now()
	order.DeliveredAt = &now
}
