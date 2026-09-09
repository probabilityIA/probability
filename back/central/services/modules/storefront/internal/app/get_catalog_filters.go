package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
)

func (uc *UseCase) GetCatalogFilters(ctx context.Context, businessID uint) (entities.StorefrontFilters, error) {
	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return entities.StorefrontFilters{}, err
	}
	if !active {
		return entities.StorefrontFilters{}, domainerrors.ErrStorefrontNotActive
	}

	filters, err := uc.repo.GetCatalogFilters(ctx, businessID)
	if err != nil {
		return entities.StorefrontFilters{}, err
	}

	layout, err := uc.repo.GetCatalogLayout(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return entities.StorefrontFilters{}, err
	}
	filters.Layout = layout

	banner, err := uc.repo.GetCatalogBanner(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return entities.StorefrontFilters{}, err
	}
	filters.Banner = banner

	return filters, nil
}
