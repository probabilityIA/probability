package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orderstatus/internal/domain/entities"
)

func (uc *useCase) GetChannelStatus(ctx context.Context, id uint) (*entities.ChannelStatusInfo, error) {
	return uc.repo.GetChannelStatusByID(ctx, id)
}
