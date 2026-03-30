package app

import (
	"context"

	domainErrors "github.com/secamc93/probability/back/monitoring/internal/domain/errors"
)

func (uc *useCase) ContainerAction(ctx context.Context, id string, action string) error {
	switch action {
	case "restart":
		return uc.docker.RestartContainer(ctx, id)
	case "stop":
		return uc.docker.StopContainer(ctx, id)
	case "start":
		return uc.docker.StartContainer(ctx, id)
	default:
		return domainErrors.ErrInvalidAction
	}
}
