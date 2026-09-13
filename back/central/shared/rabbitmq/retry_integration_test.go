package rabbitmq_test

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/stretchr/testify/require"
)

func brokerLocal(t *testing.T) (rabbitmq.IQueue, *amqp.Channel) {
	t.Helper()

	usuario := valorODefecto("RABBITMQ_USER", "admin")
	clave := valorODefecto("RABBITMQ_PASS", "admin")
	url := fmt.Sprintf("amqp://%s:%s@localhost:5672/", usuario, clave)

	conn, err := amqp.DialConfig(url, amqp.Config{Dial: amqp.DefaultDial(2 * time.Second)})
	if err != nil {
		t.Skipf("sin RabbitMQ en localhost:5672, se omite la prueba de integracion: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	ch, err := conn.Channel()
	require.NoError(t, err)
	t.Cleanup(func() { ch.Close() })

	os.Setenv("RABBITMQ_HOST", "localhost")
	os.Setenv("RABBITMQ_PORT", "5672")
	os.Setenv("RABBITMQ_USER", usuario)
	os.Setenv("RABBITMQ_PASS", clave)
	os.Setenv("RABBITMQ_VHOST", "/")

	lg := log.New()
	q, err := rabbitmq.New(lg, env.New(lg))
	require.NoError(t, err)
	t.Cleanup(func() { q.Close() })

	return q, ch
}

func valorODefecto(clave, defecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return defecto
}

func TestIntegracionUnMensajeVenenosoDejaDeGirar(t *testing.T) {
	q, ch := brokerLocal(t)

	cola := fmt.Sprintf("test.veneno.%d", time.Now().UnixNano())
	require.NoError(t, q.DeclareQueue(cola, true))
	t.Cleanup(func() { ch.QueueDelete(cola, false, false, false) })

	var intentos int64
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	require.NoError(t, q.Consume(ctx, cola, func([]byte) error {
		atomic.AddInt64(&intentos, 1)
		return fmt.Errorf("order not found")
	}))

	require.NoError(t, q.Publish(ctx, cola, []byte(`{"order_number":"PRUEBA-ISABEL"}`)))
	time.Sleep(20 * time.Second)

	n := atomic.LoadInt64(&intentos)
	require.LessOrEqual(t, n, int64(3),
		"el bucle sin backoff hacia 657 entregas POR SEGUNDO; con backoff caben 2 en 20 s")
	require.GreaterOrEqual(t, n, int64(2),
		"el tramo de 10 s debe devolver el mensaje a su cola original")
}

func TestIntegracionElMensajeBuenoNoQuedaAtrapadoDetrasDelVenenoso(t *testing.T) {
	q, ch := brokerLocal(t)

	cola := fmt.Sprintf("test.mixto.%d", time.Now().UnixNano())
	require.NoError(t, q.DeclareQueue(cola, true))
	t.Cleanup(func() { ch.QueueDelete(cola, false, false, false) })

	var buenos int64
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	require.NoError(t, q.Consume(ctx, cola, func(b []byte) error {
		if string(b) == "VENENO" {
			return fmt.Errorf("order not found")
		}
		atomic.AddInt64(&buenos, 1)
		return nil
	}))

	require.NoError(t, q.Publish(ctx, cola, []byte("VENENO")))
	require.NoError(t, q.Publish(ctx, cola, []byte("VIG-0158")))
	time.Sleep(5 * time.Second)

	require.Equal(t, int64(1), atomic.LoadInt64(&buenos))
}

func TestIntegracionAlAgotarLosIntentosCaeEnLaDLQConSuOrigen(t *testing.T) {
	q, ch := brokerLocal(t)

	cola := fmt.Sprintf("test.agotado.%d", time.Now().UnixNano())
	require.NoError(t, q.DeclareQueue(cola, true))
	t.Cleanup(func() { ch.QueueDelete(cola, false, false, false) })

	var intentos int64
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	require.NoError(t, q.Consume(ctx, cola, func([]byte) error {
		atomic.AddInt64(&intentos, 1)
		return fmt.Errorf("order not found")
	}))

	antes, err := ch.QueueInspect(rabbitmq.DeadLetterQueue)
	require.NoError(t, err)

	cuerpo := fmt.Sprintf(`{"order_number":"PRUEBA-ISABEL","t":%d}`, time.Now().UnixNano())
	require.NoError(t, ch.PublishWithContext(ctx, "", cola, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        []byte(cuerpo),
		Headers:     amqp.Table{"x-retry-count": int32(rabbitmq.MaxDeliveryAttempts - 1)},
	}))
	time.Sleep(3 * time.Second)

	require.Equal(t, int64(1), atomic.LoadInt64(&intentos),
		"agotados los intentos no hay una entrega mas: ahi termina el bucle")

	despues, err := ch.QueueInspect(rabbitmq.DeadLetterQueue)
	require.NoError(t, err)
	require.Equal(t, antes.Messages+1, despues.Messages, "el mensaje se conserva, no se pierde")

	entregado, ok, err := ch.Get(rabbitmq.DeadLetterQueue, true)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, cuerpo, string(entregado.Body))
	require.Equal(t, cola, entregado.Headers["x-origin-queue"])
	require.Contains(t, entregado.Headers["x-last-error"], "order not found")
}
