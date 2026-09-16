package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
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

type statusFake struct {
	pedidos []string
	err     error
}

func (s *statusFake) ChangeStatus(ctx context.Context, orderID string, req *dtos.ChangeStatusRequest) (*dtos.OrderResponse, error) {
	s.pedidos = append(s.pedidos, req.Status)
	return nil, s.err
}

type plantillaFake struct {
	nombres []string
}

func (p *plantillaFake) RequestTemplate(ctx context.Context, businessID uint, phone, templateName string, parameters []string) error {
	p.nombres = append(p.nombres, templateName)
	return nil
}

func eventoCancelacion(t *testing.T) []byte {
	t.Helper()
	cuerpo, err := json.Marshal(WhatsAppCancelledEvent{EventType: "cancelled", OrderNumber: "DEMO-1", PhoneNumber: "573001234567", BusinessID: 26})
	require.NoError(t, err)
	return cuerpo
}

func consumerCancelacion(orden *entities.ProbabilityOrder) (*WhatsAppConsumer, *statusFake, *plantillaFake) {
	status := &statusFake{}
	plantillas := &plantillaFake{}
	return &WhatsAppConsumer{repository: &repoFake{orden: orden}, statusUseCase: status, templateRequester: plantillas, log: log.New()}, status, plantillas
}

func TestCancelacionSinGuiaCancelaLaOrden(t *testing.T) {
	c, status, plantillas := consumerCancelacion(&entities.ProbabilityOrder{ID: "ord-1", OrderNumber: "DEMO-1", Status: "pending"})

	require.NoError(t, c.handleCancelled(eventoCancelacion(t)))

	assert.Equal(t, []string{"cancelled"}, status.pedidos)
	assert.Equal(t, []string{"pedido_cancelado"}, plantillas.nombres)
}

func TestCancelacionConGuiaQuedaEnSolicitud(t *testing.T) {
	guia := "240012345678"
	orden := &entities.ProbabilityOrder{ID: "ord-1", OrderNumber: "DEMO-1", Status: "ready_to_ship", Shipments: []entities.ProbabilityShipment{{TrackingNumber: &guia, Status: "pending"}}}
	c, status, plantillas := consumerCancelacion(orden)

	require.NoError(t, c.handleCancelled(eventoCancelacion(t)))

	assert.Equal(t, []string{"cancel_requested"}, status.pedidos)
	assert.Equal(t, []string{"solicitud_cancelacion_recibida"}, plantillas.nombres)
}

func TestCancelacionConGuiaAnuladaCancelaLaOrden(t *testing.T) {
	guia := "240012345678"
	orden := &entities.ProbabilityOrder{ID: "ord-1", OrderNumber: "DEMO-1", Status: "pending", Shipments: []entities.ProbabilityShipment{{TrackingNumber: &guia, Status: "cancelled"}}}
	c, status, _ := consumerCancelacion(orden)

	require.NoError(t, c.handleCancelled(eventoCancelacion(t)))

	assert.Equal(t, []string{"cancelled"}, status.pedidos)
}

func TestCancelacionRepetidaNoHaceNada(t *testing.T) {
	c, status, plantillas := consumerCancelacion(&entities.ProbabilityOrder{ID: "ord-1", OrderNumber: "DEMO-1", Status: "cancel_requested"})

	require.NoError(t, c.handleCancelled(eventoCancelacion(t)))

	assert.Empty(t, status.pedidos)
	assert.Empty(t, plantillas.nombres)
}
