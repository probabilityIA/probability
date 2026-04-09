package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

func (uc *UseCase) ListCustomerOrderItems(ctx context.Context, params dtos.ListCustomerOrderItemsParams) ([]entities.CustomerOrderItem, int64, error) {
	return uc.repo.ListCustomerOrderItems(ctx, params)
}
