package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
)

func (uc *UseCase) GetBusinessPage(ctx context.Context, slug string) (*entities.BusinessPage, error) {
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

	return business, nil
}
