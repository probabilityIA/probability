package app

import (
	"context"

	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
)

func (uc *useCase) GetSystemStats(ctx context.Context) (*entities.SystemStats, error) {
	return uc.docker.GetSystemStats(ctx)
}
