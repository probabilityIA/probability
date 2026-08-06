package queue

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/testkit"
)

func nuevoWebhook(q *testkit.QueueMock) ports.IEventPublisher {
	return NewWebhookPublisher(q, testkit.NewSilentLogger())
}

func TestPublishOrderConfirmed_VaASuColaPropia(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoWebhook(q).PublishOrderConfirmed(context.Background(), "1001", "573001112233", 26))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	assert.Equal(t, rabbitmq.QueueWhatsAppOrderConfirmed, pubs[0].Queue,
		"cada intencion del cliente tiene su propia cola: a diferencia de facturacion, aca el consumidor no filtra por event_type")

	ev := evento(t, q)
	assert.Equal(t, "order.confirmed", ev["event_type"])
	assert.Equal(t, "1001", ev["order_number"])
	assert.Equal(t, "573001112233", ev["phone_number"])
	assert.EqualValues(t, 26, ev["business_id"])
	assert.Equal(t, "whatsapp", ev["source"],
		"el source distingue una confirmacion por bot de una hecha a mano en el panel")
}

func TestPublishOrderCancelled_LlevaElMotivoDeCancelacion(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoWebhook(q).PublishOrderCancelled(context.Background(),
		"1001", "ya no lo quiero", "573001112233", 26))

	pubs := q.Publicaciones()
	assert.Equal(t, rabbitmq.QueueWhatsAppOrderCancelled, pubs[0].Queue)
	ev := evento(t, q)
	assert.Equal(t, "order.cancelled", ev["event_type"])
	assert.Equal(t, "ya no lo quiero", ev["cancellation_reason"],
		"el motivo que dio el cliente queda en la orden para el informe de cancelaciones")
}

func TestPublishNoveltyRequested_LlevaElTipoDeNovedad(t *testing.T) {
	casos := []string{"change_address", "change_products", "change_payment"}

	for _, tipo := range casos {
		t.Run(tipo, func(t *testing.T) {
			q := &testkit.QueueMock{}

			require.NoError(t, nuevoWebhook(q).PublishNoveltyRequested(context.Background(),
				"1001", tipo, "573001112233", 26))

			pubs := q.Publicaciones()
			assert.Equal(t, rabbitmq.QueueWhatsAppOrderNovelty, pubs[0].Queue)
			ev := evento(t, q)
			assert.Equal(t, "order.novelty_requested", ev["event_type"])
			assert.Equal(t, tipo, ev["novelty_type"],
				"el tipo decide que flujo de correccion se dispara del otro lado")
		})
	}
}

func TestPublishNoveltyRequested_NoValidaElTipoDeNovedad(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoWebhook(q).PublishNoveltyRequested(context.Background(),
		"1001", "tipo_inventado", "573001112233", 26))

	assert.Equal(t, "tipo_inventado", evento(t, q)["novelty_type"],
		"el publisher no valida contra la lista de tipos conocidos: un tipo nuevo llega igual y el consumidor decide")
}

func TestPublishHandoffRequested_LlevaLaConversacionYQuedaPendienteDeAgente(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoWebhook(q).PublishHandoffRequested(context.Background(),
		"1001", "573001112233", 26, "conv-1"))

	pubs := q.Publicaciones()
	assert.Equal(t, rabbitmq.QueueWhatsAppCustomerHandoff, pubs[0].Queue)
	ev := evento(t, q)
	assert.Equal(t, "customer.handoff_requested", ev["event_type"])
	assert.Equal(t, "conv-1", ev["conversation_id"],
		"el agente humano necesita la conversacion para retomar el hilo donde lo dejo el bot")
	assert.Equal(t, "pending_human_agent", ev["status"])
}

func TestTodosLosEventos_LlevanFechaUnixYNegocio(t *testing.T) {
	ctx := context.Background()
	casos := map[string]func(ports.IEventPublisher) error{
		"confirmed": func(p ports.IEventPublisher) error { return p.PublishOrderConfirmed(ctx, "1", "57300", 26) },
		"cancelled": func(p ports.IEventPublisher) error { return p.PublishOrderCancelled(ctx, "1", "r", "57300", 26) },
		"novelty":   func(p ports.IEventPublisher) error { return p.PublishNoveltyRequested(ctx, "1", "t", "57300", 26) },
		"handoff":   func(p ports.IEventPublisher) error { return p.PublishHandoffRequested(ctx, "1", "57300", 26, "c") },
	}

	for nombre, publicar := range casos {
		t.Run(nombre, func(t *testing.T) {
			q := &testkit.QueueMock{}

			require.NoError(t, publicar(nuevoWebhook(q)))

			ev := evento(t, q)
			assert.Greater(t, ev["timestamp"], float64(0))
			assert.EqualValues(t, 26, ev["business_id"],
				"sin business_id el evento se aplicaria sobre la orden de otro negocio con el mismo numero")
			assert.Equal(t, "whatsapp", ev["source"])
		})
	}
}

func TestPublish_PropagaElErrorDelBrokerEnLosCuatroEventos(t *testing.T) {
	ctx := context.Background()
	casos := map[string]func(ports.IEventPublisher) error{
		"confirmed": func(p ports.IEventPublisher) error { return p.PublishOrderConfirmed(ctx, "1", "57300", 26) },
		"cancelled": func(p ports.IEventPublisher) error { return p.PublishOrderCancelled(ctx, "1", "r", "57300", 26) },
		"novelty":   func(p ports.IEventPublisher) error { return p.PublishNoveltyRequested(ctx, "1", "t", "57300", 26) },
		"handoff":   func(p ports.IEventPublisher) error { return p.PublishHandoffRequested(ctx, "1", "57300", 26, "c") },
	}

	for nombre, publicar := range casos {
		t.Run(nombre, func(t *testing.T) {
			err := publicar(nuevoWebhook(brokerCaido()))

			require.Error(t, err)
			assert.Contains(t, err.Error(), "error al publicar evento en RabbitMQ",
				"si el broker falla, la confirmacion del cliente se pierde: el bot ya le respondio que quedo confirmado")
		})
	}
}

func TestPublish_NoDeclaraLasColasNiEsperaConfirmacion(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoWebhook(q).PublishOrderConfirmed(context.Background(), "1001", "57300", 26))

	assert.Empty(t, q.Declaradas)
}

func TestPublish_ElCuerpoEsJSONPlanoSinAnidamiento(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoWebhook(q).PublishOrderConfirmed(context.Background(), "1001", "57300", 26))

	var ev map[string]any
	require.NoError(t, json.Unmarshal(q.Publicaciones()[0].Body, &ev))
	for clave, valor := range ev {
		_, esMapa := valor.(map[string]any)
		assert.False(t, esMapa,
			"el evento es plano a proposito; la clave %q trae un objeto anidado", clave)
	}
}
