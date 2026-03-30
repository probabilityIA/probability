package app

import (
	"context"

	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
)

func (uc *useCase) GetContainerStats(ctx context.Context, id string) (*entities.ContainerStats, error) {
	return uc.docker.GetContainerStats(ctx, id)
}
