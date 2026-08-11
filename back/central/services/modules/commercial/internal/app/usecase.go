package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/commercial/internal/domain/entities"
)

func (uc *UseCase) ListProspects(ctx context.Context, filters dtos.ListProspectsFilters) ([]entities.Prospect, dtos.ListResult, error) {
	filters.Page, filters.PageSize = dtos.NormalizePagination(filters.Page, filters.PageSize)

	prospects, total, err := uc.repo.ListProspects(ctx, filters)
	if err != nil {
		return nil, dtos.ListResult{}, err
	}

	return prospects, dtos.BuildListResult(total, filters.Page, filters.PageSize), nil
}

func (uc *UseCase) GetStats(ctx context.Context) (*entities.ProspectStats, error) {
	return uc.repo.GetStats(ctx)
}

func (uc *UseCase) UpdateSeen(ctx context.Context, id uint, seen bool) error {
	return uc.repo.UpdateSeen(ctx, id, seen)
}
