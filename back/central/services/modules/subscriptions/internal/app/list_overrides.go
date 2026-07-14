package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

func (uc *UseCase) ListOverrides(ctx context.Context, businessID uint) ([]entities.BusinessModuleOverride, error) {
	return uc.repo.ListOverridesByBusiness(ctx, businessID)
}
