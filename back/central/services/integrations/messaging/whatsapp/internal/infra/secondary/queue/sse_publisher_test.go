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

func nuevoSSE(q *testkit.QueueMock) ports.ISSEEventPublisher {
	return NewSSEPublisher(q, testkit.NewSilentLogger())
}

func sobreSSE(t *testing.T, q *testkit.QueueMock) rabbitmq.EventEnvelope {
	t.Helper()
	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	var env rabbitmq.EventEnvelope
	require.NoError(t, json.Unmarshal(pubs[0].Body, &env))
	return env
}

func datosSSE(t *testing.T, q *testkit.QueueMock) map[string]any {
	t.Helper()
	crudo, err := json.Marshal(sobreSSE(t, q).Data)
	require.NoError(t, err)
	var out map[string]any
	require.NoError(t, json.Unmarshal(crudo, &out))
	return out
}

func TestPublishMessageReceived_VaAlExchangeDeEventosConCategoriaWhatsApp(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoSSE(q).PublishMessageReceived(context.Background(),
		26, "conv-1", "573001112233", "wamid.1", "hola"))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	assert.Equal(t, rabbitmq.EventsExchangeName, pubs[0].Exchange,
		"los eventos de UI van al exchange compartido de eventos, no a una cola propia")

	env := sobreSSE(t, q)
	assert.Equal(t, "whatsapp", env.Category)
	assert.EqualValues(t, 26, env.BusinessID,
		"sin business_id el mensaje de un cliente se le mostraria al panel de otro negocio")
	assert.NotEmpty(t, env.ID)
	assert.False(t, env.Timestamp.IsZero())
}

func TestPublishMessageReceived_MarcaLaDireccionComoEntrante(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoSSE(q).PublishMessageReceived(context.Background(),
		26, "conv-1", "573001112233", "wamid.1", "hola"))

	d := datosSSE(t, q)
	assert.Equal(t, "inbound", d["direction"],
		"el panel pinta a la izquierda o derecha segun la direccion; este publisher solo emite entrantes")
	assert.Equal(t, "hola", d["content"])
	assert.Equal(t, "wamid.1", d["message_id"])
	assert.Equal(t, "conv-1", d["conversation_id"])
	assert.Equal(t, "573001112233", d["phone_number"])
}

func TestPublishConversationStarted_LlevaLaConversacionYElTelefono(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoSSE(q).PublishConversationStarted(context.Background(),
		26, "conv-1", "573001112233"))

	d := datosSSE(t, q)
	assert.Equal(t, "conv-1", d["conversation_id"])
	assert.Equal(t, "573001112233", d["phone_number"])
	assert.NotContains(t, d, "content",
		"iniciar conversacion no lleva contenido: el primer mensaje viaja en su propio evento")
}

func TestPublishMessageStatusUpdated_LlevaElMensajeYSuEstadoNuevo(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoSSE(q).PublishMessageStatusUpdated(context.Background(),
		26, "wamid.1", "read"))

	d := datosSSE(t, q)
	assert.Equal(t, "wamid.1", d["message_id"])
	assert.Equal(t, "read", d["status"],
		"es lo que pinta el doble check azul en el panel")
}

func TestLosTresEventos_TienenTiposDistintosYUsanLaRoutingKeyDelTipo(t *testing.T) {
	ctx := context.Background()
	casos := []func(ports.ISSEEventPublisher) error{
		func(p ports.ISSEEventPublisher) error {
			return p.PublishMessageReceived(ctx, 26, "c", "57300", "w", "hola")
		},
		func(p ports.ISSEEventPublisher) error {
			return p.PublishConversationStarted(ctx, 26, "c", "57300")
		},
		func(p ports.ISSEEventPublisher) error {
			return p.PublishMessageStatusUpdated(ctx, 26, "w", "read")
		},
	}

	vistos := map[string]bool{}
	for _, publicar := range casos {
		q := &testkit.QueueMock{}

		require.NoError(t, publicar(nuevoSSE(q)))

		env := sobreSSE(t, q)
		assert.False(t, vistos[env.Type], "el tipo %q se repitio entre eventos distintos", env.Type)
		vistos[env.Type] = true
		assert.Equal(t, env.Type, q.Publicaciones()[0].RoutingKey,
			"la routing key es el tipo: los suscriptores del panel filtran por ahi")
	}
	assert.Len(t, vistos, 3)
}

func TestCadaEventoLlevaUnIDDistinto(t *testing.T) {
	q1 := &testkit.QueueMock{}
	q2 := &testkit.QueueMock{}

	require.NoError(t, nuevoSSE(q1).PublishMessageReceived(context.Background(), 26, "c", "57300", "w", "hola"))
	require.NoError(t, nuevoSSE(q2).PublishMessageReceived(context.Background(), 26, "c", "57300", "w", "hola"))

	assert.NotEqual(t, sobreSSE(t, q1).ID, sobreSSE(t, q2).ID)
}

func TestPublishSSE_PropagaElErrorDelBroker(t *testing.T) {
	ctx := context.Background()
	casos := map[string]func(ports.ISSEEventPublisher) error{
		"message_received":       func(p ports.ISSEEventPublisher) error { return p.PublishMessageReceived(ctx, 26, "c", "p", "m", "x") },
		"conversation_started":   func(p ports.ISSEEventPublisher) error { return p.PublishConversationStarted(ctx, 26, "c", "p") },
		"message_status_updated": func(p ports.ISSEEventPublisher) error { return p.PublishMessageStatusUpdated(ctx, 26, "m", "read") },
	}

	for nombre, publicar := range casos {
		t.Run(nombre, func(t *testing.T) {
			q := &testkit.QueueMock{
				PublishToExchangeFn: func(ctx context.Context, exchangeName, routingKey string, message []byte) error {
					return assert.AnError
				},
			}

			err := publicar(nuevoSSE(q))

			require.Error(t, err,
				"a diferencia del publisher SSE de shipments, que es fire-and-forget en goroutine, este si devuelve el error")
		})
	}
}

func TestPublishSSE_SinBrokerNoRevienta(t *testing.T) {
	p := NewSSEPublisher(nil, testkit.NewSilentLogger())

	err := p.PublishMessageReceived(context.Background(), 26, "c", "57300", "w", "hola")

	assert.NoError(t, err,
		"rabbitmq.PublishEvent devuelve nil cuando la cola es nil: el backend puede arrancar sin broker y el panel simplemente no recibe eventos")
}

func TestNoopSSEPublisher_NoPublicaYDevuelveNil(t *testing.T) {
	p := NewNoopSSEPublisher()
	ctx := context.Background()

	assert.NoError(t, p.PublishMessageReceived(ctx, 26, "c", "p", "m", "x"))
	assert.NoError(t, p.PublishConversationStarted(ctx, 26, "c", "p"))
	assert.NoError(t, p.PublishMessageStatusUpdated(ctx, 26, "m", "read"))
}
