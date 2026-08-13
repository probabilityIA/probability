package usecases

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
)

type credencialesTienda struct {
	StoreURL       string
	ConsumerKey    string
	ConsumerSecret string
}

func (uc *wooCommerceUseCase) credenciales(ctx context.Context, integrationID string, integration *domain.Integration) (credencialesTienda, error) {
	storeURL, err := extractString(integration.Config, "store_url")
	if err != nil {
		return credencialesTienda{}, domain.ErrMissingStoreURL
	}

	consumerKey, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_key")
	if err != nil {
		return credencialesTienda{}, fmt.Errorf("decrypting consumer_key: %w", err)
	}
	consumerSecret, err := uc.service.DecryptCredential(ctx, integrationID, "consumer_secret")
	if err != nil {
		return credencialesTienda{}, fmt.Errorf("decrypting consumer_secret: %w", err)
	}

	return credencialesTienda{
		StoreURL:       resolveEffectiveStoreURL(integration, storeURL),
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
	}, nil
}

func filtrarPorSKU(mapped []domain.MappedItem, skus []string) []domain.MappedItem {
	if len(skus) == 0 {
		return mapped
	}
	buscados := make(map[string]bool, len(skus))
	for _, sku := range skus {
		limpio := strings.ToUpper(strings.TrimSpace(sku))
		if limpio != "" {
			buscados[limpio] = true
		}
	}
	if len(buscados) == 0 {
		return mapped
	}
	filtrados := make([]domain.MappedItem, 0, len(buscados))
	for _, m := range mapped {
		if buscados[strings.ToUpper(strings.TrimSpace(m.SKU))] {
			filtrados = append(filtrados, m)
		}
	}
	return filtrados
}

func (uc *wooCommerceUseCase) CompareInventory(ctx context.Context, integrationID string, businessID uint, page, pageSize int, skus ...string) (*inventorycompare.Page, error) {
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

	creds, err := uc.credenciales(ctx, integrationID, integration)
	if err != nil {
		return nil, err
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

	stockCanal, err := uc.client.GetProductsStock(ctx, creds.StoreURL, creds.ConsumerKey, creds.ConsumerSecret, externos)
	if err != nil {
		return nil, fmt.Errorf("reading channel stock: %w", err)
	}

	porPublicacion := make(map[string]domain.ChannelStock, len(stockCanal))
	for _, s := range stockCanal {
		porPublicacion[s.ExternalID] = s
	}

	items := make([]inventorycompare.Item, 0, len(grupo))
	for _, m := range grupo {
		item := inventorycompare.Item{
			ProductID:         m.ProductID,
			SKU:               m.SKU,
			Name:              m.Name,
			ExternalItemID:    m.ExternalItemID,
			ExternalVariantID: m.ExternalVariantID,
		}
		if cantidad, ok := stock[m.ProductID]; ok {
			item.ProbabilityQty = cantidad
			item.HasStock = true
		}
		if canal, ok := porPublicacion[m.ExternalItemID]; ok {
			item.ChannelQty = canal.AvailableQuantity
			item.ChannelFound = true
			item.ChannelManaged = true
		}
		items = append(items, item)
	}

	resultado.Rows = inventorycompare.Build(items)
	resultado.Totals = inventorycompare.Summarize(resultado.Rows)
	return resultado, nil
}
