package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func syncSetup(mapped []domain.MappedItem, stock map[string]int) (*mockClient, *mockInventoryRepo, *meliUseCase) {
	cli := &mockClient{}
	svc := &mockService{integration: &domain.Integration{ID: 254, Config: futureToken()}, token: "tok"}
	repo := &mockInventoryRepo{mapped: mapped, stock: stock}
	return cli, repo, newUseCase(cli, svc, repo)
}

func runSync(t *testing.T, uc *meliUseCase, mapped []domain.MappedItem, stock map[string]int) {
	t.Helper()
	cfg := domain.InventoryConfig{Mode: "single", SingleWarehouseID: 22, Enabled: true}
	if err := uc.syncInventorySingle(context.Background(), uc.client, 46, 254, "254", "tok", cfg, mapped, "corr", false); err != nil {
		t.Fatalf("error inesperado en el sync: %v", err)
	}
}

func TestSyncEmpujaLoQueCambioYRegistra(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10"},
	}
	stock := map[string]int{"P1": 5}
	cli, repo, uc := syncSetup(mapped, stock)

	runSync(t, uc, mapped, stock)

	puts := cli.puts()
	if len(puts) != 1 || puts[0].VariantID != "10" || puts[0].Quantity != 5 {
		t.Fatalf("esperaba un push de 5 a la variante 10, obtuve %+v", puts)
	}
	marks := repo.markedQty()
	if len(marks) != 1 || marks[0].VariantID != "10" || marks[0].Quantity != 5 {
		t.Fatalf("debia registrar por variante, obtuve %+v", marks)
	}
}

func TestSyncSaltaLoQueYaSeHabiaEmpujado(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10", LastPushedQty: intPtr(5)},
	}
	stock := map[string]int{"P1": 5}
	cli, repo, uc := syncSetup(mapped, stock)

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 0 {
		t.Fatal("si la cantidad es la misma que la ultima empujada no debe haber llamada")
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("no debia reescribir el registro")
	}
}

func TestSyncEmpujaSiLaCantidadCambioDesdeElUltimoPush(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10", LastPushedQty: intPtr(5)},
	}
	stock := map[string]int{"P1": 2}
	cli, _, uc := syncSetup(mapped, stock)

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 1 || cli.puts()[0].Quantity != 2 {
		t.Fatalf("una cantidad distinta si debe empujarse, obtuve %+v", cli.puts())
	}
}

func TestSyncNoApagaUnProductoSinRegistroDeInventario(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "SIN-NIVEL", SKU: "SDH22-S", ExternalItemID: "MCO1", ExternalVariantID: "10"},
	}
	stock := map[string]int{}
	cli, repo, uc := syncSetup(mapped, stock)

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 0 {
		t.Fatal("sin fila en inventory_levels no se debe mandar 0: apagaria la publicacion")
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("no hubo push, no hay nada que registrar")
	}
}

func TestSyncNoConfundeCeroRealConAusenciaDeRegistro(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10"},
	}
	stock := map[string]int{"P1": 0}
	cli, _, uc := syncSetup(mapped, stock)

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 1 || cli.puts()[0].Quantity != 0 {
		t.Fatalf("un cero real si debe empujarse, obtuve %+v", cli.puts())
	}
}

func TestSyncNoTocaPublicacionesDeFulfillment(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10", LogisticType: domain.LogisticFulfillment},
	}
	stock := map[string]int{"P1": 5}
	cli, _, uc := syncSetup(mapped, stock)

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 0 {
		t.Fatal("en fulfillment el stock lo administra MercadoLibre")
	}
}

func TestSyncSaltaFulfillmentAntesDeMirarElInventario(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "SIN-NIVEL", SKU: "A", ExternalItemID: "MCO1", ExternalVariantID: "10", LogisticType: domain.LogisticFulfillment},
	}
	stock := map[string]int{}
	cli, _, uc := syncSetup(mapped, stock)

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 0 {
		t.Fatal("ninguno de los dos motivos debe terminar en llamada")
	}
}

func TestSyncProcesaCadaVarianteDeLaMismaPublicacion(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10"},
		{ProductID: "P2", SKU: "A-3XL", ExternalItemID: "MCO1", ExternalVariantID: "11"},
		{ProductID: "P3", SKU: "A-4XL", ExternalItemID: "MCO1", ExternalVariantID: "12"},
	}
	stock := map[string]int{"P1": 1, "P2": 2, "P3": 3}
	cli, _, uc := syncSetup(mapped, stock)

	runSync(t, uc, mapped, stock)

	puts := cli.puts()
	if len(puts) != 3 {
		t.Fatalf("cada variante se actualiza por separado, esperaba 3 y hubo %d", len(puts))
	}
	got := map[string]int{}
	for _, p := range puts {
		got[p.VariantID] = p.Quantity
	}
	if got["10"] != 1 || got["11"] != 2 || got["12"] != 3 {
		t.Fatalf("cada variante debe recibir su propia cantidad, obtuve %+v", got)
	}
}

func TestSyncMezclaSimplesYVariantes(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10"},
		{ProductID: "P2", SKU: "B", ExternalItemID: "MCO2"},
	}
	stock := map[string]int{"P1": 4, "P2": 9}
	cli, _, uc := syncSetup(mapped, stock)

	runSync(t, uc, mapped, stock)

	puts := cli.puts()
	if len(puts) != 2 {
		t.Fatalf("esperaba 2 llamadas, hubo %d", len(puts))
	}
	for _, p := range puts {
		if p.ItemID == "MCO2" && p.VariantID != "" {
			t.Fatal("la publicacion simple no debe llevar variation_id")
		}
	}
}

func TestSyncNoRegistraLaCantidadSiElPushFallo(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A-2XL", ExternalItemID: "MCO1", ExternalVariantID: "10"},
	}
	stock := map[string]int{"P1": 5}
	cli, repo, uc := syncSetup(mapped, stock)
	cli.putErrs = []error{errors.New("meli client: update stock status 400")}

	runSync(t, uc, mapped, stock)

	if len(repo.markedQty()) != 0 {
		t.Fatal("un push fallido no debe quedar registrado como empujado")
	}
}

func TestSyncPropagaElErrorAlLeerElStock(t *testing.T) {
	cli := &mockClient{}
	svc := &mockService{integration: &domain.Integration{ID: 254, Config: futureToken()}, token: "tok"}
	repo := &mockInventoryRepo{stockErr: errors.New("base caida")}
	uc := newUseCase(cli, svc, repo)

	cfg := domain.InventoryConfig{Mode: "single", SingleWarehouseID: 22}
	err := uc.syncInventorySingle(context.Background(), cli, 46, 254, "254", "tok", cfg,
		[]domain.MappedItem{{ProductID: "P1"}}, "corr", false)
	if err == nil {
		t.Fatal("si no se puede leer el stock hay que abortar, no empujar a ciegas")
	}
	if len(cli.puts()) != 0 {
		t.Fatal("no debia empujar nada")
	}
}

func TestSyncSinMapeosNoHaceNada(t *testing.T) {
	cli, _, uc := syncSetup(nil, map[string]int{})
	runSync(t, uc, nil, map[string]int{})

	if len(cli.puts()) != 0 {
		t.Fatal("sin mapeos no hay a que empujarle")
	}
}

func TestSyncReintentaTrasRefrescarElToken(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A", ExternalItemID: "MCO1", ExternalVariantID: "10"},
	}
	stock := map[string]int{"P1": 5}
	cli, repo, uc := syncSetup(mapped, stock)
	cli.putErrs = []error{domain.ErrTokenExpired}

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 2 {
		t.Fatalf("esperaba el intento fallido y el reintento, hubo %d", len(cli.puts()))
	}
	if len(repo.markedQty()) != 1 {
		t.Fatal("el reintento exitoso debe registrar la cantidad")
	}
}

func TestSyncCuentaComoFallidoSiElReintentoNoSirve(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A", ExternalItemID: "MCO1", ExternalVariantID: "10"},
	}
	stock := map[string]int{"P1": 5}
	cli, repo, uc := syncSetup(mapped, stock)
	cli.putErrs = []error{domain.ErrTokenExpired, errors.New("sigue fallando")}

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 2 {
		t.Fatalf("esperaba dos intentos, hubo %d", len(cli.puts()))
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("nada llego al canal, nada que registrar")
	}
}

func TestSyncSigueConLasDemasAunqueUnaFalle(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A", ExternalItemID: "MCO1", ExternalVariantID: "10"},
		{ProductID: "P2", SKU: "B", ExternalItemID: "MCO1", ExternalVariantID: "11"},
	}
	stock := map[string]int{"P1": 5, "P2": 6}
	cli, repo, uc := syncSetup(mapped, stock)
	cli.putErrs = []error{errors.New("400 en la primera")}

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 2 {
		t.Fatalf("un fallo no debe abortar la corrida, hubo %d intentos", len(cli.puts()))
	}
	marks := repo.markedQty()
	if len(marks) != 1 || marks[0].ProductID != "P2" {
		t.Fatalf("solo la que si llego debe registrarse, obtuve %+v", marks)
	}
}

func TestSyncSigueAunqueFalleElRegistroDeLaCantidad(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", SKU: "A", ExternalItemID: "MCO1", ExternalVariantID: "10"},
	}
	stock := map[string]int{"P1": 5}
	cli, repo, uc := syncSetup(mapped, stock)
	repo.markErr = errors.New("no se pudo escribir")

	runSync(t, uc, mapped, stock)

	if len(cli.puts()) != 1 {
		t.Fatal("el push debia hacerse igual")
	}
}

func TestSeleccionManualCorrigeElCanalACeroCuandoNoHayRegistro(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P9", SKU: "SIN-STOCK", ExternalItemID: "MCO9", ExternalVariantID: "90"},
	}
	cli, _, uc := syncSetup(mapped, map[string]int{})

	cfg := domain.InventoryConfig{Mode: "single", SingleWarehouseID: 22, Enabled: true}
	if err := uc.syncInventorySingle(context.Background(), cli, 46, 254, "254", "tok", cfg, mapped, "corr", true); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	puts := cli.puts()
	if len(puts) != 1 || puts[0].Quantity != 0 {
		t.Fatalf("al elegirlo a mano se debe corregir el canal a 0, obtuve %+v", puts)
	}
}

func TestElSyncAutomaticoNoTocaLoQueNoTieneRegistro(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P9", SKU: "SIN-STOCK", ExternalItemID: "MCO9", ExternalVariantID: "90"},
	}
	cli, _, uc := syncSetup(mapped, map[string]int{})

	runSync(t, uc, mapped, map[string]int{})

	if len(cli.puts()) != 0 {
		t.Fatal("sin seleccion manual no se debe apagar el stock del canal")
	}
}
