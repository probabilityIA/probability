package app

import (
	"context"
	"fmt"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
)

func (uc *UseCase) DeleteClient(ctx context.Context, businessID, clientID uint) error {
	_, err := uc.repo.GetByID(ctx, businessID, clientID)
	if err != nil {
		return err
	}

	summary, err := uc.repo.GetCustomerSummary(ctx, businessID, clientID)
	if err != nil {
		return err
	}
	if summary != nil && summary.TotalOrders > 0 {
		return fmt.Errorf("%w: tiene %d orden(es)", domainerrors.ErrClientHasOrders, summary.TotalOrders)
	}

	return uc.repo.Delete(ctx, businessID, clientID)
}
