package queue

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/app/usecases"
	"github.com/secamc93/probability/back/central/shared/inventorycompare"
	"github.com/secamc93/probability/back/central/shared/log"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

type usoFalso struct {
	usecases.ITiendanubeUseCase
	err      error
	llamadas int
}

func (u *usoFalso) UpdateInventory(ctx context.Context, integrationID string, productExternalID string, quantity int) error {
	u.llamadas++
	return u.err
}

func (u *usoFalso) CompareInventory(ctx context.Context, integrationID string, businessID uint, page, pageSize int, skus ...string) (*inventorycompare.Page, error) {
	return nil, nil
}

func nuevoConsumidor(err error) (*InventoryPushConsumer, *usoFalso) {
	uso := &usoFalso{err: err}
	return NewInventoryPushConsumer(nil, uso, log.New()), uso
}

func TestUnErrorPermanenteSeDescartaParaNoDejarLaColaEnBucle(t *testing.T) {
	permanentes := []error{
		domain.ErrIntegrationNotFound,
		domain.ErrInvalidCredentials,
		domain.ErrMissingAccessToken,
		domain.ErrMissingStoreID,
		errors.New("tiendanube: el external_id 2001 no incluye variante, no se puede escribir stock"),
		errors.New("tiendanube client: PUT /products returned 404: not found"),
		errors.New("tiendanube client: PUT /products returned 422: sku taken"),
	}

	for _, err := range permanentes {
		if !isPermanent(err) {
			t.Fatalf("%q es permanente: reintentarlo deja la cola girando sin fin", err)
		}
	}
}

func TestUnErrorTransitorioSeReintenta(t *testing.T) {
	transitorios := []error{
		errors.New("tiendanube client: request failed: connection refused"),
		errors.New("tiendanube client: GET /products returned 500: server error"),
		domain.ErrRateLimited,
		errors.New("context deadline exceeded"),
	}

	for _, err := range transitorios {
		if isPermanent(err) {
			t.Fatalf("%q puede recuperarse solo: descartarlo pierde el evento", err)
		}
	}

	if isPermanent(nil) {
		t.Fatal("sin error no hay nada que clasificar")
	}
}

func TestElMensajeIncompletoSeDescartaSinReintentar(t *testing.T) {
	c, uso := nuevoConsumidor(nil)

	if err := c.handle(context.Background(), []byte(`{"integration_id": 0, "external_product_id": ""}`)); err != nil {
		t.Fatalf("un mensaje incompleto nunca va a poder procesarse: se descarta, no se reintenta. Llego %v", err)
	}
	if uso.llamadas != 0 {
		t.Fatal("no se debe llamar al canal con un mensaje incompleto")
	}
}

func TestUnMensajeIlegibleSeDescarta(t *testing.T) {
	c, _ := nuevoConsumidor(nil)

	if err := c.handle(context.Background(), []byte(`{esto no es json`)); err != nil {
		t.Fatalf("un mensaje corrupto no mejora con reintentos: se descarta. Llego %v", err)
	}
}

func TestElMensajeValidoEmpujaElStock(t *testing.T) {
	c, uso := nuevoConsumidor(nil)

	err := c.handle(context.Background(), []byte(`{"integration_id": 77, "external_product_id": "2001:9001", "quantity": 5}`))
	if err != nil {
		t.Fatalf("un mensaje valido no debe fallar: %v", err)
	}
	if uso.llamadas != 1 {
		t.Fatalf("se esperaba una escritura de stock, hubo %d", uso.llamadas)
	}
}

func TestElFalloTransitorioDevuelveErrorParaQueSeReintente(t *testing.T) {
	c, _ := nuevoConsumidor(errors.New("tiendanube client: request failed: connection refused"))

	err := c.handle(context.Background(), []byte(`{"integration_id": 77, "external_product_id": "2001:9001", "quantity": 5}`))
	if err == nil {
		t.Fatal("un fallo transitorio debe devolver error para que el mensaje se reencole")
	}
}

func TestElFalloPermanenteSeAceptaYNoSeReencola(t *testing.T) {
	c, _ := nuevoConsumidor(domain.ErrInvalidCredentials)

	err := c.handle(context.Background(), []byte(`{"integration_id": 77, "external_product_id": "2001:9001", "quantity": 5}`))
	if err != nil {
		t.Fatalf("credenciales invalidas no se arreglan reintentando: debe ACKear. Llego %v", err)
	}
}
