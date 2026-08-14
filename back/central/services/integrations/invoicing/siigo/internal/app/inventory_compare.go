package app

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
)

const vigenciaFotoStock = 2 * time.Minute

type fotoStockSiigo struct {
	productos []dtos.ProductItem
	tomadaA   time.Time
}

func (uc *invoicingUseCase) productosSiigoRecientes(
	ctx context.Context,
	integrationID string,
	credentials dtos.Credentials,
) ([]dtos.ProductItem, time.Time, error) {
	if guardado, ok := uc.fotoStock.Load(integrationID); ok {
		foto := guardado.(fotoStockSiigo)
		if time.Since(foto.tomadaA) < vigenciaFotoStock {
			return foto.productos, foto.tomadaA, nil
		}
	}

	productos, err := uc.listAllSiigoProducts(ctx, credentials)
	if err != nil {
		return nil, time.Time{}, err
	}

	tomadaA := time.Now()
	uc.fotoStock.Store(integrationID, fotoStockSiigo{productos: productos, tomadaA: tomadaA})
	return productos, tomadaA, nil
}

func normalizarSKU(sku string) string {
	return strings.ToLower(strings.TrimSpace(sku))
}

func (uc *invoicingUseCase) CompareInventory(
	ctx context.Context,
	integrationID string,
	businessID uint,
	page, pageSize int,
	skus ...string,
) (*inventorycompare.Page, error) {
	page, pageSize = inventorycompare.NormalizePaging(page, pageSize)

	integration, err := uc.integrationCore.GetIntegrationByID(ctx, integrationID)
	if err != nil || integration == nil {
		return nil, fmt.Errorf("integracion no encontrada")
	}
	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, fmt.Errorf("la integracion no pertenece al negocio")
	}

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	asociados, err := uc.productRepo.ListAssociatedSKUs(ctx, businessID, uint(integIDUint))
	if err != nil {
		return nil, fmt.Errorf("listing associated skus: %w", err)
	}

	propios, err := uc.productRepo.ListProductsByBusiness(ctx, businessID)
	if err != nil {
		return nil, fmt.Errorf("listing probability products: %w", err)
	}

	pedidos := make(map[string]bool, len(skus))
	for _, sku := range skus {
		if limpio := normalizarSKU(sku); limpio != "" {
			pedidos[limpio] = true
		}
	}

	candidatos := make([]dtos.ProductForSync, 0, len(propios))
	for _, producto := range propios {
		clave := normalizarSKU(producto.SKU)
		if clave == "" || !asociados[clave] {
			continue
		}
		if len(pedidos) > 0 && !pedidos[clave] {
			continue
		}
		candidatos = append(candidatos, producto)
	}

	total := len(candidatos)
	grupo := inventorycompare.Slice(candidatos, page, pageSize)
	if len(pedidos) > 0 {
		grupo = candidatos
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

	credentials, err := uc.resolveWebhookCredentials(ctx, integrationID)
	if err != nil {
		return nil, err
	}
	productosSiigo, tomadaA, err := uc.productosSiigoRecientes(ctx, integrationID, credentials)
	if err != nil {
		return nil, fmt.Errorf("listing siigo products: %w", err)
	}
	resultado.CheckedAt = tomadaA

	enSiigo := make(map[string]dtos.ProductItem, len(productosSiigo))
	for _, producto := range productosSiigo {
		if clave := normalizarSKU(producto.Code); clave != "" {
			enSiigo[clave] = producto
		}
	}

	items := make([]inventorycompare.Item, 0, len(grupo))
	for _, producto := range grupo {
		item := inventorycompare.Item{
			ProductID:      producto.ID,
			SKU:            producto.SKU,
			Name:           producto.Name,
			ImageURL:       producto.ImageURL,
			ProbabilityQty: producto.StockQuantity,
			HasStock:       true,
		}
		if siigo, ok := enSiigo[normalizarSKU(producto.SKU)]; ok {
			item.ExternalItemID = siigo.ID
			item.ChannelQty = int(siigo.AvailableQuantity)
			item.ChannelFound = true
			item.ChannelManaged = siigo.StockControl
			if !siigo.StockControl {
				item.SkipReason = "el producto no controla inventario en Siigo"
			}
		}
		items = append(items, item)
	}

	resultado.Rows = inventorycompare.Build(items)
	resultado.Totals = inventorycompare.Summarize(resultado.Rows)

	if err := uc.productRepo.SaveCompareSnapshot(ctx, businessID, uint(integIDUint), resultado.Rows, resultado.CheckedAt); err != nil {
		uc.log.Warn(ctx).Err(err).
			Uint64("integration_id", integIDUint).
			Msg("No se pudo guardar la foto del comparativo de inventario")
	}

	return resultado, nil
}

func (uc *invoicingUseCase) LoadInventoryCompare(
	ctx context.Context,
	integrationID string,
	businessID uint,
	opts inventorycompare.LoadOptions,
) (*inventorycompare.Page, error) {
	integration, err := uc.integrationCore.GetIntegrationByID(ctx, integrationID)
	if err != nil || integration == nil {
		return nil, fmt.Errorf("integracion no encontrada")
	}
	if integration.BusinessID == nil || *integration.BusinessID != businessID {
		return nil, fmt.Errorf("la integracion no pertenece al negocio")
	}

	integIDUint, _ := strconv.ParseUint(integrationID, 10, 64)
	return uc.productRepo.LoadCompareSnapshot(ctx, businessID, uint(integIDUint), opts)
}
