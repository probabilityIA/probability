package usecases

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/shopify/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

// legacyResolveInventoryItemID es la implementacion anterior al match configurable.
// Se conserva solo para verificar que el comportamiento no cambio con las reglas por defecto.
func legacyResolveInventoryItemID(product *domain.ShopifyProduct, sku string) int64 {
	if product == nil {
		return 0
	}
	if len(product.Variants) > 1 && sku != "" {
		for _, v := range product.Variants {
			if v.SKU == sku && v.InventoryItemID != 0 {
				return v.InventoryItemID
			}
		}
	}
	for _, v := range product.Variants {
		if v.InventoryItemID != 0 {
			return v.InventoryItemID
		}
	}
	return 0
}

func TestResolveInventoryItemIDIgualQueLegacyConReglasPorDefecto(t *testing.T) {
	casos := []struct {
		nombre  string
		product *domain.ShopifyProduct
		sku     string
	}{
		{
			nombre:  "producto nulo",
			product: nil,
			sku:     "ABC",
		},
		{
			nombre: "una sola variante que casa",
			product: &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
				{ID: 10, SKU: "ABC", InventoryItemID: 100},
			}},
			sku: "ABC",
		},
		{
			nombre: "una sola variante que NO casa: igual devuelve la unica",
			product: &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
				{ID: 10, SKU: "OTRO", InventoryItemID: 100},
			}},
			sku: "ABC",
		},
		{
			nombre: "varias variantes, una casa",
			product: &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
				{ID: 10, SKU: "XXX", InventoryItemID: 100},
				{ID: 11, SKU: "ABC", InventoryItemID: 200},
				{ID: 12, SKU: "YYY", InventoryItemID: 300},
			}},
			sku: "ABC",
		},
		{
			nombre: "varias variantes, ninguna casa: cae a la primera con inventory_item_id",
			product: &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
				{ID: 10, SKU: "XXX", InventoryItemID: 0},
				{ID: 11, SKU: "YYY", InventoryItemID: 200},
			}},
			sku: "ABC",
		},
		{
			nombre: "la variante que casa no tiene inventory_item_id",
			product: &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
				{ID: 10, SKU: "ABC", InventoryItemID: 0},
				{ID: 11, SKU: "YYY", InventoryItemID: 200},
			}},
			sku: "ABC",
		},
		{
			nombre: "sku duplicado entre variantes",
			product: &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
				{ID: 10, SKU: "ABC", InventoryItemID: 100},
				{ID: 11, SKU: "ABC", InventoryItemID: 200},
			}},
			sku: "ABC",
		},
		{
			nombre: "variantes sin sku",
			product: &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
				{ID: 10, SKU: "", InventoryItemID: 100},
				{ID: 11, SKU: "", InventoryItemID: 200},
			}},
			sku: "",
		},
		{
			nombre: "ninguna variante tiene inventory_item_id",
			product: &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
				{ID: 10, SKU: "ABC", InventoryItemID: 0},
			}},
			sku: "ABC",
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			esperado := legacyResolveInventoryItemID(c.product, c.sku)
			// Mapeo tal como existe hoy en la base: sin external_variant_id ni external_sku.
			m := domain.MappedItem{SKU: c.sku}
			obtenido := resolveInventoryItemID(c.product, m, productmatch.DefaultRules())
			if obtenido != esperado {
				t.Fatalf("cambio de comportamiento: legacy=%d nuevo=%d", esperado, obtenido)
			}
		})
	}
}

func TestResolveInventoryItemIDPrefiereElVariantIDGuardado(t *testing.T) {
	product := &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
		{ID: 10, SKU: "ABC", InventoryItemID: 100},
		{ID: 11, SKU: "ABC", InventoryItemID: 200},
	}}
	m := domain.MappedItem{SKU: "ABC", ExternalVariantID: "11"}

	if got := resolveInventoryItemID(product, m, productmatch.DefaultRules()); got != 200 {
		t.Fatalf("con external_variant_id guardado debia resolver 200, obtuve %d", got)
	}
}

func TestResolveInventoryItemIDPorBarcodeCuandoSeConfigura(t *testing.T) {
	product := &domain.ShopifyProduct{ID: 1, Variants: []domain.ShopifyVariant{
		{ID: 10, SKU: "XXX", Barcode: "7701", InventoryItemID: 100},
		{ID: 11, SKU: "YYY", Barcode: "7702", InventoryItemID: 200},
	}}
	reglas := []productmatch.Rule{{Probability: productmatch.FieldBarcode, Channel: productmatch.FieldBarcode}}
	m := domain.MappedItem{SKU: "SIN-MATCH", Barcode: "7702"}

	if got := resolveInventoryItemID(product, m, reglas); got != 200 {
		t.Fatalf("con la regla barcode debia resolver 200, obtuve %d", got)
	}
}
