package usecases

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
)

func filtrarPorSKU(mapped []domain.MappedItem, skus []string) []domain.MappedItem {
	if len(skus) == 0 {
		return mapped
	}
	buscados := make(map[string]bool, len(skus))
	for _, sku := range skus {
		if limpio := normalizeSKU(sku); limpio != "" {
			buscados[limpio] = true
		}
	}
	if len(buscados) == 0 {
		return mapped
	}
	filtrados := make([]domain.MappedItem, 0, len(buscados))
	for _, m := range mapped {
		if buscados[normalizeSKU(m.SKU)] {
			filtrados = append(filtrados, m)
		}
	}
	return filtrados
}

func (uc *tiendanubeUseCase) CompareInventory(ctx context.Context, integrationID string, businessID uint, page, pageSize int, skus ...string) (*inventorycompare.Page, error) {
	page, pageSize = inventorycompare.NormalizePaging(page, pageSize)

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	integration, cred, err := uc.resolveIntegrationForBusiness(ctx, integrationID, businessID)
	if err != nil {
		return nil, err
	}

	mapped, err := uc.productRepo.ListMappedItems(ctx, uint(integIDUint))
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
	externos := make([]string, 0, len(grupo))
	for _, m := range grupo {
		productIDs = append(productIDs, m.ProductID)
		externos = append(externos, m.ExternalItemID)
	}

	stock, err := uc.productRepo.GetStockForProducts(ctx, productIDs, resolveWarehouseIDs(cfg))
	if err != nil {
		return nil, fmt.Errorf("getting stock: %w", err)
	}

	stockCanal, err := uc.client.GetProductsStock(ctx, cred, externos)
	if err != nil {
		return nil, fmt.Errorf("reading channel stock: %w", err)
	}

	porPublicacion := make(map[string]domain.ChannelStock, len(stockCanal))
	for _, s := range stockCanal {
		porPublicacion[strings.TrimSpace(s.ExternalID)] = s
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
		}
		if cantidad, ok := stock[m.ProductID]; ok {
			item.ProbabilityQty = cantidad
			item.HasStock = true
		}
		if canal, ok := porPublicacion[strings.TrimSpace(m.ExternalItemID)]; ok && canal.Found {
			item.ChannelQty = canal.Quantity
			item.ChannelFound = true
			item.ChannelManaged = canal.ManageStock
		}
		items = append(items, item)
	}

	resultado.Rows = inventorycompare.Build(items)
	resultado.Totals = inventorycompare.Summarize(resultado.Rows)

	if err := uc.productRepo.SaveCompareSnapshot(ctx, businessID, uint(integIDUint), resultado.Rows, resultado.CheckedAt); err != nil {
		uc.logger.Warn(ctx).Err(err).
			Uint64("integration_id", integIDUint).
			Msg("No se pudo guardar la foto del comparativo de inventario")
	}

	return resultado, nil
}

func (uc *tiendanubeUseCase) LoadInventoryCompare(ctx context.Context, integrationID string, businessID uint, opts inventorycompare.LoadOptions) (*inventorycompare.Page, error) {
	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)

	integration, err := uc.service.GetIntegrationByID(ctx, integrationID)
	if err != nil {
		return nil, fmt.Errorf("getting integration: %w", err)
	}
	if integration == nil || integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, domain.ErrIntegrationNotFound
	}

	return uc.productRepo.LoadCompareSnapshot(ctx, businessID, uint(integIDUint), opts)
}
