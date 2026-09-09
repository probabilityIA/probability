package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
)

func (uc *UseCase) ListClients(ctx context.Context, businessID uint, page, pageSize int) ([]entities.StorefrontClient, int64, error) {
	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return nil, 0, err
	}
	if !active {
		return nil, 0, domainerrors.ErrStorefrontNotActive
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	return uc.repo.ListClientsByBusiness(ctx, businessID, page, pageSize)
}
