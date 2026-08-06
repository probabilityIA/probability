package queue

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/testkit"
)

func nuevoForwarder(q *testkit.QueueMock, r *testkit.RedisMock) ports.IAIForwarder {
	return NewAIForwarder(q, r, testkit.NewSilentLogger())
}

func credsAIActiva() *testkit.RedisMock {
	return testkit.NewRedis(platformCredsRedisKey, `{"ai_sales_enabled":true,"ai_sales_demo_business_id":26}`)
}

func TestForwardToAI_ConAIActivaPublicaEnLaColaDeEntrada(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoForwarder(q, credsAIActiva()).
		ForwardToAI(context.Background(), "573001112233", "hola", "wamid.1", "text"))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	assert.Equal(t, rabbitmq.QueueWhatsAppAIIncoming, pubs[0].Queue)

	var dto map[string]any
	require.NoError(t, json.Unmarshal(pubs[0].Body, &dto))
	assert.Equal(t, "573001112233", dto["PhoneNumber"])
	assert.Equal(t, "hola", dto["MessageText"])
	assert.Equal(t, "wamid.1", dto["MessageID"])
	assert.EqualValues(t, 26, dto["BusinessID"])
	assert.Greater(t, dto["Timestamp"], float64(0))
}

func TestForwardToAI_SoloReenviaMensajesDeTexto(t *testing.T) {
	casos := []string{"image", "audio", "video", "document", "location", "sticker", ""}

	for _, tipo := range casos {
		t.Run(tipo, func(t *testing.T) {
			q := &testkit.QueueMock{}
			r := credsAIActiva()

			require.NoError(t, nuevoForwarder(q, r).
				ForwardToAI(context.Background(), "573001112233", "", "wamid.1", tipo))

			assert.Empty(t, q.Publicaciones(),
				"la IA solo sabe responder texto; un audio o una foto se descartan antes de gastar tokens")
			assert.Empty(t, r.Leidas,
				"ni siquiera se consulta Redis: el filtro por tipo va primero")
		})
	}
}

func TestForwardToAI_SinPlatformCredsEnRedisNoReenvia(t *testing.T) {
	q := &testkit.QueueMock{}
	r := testkit.NewRedis()

	err := nuevoForwarder(q, r).ForwardToAI(context.Background(), "573001112233", "hola", "wamid.1", "text")

	assert.NoError(t, err,
		"no poder leer las credenciales no es un error del webhook: se responde ok y simplemente no hay IA")
	assert.Empty(t, q.Publicaciones())
}

func TestForwardToAI_ConPlatformCredsCorruptasNoReenvia(t *testing.T) {
	q := &testkit.QueueMock{}
	r := testkit.NewRedis(platformCredsRedisKey, `{no es json`)

	err := nuevoForwarder(q, r).ForwardToAI(context.Background(), "573001112233", "hola", "wamid.1", "text")

	assert.NoError(t, err)
	assert.Empty(t, q.Publicaciones(),
		"falla cerrado hacia la IA: si no se puede confirmar que esta habilitada, no se reenvia")
}

func TestForwardToAI_ConAIDeshabilitadaNoReenvia(t *testing.T) {
	casos := map[string]string{
		"bandera en false":      `{"ai_sales_enabled":false}`,
		"bandera ausente":       `{}`,
		"bandera con tipo raro": `{"ai_sales_enabled":"si"}`,
		"bandera numerica":      `{"ai_sales_enabled":1}`,
	}

	for nombre, creds := range casos {
		t.Run(nombre, func(t *testing.T) {
			q := &testkit.QueueMock{}
			r := testkit.NewRedis(platformCredsRedisKey, creds)

			require.NoError(t, nuevoForwarder(q, r).
				ForwardToAI(context.Background(), "573001112233", "hola", "wamid.1", "text"))

			assert.Empty(t, q.Publicaciones(),
				"solo un booleano true habilita la IA: cualquier otro valor la deja apagada")
		})
	}
}

func TestForwardToAI_SinBusinessIDEnLasCredencialesUsaElUno(t *testing.T) {
	casos := map[string]string{
		"sin la clave":   `{"ai_sales_enabled":true}`,
		"con tipo texto": `{"ai_sales_enabled":true,"ai_sales_demo_business_id":"26"}`,
		"con valor nulo": `{"ai_sales_enabled":true,"ai_sales_demo_business_id":null}`,
	}

	for nombre, creds := range casos {
		t.Run(nombre, func(t *testing.T) {
			q := &testkit.QueueMock{}
			r := testkit.NewRedis(platformCredsRedisKey, creds)

			require.NoError(t, nuevoForwarder(q, r).
				ForwardToAI(context.Background(), "573001112233", "hola", "wamid.1", "text"))

			var dto map[string]any
			require.NoError(t, json.Unmarshal(q.Publicaciones()[0].Body, &dto))
			assert.EqualValues(t, 1, dto["BusinessID"],
				"el default es el negocio 1 de la fase 1 de la demo; si la clave viene mal tipada la conversacion se le atribuye a ese negocio, no al real")
		})
	}
}

func TestForwardToAI_LeElBusinessIDDeLasCredenciales(t *testing.T) {
	q := &testkit.QueueMock{}
	r := testkit.NewRedis(platformCredsRedisKey, `{"ai_sales_enabled":true,"ai_sales_demo_business_id":217}`)

	require.NoError(t, nuevoForwarder(q, r).
		ForwardToAI(context.Background(), "573001112233", "hola", "wamid.1", "text"))

	var dto map[string]any
	require.NoError(t, json.Unmarshal(q.Publicaciones()[0].Body, &dto))
	assert.EqualValues(t, 217, dto["BusinessID"])
}

func TestForwardToAI_LeeLaClaveDePlatformCredsDelTipoDeIntegracionDos(t *testing.T) {
	q := &testkit.QueueMock{}
	r := credsAIActiva()

	require.NoError(t, nuevoForwarder(q, r).
		ForwardToAI(context.Background(), "573001112233", "hola", "wamid.1", "text"))

	assert.Equal(t, []string{"integration:platform_creds:2"}, r.Leidas,
		"el 2 es el id del tipo whatsapp; la configuracion de IA vive en las credenciales de plataforma de ese tipo, no en un .env")
}

func TestForwardToAI_PropagaElErrorDelBroker(t *testing.T) {
	q := &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("canal cerrado")
		},
	}

	err := nuevoForwarder(q, credsAIActiva()).
		ForwardToAI(context.Background(), "573001112233", "hola", "wamid.1", "text")

	require.Error(t, err,
		"aca si sube el error: el cliente escribio y nadie le va a responder si el mensaje no llega a la IA")
	assert.Contains(t, err.Error(), rabbitmq.QueueWhatsAppAIIncoming)
}

func TestForwardToAI_UnErrorDeRedisNoRompeElWebhook(t *testing.T) {
	q := &testkit.QueueMock{}
	r := testkit.NewRedis()
	r.GetFn = func(ctx context.Context, key string) (string, error) {
		return "", errors.New("redis caido")
	}

	err := nuevoForwarder(q, r).ForwardToAI(context.Background(), "573001112233", "hola", "wamid.1", "text")

	assert.NoError(t, err,
		"con Redis caido el webhook de Meta igual responde 200; se pierde el reenvio a la IA pero no la recepcion del mensaje")
	assert.Empty(t, q.Publicaciones())
}
