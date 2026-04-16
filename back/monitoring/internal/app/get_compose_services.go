package app

import (
	"context"

	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
)

func (uc *useCase) GetComposeServices(ctx context.Context) ([]entities.ComposeService, error) {
	return uc.docker.GetComposeServices(ctx)
}
