package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

func (uc *UseCase) ListCustomerAddresses(ctx context.Context, params dtos.ListCustomerAddressesParams) ([]entities.CustomerAddress, int64, error) {
	return uc.repo.ListCustomerAddresses(ctx, params)
}
