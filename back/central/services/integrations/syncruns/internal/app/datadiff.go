package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

func (uc *useCase) DataSummary(ctx context.Context, businessID uint) (*domain.DataSummary, error) {
	if businessID == 0 {
		return nil, domain.ErrIntegrationNotOwned
	}
	return uc.repo.DataSummary(ctx, businessID)
}

func (uc *useCase) DataPreview(ctx context.Context, query domain.DataPreviewQuery) (*domain.DataPreview, error) {
	if query.BusinessID == 0 || query.IntegrationID == 0 {
		return nil, domain.ErrIntegrationNotOwned
	}
	if !domain.IsDataField(query.Field) {
		return nil, domain.ErrInvalidDataField
	}
	if query.Mode != domain.ModeOverwrite {
		query.Mode = domain.ModeFillEmpty
	}

	owned, err := uc.repo.IntegrationBelongsToBusiness(ctx, query.IntegrationID, query.BusinessID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, domain.ErrIntegrationNotOwned
	}

	return uc.repo.DataPreview(ctx, query)
}
