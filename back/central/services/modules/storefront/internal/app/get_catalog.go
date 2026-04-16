package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
)

func (uc *UseCase) ListCatalog(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error) {
	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return nil, 0, err
	}
	if !active {
		return nil, 0, domainerrors.ErrStorefrontNotActive
	}

	filters.Normalize()
	return uc.repo.ListActiveProducts(ctx, businessID, filters)
}
