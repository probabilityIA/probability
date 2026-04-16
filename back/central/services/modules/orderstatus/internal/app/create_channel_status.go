package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) CreateChannelStatus(ctx context.Context, status *entities.ChannelStatusInfo) (*entities.ChannelStatusInfo, error) {
	return uc.repo.CreateChannelStatus(ctx, status)
}
