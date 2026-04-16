package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

func (uc *UseCase) ListCustomerProducts(ctx context.Context, params dtos.ListCustomerProductsParams) ([]entities.CustomerProductHistory, int64, error) {
	return uc.repo.ListCustomerProducts(ctx, params)
}
