package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) ListChannelStatuses(ctx context.Context, integrationTypeID uint, isActive *bool) ([]entities.ChannelStatusInfo, error) {
	return uc.repo.ListChannelStatuses(ctx, integrationTypeID, isActive)
}
