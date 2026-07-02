package usecases

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
)

func (uc *UseCases) UpsertFromProvider(ctx context.Context, dto *domain.ProductProviderUpsertDTO) error {
	if dto == nil || dto.BusinessID == 0 || dto.SKU == "" {
		return domain.ErrInvalidProductData
	}

	existing, err := uc.repo.GetProductBySKU(ctx, dto.BusinessID, dto.SKU)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			_, cerr := uc.ProductCRUD.CreateProduct(ctx, &domain.CreateProductRequest{
				BusinessID:     dto.BusinessID,
				SKU:            dto.SKU,
				Name:           dto.Name,
				ExternalID:     dto.ExternalID,
				Price:          dto.Price,
				Currency:       "COP",
				TrackInventory: dto.TrackInventory,
				Status:         "active",
				IsActive:       true,
			})
			return cerr
		}
		return err
	}

	name := dto.Name
	price := dto.Price
	track := dto.TrackInventory
	_, uerr := uc.ProductCRUD.UpdateProduct(ctx, dto.BusinessID, existing.ID, &domain.UpdateProductRequest{
		Name:           &name,
		Price:          &price,
		TrackInventory: &track,
	})
	return uerr
}
