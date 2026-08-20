package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

func TestElStoreIDSaleDelConfigODeLaURLLegacy(t *testing.T) {
	casos := []struct {
		nombre   string
		config   map[string]interface{}
		fallback string
		esperado string
	}{
		{"config directo", map[string]interface{}{"store_id": "1234567"}, "", "1234567"},
		{"store_url legacy", map[string]interface{}{"store_url": "https://api.tiendanube.com/v1/987654"}, "", "987654"},
		{"store_url con barra final", map[string]interface{}{"store_url": "https://api.tiendanube.com/v1/555/"}, "", "555"},
		{"store_url sin id numerico", map[string]interface{}{"store_url": "https://mitienda.mitiendanube.com"}, "42", "42"},
		{"sin nada, cae al de la integracion", map[string]interface{}{}, "77", "77"},
		{"sin nada en absoluto", map[string]interface{}{}, "", ""},
	}

	for _, c := range casos {
		if got := resolveStoreID(c.config, c.fallback); got != c.esperado {
			t.Fatalf("%s: se esperaba %q y llego %q", c.nombre, c.esperado, got)
		}
	}
}

func TestExtraerTextoRechazaVacioYTipoEquivocado(t *testing.T) {
	if _, err := extractString(map[string]interface{}{}, "falta"); err == nil {
		t.Fatal("un campo ausente debe fallar")
	}
	if _, err := extractString(map[string]interface{}{"x": "   "}, "x"); err == nil {
		t.Fatal("un campo en blanco debe fallar: pasaria como valido y romperia mas adelante")
	}
	if _, err := extractString(map[string]interface{}{"x": 5}, "x"); err == nil {
		t.Fatal("un campo que no es texto debe fallar")
	}
	v, err := extractString(map[string]interface{}{"x": "  ok  "}, "x")
	if err != nil || v != "ok" {
		t.Fatalf("debe recortar espacios: %q %v", v, err)
	}
}

func TestLaURLEfectivaCambiaEnModoPruebas(t *testing.T) {
	prod := &domain.Integration{BaseURL: "https://api.tiendanube.com/v1", BaseURLTest: "http://mock:9102/v1"}
	if url, err := resolveEffectiveBaseURL(prod); err != nil || url != "https://api.tiendanube.com/v1" {
		t.Fatalf("sin modo pruebas debe usar la URL real: %q %v", url, err)
	}

	prod.IsTesting = true
	if url, err := resolveEffectiveBaseURL(prod); err != nil || url != "http://mock:9102/v1" {
		t.Fatalf("en modo pruebas debe usar el sandbox: %q %v", url, err)
	}

	sinTest := &domain.Integration{IsTesting: true, BaseURL: "https://api.tiendanube.com/v1"}
	if _, err := resolveEffectiveBaseURL(sinTest); !errors.Is(err, domain.ErrMissingBaseURLTest) {
		t.Fatalf("modo pruebas sin sandbox configurado debe fallar claro, llego %v", err)
	}

	sinProd := &domain.Integration{}
	if _, err := resolveEffectiveBaseURL(sinProd); !errors.Is(err, domain.ErrMissingBaseURL) {
		t.Fatalf("sin URL de produccion debe fallar claro, llego %v", err)
	}
}

func TestNoSeOperaUnaIntegracionDeOtroNegocio(t *testing.T) {
	uc := newTestUseCase(&fakeClient{}, &fakeRepo{})

	_, _, err := uc.resolveIntegrationForBusiness(context.Background(), testIntegrationID, testBusinessID+1)
	if !errors.Is(err, domain.ErrIntegrationNotFound) {
		t.Fatalf("una integracion de otro negocio no se debe poder usar, llego %v", err)
	}
}

func TestLaCredencialSeArmaConTokenStoreYURL(t *testing.T) {
	uc := newTestUseCase(&fakeClient{}, &fakeRepo{})

	_, cred, err := uc.resolveIntegrationForBusiness(context.Background(), testIntegrationID, testBusinessID)
	if err != nil {
		t.Fatalf("resolveIntegrationForBusiness fallo: %v", err)
	}
	if cred.AccessToken == "" || cred.StoreID != "1234567" || cred.BaseURL == "" {
		t.Fatalf("la credencial quedo incompleta: %+v", cred)
	}
}

func TestSinStoreIDNoSePuedeArmarLaCredencial(t *testing.T) {
	uc := newTestUseCase(&fakeClient{}, &fakeRepo{})
	servicio := uc.service.(*fakeService)
	servicio.integration.Config = map[string]interface{}{}
	servicio.integration.StoreID = ""

	_, _, err := uc.resolveIntegrationForBusiness(context.Background(), testIntegrationID, testBusinessID)
	if !errors.Is(err, domain.ErrMissingStoreID) {
		t.Fatalf("sin store_id la ruta de la API no existe: debe fallar claro, llego %v", err)
	}
}

func TestElTestDeConexionExigeTokenStoreYURL(t *testing.T) {
	uc := newTestUseCase(&fakeClient{}, &fakeRepo{})
	ctx := context.Background()

	err := uc.TestConnection(ctx, map[string]interface{}{"store_id": "1"}, map[string]interface{}{"access_token": "t"})
	if !errors.Is(err, domain.ErrMissingBaseURL) {
		t.Fatalf("sin base_url debe fallar claro, llego %v", err)
	}

	base := map[string]interface{}{"base_url": "https://api.tiendanube.com/v1", "store_id": "1"}
	err = uc.TestConnection(ctx, base, map[string]interface{}{})
	if !errors.Is(err, domain.ErrMissingAccessToken) {
		t.Fatalf("sin token debe fallar claro, llego %v", err)
	}

	err = uc.TestConnection(ctx, map[string]interface{}{"base_url": "https://api.tiendanube.com/v1"}, map[string]interface{}{"access_token": "t"})
	if !errors.Is(err, domain.ErrMissingStoreID) {
		t.Fatalf("sin store_id debe fallar claro, llego %v", err)
	}
}

func TestElTestDeConexionEnModoPruebasUsaElSandbox(t *testing.T) {
	uc := newTestUseCase(&fakeClient{}, &fakeRepo{})

	err := uc.TestConnection(context.Background(), map[string]interface{}{
		"is_testing": true,
		"base_url":   "https://api.tiendanube.com/v1",
		"store_id":   "1",
	}, map[string]interface{}{"access_token": "t"})

	if !errors.Is(err, domain.ErrMissingBaseURLTest) {
		t.Fatalf("en modo pruebas no debe caer a la API real: %v", err)
	}
}
