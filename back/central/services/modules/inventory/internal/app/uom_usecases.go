package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
)

func (uc *useCase) ListUoMs(ctx context.Context) ([]entities.UnitOfMeasure, error) {
	return uc.repo.ListUoMs(ctx)
}

func (uc *useCase) ListProductUoMs(ctx context.Context, businessID uint, productID string) ([]entities.ProductUoM, error) {
	return uc.repo.ListProductUoMs(ctx, dtos.ListProductUoMParams{BusinessID: businessID, ProductID: productID})
}

func (uc *useCase) CreateProductUoM(ctx context.Context, dto request.CreateProductUoMDTO) (*entities.ProductUoM, error) {
	if _, _, _, err := uc.repo.GetProductByID(ctx, dto.ProductID, dto.BusinessID); err != nil {
		return nil, domainerrors.ErrProductNotFound
	}
	uom, err := uc.repo.GetUoMByCode(ctx, dto.UomCode)
	if err != nil {
		return nil, err
	}
	factor := dto.ConversionFactor
	if factor <= 0 {
		factor = 1
	}
	pu := &entities.ProductUoM{
		ProductID:        dto.ProductID,
		UomID:            uom.ID,
		BusinessID:       dto.BusinessID,
		ConversionFactor: factor,
		IsBase:           dto.IsBase,
		Barcode:          dto.Barcode,
		IsActive:         true,
	}
	return uc.repo.CreateProductUoM(ctx, pu)
}

func (uc *useCase) DeleteProductUoM(ctx context.Context, businessID, id uint) error {
	return uc.repo.DeleteProductUoM(ctx, businessID, id)
}

func (uc *useCase) ConvertUoM(ctx context.Context, dto request.ConvertUoMDTO) (*response.ConvertUoMResult, error) {
	uoms, err := uc.repo.ListProductUoMs(ctx, dtos.ListProductUoMParams{BusinessID: dto.BusinessID, ProductID: dto.ProductID})
	if err != nil {
		return nil, err
	}
	if len(uoms) == 0 {
		return nil, domainerrors.ErrProductUoMNotFound
	}

	var fromPU, toPU, basePU *entities.ProductUoM
	for i := range uoms {
		if uoms[i].IsBase {
			basePU = &uoms[i]
		}
		if uoms[i].UomCode == dto.FromUomCode {
			fromPU = &uoms[i]
		}
		if uoms[i].UomCode == dto.ToUomCode {
			toPU = &uoms[i]
		}
	}
	if fromPU == nil || toPU == nil || basePU == nil {
		return nil, domainerrors.ErrUomConversion
	}

	baseQty := dto.Quantity * fromPU.ConversionFactor
	converted := baseQty / toPU.ConversionFactor

	return &response.ConvertUoMResult{
		FromUomCode:      dto.FromUomCode,
		ToUomCode:        dto.ToUomCode,
		InputQuantity:    dto.Quantity,
		ConvertedQty:     converted,
		BaseUnitQuantity: baseQty,
		BaseUomCode:      basePU.UomCode,
	}, nil
}
