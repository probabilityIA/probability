package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/storefront/internal/domain/errors"
)

func (uc *UseCase) UpdateCatalogLayout(ctx context.Context, businessID, requesterUserID uint, columns, rows int) (entities.StorefrontCatalogLayout, error) {
	active, err := uc.repo.IsIntegrationActiveOrMissing(ctx, businessID, tiendaIntegrationTypeID)
	if err != nil {
		return entities.StorefrontCatalogLayout{}, err
	}
	if !active {
		return entities.StorefrontCatalogLayout{}, domainerrors.ErrStorefrontNotActive
	}

	level, err := uc.repo.GetRoleLevelByUserAndBusiness(ctx, requesterUserID, businessID)
	if err != nil {
		return entities.StorefrontCatalogLayout{}, err
	}
	if level > clienteFinalMaxRoleLevel {
		return entities.StorefrontCatalogLayout{}, domainerrors.ErrRoleNotAllowed
	}

	if columns < 1 {
		columns = 1
	}
	if columns > 6 {
		columns = 6
	}
	if rows < 1 {
		rows = 1
	}
	if rows > 12 {
		rows = 12
	}

	layout := entities.StorefrontCatalogLayout{Columns: columns, Rows: rows}
	if err := uc.repo.UpdateCatalogLayout(ctx, businessID, tiendaIntegrationTypeID, requesterUserID, layout); err != nil {
		return entities.StorefrontCatalogLayout{}, err
	}
	return layout, nil
}
