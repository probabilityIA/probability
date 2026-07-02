package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/response"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
)

const providerSyncProgressBatch = 25

func (uc *useCase) SyncProviderStock(ctx context.Context, dto request.ProviderStockSyncDTO) (*response.ProviderSyncResult, error) {
	total := len(dto.Items)
	result := &response.ProviderSyncResult{Total: total}

	uc.publishInventorySyncEvent(dto.BusinessID, "inventory.sync.started", map[string]interface{}{
		"correlation_id": dto.CorrelationID,
		"total":          total,
		"integration_id": dto.IntegrationID,
		"provider":       dto.Provider,
	})

	inboundID, errIn := uc.repo.GetMovementTypeIDByCode(ctx, "inbound")
	outboundID, errOut := uc.repo.GetMovementTypeIDByCode(ctx, "outbound")

	for i, item := range dto.Items {
		uc.applyProviderStockItem(ctx, dto.BusinessID, item, inboundID, outboundID, errIn, errOut, result)

		if (i+1)%providerSyncProgressBatch == 0 || i+1 == total {
			uc.publishInventorySyncEvent(dto.BusinessID, "inventory.sync.progress", map[string]interface{}{
				"correlation_id": dto.CorrelationID,
				"processed":      i + 1,
				"total":          total,
				"updated":        result.Updated,
				"unchanged":      result.Unchanged,
				"skipped":        result.Skipped,
				"failed":         result.Failed,
			})
		}
	}

	uc.publishInventorySyncEvent(dto.BusinessID, "inventory.sync.completed", map[string]interface{}{
		"correlation_id": dto.CorrelationID,
		"total":          total,
		"updated":        result.Updated,
		"unchanged":      result.Unchanged,
		"skipped":        result.Skipped,
		"failed":         result.Failed,
	})

	return result, nil
}

func (uc *useCase) applyProviderStockItem(
	ctx context.Context,
	businessID uint,
	item request.ProviderStockSyncItem,
	inboundID, outboundID uint,
	errIn, errOut error,
	result *response.ProviderSyncResult,
) {
	productID, _, track, err := uc.repo.GetProductBySKU(ctx, item.SKU, businessID)
	if err != nil || productID == "" || !track {
		result.Skipped++
		return
	}

	warehouseID := item.WarehouseID
	if warehouseID == 0 {
		if def, derr := uc.repo.GetDefaultWarehouseID(ctx, businessID); derr == nil {
			warehouseID = def
		}
	}
	if warehouseID == 0 {
		result.Failed++
		return
	}

	current := uc.currentWarehouseQty(ctx, productID, businessID, warehouseID)
	delta := item.Quantity - current
	if delta == 0 {
		result.Unchanged++
		return
	}

	movTypeID := inboundID
	movErr := errIn
	if delta < 0 {
		movTypeID = outboundID
		movErr = errOut
	}
	if movErr != nil || movTypeID == 0 {
		result.Failed++
		return
	}

	_, aerr := uc.repo.AdjustStockTx(ctx, dtos.AdjustStockTxParams{
		ProductID:      productID,
		WarehouseID:    warehouseID,
		BusinessID:     businessID,
		Quantity:       delta,
		MovementTypeID: movTypeID,
		Reason:         "Sincronizacion inventario Siigo",
		ReferenceType:  "provider_sync",
	})
	if aerr != nil {
		result.Failed++
		return
	}

	result.Updated++
	uc.updateProductTotalStock(ctx, productID, businessID)
}

func (uc *useCase) currentWarehouseQty(ctx context.Context, productID string, businessID, warehouseID uint) int {
	levels, err := uc.repo.GetProductInventory(ctx, dtos.GetProductInventoryParams{
		ProductID:  productID,
		BusinessID: businessID,
	})
	if err != nil {
		return 0
	}
	for _, l := range levels {
		if l.WarehouseID == warehouseID {
			return l.Quantity
		}
	}
	return 0
}

func (uc *useCase) publishInventorySyncEvent(businessID uint, eventType string, data map[string]interface{}) {
	if uc.eventPublisher == nil {
		return
	}
	go func() {
		_ = uc.eventPublisher.PublishInventoryEvent(context.Background(), ports.InventoryEvent{
			EventType:  eventType,
			BusinessID: businessID,
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			Data:       data,
		})
	}()
}
