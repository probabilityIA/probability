package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (uc *UseCase) GetProduct(ctx context.Context, slug string, productID string) (*entities.PublicProduct, error) {
	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, domainerrors.ErrBusinessNotFound
	}

	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, business.ID, tiendaWebIntegrationTypeID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, domainerrors.ErrPublicSiteNotActive
	}

	return uc.repo.GetProductByID(ctx, business.ID, productID)
}

func (uc *UseCase) GetFeaturedProducts(ctx context.Context, slug string, limit int) ([]entities.PublicProduct, error) {
	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, domainerrors.ErrBusinessNotFound
	}

	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, business.ID, tiendaWebIntegrationTypeID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, domainerrors.ErrPublicSiteNotActive
	}

	if limit < 1 || limit > 20 {
		limit = 8
	}
	return uc.repo.GetFeaturedProducts(ctx, business.ID, limit)
}
