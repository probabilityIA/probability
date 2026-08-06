package rabbitmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecidirRecuperacion_ConLaConexionVivaYLaMismaEpocaReiniciaElConsumidor(t *testing.T) {
	decision := decidirRecuperacion(false, false, true, 7, 7)

	assert.Equal(t, recuperacionReiniciar, decision,
		"es el caso del incidente del 2026-08-05: el broker cerro el canal por consumer_timeout y la conexion siguio viva, asi que nadie mas va a levantar este consumidor")
}

func TestDecidirRecuperacion_ConLaConexionCaidaDejaQueReconnectSeEncargue(t *testing.T) {
	decision := decidirRecuperacion(false, false, false, 7, 7)

	assert.Equal(t, recuperacionConexionCaida, decision,
		"si la conexion tambien murio, reregisterConsumers va a recrear todos los consumidores; reiniciar aqui los duplicaria")
}

func TestDecidirRecuperacion_ConUnaEpocaViejaNoDuplicaElConsumidor(t *testing.T) {
	decision := decidirRecuperacion(false, false, true, 8, 7)

	assert.Equal(t, recuperacionEpocaVieja, decision,
		"la conexion ya se reemplazo entre que el canal murio y que se evaluo la recuperacion: el consumidor nuevo ya existe")
}

func TestDecidirRecuperacion_ElApagadoIntencionalGanaSobreTodoLoDemas(t *testing.T) {
	decision := decidirRecuperacion(true, false, true, 7, 7)

	assert.Equal(t, recuperacionApagando, decision,
		"al cerrar el proceso los canales se cierran solos; eso no es una falla que haya que recuperar")
}

func TestDecidirRecuperacion_ElContextoCanceladoNoRestauraElConsumidor(t *testing.T) {
	decision := decidirRecuperacion(false, true, true, 7, 7)

	assert.Equal(t, recuperacionContextoCancelado, decision,
		"si el modulo dueno del consumidor cancelo su contexto, no hay que revivirlo")
}

func TestDecidirRecuperacion_ElOrdenDePrioridadEsApagadoContextoConexionEpoca(t *testing.T) {
	assert.Equal(t, recuperacionApagando, decidirRecuperacion(true, true, false, 8, 7))
	assert.Equal(t, recuperacionContextoCancelado, decidirRecuperacion(false, true, false, 8, 7))
	assert.Equal(t, recuperacionConexionCaida, decidirRecuperacion(false, false, false, 8, 7),
		"la conexion caida se evalua antes que la epoca: si no hay conexion, la epoca no importa")
}

func TestConsumerPrefetch_EntregaUnMensajePorTrabajador(t *testing.T) {
	casos := map[int]int{
		1:  1,
		3:  3,
		10: 10,
	}

	for workers, esperado := range casos {
		assert.Equal(t, esperado, consumerPrefetch(workers),
			"el prefetch acota cuantos mensajes quedan sin ACK: con uno por trabajador, la espera de un mensaje nunca depende de cuantos haya delante")
	}
}

func TestConsumerPrefetch_NuncaDevuelveCeroNiNegativo(t *testing.T) {
	for _, workers := range []int{0, -1, -50} {
		assert.Equal(t, 1, consumerPrefetch(workers),
			"un prefetch de 0 en AMQP significa ilimitado, que es justo lo contrario de lo que queremos")
	}
}

func TestConsumerPrefetch_ElValorViejoDe50EraLaCausaDelTimeout(t *testing.T) {
	const prefetchViejo = 50
	const consumerTimeoutMs = 1800000

	prefetchNuevo := consumerPrefetch(1)

	assert.Less(t, prefetchNuevo, prefetchViejo,
		"con 1 worker y prefetch 50, el mensaje 50 esperaba 49 facturas antes de su ACK y superaba los %d ms de consumer_timeout; asi murio el canal de invoicing.softpymes.requests", consumerTimeoutMs)
	assert.Equal(t, 1, prefetchNuevo)
}
