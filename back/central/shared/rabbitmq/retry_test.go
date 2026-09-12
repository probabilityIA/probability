package rabbitmq

import (
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestElBackoffCreceYSeEstancaEnElUltimoTramo(t *testing.T) {
	assert.Equal(t, 10*time.Second, nextRetryDelay(1))
	assert.Equal(t, 60*time.Second, nextRetryDelay(2))
	assert.Equal(t, 300*time.Second, nextRetryDelay(3))
	assert.Equal(t, 300*time.Second, nextRetryDelay(4),
		"pasado el ultimo tramo se repite, no se desborda el arreglo")
}

func TestElBackoffNoSeCaeConIntentoCero(t *testing.T) {
	assert.Equal(t, 10*time.Second, nextRetryDelay(0))
	assert.Equal(t, 10*time.Second, nextRetryDelay(-3))
}

func TestElTotalDeEsperaEsMenorAMediaHora(t *testing.T) {
	var total time.Duration
	for i := 1; i < MaxDeliveryAttempts; i++ {
		total += nextRetryDelay(i)
	}
	assert.Less(t, total, 30*time.Minute,
		"un mensaje muerto no puede tardar horas en llegar a la DLQ")
}

func TestElContadorDeIntentosArrancaEnCeroSinCabeceras(t *testing.T) {
	assert.Equal(t, 0, retryCountFromHeaders(nil))
	assert.Equal(t, 0, retryCountFromHeaders(amqp.Table{}))
}

func TestElContadorLeeLosTiposQueMandaAMQP(t *testing.T) {
	assert.Equal(t, 3, retryCountFromHeaders(amqp.Table{headerRetryCount: int32(3)}))
	assert.Equal(t, 3, retryCountFromHeaders(amqp.Table{headerRetryCount: int64(3)}))
	assert.Equal(t, 3, retryCountFromHeaders(amqp.Table{headerRetryCount: 3}))
}

func TestElContadorIgnoraUnaCabeceraBasura(t *testing.T) {
	assert.Equal(t, 0, retryCountFromHeaders(amqp.Table{headerRetryCount: "tres"}),
		"una cabecera corrupta no puede hacer que el mensaje gire para siempre")
}

func TestLasCabecerasDeReintentoConservanLasOriginales(t *testing.T) {
	msg := amqp.Delivery{Headers: amqp.Table{"trace_id": "abc"}}

	h := headersForRetry("orders.whatsapp.confirmed", msg, assert.AnError, 1)

	assert.Equal(t, "abc", h["trace_id"])
	assert.Equal(t, int32(1), h[headerRetryCount])
	assert.Equal(t, "orders.whatsapp.confirmed", h[headerOriginQueue])
	assert.NotEmpty(t, h[headerFirstFailed])
}

func TestLaMarcaDelPrimerFalloNoSePisaEnLosSiguientesIntentos(t *testing.T) {
	msg := amqp.Delivery{Headers: amqp.Table{headerFirstFailed: "2026-09-11T00:00:00Z"}}

	h := headersForRetry("q", msg, assert.AnError, 2)

	assert.Equal(t, "2026-09-11T00:00:00Z", h[headerFirstFailed],
		"sirve para saber cuanto lleva atascado el mensaje")
}

func TestElNombreDeLaColaDeReintentoSaleDelTramo(t *testing.T) {
	assert.Equal(t, "probability.retry.10s", retryQueueName(10*time.Second))
	assert.Equal(t, "probability.retry.300s", retryQueueName(300*time.Second))
}

func TestElErrorSeRecortaParaNoInflarLaCabecera(t *testing.T) {
	largo := make([]byte, 900)
	for i := range largo {
		largo[i] = 'x'
	}
	assert.Len(t, truncar(string(largo), 500), 500)
	assert.Equal(t, "corto", truncar("corto", 500))
}
