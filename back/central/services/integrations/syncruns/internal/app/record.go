package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
)

func (uc *useCase) Record(ctx context.Context, run *domain.SyncRun) error {
	if run == nil || run.BusinessID == 0 || run.IntegrationID == 0 {
		return domain.ErrIntegrationNotOwned
	}

	owned, err := uc.repo.IntegrationBelongsToBusiness(ctx, run.IntegrationID, run.BusinessID)
	if err != nil {
		return err
	}
	if !owned {
		return domain.ErrIntegrationNotOwned
	}

	run.Normalize()

	if err := uc.repo.Save(ctx, run); err != nil {
		uc.logger.Error(ctx).Err(err).
			Uint("integration_id", run.IntegrationID).
			Str("kind", run.Kind).
			Msg("No se pudo guardar el resultado de la sincronizacion")
		return err
	}

	return nil
}

func (uc *useCase) ListLast(ctx context.Context, businessID uint) ([]domain.SyncRun, error) {
	return uc.repo.ListLastByBusiness(ctx, businessID)
}

func (uc *useCase) ListDetail(ctx context.Context, query domain.DetailQuery) (*domain.DetailPage, error) {
	if query.BusinessID == 0 || query.IntegrationID == 0 {
		return nil, domain.ErrIntegrationNotOwned
	}

	owned, err := uc.repo.IntegrationBelongsToBusiness(ctx, query.IntegrationID, query.BusinessID)
	if err != nil {
		return nil, err
	}
	if !owned {
		return nil, domain.ErrIntegrationNotOwned
	}

	query.Normalize()
	return uc.repo.ListDetail(ctx, query)
}
