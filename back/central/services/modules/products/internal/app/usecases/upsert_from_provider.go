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
			created, cerr := uc.ProductCRUD.CreateProduct(ctx, &domain.CreateProductRequest{
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
			if cerr != nil {
				return cerr
			}
			return uc.ensureIntegrationMapping(ctx, created.ID, dto.IntegrationID, dto.ExternalID)
		}
		return err
	}

	name := dto.Name
	price := dto.Price
	track := dto.TrackInventory
	if _, uerr := uc.ProductCRUD.UpdateProduct(ctx, dto.BusinessID, existing.ID, &domain.UpdateProductRequest{
		Name:           &name,
		Price:          &price,
		TrackInventory: &track,
	}); uerr != nil {
		return uerr
	}
	return uc.ensureIntegrationMapping(ctx, existing.ID, dto.IntegrationID, dto.ExternalID)
}

func (uc *UseCases) ensureIntegrationMapping(ctx context.Context, productID string, integrationID uint, externalID string) error {
	if integrationID == 0 || productID == "" {
		return nil
	}
	exists, err := uc.repo.ProductIntegrationExists(ctx, productID, integrationID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = uc.repo.AddProductIntegration(ctx, productID, integrationID, externalID, nil, nil, nil)
	return err
}
