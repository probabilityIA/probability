package usecases

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func TestTraerLoQueSoloEstaEnElCanalPublicaALaColaDeUpsert(t *testing.T) {
	repo := &fakeRepo{probability: []domain.ProductForSync{{ID: "p1", SKU: "COMUN-1", Name: "Comun"}}}
	client := &fakeClient{products: []domain.TiendanubeProduct{
		productoCanal(2001, 9001, "COMUN-1", "Comun", 100, 5, true),
		productoCanal(2002, 9002, "SOLO-CANAL", "Solo en el canal", 250, 3, true),
	}}
	cola := nuevaCola()

	uc := newTestUseCaseConCola(client, repo, cola)
	if err := uc.ApplyProductsToProbability(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("ApplyProductsToProbability fallo: %v", err)
	}

	msgs := cola.mensajes(rabbitmq.QueueProductsProviderUpsert)
	if len(msgs) != 1 || msgs[0].SKU != "SOLO-CANAL" {
		t.Fatalf("solo debe publicarse lo que falta en Probability: %+v", msgs)
	}
	if msgs[0].IntegrationID == 0 {
		t.Fatal("el mensaje debe llevar integration_id: sin el se pierde la relacion producto-canal")
	}
	if msgs[0].ExternalID != "2002:9002" {
		t.Fatalf("el external_id debe traer producto:variante, llego %q", msgs[0].ExternalID)
	}
	if msgs[0].Price != 250 {
		t.Fatalf("el precio del canal no viajo: %v", msgs[0].Price)
	}
}

func TestActualizarHaciaProbabilityPublicaLoQueYaCoincide(t *testing.T) {
	repo := &fakeRepo{probability: []domain.ProductForSync{{ID: "p1", SKU: "COMUN-1", Name: "Viejo"}}}
	client := &fakeClient{products: []domain.TiendanubeProduct{
		productoCanal(2001, 9001, "COMUN-1", "Nombre nuevo del canal", 999, 5, true),
	}}
	cola := nuevaCola()

	uc := newTestUseCaseConCola(client, repo, cola)
	if err := uc.UpdateProductsToProbability(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("UpdateProductsToProbability fallo: %v", err)
	}

	msgs := cola.mensajes(rabbitmq.QueueProductsProviderUpsert)
	if len(msgs) != 1 || msgs[0].Name != "Nombre nuevo del canal" {
		t.Fatalf("debe publicarse el dato del canal para el SKU cruzado: %+v", msgs)
	}
}

func TestSinRabbitNoSePierdenProductosEnSilencio(t *testing.T) {
	repo := &fakeRepo{}
	client := &fakeClient{products: []domain.TiendanubeProduct{
		productoCanal(2002, 9002, "SOLO-CANAL", "Solo en el canal", 250, 3, true),
	}}

	uc := newTestUseCase(client, repo)
	err := uc.ApplyProductsToProbability(context.Background(), testIntegrationID, testBusinessID, "corr-1")
	if err == nil {
		t.Fatal("sin RabbitMQ la operacion debe fallar visible, nunca aparentar exito")
	}
}

func TestAsociarSoloTocaLoCruzadoYNoLoYaAsociado(t *testing.T) {
	repo := &fakeRepo{
		probability: []domain.ProductForSync{
			{ID: "p1", SKU: "YA-ASOCIADO", Name: "Uno"},
			{ID: "p2", SKU: "SIN-ASOCIAR", Name: "Dos"},
		},
		mapped: []domain.MappedItem{{ProductID: "p1", SKU: "YA-ASOCIADO", ExternalItemID: "2001:9001"}},
	}
	client := &fakeClient{products: []domain.TiendanubeProduct{
		productoCanal(2001, 9001, "YA-ASOCIADO", "Uno", 100, 5, true),
		productoCanal(2002, 9002, "SIN-ASOCIAR", "Dos", 200, 5, true),
	}}

	uc := newTestUseCase(client, repo)
	if err := uc.AssociateProducts(context.Background(), testIntegrationID, testBusinessID, "corr-1", nil); err != nil {
		t.Fatalf("AssociateProducts fallo: %v", err)
	}

	if len(repo.mappings) != 1 || repo.mappings[0].SKU != "SIN-ASOCIAR" {
		t.Fatalf("solo debe asociarse lo que aun no lo estaba: %+v", repo.mappings)
	}
}

func TestAsociarConSeleccionExplicitaIgnoraLoYaAsociado(t *testing.T) {
	repo := &fakeRepo{
		probability: []domain.ProductForSync{{ID: "p1", SKU: "YA-ASOCIADO", Name: "Uno"}},
		mapped:      []domain.MappedItem{{ProductID: "p1", SKU: "YA-ASOCIADO", ExternalItemID: "2001:9001"}},
	}
	client := &fakeClient{products: []domain.TiendanubeProduct{
		productoCanal(2001, 9001, "YA-ASOCIADO", "Uno", 100, 5, true),
	}}

	uc := newTestUseCase(client, repo)
	if err := uc.AssociateProducts(context.Background(), testIntegrationID, testBusinessID, "corr-1", []string{"ya-asociado"}); err != nil {
		t.Fatalf("AssociateProducts fallo: %v", err)
	}
	if len(repo.mappings) != 1 {
		t.Fatalf("elegir un SKU explicitamente debe re-asociarlo: %+v", repo.mappings)
	}
}

func TestActualizarHaciaElCanalTocaProductoYVariante(t *testing.T) {
	repo := &fakeRepo{probability: []domain.ProductForSync{
		{ID: "p1", SKU: "COMUN-1", Name: "Nombre en Probability", Price: 555},
	}}
	client := &fakeClient{products: []domain.TiendanubeProduct{
		productoCanal(2001, 9001, "COMUN-1", "Nombre viejo", 100, 5, true),
	}}

	uc := newTestUseCase(client, repo)
	if err := uc.UpdateProductsToTiendanube(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("UpdateProductsToTiendanube fallo: %v", err)
	}

	if len(client.actualizados) != 1 || client.actualizados[0].Name != "Nombre en Probability" {
		t.Fatalf("el producto del canal debe actualizarse con el nombre de Probability: %+v", client.actualizados)
	}
	if len(client.variantes) != 1 || client.variantes[0].Price == nil || *client.variantes[0].Price != 555 {
		t.Fatalf("el precio vive en la variante y debe actualizarse ahi: %+v", client.variantes)
	}
}

func TestElProductoDelCanalNoSeActualizaDosVecesPorSusVariantes(t *testing.T) {
	repo := &fakeRepo{probability: []domain.ProductForSync{
		{ID: "p1", SKU: "TALLA-S", Name: "Camiseta", Price: 100},
		{ID: "p2", SKU: "TALLA-M", Name: "Camiseta", Price: 110},
	}}
	client := &fakeClient{products: []domain.TiendanubeProduct{{
		ID:   2001,
		Name: "Camiseta",
		Variants: []domain.TiendanubeVariant{
			{ID: 9001, ProductID: 2001, SKU: "TALLA-S", Price: 1},
			{ID: 9002, ProductID: 2001, SKU: "TALLA-M", Price: 1},
		},
	}}}

	uc := newTestUseCase(client, repo)
	if err := uc.UpdateProductsToTiendanube(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("UpdateProductsToTiendanube fallo: %v", err)
	}

	if len(client.actualizados) != 1 {
		t.Fatalf("el producto padre se debe tocar una sola vez, se toco %d veces", len(client.actualizados))
	}
	if len(client.variantes) != 2 {
		t.Fatalf("cada variante si debe actualizarse: %d", len(client.variantes))
	}
}

func TestSincronizarInventarioEscribeElStockDeCadaVariante(t *testing.T) {
	repo := &fakeRepo{
		mapped: []domain.MappedItem{
			{ProductID: "p1", SKU: "SKU-1", ExternalItemID: "2001:9001"},
			{ProductID: "p2", SKU: "SKU-2", ExternalItemID: "2002:9002"},
		},
		stock: map[string]int{"p1": 4, "p2": 0},
	}
	client := &fakeClient{}

	uc := newTestUseCase(client, repo)
	if err := uc.SyncInventory(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("SyncInventory fallo: %v", err)
	}

	if len(client.stockEscrito) != 2 {
		t.Fatalf("se esperaba escribir stock de dos variantes: %+v", client.stockEscrito)
	}
	if client.stockEscrito["2001:9001"] != 4 {
		t.Fatalf("el stock de Probability no viajo al canal: %+v", client.stockEscrito)
	}
	if _, ok := client.stockEscrito["2002:9002"]; !ok {
		t.Fatal("un producto en cero tambien debe empujarse, si no el canal sigue vendiendo lo que no hay")
	}
}

func TestSincronizarInventarioResuelveLaVarianteCuandoElMapeoNoLaTiene(t *testing.T) {
	repo := &fakeRepo{
		mapped: []domain.MappedItem{{ProductID: "p1", SKU: "SKU-1", ExternalItemID: "2001"}},
		stock:  map[string]int{"p1": 7},
	}
	client := &fakeClient{objetivo: &domain.StockTarget{ProductID: 2001, VariantID: 9001, Found: true}}

	uc := newTestUseCase(client, repo)
	if err := uc.SyncInventory(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("SyncInventory fallo: %v", err)
	}
	if client.stockEscrito["2001:9001"] != 7 {
		t.Fatalf("debio resolverse la variante por SKU antes de escribir: %+v", client.stockEscrito)
	}
}

func TestSincronizarInventarioNoRompeSiNoEncuentraLaVariante(t *testing.T) {
	repo := &fakeRepo{
		mapped: []domain.MappedItem{{ProductID: "p1", SKU: "SKU-1", ExternalItemID: "2001"}},
		stock:  map[string]int{"p1": 7},
	}
	client := &fakeClient{objetivo: &domain.StockTarget{Found: false}}

	uc := newTestUseCase(client, repo)
	if err := uc.SyncInventory(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("una variante que no se encuentra no debe tumbar toda la sincronizacion: %v", err)
	}
	if len(client.stockEscrito) != 0 {
		t.Fatalf("no se debe escribir stock a ciegas: %+v", client.stockEscrito)
	}
}

func TestEmpujarStockExigeVariante(t *testing.T) {
	repo := &fakeRepo{}
	client := &fakeClient{}
	uc := newTestUseCase(client, repo)

	if err := uc.UpdateInventory(context.Background(), testIntegrationID, "2001", 5); err == nil {
		t.Fatal("sin variante no se puede escribir stock en Tiendanube: debe fallar explicito")
	}

	if err := uc.UpdateInventory(context.Background(), testIntegrationID, "2001:9001", 5); err != nil {
		t.Fatalf("con producto:variante debe funcionar: %v", err)
	}
	if client.stockEscrito["2001:9001"] != 5 {
		t.Fatalf("el stock no se escribio: %+v", client.stockEscrito)
	}
}

func TestElResumenDeLaComparacionSeGuardaConSusGrupos(t *testing.T) {
	repo := &fakeRepo{
		probability: []domain.ProductForSync{
			{ID: "p1", SKU: "COMUN-1", Name: "Comun"},
			{ID: "p2", SKU: "SOLO-PROB", Name: "Solo aca"},
		},
	}
	client := &fakeClient{products: []domain.TiendanubeProduct{
		productoCanal(2001, 9001, "COMUN-1", "Comun", 100, 5, true),
		productoCanal(2002, 9002, "SOLO-CANAL", "Solo alla", 200, 5, true),
	}}
	cola := nuevaCola()

	uc := newTestUseCaseConCola(client, repo, cola)
	uc.ReconcileProductsAsync(context.Background(), testIntegrationID, testBusinessID, 77, "corr-1")

	guardados := cola.publicados[rabbitmq.QueueIntegrationSyncRuns]
	if len(guardados) != 1 {
		t.Fatalf("la comparacion debe dejar un resultado guardado: %d", len(guardados))
	}

	resultado, err := uc.ReconcileProducts(context.Background(), testIntegrationID, testBusinessID)
	if err != nil {
		t.Fatalf("ReconcileProducts fallo: %v", err)
	}
	detalle := reconcileDetail(resultado)
	grupos := map[string]int{}
	for _, d := range detalle {
		grupos[d.Group]++
	}
	if grupos["only_probability"] != 1 || grupos["only_channel"] != 1 || grupos["both"] != 1 {
		t.Fatalf("el detalle debe separar los tres grupos: %+v", grupos)
	}
}

func TestSincronizarProductosHaceLosDosSentidos(t *testing.T) {
	repo := &fakeRepo{probability: []domain.ProductForSync{{ID: "p1", SKU: "SOLO-PROB", Name: "Uno", Price: 10}}}
	client := &fakeClient{products: []domain.TiendanubeProduct{
		productoCanal(2002, 9002, "SOLO-CANAL", "Dos", 20, 1, true),
	}}
	cola := nuevaCola()

	uc := newTestUseCaseConCola(client, repo, cola)
	if err := uc.SyncProducts(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err != nil {
		t.Fatalf("SyncProducts fallo: %v", err)
	}

	if len(client.created) != 1 || client.created[0].SKU != "SOLO-PROB" {
		t.Fatalf("lo que falta en el canal debe crearse alla: %+v", client.created)
	}
	msgs := cola.mensajes(rabbitmq.QueueProductsProviderUpsert)
	if len(msgs) != 1 || msgs[0].SKU != "SOLO-CANAL" {
		t.Fatalf("lo que falta en Probability debe publicarse a la cola: %+v", msgs)
	}
}
