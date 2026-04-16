package app

import (
	"context"
	"io"

	"github.com/secamc93/probability/back/monitoring/internal/domain/entities"
)

func (uc *useCase) GetContainerLogs(ctx context.Context, id string, tail int) ([]entities.LogEntry, error) {
	if tail <= 0 {
		tail = 100
	}
	return uc.docker.GetContainerLogs(ctx, id, tail)
}

func (uc *useCase) StreamContainerLogs(ctx context.Context, id string) (io.ReadCloser, error) {
	return uc.docker.StreamContainerLogs(ctx, id)
}
