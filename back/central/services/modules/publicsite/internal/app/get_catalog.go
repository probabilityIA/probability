package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (uc *UseCase) ListCatalog(ctx context.Context, slug string, filters dtos.CatalogFilters) ([]entities.PublicProduct, int64, error) {
	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return nil, 0, err
	}
	if business == nil {
		return nil, 0, domainerrors.ErrBusinessNotFound
	}

	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, business.ID, tiendaWebIntegrationTypeID)
	if err != nil {
		return nil, 0, err
	}
	if !active {
		return nil, 0, domainerrors.ErrPublicSiteNotActive
	}

	filters.Normalize()
	return uc.repo.ListActiveProducts(ctx, business.ID, filters)
}
