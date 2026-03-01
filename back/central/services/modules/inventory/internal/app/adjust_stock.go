package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
)

func (uc *useCase) AdjustStock(ctx context.Context, dto dtos.AdjustStockDTO) (*entities.StockMovement, error) {
	if dto.Quantity == 0 {
		return nil, domainerrors.ErrInvalidQuantity
	}

	// Verificar producto
	_, _, trackInventory, err := uc.repo.GetProductByID(ctx, dto.ProductID, dto.BusinessID)
	if err != nil {
		return nil, domainerrors.ErrProductNotFound
	}
	if !trackInventory {
		return nil, domainerrors.ErrProductNoTracking
	}

	// Verificar bodega
	exists, err := uc.repo.WarehouseExists(ctx, dto.WarehouseID, dto.BusinessID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domainerrors.ErrWarehouseNotFound
	}

	// Determinar tipo de movimiento
	movTypeCode := "adjustment"
	if dto.Quantity > 0 {
		movTypeCode = "inbound"
	} else {
		movTypeCode = "outbound"
	}

	movTypeID, err := uc.repo.GetMovementTypeIDByCode(ctx, movTypeCode)
	if err != nil {
		return nil, err
	}

	// Ejecutar ajuste dentro de transacción con SELECT FOR UPDATE
	txResult, err := uc.repo.AdjustStockTx(ctx, dtos.AdjustStockTxParams{
		ProductID:      dto.ProductID,
		WarehouseID:    dto.WarehouseID,
		LocationID:     dto.LocationID,
		BusinessID:     dto.BusinessID,
		Quantity:       dto.Quantity,
		MovementTypeID: movTypeID,
		Reason:         dto.Reason,
		Notes:          dto.Notes,
		ReferenceType:  "manual",
		CreatedByID:    dto.CreatedByID,
	})
	if err != nil {
		return nil, err
	}

	// Actualizar stock total del producto (best-effort, fuera de la tx)
	uc.updateProductTotalStock(ctx, dto.ProductID, dto.BusinessID)

	// Publicar sync a canales de venta (fire-and-forget)
	uc.publishSync(ctx, dto.ProductID, dto.BusinessID, txResult.NewQuantity, dto.WarehouseID, "manual_adjustment")

	return txResult.Movement, nil
}

// updateProductTotalStock recalcula y actualiza el StockQuantity del producto
func (uc *useCase) updateProductTotalStock(ctx context.Context, productID string, businessID uint) {
	levels, err := uc.repo.GetProductInventory(ctx, dtos.GetProductInventoryParams{
		ProductID:  productID,
		BusinessID: businessID,
	})
	if err != nil {
		return
	}

	total := 0
	for _, level := range levels {
		total += level.Quantity
	}

	_ = uc.repo.UpdateProductStockQuantity(ctx, productID, total)
}

// publishSync publica sync a canales de venta vinculados
func (uc *useCase) publishSync(ctx context.Context, productID string, businessID uint, newQuantity int, warehouseID uint, source string) {
	if uc.publisher == nil {
		return
	}

	integrations, err := uc.repo.GetProductIntegrations(ctx, productID, businessID)
	if err != nil || len(integrations) == 0 {
		return
	}

	for _, integ := range integrations {
		_ = uc.publisher.PublishInventorySync(ctx, ports.InventorySyncMessage{
			ProductID:         productID,
			ExternalProductID: integ.ExternalProductID,
			IntegrationID:     integ.IntegrationID,
			BusinessID:        businessID,
			NewQuantity:       newQuantity,
			WarehouseID:       warehouseID,
			Source:            source,
			Timestamp:         time.Now().UTC().Format(time.RFC3339),
		})
	}
}
