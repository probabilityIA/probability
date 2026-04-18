package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
)

func (uc *useCase) ExportKardex(ctx context.Context, dto request.KardexExportDTO) (*response.KardexExportResult, error) {
	entries, err := uc.repo.GetKardex(ctx, dtos.KardexQueryParams{
		BusinessID:  dto.BusinessID,
		ProductID:   dto.ProductID,
		WarehouseID: dto.WarehouseID,
		From:        dto.From,
		To:          dto.To,
	})
	if err != nil {
		return nil, err
	}

	result := &response.KardexExportResult{
		BusinessID:  dto.BusinessID,
		ProductID:   dto.ProductID,
		WarehouseID: dto.WarehouseID,
		Entries:     entries,
	}

	for _, e := range entries {
		if e.Quantity > 0 {
			result.TotalIn += e.Quantity
		} else {
			result.TotalOut += -e.Quantity
		}
		result.FinalBalance = e.RunningBalance
	}
	return result, nil
}
