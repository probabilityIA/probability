package app

import (
	"context"
	"fmt"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
)

func (uc *UseCase) DeleteClient(ctx context.Context, businessID, clientID uint) error {
	// Verificar que existe
	_, err := uc.repo.GetByID(ctx, businessID, clientID)
	if err != nil {
		return err
	}

	// Verificar si tiene órdenes
	orderCount, _, _, err := uc.repo.GetOrderStats(ctx, clientID)
	if err != nil {
		return err
	}
	if orderCount > 0 {
		return fmt.Errorf("%w: tiene %d orden(es)", domainerrors.ErrClientHasOrders, orderCount)
	}

	return uc.repo.Delete(ctx, businessID, clientID)
}
