package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/routes/internal/domain/entities"
)

func (uc *UseCase) ListRoutes(ctx context.Context, params dtos.ListRoutesParams) ([]entities.Route, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}
	return uc.repo.ListRoutes(ctx, params)
}
