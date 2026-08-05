package usecaseupdateorder

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type eventosCapturados struct {
	mu      sync.Mutex
	eventos []*entities.OrderEvent
}

func (e *eventosCapturados) registrar(event *entities.OrderEvent) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventos = append(e.eventos, event)
}

func (e *eventosCapturados) total() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return len(e.eventos)
}

func (e *eventosCapturados) copia() []*entities.OrderEvent {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]*entities.OrderEvent, len(e.eventos))
	copy(out, e.eventos)
	return out
}

func publisherQueCaptura(capturados *eventosCapturados) *mockRabbitPublisher {
	return &mockRabbitPublisher{
		PublishOrderEventFn: func(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error {
			capturados.registrar(event)
			return nil
		},
	}
}

func esperarEventos(t *testing.T, capturados *eventosCapturados, cantidad int) {
	t.Helper()
	assert.Eventually(t, func() bool {
		return capturados.total() >= cantidad
	}, 2*time.Second, 5*time.Millisecond)
}

func ordenParaEventos() *entities.ProbabilityOrder {
	businessID := uint(36)
	return &entities.ProbabilityOrder{
		ID:             "order-uuid",
		OrderNumber:    "MYS-0302",
		InternalNumber: "INT-1",
		ExternalID:     "EXT-1",
		Status:         "confirmed",
		CustomerEmail:  "ana@test.com",
		TotalAmount:    126000,
		Currency:       "COP",
		Platform:       "Shopify",
		BusinessID:     &businessID,
		IntegrationID:  197,
	}
}

func TestPublishOrderUpdatedEvent_PublicaEventoDeActualizacion(t *testing.T) {
	capturados := &eventosCapturados{}
	uc := newTestUpdateUseCase(&mockRepository{}, publisherQueCaptura(capturados), nil)

	uc.publishOrderUpdatedEvent(context.Background(), ordenParaEventos())

	esperarEventos(t, capturados, 1)
	eventos := capturados.copia()
	require.Len(t, eventos, 1)

	evento := eventos[0]
	assert.Equal(t, entities.OrderEventTypeUpdated, evento.Type)
	assert.Equal(t, "order-uuid", evento.OrderID)
	require.NotNil(t, evento.BusinessID)
	assert.Equal(t, uint(36), *evento.BusinessID)
	require.NotNil(t, evento.IntegrationID)
	assert.Equal(t, uint(197), *evento.IntegrationID)
	assert.Equal(t, "MYS-0302", evento.Data.OrderNumber)
	assert.Equal(t, "confirmed", evento.Data.CurrentStatus)
	assert.Equal(t, "COP", evento.Data.Currency)
}

func TestPublishOrderUpdatedEvent_SinIntegracion_NoAsignaIntegrationID(t *testing.T) {
	capturados := &eventosCapturados{}
	uc := newTestUpdateUseCase(&mockRepository{}, publisherQueCaptura(capturados), nil)

	order := ordenParaEventos()
	order.IntegrationID = 0

	uc.publishOrderUpdatedEvent(context.Background(), order)

	esperarEventos(t, capturados, 1)
	assert.Nil(t, capturados.copia()[0].IntegrationID)
}

func TestPublishOrderUpdatedEvent_SinPublisher_NoRompe(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	assert.NotPanics(t, func() {
		uc.publishOrderUpdatedEvent(context.Background(), ordenParaEventos())
	})
}

func TestPublishOrderUpdatedEvent_FallaPublicacion_NoPropaga(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	publisher := &mockRabbitPublisher{
		PublishOrderEventFn: func(ctx context.Context, event *entities.OrderEvent, order *entities.ProbabilityOrder) error {
			defer wg.Done()
			return errors.New("rabbitmq caido")
		},
	}
	uc := newTestUpdateUseCase(&mockRepository{}, publisher, nil)

	assert.NotPanics(t, func() {
		uc.publishOrderUpdatedEvent(context.Background(), ordenParaEventos())
	})

	wg.Wait()
}

func TestPublishOrderStatusChangedEvent_IncluyeEstadoAnterior(t *testing.T) {
	capturados := &eventosCapturados{}
	uc := newTestUpdateUseCase(&mockRepository{}, publisherQueCaptura(capturados), nil)

	uc.publishOrderStatusChangedEvent(context.Background(), ordenParaEventos(), "pending")

	esperarEventos(t, capturados, 1)
	evento := capturados.copia()[0]

	assert.Equal(t, entities.OrderEventTypeStatusChanged, evento.Type)
	assert.Equal(t, "pending", evento.Data.PreviousStatus)
	assert.Equal(t, "confirmed", evento.Data.CurrentStatus)
}

func TestPublishOrderStatusChangedEvent_SinPublisher_NoRompe(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	assert.NotPanics(t, func() {
		uc.publishOrderStatusChangedEvent(context.Background(), ordenParaEventos(), "pending")
	})
}

func TestPublishSyncOrderUpdated_NotificaAIntegraciones(t *testing.T) {
	var recibidoIntegracion uint
	var recibidoBusiness *uint
	var recibidoData map[string]interface{}

	integrationPub := &mockIntegrationEventPublisher{
		PublishSyncOrderUpdatedFn: func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
			recibidoIntegracion = integrationID
			recibidoBusiness = businessID
			recibidoData = data
		},
	}
	uc := newTestUpdateUseCase(&mockRepository{}, nil, integrationPub)

	uc.publishSyncOrderUpdated(context.Background(), ordenParaEventos())

	assert.Equal(t, uint(197), recibidoIntegracion)
	require.NotNil(t, recibidoBusiness)
	assert.Equal(t, uint(36), *recibidoBusiness)
	assert.Equal(t, "order-uuid", recibidoData["order_id"])
	assert.Equal(t, "MYS-0302", recibidoData["order_number"])
	assert.Equal(t, "confirmed", recibidoData["status"])
	assert.Equal(t, "ana@test.com", recibidoData["customer_email"])
}

func TestPublishSyncOrderUpdated_SinPublisher_NoRompe(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	assert.NotPanics(t, func() {
		uc.publishSyncOrderUpdated(context.Background(), ordenParaEventos())
	})
}

func TestPublishSyncOrderSkipped_MarcaLaOrdenComoOmitida(t *testing.T) {
	var recibidoData map[string]interface{}
	integrationPub := &mockIntegrationEventPublisher{
		PublishSyncOrderRejectedFn: func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
			recibidoData = data
		},
	}
	uc := newTestUpdateUseCase(&mockRepository{}, nil, integrationPub)

	uc.publishSyncOrderSkipped(context.Background(), ordenParaEventos())

	require.NotNil(t, recibidoData)
	assert.Equal(t, true, recibidoData["skipped"])
	assert.Equal(t, "Ya existe (sin cambios)", recibidoData["reason"])
	assert.Equal(t, "order-uuid", recibidoData["order_id"])
}

func TestPublishSyncOrderSkipped_SinIntegracion_NoNotifica(t *testing.T) {
	llamado := false
	integrationPub := &mockIntegrationEventPublisher{
		PublishSyncOrderRejectedFn: func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
			llamado = true
		},
	}
	uc := newTestUpdateUseCase(&mockRepository{}, nil, integrationPub)

	order := ordenParaEventos()
	order.IntegrationID = 0

	uc.publishSyncOrderSkipped(context.Background(), order)

	assert.False(t, llamado)
}

func TestPublishSyncOrderSkipped_SinPublisher_NoRompe(t *testing.T) {
	uc := newTestUpdateUseCase(&mockRepository{}, nil, nil)

	assert.NotPanics(t, func() {
		uc.publishSyncOrderSkipped(context.Background(), ordenParaEventos())
	})
}

func TestPublishUpdateEvents_OrdenDeIntegracionConCambioDeEstado(t *testing.T) {
	capturados := &eventosCapturados{}
	syncLlamado := false
	integrationPub := &mockIntegrationEventPublisher{
		PublishSyncOrderUpdatedFn: func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
			syncLlamado = true
		},
	}
	uc := newTestUpdateUseCase(&mockRepository{}, publisherQueCaptura(capturados), integrationPub)

	uc.publishUpdateEvents(context.Background(), ordenParaEventos(), "pending", false)

	esperarEventos(t, capturados, 2)
	assert.True(t, syncLlamado)

	tipos := map[entities.OrderEventType]bool{}
	for _, e := range capturados.copia() {
		tipos[e.Type] = true
	}
	assert.True(t, tipos[entities.OrderEventTypeUpdated])
	assert.True(t, tipos[entities.OrderEventTypeStatusChanged])
}

func TestPublishUpdateEvents_SinCambioDeEstado_NoPublicaStatusChanged(t *testing.T) {
	capturados := &eventosCapturados{}
	uc := newTestUpdateUseCase(&mockRepository{}, publisherQueCaptura(capturados), nil)

	uc.publishUpdateEvents(context.Background(), ordenParaEventos(), "confirmed", false)

	esperarEventos(t, capturados, 1)
	time.Sleep(50 * time.Millisecond)

	eventos := capturados.copia()
	require.Len(t, eventos, 1)
	assert.Equal(t, entities.OrderEventTypeUpdated, eventos[0].Type)
}

func TestPublishUpdateEvents_OrdenManual_NoNotificaIntegraciones(t *testing.T) {
	capturados := &eventosCapturados{}
	syncLlamado := false
	integrationPub := &mockIntegrationEventPublisher{
		PublishSyncOrderUpdatedFn: func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
			syncLlamado = true
		},
	}
	uc := newTestUpdateUseCase(&mockRepository{}, publisherQueCaptura(capturados), integrationPub)

	uc.publishUpdateEvents(context.Background(), ordenParaEventos(), "pending", true)

	esperarEventos(t, capturados, 2)
	assert.False(t, syncLlamado, "una orden manual no notifica sincronizacion a integraciones")
}

func TestPublishUpdateEvents_SinIntegracion_NoNotificaSync(t *testing.T) {
	capturados := &eventosCapturados{}
	syncLlamado := false
	integrationPub := &mockIntegrationEventPublisher{
		PublishSyncOrderUpdatedFn: func(ctx context.Context, integrationID uint, businessID *uint, data map[string]interface{}) {
			syncLlamado = true
		},
	}
	uc := newTestUpdateUseCase(&mockRepository{}, publisherQueCaptura(capturados), integrationPub)

	order := ordenParaEventos()
	order.IntegrationID = 0

	uc.publishUpdateEvents(context.Background(), order, "confirmed", false)

	esperarEventos(t, capturados, 1)
	assert.False(t, syncLlamado)
}
