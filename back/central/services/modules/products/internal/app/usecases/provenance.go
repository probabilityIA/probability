package usecases

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (uc *UseCases) stampOrigin(ctx context.Context, businessID uint, before *domain.Product, origin domain.FieldOrigin, batchID string) {
	if before == nil {
		return
	}
	after, err := uc.repo.GetProductByID(ctx, businessID, before.ID)
	if err != nil || after == nil {
		return
	}
	changes := domain.DiffProduct(before, after, origin, batchID, time.Now())
	if len(changes) == 0 {
		return
	}
	_ = uc.repo.RecordFieldChanges(ctx, changes)
}

func (uc *UseCases) UpdateProductWithOrigin(
	ctx context.Context,
	businessID uint,
	id string,
	req *domain.UpdateProductRequest,
	origin domain.FieldOrigin,
) (*domain.ProductResponse, error) {
	before, err := uc.repo.GetProductByID(ctx, businessID, id)
	if err != nil {
		return nil, err
	}
	if before == nil {
		return nil, domain.ErrProductNotFound
	}
	previo := *before

	updated, err := uc.ProductCRUD.UpdateProduct(ctx, businessID, id, req)
	if err != nil {
		return nil, err
	}

	uc.stampOrigin(ctx, businessID, &previo, origin, "")
	return updated, nil
}

func (uc *UseCases) GetProductFieldOrigins(ctx context.Context, businessID uint, productID string) ([]domain.FieldOriginView, error) {
	return uc.repo.GetProductFieldOrigins(ctx, businessID, productID)
}
