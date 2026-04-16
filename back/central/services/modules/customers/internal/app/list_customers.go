package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
)

func (uc *UseCase) ListClients(ctx context.Context, params dtos.ListClientsParams) ([]entities.Client, int64, error) {
	// Validar paginación
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	return uc.repo.List(ctx, params)
}
