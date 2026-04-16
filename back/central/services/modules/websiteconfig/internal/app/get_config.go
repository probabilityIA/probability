package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
)

func (uc *UseCase) GetConfig(ctx context.Context, businessID uint) (*entities.WebsiteConfig, error) {
	config, err := uc.repo.GetConfig(ctx, businessID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		// Return defaults
		return &entities.WebsiteConfig{
			BusinessID:           businessID,
			Template:             "default",
			ShowHero:             true,
			ShowFeaturedProducts: true,
			ShowFullCatalog:      true,
			ShowContact:          true,
		}, nil
	}
	return config, nil
}
