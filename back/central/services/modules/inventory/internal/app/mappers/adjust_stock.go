package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
)

func AdjustStockDTOToTxParams(dto request.AdjustStockDTO, movTypeID uint, referenceType string) dtos.AdjustStockTxParams {
	return dtos.AdjustStockTxParams{
		ProductID:      dto.ProductID,
		WarehouseID:    dto.WarehouseID,
		LocationID:     dto.LocationID,
		BusinessID:     dto.BusinessID,
		Quantity:       dto.Quantity,
		MovementTypeID: movTypeID,
		Reason:         dto.Reason,
		Notes:          dto.Notes,
		ReferenceType:  referenceType,
		CreatedByID:    dto.CreatedByID,
	}
}

func TransferStockDTOToTxParams(dto request.TransferStockDTO, movTypeID uint, referenceType string) dtos.TransferStockTxParams {
	return dtos.TransferStockTxParams{
		ProductID:       dto.ProductID,
		FromWarehouseID: dto.FromWarehouseID,
		ToWarehouseID:   dto.ToWarehouseID,
		FromLocationID:  dto.FromLocationID,
		ToLocationID:    dto.ToLocationID,
		BusinessID:      dto.BusinessID,
		Quantity:        dto.Quantity,
		MovementTypeID:  movTypeID,
		Reason:          dto.Reason,
		Notes:           dto.Notes,
		ReferenceType:   referenceType,
		CreatedByID:     dto.CreatedByID,
	}
}
