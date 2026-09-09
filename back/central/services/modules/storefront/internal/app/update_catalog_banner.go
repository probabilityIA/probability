package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
)

func (uc *UseCase) UpdateCatalogBanner(ctx context.Context, businessID, requesterUserID uint, enabled *bool, imageURL *string) (entities.StorefrontBanner, error) {
	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return entities.StorefrontBanner{}, err
	}
	if !active {
		return entities.StorefrontBanner{}, domainerrors.ErrStorefrontNotActive
	}

	level, err := uc.repo.GetRoleLevelByUserAndBusiness(ctx, requesterUserID, businessID)
	if err != nil {
		return entities.StorefrontBanner{}, err
	}
	if level > clienteFinalMaxRoleLevel {
		return entities.StorefrontBanner{}, domainerrors.ErrRoleNotAllowed
	}

	current, err := uc.repo.GetCatalogBanner(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return entities.StorefrontBanner{}, err
	}

	if enabled != nil {
		current.Enabled = *enabled
	}
	if imageURL != nil {
		current.ImageURL = *imageURL
	}

	if err := uc.repo.UpdateCatalogBanner(ctx, businessID, tiendaIntegrationTypeID, requesterUserID, current); err != nil {
		return entities.StorefrontBanner{}, err
	}
	return current, nil
}
