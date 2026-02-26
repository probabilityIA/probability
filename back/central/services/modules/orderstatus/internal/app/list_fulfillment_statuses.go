package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

// ListFulfillmentStatuses obtiene el catálogo de estados de fulfillment
func (uc *useCase) ListFulfillmentStatuses(ctx context.Context, isActive *bool) ([]entities.FulfillmentStatusInfo, error) {
	return uc.repo.ListFulfillmentStatuses(ctx, isActive)
}
