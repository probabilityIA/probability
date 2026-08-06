package queue

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/testkit"
)

func nuevoPersistence(q *testkit.QueueMock) ports.IPersistencePublisher {
	return NewPersistencePublisher(q, testkit.NewSilentLogger())
}

func brokerCaido() *testkit.QueueMock {
	return &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("canal cerrado")
		},
	}
}

func evento(t *testing.T, q *testkit.QueueMock) map[string]any {
	t.Helper()
	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	var out map[string]any
	require.NoError(t, json.Unmarshal(pubs[0].Body, &out))
	return out
}

func conversacion() *entities.Conversation {
	return &entities.Conversation{
		ID:               "conv-1",
		PhoneNumber:      "573001112233",
		OrderNumber:      "1001",
		ConversationType: entities.ConversationTypeOrder,
		BusinessID:       26,
		LastMessageID:    "wamid.1",
		CreatedAt:        time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestPublishConversationCreated_VaALaColaDePersistenciaConSuTipo(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPersistence(q).PublishConversationCreated(context.Background(), conversacion()))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	assert.Equal(t, rabbitmq.QueueWhatsAppPersistenceEvents, pubs[0].Queue,
		"la persistencia se delega por cola en vez de escribir la base desde el bot: asi el bot responde rapido y la escritura va aparte")
	assert.Equal(t, "conversation.created", evento(t, q)["event_type"])
}

func TestPublishConversation_LosTresEventosSoloCambianElTipo(t *testing.T) {
	casos := map[string]func(ports.IPersistencePublisher) error{
		"conversation.created": func(p ports.IPersistencePublisher) error {
			return p.PublishConversationCreated(context.Background(), conversacion())
		},
		"conversation.updated": func(p ports.IPersistencePublisher) error {
			return p.PublishConversationUpdated(context.Background(), conversacion())
		},
		"conversation.expired": func(p ports.IPersistencePublisher) error {
			return p.PublishConversationExpired(context.Background(), "conv-1")
		},
	}

	for tipo, publicar := range casos {
		t.Run(tipo, func(t *testing.T) {
			q := &testkit.QueueMock{}

			require.NoError(t, publicar(nuevoPersistence(q)))

			ev := evento(t, q)
			assert.Equal(t, tipo, ev["event_type"])
			conv := ev["conversation"].(map[string]any)
			assert.Equal(t, "conv-1", conv["id"])
		})
	}
}

func TestPublishConversationCreated_LlevaTelefonoOrdenYNegocio(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPersistence(q).PublishConversationCreated(context.Background(), conversacion()))

	conv := evento(t, q)["conversation"].(map[string]any)
	assert.Equal(t, "573001112233", conv["phone_number"],
		"el telefono es la llave natural de la conversacion en WhatsApp")
	assert.Equal(t, "1001", conv["order_number"])
	assert.EqualValues(t, 26, conv["business_id"],
		"sin business_id la conversacion se guardaria sin dueno y otro negocio podria verla")
}

func TestToConversationPayload_SinTipoDeConversacionAsumeQueEsDeOrden(t *testing.T) {
	q := &testkit.QueueMock{}
	c := conversacion()
	c.ConversationType = ""

	require.NoError(t, nuevoPersistence(q).PublishConversationCreated(context.Background(), c))

	conv := evento(t, q)["conversation"].(map[string]any)
	assert.Equal(t, string(entities.ConversationTypeOrder), conv["conversation_type"],
		"conversaciones viejas quedaron sin tipo; se asume 'order' al publicar para no romper el consumidor")
}

func TestToConversationPayload_RespetaElTipoQueYaViene(t *testing.T) {
	q := &testkit.QueueMock{}
	c := conversacion()
	c.ConversationType = entities.ConversationTypeSystemAlert

	require.NoError(t, nuevoPersistence(q).PublishConversationCreated(context.Background(), c))

	conv := evento(t, q)["conversation"].(map[string]any)
	assert.Equal(t, string(entities.ConversationTypeSystemAlert), conv["conversation_type"])
}

func TestPublishConversationExpired_SoloMandaElIDSinElRestoDeLaConversacion(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPersistence(q).PublishConversationExpired(context.Background(), "conv-1"))

	conv := evento(t, q)["conversation"].(map[string]any)
	assert.Equal(t, "conv-1", conv["id"])
	assert.Equal(t, "", conv["phone_number"],
		"expirar solo necesita el id; el resto del payload viaja vacio y el consumidor no debe sobreescribir con eso")
	assert.EqualValues(t, 0, conv["business_id"])
}

func TestPublishConversation_PropagaElErrorDelBroker(t *testing.T) {
	err := nuevoPersistence(brokerCaido()).PublishConversationCreated(context.Background(), conversacion())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error publicando conversation event",
		"a diferencia de los publishers de facturacion, aca si se devuelve error siempre")
}

func TestPublishMessageLogCreated_LlevaElMensajeConSuDireccionYEstado(t *testing.T) {
	q := &testkit.QueueMock{}
	log := &entities.MessageLog{
		ID:             "log-1",
		ConversationID: "conv-1",
		Direction:      entities.MessageDirectionOutbound,
		MessageID:      "wamid.1",
		TemplateName:   "guia_generada",
		Content:        "Tu guia esta lista",
		Status:         entities.MessageStatusSent,
	}

	require.NoError(t, nuevoPersistence(q).PublishMessageLogCreated(context.Background(), log))

	ev := evento(t, q)
	assert.Equal(t, "messagelog.created", ev["event_type"])
	ml := ev["message_log"].(map[string]any)
	assert.Equal(t, "conv-1", ml["conversation_id"])
	assert.Equal(t, "wamid.1", ml["message_id"],
		"el wamid es lo que permite casar despues el webhook de estado de Meta con este mensaje")
	assert.Equal(t, "guia_generada", ml["template_name"])
	assert.Equal(t, string(entities.MessageDirectionOutbound), ml["direction"])
}

func TestPublishMessageLogCreated_LasFechasDeEntregaYLecturaSonOpcionales(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPersistence(q).PublishMessageLogCreated(context.Background(),
		&entities.MessageLog{ID: "log-1", MessageID: "wamid.1"}))

	ml := evento(t, q)["message_log"].(map[string]any)
	assert.NotContains(t, ml, "delivered_at",
		"un mensaje recien enviado no tiene entrega ni lectura: se omiten en vez de mandar la fecha cero")
	assert.NotContains(t, ml, "read_at")
}

func TestPublishMessageLogCreated_ConFechasDeEntregaLasIncluye(t *testing.T) {
	momento := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPersistence(q).PublishMessageLogCreated(context.Background(),
		&entities.MessageLog{ID: "log-1", MessageID: "wamid.1", DeliveredAt: &momento, ReadAt: &momento}))

	ml := evento(t, q)["message_log"].(map[string]any)
	assert.Contains(t, ml, "delivered_at")
	assert.Contains(t, ml, "read_at")
}

func TestPublishMessageLogCreated_PropagaElErrorDelBroker(t *testing.T) {
	err := nuevoPersistence(brokerCaido()).PublishMessageLogCreated(context.Background(),
		&entities.MessageLog{ID: "log-1"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error publicando messagelog event")
}

func TestPublishMessageStatusUpdated_LlevaElEstadoNuevoDelMensaje(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPersistence(q).PublishMessageStatusUpdated(context.Background(),
		"wamid.1", entities.MessageStatusRead, nil))

	ev := evento(t, q)
	assert.Equal(t, "messagelog.status_updated", ev["event_type"])
	upd := ev["update"].(map[string]any)
	assert.Equal(t, "wamid.1", upd["message_id"])
	assert.Equal(t, string(entities.MessageStatusRead), upd["status"])
}

func TestPublishMessageStatusUpdated_FormateaLasMarcasDeTiempoEnRFC3339(t *testing.T) {
	momento := time.Date(2026, 8, 5, 12, 30, 0, 0, time.UTC)
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPersistence(q).PublishMessageStatusUpdated(context.Background(),
		"wamid.1", entities.MessageStatusDelivered,
		map[string]time.Time{"delivered_at": momento}))

	upd := evento(t, q)["update"].(map[string]any)
	ts := upd["timestamps"].(map[string]any)
	assert.Equal(t, "2026-08-05T12:30:00Z", ts["delivered_at"],
		"las marcas viajan como texto RFC3339, no como time.Time, para que el consumidor no dependa del formato interno de Go")
}

func TestPublishMessageStatusUpdated_SinMarcasMandaUnMapaVacioNoNulo(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPersistence(q).PublishMessageStatusUpdated(context.Background(),
		"wamid.1", entities.MessageStatusSent, nil))

	upd := evento(t, q)["update"].(map[string]any)
	assert.NotNil(t, upd["timestamps"])
	assert.Empty(t, upd["timestamps"])
}

func TestPublishMessageStatusUpdated_PropagaElErrorDelBroker(t *testing.T) {
	err := nuevoPersistence(brokerCaido()).PublishMessageStatusUpdated(context.Background(),
		"wamid.1", entities.MessageStatusRead, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error publicando message status update")
}

func TestTodosLosEventos_CompartenLaMismaColaYLlevanFechaUnix(t *testing.T) {
	q := &testkit.QueueMock{}
	p := nuevoPersistence(q)
	ctx := context.Background()

	require.NoError(t, p.PublishConversationCreated(ctx, conversacion()))
	require.NoError(t, p.PublishMessageLogCreated(ctx, &entities.MessageLog{ID: "log-1"}))
	require.NoError(t, p.PublishMessageStatusUpdated(ctx, "wamid.1", entities.MessageStatusSent, nil))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 3)
	for _, pub := range pubs {
		assert.Equal(t, rabbitmq.QueueWhatsAppPersistenceEvents, pub.Queue,
			"el consumidor discrimina por event_type dentro del cuerpo, no por la cola")
		var ev map[string]any
		require.NoError(t, json.Unmarshal(pub.Body, &ev))
		assert.Greater(t, ev["timestamp"], float64(0),
			"todos llevan fecha unix propia del publisher, no del contenido")
	}
	assert.Empty(t, q.Declaradas,
		"no declara la cola ni usa publisher confirms: si el broker rechaza, la conversacion no se persiste")
}
