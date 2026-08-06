package queue

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/testkit"
)

func evento() entities.AlertEvent {
	return entities.AlertEvent{
		AlertType: "RAM",
		Summary:   "RAM al 95% en back-central",
		Status:    "firing",
		FiredAt:   time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC),
	}
}

func TestPublish_PublicaEnLaColaDeAlertas(t *testing.T) {
	q := &testkit.QueueMock{}

	err := New(q, testkit.NewSilentLogger()).Publish(context.Background(), evento())

	require.NoError(t, err)
	require.Len(t, q.Publicados, 1)
	assert.Equal(t, rabbitmq.QueueMonitoringAlerts, q.Publicados[0].Queue,
		"el nombre de la cola sale de la constante compartida: si cambia ahi, el consumidor y el publisher no se desincronizan")
}

func TestPublish_DeclaraLaColaComoDurableAntesDePublicar(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, New(q, testkit.NewSilentLogger()).Publish(context.Background(), evento()))

	require.Len(t, q.Declaradas, 1)
	assert.Equal(t, rabbitmq.QueueMonitoringAlerts, q.Declaradas[0].Name)
	assert.True(t, q.Declaradas[0].Durable,
		"la cola es durable, asi sobrevive un reinicio del broker. Ojo: eso NO hace persistente al mensaje, que es un problema aparte")
}

func TestPublish_SerializaElEventoComoJSONConTodosSusCampos(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, New(q, testkit.NewSilentLogger()).Publish(context.Background(), evento()))

	var recibido entities.AlertEvent
	require.NoError(t, json.Unmarshal(q.Publicados[0].Body, &recibido))
	assert.Equal(t, "RAM", recibido.AlertType)
	assert.Equal(t, "RAM al 95% en back-central", recibido.Summary)
	assert.Equal(t, "firing", recibido.Status)
	assert.True(t, recibido.FiredAt.Equal(evento().FiredAt),
		"la fecha del disparo debe sobrevivir el viaje por la cola sin perder la zona horaria")
}

func TestPublish_SiFallaLaDeclaracionNoPublicaNada(t *testing.T) {
	fallo := errors.New("canal cerrado")
	q := &testkit.QueueMock{
		DeclareQueueFn: func(queueName string, durable bool) error { return fallo },
	}

	err := New(q, testkit.NewSilentLogger()).Publish(context.Background(), evento())

	assert.ErrorIs(t, err, fallo)
	assert.Empty(t, q.Publicados)
}

func TestPublish_PropagaElErrorDePublicacion(t *testing.T) {
	fallo := errors.New("broker caido")
	q := &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error { return fallo },
	}

	err := New(q, testkit.NewSilentLogger()).Publish(context.Background(), evento())

	assert.ErrorIs(t, err, fallo,
		"el caso de uso se traga este error a proposito para que Grafana no reintente; la alerta se pierde aqui")
}

func TestPublish_NoHayConfirmacionDelBroker(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, New(q, testkit.NewSilentLogger()).Publish(context.Background(), evento()))

	assert.Len(t, q.Publicados, 1,
		"Publish devuelve nil apenas entrega el mensaje al canal: sin publisher confirms, un nil aqui no garantiza que el broker lo haya aceptado")
}
