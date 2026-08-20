package usecases

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

func briefTieneSKU(items []domain.ProductBrief, sku string) bool {
	for _, item := range items {
		if item.SKU == sku {
			return true
		}
	}
	return false
}

func TestElAplanadoUsaLaVarianteComoUnidadYArmaElExternalID(t *testing.T) {
	plano := flattenProductSKUs([]domain.TiendanubeProduct{
		{
			ID:   2001,
			Name: "Camiseta",
			Variants: []domain.TiendanubeVariant{
				{ID: 9001, SKU: "TALLA-S", Price: 100},
				{ID: 9002, SKU: "TALLA-M", Price: 110},
			},
		},
	})

	if len(plano) != 2 {
		t.Fatalf("cada variante es un SKU vendible, se esperaban 2 y llegaron %d", len(plano))
	}
	if plano[0].ExternalID != "2001:9001" || plano[1].ExternalID != "2001:9002" {
		t.Fatalf("el external_id debe ser producto:variante para poder escribir stock: %q y %q", plano[0].ExternalID, plano[1].ExternalID)
	}
	if plano[0].Name != "Camiseta" {
		t.Fatalf("la variante hereda el nombre del producto: %q", plano[0].Name)
	}
}

func TestElAplanadoIgnoraElProductoSinVariantes(t *testing.T) {
	plano := flattenProductSKUs([]domain.TiendanubeProduct{{ID: 2001, Name: "Sin variantes"}})
	if len(plano) != 0 {
		t.Fatalf("en Tiendanube el SKU vive en la variante: sin variantes no hay nada que cruzar, llegaron %d", len(plano))
	}
}

func TestElReconcileSeparaLoQueCoincideDeLoQueFaltaEnCadaLado(t *testing.T) {
	repo := &fakeRepo{
		probability: []domain.ProductForSync{
			{ID: "p1", SKU: "COMUN-1", Name: "Comun"},
			{ID: "p2", SKU: "SOLO-PROB", Name: "Solo en Probability"},
		},
	}
	client := &fakeClient{
		products: []domain.TiendanubeProduct{
			productoCanal(2001, 9001, "COMUN-1", "Comun en el canal", 100, 5, true),
			productoCanal(2002, 9002, "SOLO-CANAL", "Solo en Tiendanube", 200, 3, true),
		},
	}

	result, err := newTestUseCase(client, repo).ReconcileProducts(context.Background(), testIntegrationID, testBusinessID)
	if err != nil {
		t.Fatalf("ReconcileProducts fallo: %v", err)
	}

	if len(result.MatchedItems) != 1 || result.MatchedItems[0].SKU != "COMUN-1" {
		t.Fatalf("el SKU que existe en ambos lados no se cruzo: %+v", result.MatchedItems)
	}
	if !briefTieneSKU(result.OnlyInProbability, "SOLO-PROB") {
		t.Fatalf("falta el que solo esta en Probability: %+v", result.OnlyInProbability)
	}
	if !briefTieneSKU(result.OnlyInTiendanube, "SOLO-CANAL") {
		t.Fatalf("falta el que solo esta en Tiendanube: %+v", result.OnlyInTiendanube)
	}
}

func TestElReconcileDistingueLoCruzadoDeLoYaAsociado(t *testing.T) {
	repo := &fakeRepo{
		probability: []domain.ProductForSync{{ID: "p1", SKU: "COMUN-1", Name: "Comun"}},
	}
	client := &fakeClient{
		products: []domain.TiendanubeProduct{productoCanal(2001, 9001, "COMUN-1", "Comun", 100, 5, true)},
	}

	result, err := newTestUseCase(client, repo).ReconcileProducts(context.Background(), testIntegrationID, testBusinessID)
	if err != nil {
		t.Fatalf("ReconcileProducts fallo: %v", err)
	}
	if result.Matched != 0 || len(result.MatchedNotAssociated) != 1 {
		t.Fatalf("cruzar por SKU no es lo mismo que estar asociado: matched=%d sin_asociar=%d", result.Matched, len(result.MatchedNotAssociated))
	}

	repo.mapped = []domain.MappedItem{{ProductID: "p1", SKU: "COMUN-1", ExternalItemID: "2001:9001"}}
	result, err = newTestUseCase(client, repo).ReconcileProducts(context.Background(), testIntegrationID, testBusinessID)
	if err != nil {
		t.Fatalf("ReconcileProducts fallo: %v", err)
	}
	if result.Matched != 1 || len(result.MatchedNotAssociated) != 0 {
		t.Fatalf("ya asociado debe contar como matched: matched=%d sin_asociar=%d", result.Matched, len(result.MatchedNotAssociated))
	}
}

func TestAplicarHaciaTiendanubeCreaElFaltanteYMapeaElQueYaExiste(t *testing.T) {
	repo := &fakeRepo{
		probability: []domain.ProductForSync{
			{ID: "p1", SKU: "COMUN-1", Name: "Comun", Price: 100},
			{ID: "p2", SKU: "NUEVO-1", Name: "Nuevo", Price: 250, TrackInventory: true, StockQuantity: 4},
		},
	}
	client := &fakeClient{
		products: []domain.TiendanubeProduct{productoCanal(2001, 9001, "COMUN-1", "Comun", 100, 5, true)},
	}

	uc := newTestUseCase(client, repo)
	if err := uc.ApplyProductsToTiendanube(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("ApplyProductsToTiendanube fallo: %v", err)
	}

	if len(client.created) != 1 || client.created[0].SKU != "NUEVO-1" {
		t.Fatalf("solo el que falta en el canal debe crearse: %+v", client.created)
	}
	if len(repo.mappings) != 2 {
		t.Fatalf("ambos productos deben quedar mapeados (el existente y el creado): %d", len(repo.mappings))
	}

	var refCreado string
	for _, ref := range repo.mappings {
		if ref.SKU == "NUEVO-1" {
			refCreado = ref.ProductID
		}
	}
	if refCreado != "1001:6001" {
		t.Fatalf("el producto creado debe mapearse con producto:variante para poder escribir stock, quedo %q", refCreado)
	}
}

func TestAplicarHaciaTiendanubeRespetaLaSeleccionDeSKUs(t *testing.T) {
	repo := &fakeRepo{
		probability: []domain.ProductForSync{
			{ID: "p1", SKU: "NUEVO-1", Name: "Uno", Price: 100},
			{ID: "p2", SKU: "NUEVO-2", Name: "Dos", Price: 200},
		},
	}
	client := &fakeClient{}

	uc := newTestUseCase(client, repo)
	if err := uc.ApplyProductsToTiendanube(context.Background(), testIntegrationID, testBusinessID, "corr-1", "nuevo-2"); err != nil {
		t.Fatalf("ApplyProductsToTiendanube fallo: %v", err)
	}

	if len(client.created) != 1 || client.created[0].SKU != "NUEVO-2" {
		t.Fatalf("elegir un SKU no debe arrastrar los demas al canal: %+v", client.created)
	}
}

func TestElPesoSeConvierteAKilosAntesDeMandarloAlCanal(t *testing.T) {
	gramos := 2500.0
	repo := &fakeRepo{
		probability: []domain.ProductForSync{
			{ID: "p1", SKU: "NUEVO-1", Name: "Uno", Price: 100, Weight: &gramos, WeightUnit: "g"},
		},
	}
	client := &fakeClient{}

	uc := newTestUseCase(client, repo)
	if err := uc.ApplyProductsToTiendanube(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("ApplyProductsToTiendanube fallo: %v", err)
	}

	if len(client.created) != 1 || client.created[0].Weight == nil {
		t.Fatalf("el peso no se envio: %+v", client.created)
	}
	if got := *client.created[0].Weight; got != 2.5 {
		t.Fatalf("2500 g son 2.5 kg, se envio %v", got)
	}
}

func TestElPesoConUnidadDesconocidaNoSeInventa(t *testing.T) {
	peso := 3.0
	repo := &fakeRepo{
		probability: []domain.ProductForSync{
			{ID: "p1", SKU: "NUEVO-1", Name: "Uno", Price: 100, Weight: &peso, WeightUnit: "quintales"},
		},
	}
	client := &fakeClient{}

	uc := newTestUseCase(client, repo)
	if err := uc.ApplyProductsToTiendanube(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("ApplyProductsToTiendanube fallo: %v", err)
	}

	if len(client.created) != 1 {
		t.Fatalf("el producto debe crearse igual: %+v", client.created)
	}
	if client.created[0].Weight != nil {
		t.Fatalf("con una unidad que no conocemos es mejor no mandar peso que mandar uno falso: %v", *client.created[0].Weight)
	}
}
