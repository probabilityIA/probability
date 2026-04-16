package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) UpdateChannelStatus(ctx context.Context, id uint, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
	return uc.repo.UpdateChannelStatus(ctx, id, status)
}
