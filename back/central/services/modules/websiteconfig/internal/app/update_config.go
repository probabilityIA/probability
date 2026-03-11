package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
)

func (uc *UseCase) UpdateConfig(ctx context.Context, businessID uint, dto *dtos.UpdateConfigDTO) (*entities.WebsiteConfig, error) {
	uc.logger.Info(ctx).Uint("business_id", businessID).Msg("Updating website config")
	return uc.repo.UpsertConfig(ctx, businessID, dto)
}
