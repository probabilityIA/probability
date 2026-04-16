package app

import (
	"context"

	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
)

func (uc *useCase) ListContainers(ctx context.Context) ([]entities.Container, error) {
	return uc.docker.ListContainers(ctx)
}
