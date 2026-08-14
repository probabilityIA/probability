package usecases

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
)

func (uc *meliUseCase) CompareInventory(ctx context.Context, integrationID string, businessID uint, page, pageSize int, skus ...string) (*inventorycompare.Page, error) {
	page, pageSize = inventorycompare.NormalizePaging(page, pageSize)

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil {
		return nil, domain.ErrIntegrationNotFound
	}
	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, domain.ErrIntegrationNotFound
	}

	accessToken, err := uc.EnsureValidToken(ctx, integrationID)
	if err != nil {
		return nil, err
	}

	mapped, err := uc.inventoryRepo.ListMappedItems(ctx, uint(integIDUint))
	if err != nil {
		return nil, fmt.Errorf("listing mapped items: %w", err)
	}

	mapped = filtrarPorSKU(mapped, skus)

	total := len(mapped)
	grupo := inventorycompare.Slice(mapped, page, pageSize)
	if len(skus) > 0 {
		grupo = mapped
		page, pageSize = 1, total
	}
	resultado := &inventorycompare.Page{
		Rows:       []inventorycompare.Row{},
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: inventorycompare.TotalPages(total, pageSize),
		CheckedAt:  time.Now(),
	}
	if len(grupo) == 0 {
		return resultado, nil
	}

	cfg := parseInventoryConfig(integration.Config)
	productIDs := make([]string, 0, len(grupo))
	itemIDs := make([]string, 0, len(grupo))
	for _, m := range grupo {
		productIDs = append(productIDs, m.ProductID)
		itemIDs = append(itemIDs, m.ExternalItemID)
	}

	stock, err := uc.inventoryRepo.GetStockForProducts(ctx, productIDs, resolveWarehouseIDs(cfg))
	if err != nil {
		return nil, fmt.Errorf("getting stock: %w", err)
	}

	cli := uc.clientFor(ctx, integration)
	stockCanal, err := cli.GetItemsStock(ctx, accessToken, itemIDs)
	if err != nil {
		if err == domain.ErrTokenExpired {
			nuevo, rerr := uc.EnsureValidToken(ctx, integrationID)
			if rerr != nil {
				return nil, rerr
			}
			stockCanal, err = cli.GetItemsStock(ctx, nuevo, itemIDs)
		}
		if err != nil {
			return nil, fmt.Errorf("reading channel stock: %w", err)
		}
	}

	porPublicacion := make(map[string]int, len(stockCanal))
	for _, s := range stockCanal {
		porPublicacion[claveMeli(s.ExternalItemID, s.ExternalVariantID)] = s.AvailableQuantity
	}

	items := make([]inventorycompare.Item, 0, len(grupo))
	for _, m := range grupo {
		item := inventorycompare.Item{
			ProductID:         m.ProductID,
			SKU:               m.SKU,
			Name:              m.Name,
			ImageURL:          m.ImageURL,
			ExternalItemID:    m.ExternalItemID,
			ExternalVariantID: m.ExternalVariantID,
			ChannelManaged:    true,
		}
		if m.LogisticType == domain.LogisticFulfillment {
			item.SkipReason = "publicacion de fulfillment, el stock lo administra MercadoLibre"
		}
		if cantidad, ok := stock[m.ProductID]; ok {
			item.ProbabilityQty = cantidad
			item.HasStock = true
		}
		if cantidad, ok := porPublicacion[claveMeli(m.ExternalItemID, m.ExternalVariantID)]; ok {
			item.ChannelQty = cantidad
			item.ChannelFound = true
		}
		items = append(items, item)
	}

	resultado.Rows = inventorycompare.Build(items)
	resultado.Totals = inventorycompare.Summarize(resultado.Rows)

	if err := uc.inventoryRepo.SaveCompareSnapshot(ctx, businessID, uint(integIDUint), resultado.Rows, resultado.CheckedAt); err != nil {
		uc.logger.Warn(ctx).Err(err).
			Uint64("integration_id", integIDUint).
			Msg("No se pudo guardar la foto del comparativo de inventario")
	}

	return resultado, nil
}

func (uc *meliUseCase) LoadInventoryCompare(
	ctx context.Context,
	integrationID string,
	businessID uint,
	opts inventorycompare.LoadOptions,
) (*inventorycompare.Page, error) {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil || integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, domain.ErrIntegrationNotFound
	}

	return uc.inventoryRepo.LoadCompareSnapshot(ctx, businessID, uint(integIDUint), opts)
}

func claveMeli(itemID, variantID string) string {
	return itemID + "|" + variantID
}
