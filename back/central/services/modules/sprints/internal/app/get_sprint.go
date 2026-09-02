package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/sprints/internal/domain/entities"
)

func (uc *UseCase) Get(ctx context.Context, id uint) (*entities.Sprint, error) {
	return uc.repo.GetByID(ctx, id)
}
