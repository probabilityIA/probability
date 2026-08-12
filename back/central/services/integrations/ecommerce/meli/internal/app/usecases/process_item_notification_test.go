package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func notifSetup(mapped []domain.MappedItem, stock map[string]int, item *domain.MeliItemDetail) (*mockClient, *mockInventoryRepo, *meliUseCase, *domain.Integration) {
	cli := &mockClient{item: item}
	integ := &domain.Integration{ID: 254, Config: futureToken()}
	svc := &mockService{integration: integ, token: "tok"}
	repo := &mockInventoryRepo{mapped: mapped, stock: stock}
	return cli, repo, newUseCase(cli, svc, repo), integ
}

func TestWebhookRefrescaElTipoLogistico(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1", ExternalVariantID: "10"}}
	item := &domain.MeliItemDetail{
		ID:           "MCO1",
		LogisticType: domain.LogisticFulfillment,
		Variations:   []domain.MeliItemVariation{{ID: 10, AvailableQuantity: 3}},
	}
	_, repo, uc, integ := notifSetup(mapped, map[string]int{"P1": 3}, item)

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(repo.logisticSets) != 1 || repo.logisticSets[0] != "MCO1=fulfillment" {
		t.Fatalf("debia guardar el tipo logistico que trajo el webhook, obtuve %v", repo.logisticSets)
	}
}

func TestWebhookCorrigeLaDesviacionYRegistra(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", SKU: "A", ExternalItemID: "MCO1", ExternalVariantID: "10"}}
	item := &domain.MeliItemDetail{
		ID:         "MCO1",
		Variations: []domain.MeliItemVariation{{ID: 10, AvailableQuantity: 50}},
	}
	cli, repo, uc, integ := notifSetup(mapped, map[string]int{"P1": 13}, item)

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	puts := cli.puts()
	if len(puts) != 1 || puts[0].Quantity != 13 {
		t.Fatalf("si el canal tiene 50 y aqui hay 13, hay que corregirlo, obtuve %+v", puts)
	}
	if len(repo.markedQty()) != 1 || repo.markedQty()[0].Quantity != 13 {
		t.Fatalf("la correccion tambien debe quedar registrada, obtuve %+v", repo.markedQty())
	}
}

func TestWebhookNoTocaNadaSiLaVarianteYaCoincide(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1", ExternalVariantID: "10"}}
	item := &domain.MeliItemDetail{
		ID:         "MCO1",
		Variations: []domain.MeliItemVariation{{ID: 10, AvailableQuantity: 13}},
	}
	cli, repo, uc, integ := notifSetup(mapped, map[string]int{"P1": 13}, item)

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 0 {
		t.Fatal("si el canal ya tiene la cantidad correcta no hay nada que hacer")
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("sin push no hay registro")
	}
}

func TestWebhookIgnoraUnItemQueNoTenemosMapeado(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "OTRO", ExternalVariantID: "10"}}
	cli, _, uc, integ := notifSetup(mapped, map[string]int{"P1": 3}, &domain.MeliItemDetail{ID: "MCO1"})

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 0 {
		t.Fatal("un item sin mapeo local no se toca")
	}
}

func TestWebhookRespetaElSyncDesactivado(t *testing.T) {
	cli := &mockClient{}
	integ := &domain.Integration{ID: 254, Config: map[string]interface{}{"inventory_sync_enabled": false}}
	repo := &mockInventoryRepo{mapped: []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1"}}}
	uc := newUseCase(cli, &mockService{integration: integ}, repo)

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 0 {
		t.Fatal("con el sync desactivado el webhook no empuja")
	}
}

func TestWebhookPropagaElErrorDeMapeos(t *testing.T) {
	integ := &domain.Integration{ID: 254, Config: futureToken()}
	repo := &mockInventoryRepo{mappedErr: errors.New("base caida")}
	uc := newUseCase(&mockClient{}, &mockService{integration: integ}, repo)

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err == nil {
		t.Fatal("debia propagar el error")
	}
}

func TestWebhookEmpujaAunqueNoPuedaLeerLaPublicacion(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1", ExternalVariantID: "10"}}
	cli, _, uc, integ := notifSetup(mapped, map[string]int{"P1": 3}, nil)
	cli.getErr = errors.New("no se pudo leer el item")

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 1 {
		t.Fatal("sin poder comparar hay que empujar por seguridad")
	}
}

func TestWebhookAtiendeTodasLasTallasDeLaPublicacion(t *testing.T) {
	mapped := []domain.MappedItem{
		{ProductID: "P1", ExternalItemID: "MCO1", ExternalVariantID: "10"},
		{ProductID: "P2", ExternalItemID: "MCO1", ExternalVariantID: "11"},
	}
	item := &domain.MeliItemDetail{
		ID: "MCO1",
		Variations: []domain.MeliItemVariation{
			{ID: 10, AvailableQuantity: 0},
			{ID: 11, AvailableQuantity: 0},
		},
	}
	cli, _, uc, integ := notifSetup(mapped, map[string]int{"P1": 4, "P2": 6}, item)

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 2 {
		t.Fatalf("debia corregir las dos tallas, obtuvo %d", len(cli.puts()))
	}
}

func TestWebhookNoTocaLaPublicacionSimpleQueYaCoincide(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1"}}
	item := &domain.MeliItemDetail{ID: "MCO1", AvailableQuantity: 13}
	cli, _, uc, integ := notifSetup(mapped, map[string]int{"P1": 13}, item)

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 0 {
		t.Fatal("una publicacion sin variantes que ya coincide no se toca")
	}
}

func TestWebhookReintentaTrasRefrescarElToken(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1", ExternalVariantID: "10"}}
	item := &domain.MeliItemDetail{ID: "MCO1", Variations: []domain.MeliItemVariation{{ID: 10, AvailableQuantity: 0}}}
	cli, repo, uc, integ := notifSetup(mapped, map[string]int{"P1": 4}, item)
	cli.putErrs = []error{domain.ErrTokenExpired}

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("tras refrescar debia salir bien: %v", err)
	}
	if len(cli.puts()) != 2 {
		t.Fatalf("esperaba el intento fallido y el reintento, hubo %d", len(cli.puts()))
	}
	if len(repo.markedQty()) != 1 {
		t.Fatal("el reintento exitoso debe registrar la cantidad")
	}
}

func TestWebhookPropagaSiElReintentoTambienFalla(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1", ExternalVariantID: "10"}}
	item := &domain.MeliItemDetail{ID: "MCO1", Variations: []domain.MeliItemVariation{{ID: 10, AvailableQuantity: 0}}}
	cli, repo, uc, integ := notifSetup(mapped, map[string]int{"P1": 4}, item)
	cli.putErrs = []error{domain.ErrTokenExpired, errors.New("sigue fallando")}

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err == nil {
		t.Fatal("debia propagar el error del reintento")
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("nada llego al canal")
	}
}

func TestWebhookPropagaUnErrorQueNoEsDeToken(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1", ExternalVariantID: "10"}}
	item := &domain.MeliItemDetail{ID: "MCO1", Variations: []domain.MeliItemVariation{{ID: 10, AvailableQuantity: 0}}}
	cli, _, uc, integ := notifSetup(mapped, map[string]int{"P1": 4}, item)
	cli.putErrs = []error{errors.New("meli client: update stock status 500")}

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err == nil {
		t.Fatal("debia propagar el error del canal")
	}
}

func TestWebhookSigueAunqueFalleElRegistroDeLaCantidad(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1", ExternalVariantID: "10"}}
	item := &domain.MeliItemDetail{ID: "MCO1", Variations: []domain.MeliItemVariation{{ID: 10, AvailableQuantity: 0}}}
	cli, repo, uc, integ := notifSetup(mapped, map[string]int{"P1": 4}, item)
	repo.markErr = errors.New("no se pudo escribir")

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("el stock ya se corrigio, el registro es secundario: %v", err)
	}
	if len(cli.puts()) != 1 {
		t.Fatal("la correccion debia hacerse")
	}
}

func TestWebhookSigueAunqueFalleGuardarElTipoLogistico(t *testing.T) {
	mapped := []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1", ExternalVariantID: "10"}}
	item := &domain.MeliItemDetail{
		ID:           "MCO1",
		LogisticType: "xd_drop_off",
		Variations:   []domain.MeliItemVariation{{ID: 10, AvailableQuantity: 0}},
	}
	cli, repo, uc, integ := notifSetup(mapped, map[string]int{"P1": 4}, item)
	repo.logisticErr = errors.New("no se pudo escribir")

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err != nil {
		t.Fatalf("guardar el tipo logistico es secundario, el stock manda: %v", err)
	}
	if len(cli.puts()) != 1 {
		t.Fatal("el stock debia corregirse igual")
	}
}

func TestWebhookAbortaSiNoPuedeLeerElStock(t *testing.T) {
	integ := &domain.Integration{ID: 254, Config: futureToken()}
	cli := &mockClient{}
	repo := &mockInventoryRepo{
		mapped:   []domain.MappedItem{{ProductID: "P1", ExternalItemID: "MCO1"}},
		stockErr: errors.New("base caida"),
	}
	uc := newUseCase(cli, &mockService{integration: integ, token: "tok"}, repo)

	if err := uc.reconcileItemStock(context.Background(), integ, "tok", "MCO1"); err == nil {
		t.Fatal("sin stock no se puede decidir que empujar")
	}
	if len(cli.puts()) != 0 {
		t.Fatal("no debia empujar a ciegas")
	}
}
