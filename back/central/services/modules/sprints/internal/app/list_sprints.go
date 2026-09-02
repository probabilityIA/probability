package app

import (
	"context"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/errors"
)

func (uc *UseCase) List(ctx context.Context, params dtos.ListSprintsParams) ([]entities.Sprint, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}
	params.Status = strings.ToLower(strings.TrimSpace(params.Status))
	if params.Status != "" && !validSprintStatuses[params.Status] {
		return nil, 0, dom.ErrInvalidStatus
	}
	return uc.repo.List(ctx, params)
}
