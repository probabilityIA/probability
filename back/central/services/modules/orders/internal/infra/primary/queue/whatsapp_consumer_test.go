package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type repoFake struct {
	ports.IRepository
	orden    *entities.ProbabilityOrder
	guardada *entities.ProbabilityOrder
	err      error
}

func (r *repoFake) GetOrderByOrderNumberAndBusiness(ctx context.Context, orderNumber string, businessID uint) (*entities.ProbabilityOrder, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.orden, nil
}

func (r *repoFake) UpdateOrder(ctx context.Context, order *entities.ProbabilityOrder) error {
	r.guardada = order
	return nil
}

func consumerConOrden() (*WhatsAppConsumer, *repoFake) {
	repo := &repoFake{orden: &entities.ProbabilityOrder{ID: "ord-1", OrderNumber: "DEMO-1"}}
	return &WhatsAppConsumer{repository: repo, log: log.New()}, repo
}

func eventoNovedad(t *testing.T, tipo string) []byte {
	t.Helper()
	cuerpo, err := json.Marshal(WhatsAppNoveltyEvent{
		EventType:   "novelty",
		OrderNumber: "DEMO-1",
		NoveltyType: tipo,
		PhoneNumber: "573001234567",
		BusinessID:  26,
	})
	require.NoError(t, err)
	return cuerpo
}

func TestCambioDeDireccionDejaLaOrdenComoNoConfirmada(t *testing.T) {
	c, repo := consumerConOrden()

	require.NoError(t, c.handleNovelty(eventoNovedad(t, "change_address")))

	require.NotNil(t, repo.guardada)
	require.NotNil(t, repo.guardada.IsConfirmed,
		"una novedad es una no confirmacion: dejarla en nulo la hace ver como pendiente de respuesta")
	assert.False(t, *repo.guardada.IsConfirmed)
}

func TestLaNovedadDeDireccionNoTocaLaDireccionDeLaOrden(t *testing.T) {
	c, repo := consumerConOrden()
	repo.orden.ShippingStreet = "calle 100 #10-20"

	require.NoError(t, c.handleNovelty(eventoNovedad(t, "change_address")))

	assert.Equal(t, "calle 100 #10-20", repo.guardada.ShippingStreet,
		"el cliente no puede cambiar su direccion solo: se presta para pedir a una direccion cercana y desviar el envio a una lejana")
}

func TestLaNovedadDejaElMotivoVisibleParaElNegocio(t *testing.T) {
	c, repo := consumerConOrden()

	require.NoError(t, c.handleNovelty(eventoNovedad(t, "change_address")))

	require.NotNil(t, repo.guardada.Novelty)
	assert.Contains(t, *repo.guardada.Novelty, "cambio de direccion")
	assert.Contains(t, *repo.guardada.Novelty, "573001234567")
}

func TestLasNovedadesSeAcumulanSinPisarLaAnterior(t *testing.T) {
	c, repo := consumerConOrden()
	previa := "Novedad anterior"
	repo.orden.Novelty = &previa

	require.NoError(t, c.handleNovelty(eventoNovedad(t, "change_products")))

	assert.Contains(t, *repo.guardada.Novelty, "Novedad anterior")
	assert.Contains(t, *repo.guardada.Novelty, "cambio de productos")
}

func consumerConError(err error) (*WhatsAppConsumer, *repoFake) {
	repo := &repoFake{err: err}
	return &WhatsAppConsumer{repository: repo, log: log.New()}, repo
}

func eventoConfirmacion(t *testing.T) []byte {
	t.Helper()
	cuerpo, err := json.Marshal(WhatsAppConfirmedEvent{
		EventType:   "confirmed",
		OrderNumber: "PRUEBA-ISABEL",
		PhoneNumber: "573192611891",
		BusinessID:  66,
	})
	require.NoError(t, err)
	return cuerpo
}

func TestEsPermanenteConOrdenInexistente(t *testing.T) {
	assert.True(t, esPermanente(domainerrors.ErrOrderNotFound),
		"una orden que no existe no aparece por reintentar: reencolar produce el bucle caliente")
}

func TestEsPermanenteConOrdenInexistenteEnvuelta(t *testing.T) {
	assert.True(t, esPermanente(fmt.Errorf("error getting order: %w", domainerrors.ErrOrderNotFound)))
}

func TestEsPermanenteConNegocioEliminado(t *testing.T) {
	assert.True(t, esPermanente(domainerrors.ErrOrderBusinessDeleted))
}

func TestEsPermanenteEsFalsoConErrorTransitorio(t *testing.T) {
	assert.False(t, esPermanente(errors.New("dial tcp 172.31.96.15:5432: connect: connection refused")),
		"la base caida se recupera sola: ese mensaje si debe reencolarse")
}

func TestEsPermanenteEsFalsoConNil(t *testing.T) {
	assert.False(t, esPermanente(nil))
}

func TestLaConfirmacionDeUnaOrdenInexistenteSeDescarta(t *testing.T) {
	c, _ := consumerConError(domainerrors.ErrOrderNotFound)

	assert.NoError(t, c.handleConfirmed(eventoConfirmacion(t)),
		"devolver el error hace Nack con requeue y el mensaje vuelve al instante: es el bucle de PRUEBA-ISABEL")
}

func TestLaConfirmacionConLaBaseCaidaSeReencola(t *testing.T) {
	fallo := errors.New("dial tcp 172.31.96.15:5432: connect: connection refused")
	c, _ := consumerConError(fallo)

	assert.ErrorIs(t, c.handleConfirmed(eventoConfirmacion(t)), fallo,
		"descartar un fallo transitorio pierde la confirmacion del cliente en silencio")
}

func TestUnPayloadMalformadoSeDescarta(t *testing.T) {
	c, _ := consumerConError(nil)

	assert.NoError(t, c.handleConfirmed([]byte("{esto no es json")),
		"un payload roto no se arregla por reintentar")
}
