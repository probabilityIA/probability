package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
)

func (uc *UseCase) ListCatalog(ctx context.Context, businessID uint, filters dtos.CatalogFilters) ([]entities.StorefrontProduct, int64, error) {
	filters.Normalize()
	return uc.repo.ListActiveProducts(ctx, businessID, filters)
}
