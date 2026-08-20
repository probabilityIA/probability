package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func TestUnaIntegracionQueNoExisteFallaClaro(t *testing.T) {
	uc := newTestUseCase(&fakeClient{}, &fakeRepo{})
	uc.service.(*fakeService).integration = nil

	if _, err := uc.fetchIntegration(context.Background(), testIntegrationID); !errors.Is(err, domain.ErrIntegrationNotFound) {
		t.Fatalf("una integracion inexistente debe fallar claro, llego %v", err)
	}
}

func TestSinTokenGuardadoNoSePuedeOperar(t *testing.T) {
	uc := newTestUseCase(&fakeClient{}, &fakeRepo{})
	uc.service.(*fakeService).tokenVacio = true

	_, _, err := uc.resolveIntegrationForBusiness(context.Background(), testIntegrationID, testBusinessID)
	if !errors.Is(err, domain.ErrMissingAccessToken) {
		t.Fatalf("sin access_token no se puede llamar a la API: debe fallar claro, llego %v", err)
	}
}

func TestElMensajeHaciaProbabilityLlevaPesoYMedidas(t *testing.T) {
	uc := newTestUseCase(&fakeClient{}, &fakeRepo{})

	msg := uc.upsertMsgFromTiendanube(testBusinessID, 77, tiendanubeSKU{
		SKU:        "SKU-1",
		Name:       "Producto",
		Price:      100,
		ExternalID: "2001:9001",
		Weight:     2,
		Depth:      30,
		Width:      20,
		Height:     10,
	})

	if msg.Weight == nil || *msg.Weight != 2 || msg.WeightUnit != probabilityWeightUnit {
		t.Fatalf("el peso del canal debe viajar en kilos: %+v", msg)
	}
	if msg.DimensionUnit != probabilityDimensionUnit {
		t.Fatalf("si hay medidas debe declararse la unidad: %q", msg.DimensionUnit)
	}
	if msg.Length == nil || *msg.Length != 30 {
		t.Fatalf("el largo sale del depth de la variante: %+v", msg.Length)
	}

	sinMedidas := uc.upsertMsgFromTiendanube(testBusinessID, 77, tiendanubeSKU{SKU: "SKU-2", ExternalID: "1:2"})
	if sinMedidas.DimensionUnit != "" || sinMedidas.Weight != nil {
		t.Fatalf("sin medidas no se debe inventar unidad ni peso: %+v", sinMedidas)
	}
}

func TestSiLaComparacionFallaElResultadoGuardaElMotivo(t *testing.T) {
	cola := nuevaCola()
	uc := newTestUseCaseConCola(&fakeClient{listaErr: errTiendanubeTest("Tiendanube caido")}, &fakeRepo{}, cola)

	uc.ReconcileProductsAsync(context.Background(), testIntegrationID, testBusinessID, 77, "corr-1")

	guardados := cola.publicados[rabbitmq.QueueIntegrationSyncRuns]
	if len(guardados) != 1 {
		t.Fatalf("un fallo tambien debe quedar registrado, no desaparecer: %d", len(guardados))
	}
	if !contieneTexto(string(guardados[0]), "Tiendanube caido") {
		t.Fatalf("el motivo del fallo debe quedar guardado: %s", string(guardados[0]))
	}
}

func TestSincronizarProductosCortaSiFallaElPrimerSentido(t *testing.T) {
	uc := newTestUseCase(&fakeClient{listaErr: errTiendanubeTest("Tiendanube caido")}, &fakeRepo{})

	if err := uc.SyncProducts(context.Background(), testIntegrationID, testBusinessID, "corr-1"); err == nil {
		t.Fatal("si no se puede leer el canal, la sincronizacion completa debe fallar")
	}
}

func TestElTestDeConexionPasaCuandoLaTiendaResponde(t *testing.T) {
	uc := newTestUseCase(&fakeClient{tienda: &domain.StoreInfo{Name: "Mi tienda", URL: "https://x.com"}}, &fakeRepo{})

	err := uc.TestConnection(context.Background(), map[string]interface{}{
		"base_url": "https://api.tiendanube.com/v1",
		"store_id": "1234567",
	}, map[string]interface{}{"access_token": "token"})

	if err != nil {
		t.Fatalf("con todo configurado el test debe pasar: %v", err)
	}
}

func TestElTestDeConexionPropagaElErrorDelCanal(t *testing.T) {
	uc := newTestUseCase(&fakeClient{tiendaErr: domain.ErrInvalidCredentials}, &fakeRepo{})

	err := uc.TestConnection(context.Background(), map[string]interface{}{
		"base_url": "https://api.tiendanube.com/v1",
		"store_id": "1234567",
	}, map[string]interface{}{"access_token": "malo"})

	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("un token invalido debe llegar tal cual al usuario, llego %v", err)
	}
}

func contieneTexto(texto, buscado string) bool {
	return len(texto) >= len(buscado) && (func() bool {
		for i := 0; i+len(buscado) <= len(texto); i++ {
			if texto[i:i+len(buscado)] == buscado {
				return true
			}
		}
		return false
	})()
}
