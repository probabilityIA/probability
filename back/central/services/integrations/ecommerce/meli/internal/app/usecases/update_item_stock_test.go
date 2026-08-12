package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

func pushSetup(state *domain.PushState) (*mockClient, *mockService, *mockInventoryRepo, *meliUseCase) {
	cli := &mockClient{}
	svc := &mockService{
		integration: &domain.Integration{ID: 254, Config: futureToken()},
		token:       "tok",
	}
	repo := &mockInventoryRepo{pushState: state}
	return cli, svc, repo, newUseCase(cli, svc, repo)
}

func TestUpdateItemStockEmpujaYRegistraLaCantidad(t *testing.T) {
	cli, _, repo, uc := pushSetup(&domain.PushState{})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	puts := cli.puts()
	if len(puts) != 1 {
		t.Fatalf("esperaba 1 llamada a MercadoLibre, hubo %d", len(puts))
	}
	if puts[0].ItemID != "MCO1" || puts[0].VariantID != "10" || puts[0].Quantity != 7 {
		t.Fatalf("payload equivocado: %+v", puts[0])
	}

	marks := repo.markedQty()
	if len(marks) != 1 {
		t.Fatalf("debia registrar la cantidad empujada, hubo %d registros", len(marks))
	}
	if marks[0].ProductID != "PRD_1" || marks[0].VariantID != "10" || marks[0].Quantity != 7 {
		t.Fatalf("registro equivocado: %+v", marks[0])
	}
}

func TestUpdateItemStockNoLlamaSiLaCantidadNoCambio(t *testing.T) {
	cli, _, repo, uc := pushSetup(&domain.PushState{LastPushedQty: intPtr(7)})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 0 {
		t.Fatal("no debia llamar a MercadoLibre si la cantidad ya es la ultima empujada")
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("no debia reescribir el registro si no hubo push")
	}
}

func TestUpdateItemStockSiLlamaCuandoLaCantidadCambio(t *testing.T) {
	cli, _, _, uc := pushSetup(&domain.PushState{LastPushedQty: intPtr(7)})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 3); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 1 {
		t.Fatal("una cantidad distinta si debe empujarse")
	}
}

func TestUpdateItemStockDistingueCeroDeNuncaEmpujado(t *testing.T) {
	cli, _, _, uc := pushSetup(&domain.PushState{LastPushedQty: nil})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 0); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 1 {
		t.Fatal("sin registro previo hay que empujar, aunque la cantidad sea 0")
	}
}

func TestUpdateItemStockNoEmpujaSiEsFulfillment(t *testing.T) {
	cli, _, repo, uc := pushSetup(&domain.PushState{LogisticType: domain.LogisticFulfillment})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != nil {
		t.Fatalf("fulfillment se omite, no es error: %v", err)
	}
	if len(cli.puts()) != 0 {
		t.Fatal("en fulfillment el stock lo administra MercadoLibre, no debia llamarse")
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("no debia registrar nada")
	}
}

func TestUpdateItemStockConsultaElEstadoDeLaVarianteCorrecta(t *testing.T) {
	_, _, repo, uc := pushSetup(&domain.PushState{})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "99", 4); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(repo.stateQueries) != 1 || repo.stateQueries[0] != "PRD_1|99" {
		t.Fatalf("debia consultar por producto y variante, consulto %v", repo.stateQueries)
	}
}

func TestUpdateItemStockSinProductIDEmpujaIgual(t *testing.T) {
	cli, _, repo, uc := pushSetup(&domain.PushState{LastPushedQty: intPtr(7)})

	if err := uc.UpdateItemStock(context.Background(), "254", "", "MCO1", "10", 7); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 1 {
		t.Fatal("sin product_id no se puede comparar, hay que empujar por seguridad")
	}
	if len(repo.stateQueries) != 0 {
		t.Fatal("sin product_id no tiene sentido consultar el estado")
	}
}

func TestUpdateItemStockEmpujaAunqueFalleLaConsultaDeEstado(t *testing.T) {
	cli := &mockClient{}
	svc := &mockService{integration: &domain.Integration{ID: 254, Config: futureToken()}, token: "tok"}
	repo := &mockInventoryRepo{pushStateErr: errors.New("base caida")}
	uc := newUseCase(cli, svc, repo)

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != nil {
		t.Fatalf("un fallo leyendo el estado no debe frenar el push: %v", err)
	}
	if len(cli.puts()) != 1 {
		t.Fatal("ante la duda hay que empujar")
	}
}

func TestUpdateItemStockRespetaElSyncDesactivado(t *testing.T) {
	cli := &mockClient{}
	svc := &mockService{
		integration: &domain.Integration{ID: 254, Config: map[string]interface{}{"inventory_sync_enabled": false}},
	}
	repo := &mockInventoryRepo{}
	uc := newUseCase(cli, svc, repo)

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 0 {
		t.Fatal("con el sync desactivado no se empuja")
	}
}

func TestUpdateItemStockConIntegracionInexistenteNoHaceNada(t *testing.T) {
	cli := &mockClient{}
	uc := newUseCase(cli, &mockService{integration: nil}, &mockInventoryRepo{})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(cli.puts()) != 0 {
		t.Fatal("sin integracion no hay a donde empujar")
	}
}

func TestUpdateItemStockPropagaElErrorDeLaIntegracion(t *testing.T) {
	fallo := errors.New("no se pudo leer la integracion")
	uc := newUseCase(&mockClient{}, &mockService{getErr: fallo}, &mockInventoryRepo{})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != fallo {
		t.Fatalf("debia propagar el error, obtuve %v", err)
	}
}

func TestUpdateItemStockReintentaUnaVezConTokenVencido(t *testing.T) {
	cli := &mockClient{putErrs: []error{domain.ErrTokenExpired}}
	svc := &mockService{integration: &domain.Integration{ID: 254, Config: futureToken()}, token: "tok"}
	repo := &mockInventoryRepo{pushState: &domain.PushState{}}
	uc := newUseCase(cli, svc, repo)

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != nil {
		t.Fatalf("tras refrescar el token debia salir bien: %v", err)
	}
	if len(cli.puts()) != 2 {
		t.Fatalf("esperaba el intento fallido y el reintento, hubo %d", len(cli.puts()))
	}
	if len(repo.markedQty()) != 1 {
		t.Fatal("el reintento exitoso tambien debe registrar la cantidad")
	}
}

func TestUpdateItemStockPropagaUnErrorQueNoEsDeToken(t *testing.T) {
	fallo := errors.New("meli client: update stock status 500")
	cli := &mockClient{putErrs: []error{fallo}}
	svc := &mockService{integration: &domain.Integration{ID: 254, Config: futureToken()}, token: "tok"}
	repo := &mockInventoryRepo{pushState: &domain.PushState{}}
	uc := newUseCase(cli, svc, repo)

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != fallo {
		t.Fatalf("debia propagar el error del canal, obtuve %v", err)
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("un push fallido no debe registrarse como empujado")
	}
}

func TestUpdateItemStockNoRegistraSiElReintentoTambienFalla(t *testing.T) {
	fallo := errors.New("sigue fallando")
	cli := &mockClient{putErrs: []error{domain.ErrTokenExpired, fallo}}
	svc := &mockService{integration: &domain.Integration{ID: 254, Config: futureToken()}, token: "tok"}
	repo := &mockInventoryRepo{pushState: &domain.PushState{}}
	uc := newUseCase(cli, svc, repo)

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != fallo {
		t.Fatalf("debia propagar el error del reintento, obtuve %v", err)
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("no debe registrarse una cantidad que nunca llego")
	}
}

func TestUpdateItemStockSigueAunqueFalleElRegistro(t *testing.T) {
	cli := &mockClient{}
	svc := &mockService{integration: &domain.Integration{ID: 254, Config: futureToken()}, token: "tok"}
	repo := &mockInventoryRepo{pushState: &domain.PushState{}, markErr: errors.New("no se pudo escribir")}
	uc := newUseCase(cli, svc, repo)

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err != nil {
		t.Fatalf("el stock ya se empujo, un fallo registrando no debe reportarse como error: %v", err)
	}
	if len(cli.puts()) != 1 {
		t.Fatal("el push debia haberse hecho")
	}
}

func TestUpdateItemStockEnPublicacionSimpleNoMandaVariante(t *testing.T) {
	cli, _, _, uc := pushSetup(&domain.PushState{})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "", 5); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if cli.puts()[0].VariantID != "" {
		t.Fatalf("una publicacion sin variantes no debe llevar variation_id, llevo %q", cli.puts()[0].VariantID)
	}
}

func TestUpdateItemStockFallaSiNoHayTokenUtilizable(t *testing.T) {
	svc := &mockService{
		integration: &domain.Integration{ID: 254, Config: futureToken()},
		decryptErr:  errors.New("no se pudo desencriptar"),
	}
	cli := &mockClient{}
	uc := newUseCase(cli, svc, &mockInventoryRepo{pushState: &domain.PushState{}})

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err == nil {
		t.Fatal("sin token no se puede empujar y hay que reportarlo")
	}
	if len(cli.puts()) != 0 {
		t.Fatal("no debia intentar el push")
	}
}

func TestUpdateItemStockFallaSiElRefrescoDeTokenNoSirve(t *testing.T) {
	svc := &mockService{integration: &domain.Integration{ID: 254, Config: futureToken()}, token: "tok"}
	cli := &mockClient{putErrs: []error{domain.ErrTokenExpired}}
	repo := &mockInventoryRepo{pushState: &domain.PushState{}}
	uc := newUseCase(cli, svc, repo)

	svcFalla := &mockService{integration: svc.integration, decryptErr: errors.New("credencial ilegible")}
	uc.service = svcFalla

	if err := uc.UpdateItemStock(context.Background(), "254", "PRD_1", "MCO1", "10", 7); err == nil {
		t.Fatal("si el refresco falla hay que propagarlo")
	}
	if len(repo.markedQty()) != 0 {
		t.Fatal("no debe registrarse nada")
	}
}
