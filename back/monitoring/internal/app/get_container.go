package app

import (
	"context"

	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
)

func (uc *useCase) GetContainer(ctx context.Context, id string) (*entities.Container, error) {
	return uc.docker.GetContainer(ctx, id)
}
