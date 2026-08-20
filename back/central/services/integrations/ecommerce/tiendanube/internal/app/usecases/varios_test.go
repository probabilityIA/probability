package usecases

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
)

func TestLosIdsDeConfiguracionLleganEnVariosTipos(t *testing.T) {
	casos := []struct {
		entrada  interface{}
		esperado uint
	}{
		{float64(12), 12},
		{int(13), 13},
		{int64(14), 14},
		{"15", 15},
		{"no es numero", 0},
		{nil, 0},
		{true, 0},
	}

	for _, c := range casos {
		if got := invToUint(c.entrada); got != c.esperado {
			t.Fatalf("con %v (%T) se esperaba %d y llego %d", c.entrada, c.entrada, c.esperado, got)
		}
	}
}

func TestLaConfiguracionDeInventarioSeLeeDelConfig(t *testing.T) {
	cfg := parseInventoryConfig(map[string]interface{}{
		"inventory_sync_enabled":        true,
		"inventory_single_warehouse_id": float64(9),
	})
	if !cfg.Enabled || cfg.SingleWarehouseID != 9 {
		t.Fatalf("configuracion mal leida: %+v", cfg)
	}
	if ids := resolveWarehouseIDs(cfg); len(ids) != 1 || ids[0] != 9 {
		t.Fatalf("debe filtrarse por la bodega elegida: %+v", ids)
	}

	vacia := parseInventoryConfig(map[string]interface{}{})
	if resolveWarehouseIDs(vacia) != nil {
		t.Fatal("sin bodega elegida se suma el stock de todas, no se filtra")
	}
}

func TestElConstructorDejaElCasoDeUsoUsable(t *testing.T) {
	uc := New(&fakeClient{}, &fakeService{}, nil, &fakeRepo{}, nil, log.New())
	if uc == nil {
		t.Fatal("el constructor no devolvio el caso de uso")
	}
}

func TestElExternalIDInvalidoNoSeInterpretaAMedias(t *testing.T) {
	if _, _, err := parseExternalProductID("abc"); err == nil {
		t.Fatal("un external_id que no es numerico debe fallar")
	}
	if _, _, err := parseExternalProductID("2001:xyz"); err == nil {
		t.Fatal("una variante que no es numerica debe fallar")
	}
	p, v, err := parseExternalProductID("2001:9001")
	if err != nil || p != 2001 || v != 9001 {
		t.Fatalf("producto:variante mal parseado: %d %d %v", p, v, err)
	}
	p, v, err = parseExternalProductID("2001")
	if err != nil || p != 2001 || v != 0 {
		t.Fatalf("solo producto mal parseado: %d %d %v", p, v, err)
	}
}

func TestSiElCanalNoRespondeLaComparacionFallaVisible(t *testing.T) {
	uc := newTestUseCase(&fakeClient{listaErr: errTiendanubeTest("Tiendanube caido")}, &fakeRepo{})

	if _, err := uc.ReconcileProducts(context.Background(), testIntegrationID, testBusinessID); err == nil {
		t.Fatal("si no se puede leer el catalogo del canal no se puede comparar: debe fallar")
	}
}

func TestSiElCanalNoRespondeElComparativoDeInventarioFallaVisible(t *testing.T) {
	repo := &fakeRepo{
		mapped: []domain.MappedItem{{ProductID: "p1", SKU: "SKU-1", ExternalItemID: "2001:9001"}},
		stock:  map[string]int{"p1": 1},
	}
	uc := newTestUseCase(&fakeClient{stockErr: errTiendanubeTest("Tiendanube caido")}, repo)

	if _, err := uc.CompareInventory(context.Background(), testIntegrationID, testBusinessID, 1, 100); err == nil {
		t.Fatal("sin el stock del canal no hay comparativo: debe fallar en vez de mostrar ceros")
	}
}

func TestUnFalloAlCrearEnElCanalNoDetieneElResto(t *testing.T) {
	repo := &fakeRepo{probability: []domain.ProductForSync{
		{ID: "p1", SKU: "FALLA", Name: "Uno", Price: 10},
		{ID: "p2", SKU: "PASA", Name: "Dos", Price: 20},
	}}
	client := &fakeClient{fallaCrearSKU: "FALLA"}

	uc := newTestUseCase(client, repo)
	if err := uc.ApplyProductsToTiendanube(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("un producto que falla no debe abortar el lote: %v", err)
	}
	if len(repo.mappings) != 1 || repo.mappings[0].SKU != "PASA" {
		t.Fatalf("el que si se creo debe quedar mapeado: %+v", repo.mappings)
	}
}

func TestLosProductosSinSKUNoEntranAlCruce(t *testing.T) {
	repo := &fakeRepo{probability: []domain.ProductForSync{
		{ID: "p1", SKU: "", Name: "Sin sku"},
		{ID: "p2", SKU: "CON-SKU", Name: "Con sku"},
	}}
	client := &fakeClient{}

	uc := newTestUseCase(client, repo)
	if err := uc.ApplyProductsToTiendanube(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("ApplyProductsToTiendanube fallo: %v", err)
	}
	if len(client.created) != 1 || client.created[0].SKU != "CON-SKU" {
		t.Fatalf("un producto sin SKU no se puede cruzar ni crear: %+v", client.created)
	}
}
