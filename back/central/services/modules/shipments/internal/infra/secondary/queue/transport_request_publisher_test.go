package queue

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/testkit"
)

func solicitudTransporte() *domain.TransportRequestMessage {
	return &domain.TransportRequestMessage{
		Provider:          "servientrega",
		IntegrationTypeID: 12,
		Operation:         "create_guide",
		CorrelationID:     "corr-1",
		BusinessID:        26,
		IntegrationID:     197,
		Payload:           map[string]interface{}{"peso": 1.5},
	}
}

func nuevoTransporte(q rabbitmq.IQueue) domain.ITransportRequestPublisher {
	return NewTransportRequestPublisher(q, testkit.NewSilentLogger())
}

func TestPublishTransportRequest_PublicaEnLaColaUnificadaDeTransporte(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoTransporte(q).PublishTransportRequest(context.Background(), solicitudTransporte()))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	assert.Equal(t, rabbitmq.QueueTransportRequests, pubs[0].Queue,
		"todas las transportadoras comparten una sola cola; el campo provider decide quien la atiende")
}

func TestPublishTransportRequest_ElMensajeConservaProveedorOperacionYCorrelacion(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoTransporte(q).PublishTransportRequest(context.Background(), solicitudTransporte()))

	var msg domain.TransportRequestMessage
	require.NoError(t, json.Unmarshal(q.Publicaciones()[0].Body, &msg))
	assert.Equal(t, "servientrega", msg.Provider)
	assert.Equal(t, "create_guide", msg.Operation)
	assert.Equal(t, "corr-1", msg.CorrelationID,
		"el correlation_id es lo que ata esta solicitud con la respuesta que vuelve por la cola de resultados")
	assert.EqualValues(t, 26, msg.BusinessID)
	assert.EqualValues(t, 197, msg.IntegrationID)
	assert.EqualValues(t, 12, msg.IntegrationTypeID)
}

func TestPublishTransportRequest_LeAsignaFechaSiLlegaSinElla(t *testing.T) {
	q := &testkit.QueueMock{}
	req := solicitudTransporte()
	antes := time.Now()

	require.NoError(t, nuevoTransporte(q).PublishTransportRequest(context.Background(), req))

	assert.False(t, req.Timestamp.IsZero(),
		"sin fecha el consumidor no puede descartar solicitudes viejas; se rellena en el publisher")
	assert.False(t, req.Timestamp.Before(antes.Add(-time.Second)))
}

func TestPublishTransportRequest_RespetaLaFechaQueYaTraeElMensaje(t *testing.T) {
	fijo := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	q := &testkit.QueueMock{}
	req := solicitudTransporte()
	req.Timestamp = fijo

	require.NoError(t, nuevoTransporte(q).PublishTransportRequest(context.Background(), req))

	assert.True(t, req.Timestamp.Equal(fijo),
		"si el caso de uso ya fecho la solicitud, un reintento no la vuelve a fechar y sigue viendose vieja")
}

func TestPublishTransportRequest_TrasladaElPayloadDeLaTransportadora(t *testing.T) {
	q := &testkit.QueueMock{}
	req := solicitudTransporte()
	req.Payload = map[string]interface{}{"peso": 1.5, "destino": "MEDELLIN"}

	require.NoError(t, nuevoTransporte(q).PublishTransportRequest(context.Background(), req))

	var msg domain.TransportRequestMessage
	require.NoError(t, json.Unmarshal(q.Publicaciones()[0].Body, &msg))
	assert.InDelta(t, 1.5, msg.Payload["peso"], 0.001)
	assert.Equal(t, "MEDELLIN", msg.Payload["destino"])
}

func TestPublishTransportRequest_ElModoPruebaViajaEnElMensaje(t *testing.T) {
	q := &testkit.QueueMock{}
	req := solicitudTransporte()
	req.IsTest = true

	require.NoError(t, nuevoTransporte(q).PublishTransportRequest(context.Background(), req))

	var msg domain.TransportRequestMessage
	require.NoError(t, json.Unmarshal(q.Publicaciones()[0].Body, &msg))
	assert.True(t, msg.IsTest,
		"si el flag se perdiera se generaria una guia real contra la cuenta de produccion de la transportadora")
}

func TestPublishTransportRequest_SinBrokerDevuelveErrorYNoRevienta(t *testing.T) {
	err := nuevoTransporte(nil).PublishTransportRequest(context.Background(), solicitudTransporte())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rabbitmq client is nil",
		"el backend puede arrancar sin rabbit; la solicitud falla con error claro en vez de un panic por nil")
}

func TestPublishTransportRequest_PropagaElErrorDelBroker(t *testing.T) {
	q := &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("canal cerrado")
		},
	}

	err := nuevoTransporte(q).PublishTransportRequest(context.Background(), solicitudTransporte())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish transport request",
		"a diferencia del publisher SSE, aca si se devuelve error: perder una solicitud de guia deja el envio sin despachar")
}

func TestPublishTransportRequest_UnPayloadNoSerializableFallaAntesDePublicar(t *testing.T) {
	q := &testkit.QueueMock{}
	req := solicitudTransporte()
	req.Payload = map[string]interface{}{"invalido": math.Inf(1)}

	err := nuevoTransporte(q).PublishTransportRequest(context.Background(), req)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal transport request")
	assert.Empty(t, q.Publicaciones(), "no se publica un mensaje a medias")
}

func TestPublishTransportRequest_NoDeclaraLaColaNiEsperaConfirmacion(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoTransporte(q).PublishTransportRequest(context.Background(), solicitudTransporte()))

	assert.Empty(t, q.Declaradas,
		"asume que la cola ya existe; sin publisher confirms un nil no garantiza que el broker acepto la solicitud de guia")
}
