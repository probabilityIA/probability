package app

import (
	"context"

	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (uc *UseCase) GetCategories(ctx context.Context, slug string) ([]string, error) {
	business, err := uc.repo.GetBusinessBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if business == nil {
		return nil, domainerrors.ErrBusinessNotFound
	}

	active, err := uc.repo.IsIntegrationActive(ctx, business.ID, tiendaWebIntegrationTypeID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, domainerrors.ErrPublicSiteNotActive
	}

	hidden, _ := uc.repo.GetHiddenCategories(ctx, business.ID)
	return uc.repo.ListPublicCategories(ctx, business.ID, hidden)
}
