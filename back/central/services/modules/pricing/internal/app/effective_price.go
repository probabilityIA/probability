package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
)

func (uc *UseCase) GetEffectivePrice(ctx context.Context, params dtos.EffectivePriceParams) (*entities.EffectivePrice, error) {
	return uc.repo.GetEffectivePrice(ctx, params)
}
