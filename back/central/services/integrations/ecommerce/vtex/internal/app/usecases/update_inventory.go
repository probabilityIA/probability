package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/vtex/internal/domain"
)

func (uc *vtexUseCase) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	return uc.PushStock(ctx, integrationID, "", productExternalID, quantity)
}

func (uc *vtexUseCase) PushStock(ctx context.Context, integrationID, productID, productExternalID string, quantity int) error {
	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return domain.ErrIntegrationNotFound
	}

	invCfg := uc.inventoryConfigFrom(integration.Config)
	if !invCfg.Enabled {
		uc.logger.Info(ctx).
			Str("integration_id", integrationID).
			Msg("Sync de inventario desactivado para la integracion VTEX, push omitido")
		return nil
	}

	groups := invCfg.WarehouseGroups()
	if len(groups) == 0 {
		uc.logger.Warn(ctx).
			Str("integration_id", integrationID).
			Msg("Sin bodegas mapeadas para VTEX, push de stock omitido")
		return nil
	}

	cred, err := uc.resolveCredential(ctx, integration, integrationID)
	if err != nil {
		return err
	}

	quantities, err := uc.resolveQuantities(ctx, groups, productID, quantity)
	if err != nil {
		return err
	}

	var lastErr error
	updated := 0

	for vtexWarehouseID, qty := range quantities {
		if err := uc.client.UpdateSKUInventory(ctx, cred, productExternalID, vtexWarehouseID, qty); err != nil {
			uc.logger.Error(ctx).Err(err).
				Str("vtex_sku_id", productExternalID).
				Str("vtex_warehouse", vtexWarehouseID).
				Int("quantity", qty).
				Msg("Error actualizando inventario en VTEX")
			lastErr = err
			continue
		}
		updated++
	}

	if updated == 0 && lastErr != nil {
		return lastErr
	}

	uc.logger.Info(ctx).
		Str("vtex_sku_id", productExternalID).
		Int("warehouses_updated", updated).
		Msg("Stock actualizado en VTEX")

	return nil
}

func (uc *vtexUseCase) resolveQuantities(ctx context.Context, groups map[string][]uint, productID string, quantity int) (map[string]int, error) {
	quantities := make(map[string]int, len(groups))

	if len(groups) == 1 {
		for vtexWarehouseID := range groups {
			quantities[vtexWarehouseID] = quantity
		}
		return quantities, nil
	}

	if productID == "" {
		return nil, fmt.Errorf("vtex: %d bodegas mapeadas requieren el product_id interno para repartir el stock", len(groups))
	}

	for vtexWarehouseID, internalWarehouseIDs := range groups {
		stock, err := uc.productRepo.GetStockForProducts(ctx, []string{productID}, internalWarehouseIDs)
		if err != nil {
			return nil, fmt.Errorf("obteniendo stock para la bodega %s: %w", vtexWarehouseID, err)
		}
		quantities[vtexWarehouseID] = stock[productID]
	}

	return quantities, nil
}
