package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (uc *UseCase) GetMyPrices(ctx context.Context, slug string, userID uint) (map[string]float64, error) {
	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, domainerrors.ErrBusinessNotFound
	}

	session, err := uc.repo.GetClientSession(ctx, business.ID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil || session.CustomerID == 0 {
		return map[string]float64{}, nil
	}

	return uc.repo.GetCustomerPrices(ctx, business.ID, session.CustomerID)
}
