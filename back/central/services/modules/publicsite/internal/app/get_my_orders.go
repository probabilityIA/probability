package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (uc *UseCase) GetMyOrders(ctx context.Context, slug string, userID uint, page, pageSize int) ([]entities.TiendaOrder, int64, error) {
	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return nil, 0, err
	}
	if business == nil {
		return nil, 0, domainerrors.ErrBusinessNotFound
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	var customerID uint
	session, err := uc.repo.GetClientSession(ctx, business.ID, userID)
	if err == nil && session != nil {
		customerID = session.CustomerID
	}

	return uc.repo.ListCustomerOrders(ctx, business.ID, userID, customerID, page, pageSize)
}
