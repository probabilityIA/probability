package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
)

func (uc *UseCase) List(ctx context.Context, params dtos.ListTicketsParams) ([]entities.Ticket, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	return uc.repo.List(ctx, params)
}
