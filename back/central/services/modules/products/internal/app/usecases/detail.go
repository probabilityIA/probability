package usecases

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (uc *UseCases) GetProductDetail(ctx context.Context, businessID uint, id string) (*domain.ProductDetailResponse, error) {
	return uc.ProductCRUD.GetProductDetail(ctx, businessID, id)
}
