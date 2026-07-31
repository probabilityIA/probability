package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/websiteconfig/internal/domain/entities"
)

func (uc *UseCase) GetProductCategories(ctx context.Context, businessID uint) ([]entities.ProductCategory, error) {
	return uc.repo.ListProductCategories(ctx, businessID)
}
