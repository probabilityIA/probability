package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

func (uc *UseCase) GetCustomerSummary(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error) {
	return uc.repo.GetCustomerSummary(ctx, businessID, customerID)
}
