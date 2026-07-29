package usecases

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

// legacyResolveUserProductID es la implementacion anterior al match configurable.
// Se conserva solo para verificar que el comportamiento no cambio con las reglas por defecto.
func legacyResolveUserProductID(item *domain.MeliItemDetail, sku string) string {
	if len(item.Variations) > 0 {
		for _, v := range item.Variations {
			if v.SellerSKU == sku && v.UserProductID != "" {
				return v.UserProductID
			}
		}
		for _, v := range item.Variations {
			if v.UserProductID != "" {
				return v.UserProductID
			}
		}
	}
	return item.UserProductID
}

func TestResolveUserProductIDIgualQueLegacyConReglasPorDefecto(t *testing.T) {
	casos := []struct {
		nombre string
		item   *domain.MeliItemDetail
		sku    string
	}{
		{
			nombre: "sin variaciones",
			item:   &domain.MeliItemDetail{ID: "MLA1", UserProductID: "UP-BASE"},
			sku:    "ABC",
		},
		{
			nombre: "una variacion que casa",
			item: &domain.MeliItemDetail{ID: "MLA1", UserProductID: "UP-BASE", Variations: []domain.MeliItemVariation{
				{ID: 10, SellerSKU: "ABC", UserProductID: "UP-10"},
			}},
			sku: "ABC",
		},
		{
			nombre: "varias variaciones, una casa",
			item: &domain.MeliItemDetail{ID: "MLA1", Variations: []domain.MeliItemVariation{
				{ID: 10, SellerSKU: "XXX", UserProductID: "UP-10"},
				{ID: 11, SellerSKU: "ABC", UserProductID: "UP-11"},
			}},
			sku: "ABC",
		},
		{
			nombre: "varias variaciones, ninguna casa",
			item: &domain.MeliItemDetail{ID: "MLA1", Variations: []domain.MeliItemVariation{
				{ID: 10, SellerSKU: "XXX", UserProductID: "UP-10"},
				{ID: 11, SellerSKU: "YYY", UserProductID: "UP-11"},
			}},
			sku: "ABC",
		},
		{
			nombre: "la variacion que casa no tiene user_product_id",
			item: &domain.MeliItemDetail{ID: "MLA1", Variations: []domain.MeliItemVariation{
				{ID: 10, SellerSKU: "ABC", UserProductID: ""},
				{ID: 11, SellerSKU: "YYY", UserProductID: "UP-11"},
			}},
			sku: "ABC",
		},
		{
			nombre: "variaciones sin seller_sku",
			item: &domain.MeliItemDetail{ID: "MLA1", Variations: []domain.MeliItemVariation{
				{ID: 10, SellerSKU: "", UserProductID: "UP-10"},
			}},
			sku: "",
		},
		{
			nombre: "sku duplicado entre variaciones",
			item: &domain.MeliItemDetail{ID: "MLA1", Variations: []domain.MeliItemVariation{
				{ID: 10, SellerSKU: "ABC", UserProductID: "UP-10"},
				{ID: 11, SellerSKU: "ABC", UserProductID: "UP-11"},
			}},
			sku: "ABC",
		},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			esperado := legacyResolveUserProductID(c.item, c.sku)
			m := domain.MappedItem{SKU: c.sku, ExternalItemID: c.item.ID}
			obtenido := resolveUserProductID(c.item, m, productmatch.DefaultRules())
			if obtenido != esperado {
				t.Fatalf("cambio de comportamiento: legacy=%q nuevo=%q", esperado, obtenido)
			}
		})
	}
}

func TestResolveUserProductIDPrefiereElVariantIDGuardado(t *testing.T) {
	item := &domain.MeliItemDetail{ID: "MLA1", Variations: []domain.MeliItemVariation{
		{ID: 10, SellerSKU: "ABC", UserProductID: "UP-10"},
		{ID: 11, SellerSKU: "ABC", UserProductID: "UP-11"},
	}}
	m := domain.MappedItem{SKU: "ABC", ExternalVariantID: "11"}

	if got := resolveUserProductID(item, m, productmatch.DefaultRules()); got != "UP-11" {
		t.Fatalf("con external_variant_id guardado debia resolver UP-11, obtuve %q", got)
	}
}
