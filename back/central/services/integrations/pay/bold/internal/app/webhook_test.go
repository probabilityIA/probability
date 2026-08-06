package app

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	boldErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/mocks"
)

const secretoDePrueba = "secreto"

func firmar(body []byte, secreto string) string {
	b64 := base64.StdEncoding.EncodeToString(body)
	mac := hmac.New(sha256.New, []byte(secreto))
	mac.Write([]byte(b64))
	return hex.EncodeToString(mac.Sum(nil))
}

func casoWebhook(repo *mocks.IntegrationRepositoryMock, pub *mocks.WebhookPublisherMock) ports.IWebhookUseCase {
	if repo == nil {
		repo = &mocks.IntegrationRepositoryMock{}
	}
	if pub == nil {
		pub = &mocks.WebhookPublisherMock{}
	}
	return NewWebhookUseCase(repo, pub, mocks.NewSilentLogger())
}

func TestHandleIncomingWebhook_ConFirmaValidaPublicaElEvento(t *testing.T) {
	body := []byte(`{"id":"evt-1","type":"SALE_APPROVED","data":{"payment_id":"PAY-1","amount":50000,"currency":"COP"}}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

	require.NoError(t, err)
	require.Len(t, pub.Publicados, 1)
	msg := pub.Publicados[0]
	assert.Equal(t, "evt-1", msg.BoldEventID)
	assert.Equal(t, "SALE_APPROVED", msg.Type)
	assert.Equal(t, "PAY-1", msg.PaymentID)
	assert.InDelta(t, 50000, msg.Amount, 0.001)
	assert.Equal(t, body, msg.RawPayload,
		"el payload crudo viaja a la cola para poder auditar el pago tal como lo mando Bold")
}

func TestHandleIncomingWebhook_RechazaUnaFirmaQueNoCorresponde(t *testing.T) {
	body := []byte(`{"id":"evt-1"}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, "otro-secreto"), body, false)

	assert.ErrorIs(t, err, boldErrors.ErrInvalidSignature)
	assert.Empty(t, pub.Publicados,
		"sin firma valida nadie puede inyectar un pago aprobado falso en la cola")
}

func TestHandleIncomingWebhook_RechazaCuandoNoLlegaFirma(t *testing.T) {
	body := []byte(`{"id":"evt-1"}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), "", body, false)

	assert.ErrorIs(t, err, boldErrors.ErrInvalidSignature)
	assert.Empty(t, pub.Publicados)
}

func TestHandleIncomingWebhook_ToleraEspaciosAlrededorDeLaFirma(t *testing.T) {
	body := []byte(`{"id":"evt-1"}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), "  "+firmar(body, secretoDePrueba)+"\n", body, false)

	require.NoError(t, err)
	assert.Len(t, pub.Publicados, 1)
}

func TestHandleIncomingWebhook_UnaFirmaAlteradaEnUnByteNoPasa(t *testing.T) {
	body := []byte(`{"id":"evt-1"}`)
	firma := []byte(firmar(body, secretoDePrueba))
	firma[0] ^= 0x01

	err := casoWebhook(nil, nil).HandleIncomingWebhook(context.Background(), string(firma), body, false)

	assert.ErrorIs(t, err, boldErrors.ErrInvalidSignature)
}

func TestHandleIncomingWebhook_LaFirmaDependeDelCuerpoExacto(t *testing.T) {
	original := []byte(`{"id":"evt-1","data":{"amount":1000}}`)
	manipulado := []byte(`{"id":"evt-1","data":{"amount":9000}}`)

	err := casoWebhook(nil, nil).HandleIncomingWebhook(context.Background(), firmar(original, secretoDePrueba), manipulado, false)

	assert.ErrorIs(t, err, boldErrors.ErrInvalidSignature,
		"cambiar el monto invalida la firma: no se puede aprobar un pago por mas plata de la real")
}

func TestHandleIncomingWebhook_SinSecretoConfiguradoRechazaTodo(t *testing.T) {
	body := []byte(`{"id":"evt-1"}`)
	repo := &mocks.IntegrationRepositoryMock{
		GetBoldConfigForModeFn: func(ctx context.Context, testMode bool) (*ports.BoldConfig, error) {
			return &ports.BoldConfig{APIKey: "api-key", Secret: ""}, nil
		},
	}

	err := casoWebhook(repo, nil).HandleIncomingWebhook(context.Background(), firmar(body, ""), body, false)

	assert.ErrorIs(t, err, boldErrors.ErrInvalidSignature,
		"con el secreto vacio se rechaza en vez de firmar con cadena vacia, que aceptaria cualquier payload")
}

func TestHandleIncomingWebhook_RechazaUnCuerpoVacioAntesDeConsultarLaConfig(t *testing.T) {
	repo := &mocks.IntegrationRepositoryMock{}

	err := casoWebhook(repo, nil).HandleIncomingWebhook(context.Background(), "firma", nil, false)

	require.Error(t, err)
	assert.Empty(t, repo.ModosConsultados)
}

func TestHandleIncomingWebhook_UsaLasCredencialesDelModoQueLlega(t *testing.T) {
	for _, esPrueba := range []bool{true, false} {
		body := []byte(`{"id":"evt-1"}`)
		repo := &mocks.IntegrationRepositoryMock{}
		pub := &mocks.WebhookPublisherMock{}

		err := casoWebhook(repo, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, esPrueba)

		require.NoError(t, err)
		assert.Equal(t, []bool{esPrueba}, repo.ModosConsultados,
			"sandbox y produccion tienen secretos distintos: usar el equivocado tumbaria la validacion")
		assert.Equal(t, esPrueba, pub.Publicados[0].IsTest,
			"el flag viaja al consumidor para que un pago de sandbox no marque una orden real como pagada")
	}
}

func TestHandleIncomingWebhook_SiNoSePuedeLeerLaConfigNoValidaNiPublica(t *testing.T) {
	fallo := errors.New("timeout")
	repo := &mocks.IntegrationRepositoryMock{
		GetBoldConfigForModeFn: func(ctx context.Context, testMode bool) (*ports.BoldConfig, error) {
			return nil, fallo
		},
	}
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(repo, pub).HandleIncomingWebhook(context.Background(), "firma", []byte(`{"id":"x"}`), false)

	assert.ErrorIs(t, err, fallo)
	assert.Empty(t, pub.Publicados)
}

func TestHandleIncomingWebhook_RechazaUnEventoSinID(t *testing.T) {
	body := []byte(`{"type":"SALE_APPROVED"}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing event id",
		"el id del evento es la clave de idempotencia aguas abajo: sin el, el pago podria aplicarse dos veces")
	assert.Empty(t, pub.Publicados)
}

func TestHandleIncomingWebhook_RechazaUnJSONInvalidoAunConFirmaValida(t *testing.T) {
	body := []byte(`{no es json`)

	err := casoWebhook(nil, nil).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode bold webhook envelope")
}

func TestHandleIncomingWebhook_ConvierteElTimestampEnNanosegundos(t *testing.T) {
	momento := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	body := []byte(`{"id":"evt-1","time":` + itoa(momento.UnixNano()) + `}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

	require.NoError(t, err)
	require.NotNil(t, pub.Publicados[0].OccurredAt)
	assert.True(t, pub.Publicados[0].OccurredAt.Equal(momento),
		"Bold manda el time en nanosegundos epoch; interpretarlo como segundos daria una fecha del ano 1970")
}

func TestHandleIncomingWebhook_SinTimestampNoInventaUnaFecha(t *testing.T) {
	casos := []string{`{"id":"evt-1"}`, `{"id":"evt-1","time":0}`, `{"id":"evt-1","time":-5}`}

	for _, crudo := range casos {
		body := []byte(crudo)
		pub := &mocks.WebhookPublisherMock{}

		err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

		require.NoError(t, err)
		assert.Nil(t, pub.Publicados[0].OccurredAt,
			"mejor nil que time.Now(): la fecha del evento no es la de recepcion")
	}
}

func TestHandleIncomingWebhook_BuscaLaReferenciaEnLosTresLugaresDondeBoldLaPuedeMandar(t *testing.T) {
	casos := []struct {
		nombre   string
		crudo    string
		esperado string
	}{
		{
			"merchant_reference directo",
			`{"id":"e","data":{"merchant_reference":"ORD-1","metadata":{"reference":"ORD-2","merchant_reference":"ORD-3"}}}`,
			"ORD-1",
		},
		{
			"metadata.reference como segundo intento",
			`{"id":"e","data":{"metadata":{"reference":"ORD-2","merchant_reference":"ORD-3"}}}`,
			"ORD-2",
		},
		{
			"metadata.merchant_reference como ultimo recurso",
			`{"id":"e","data":{"metadata":{"merchant_reference":"ORD-3"}}}`,
			"ORD-3",
		},
		{
			"sin referencia en ningun lado",
			`{"id":"e","data":{}}`,
			"",
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			body := []byte(caso.crudo)
			pub := &mocks.WebhookPublisherMock{}

			err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

			require.NoError(t, err)
			assert.Equal(t, caso.esperado, pub.Publicados[0].MerchantReference,
				"la referencia es lo que ata el pago a la orden; Bold la manda en sitios distintos segun el flujo")
		})
	}
}

func TestHandleIncomingWebhook_SinReferenciaIgualPublicaElEvento(t *testing.T) {
	body := []byte(`{"id":"evt-1","data":{"payment_id":"PAY-1"}}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

	require.NoError(t, err)
	assert.Len(t, pub.Publicados, 1,
		"un evento sin referencia llega igual a la cola: el consumidor decide si lo puede conciliar o lo descarta")
}

func TestHandleIncomingWebhook_PrefiereLaMonedaDeLaRaizSobreLaDelMonto(t *testing.T) {
	casos := []struct {
		nombre   string
		crudo    string
		esperado string
	}{
		{"solo en la raiz", `{"id":"e","data":{"currency":"COP"}}`, "COP"},
		{"solo dentro del monto", `{"id":"e","data":{"amount":{"value":100,"currency":"USD"}}}`, "USD"},
		{"en ambos gana la raiz", `{"id":"e","data":{"currency":"COP","amount":{"value":100,"currency":"USD"}}}`, "COP"},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			body := []byte(caso.crudo)
			pub := &mocks.WebhookPublisherMock{}

			err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

			require.NoError(t, err)
			assert.Equal(t, caso.esperado, pub.Publicados[0].Currency)
		})
	}
}

func TestHandleIncomingWebhook_LeeElMontoEnTodasLasFormasQueBoldEnvia(t *testing.T) {
	casos := []struct {
		nombre   string
		crudo    string
		esperado float64
	}{
		{"numero plano", `{"id":"e","data":{"amount":50000}}`, 50000},
		{"objeto con value", `{"id":"e","data":{"amount":{"value":50000}}}`, 50000},
		{"objeto con total", `{"id":"e","data":{"amount":{"total":50000}}}`, 50000},
		{"objeto con amount", `{"id":"e","data":{"amount":{"amount":50000}}}`, 50000},
		{"value gana sobre total", `{"id":"e","data":{"amount":{"value":100,"total":900}}}`, 100},
		{"total gana sobre amount", `{"id":"e","data":{"amount":{"total":100,"amount":900}}}`, 100},
		{"con decimales", `{"id":"e","data":{"amount":{"value":50000.75}}}`, 50000.75},
		{"nulo", `{"id":"e","data":{"amount":null}}`, 0},
		{"ausente", `{"id":"e","data":{}}`, 0},
		{"objeto sin ningun campo de monto", `{"id":"e","data":{"amount":{"currency":"COP"}}}`, 0},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			body := []byte(caso.crudo)
			pub := &mocks.WebhookPublisherMock{}

			err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

			require.NoError(t, err)
			assert.InDelta(t, caso.esperado, pub.Publicados[0].Amount, 0.001,
				"Bold cambio el formato del monto entre versiones; el parser acepta ambos para no perder pagos")
		})
	}
}

func TestHandleIncomingWebhook_UnMontoConFormatoIlegibleAbortaElEvento(t *testing.T) {
	casos := []string{
		`{"id":"e","data":{"amount":"cincuenta mil"}}`,
		`{"id":"e","data":{"amount":{"value":"abc"}}}`,
	}

	for _, crudo := range casos {
		body := []byte(crudo)
		pub := &mocks.WebhookPublisherMock{}

		err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

		require.Error(t, err,
			"un monto que no se puede leer no se publica como cero: seria conciliar un pago por el valor equivocado")
		assert.Empty(t, pub.Publicados)
	}
}

func TestHandleIncomingWebhook_CopiaElRestoDeCamposDelEvento(t *testing.T) {
	body := []byte(`{"id":"evt-1","type":"SALE_APPROVED","subject":"pago","source":"bold",
		"data":{"payment_id":"PAY-1","payment_method":"CARD","payer_email":"cliente@correo.com"}}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

	require.NoError(t, err)
	msg := pub.Publicados[0]
	assert.Equal(t, "pago", msg.Subject)
	assert.Equal(t, "bold", msg.Source)
	assert.Equal(t, "CARD", msg.PaymentMethod)
	assert.Equal(t, "cliente@correo.com", msg.PayerEmail)
}

func TestHandleIncomingWebhook_SiFallaLaPublicacionDevuelveErrorParaQueBoldReintente(t *testing.T) {
	body := []byte(`{"id":"evt-1"}`)
	pub := &mocks.WebhookPublisherMock{
		PublishWebhookEventFn: func(ctx context.Context, msg *ports.BoldWebhookMessage) error {
			return errors.New("rabbitmq caido")
		},
	}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firmar(body, secretoDePrueba), body, false)

	require.Error(t, err,
		"a diferencia del webhook de Grafana, aca si se devuelve error: perder un pago es peor que un reintento duplicado")
	assert.Contains(t, err.Error(), "publish bold webhook event")
}

func TestSecretKey_SoloEsValidoCuandoHaySecreto(t *testing.T) {
	var nulo *ports.BoldConfig
	_, ok := nulo.SecretKey()
	assert.False(t, ok, "la config nula no revienta, devuelve no-ok")

	_, ok = (&ports.BoldConfig{Secret: ""}).SecretKey()
	assert.False(t, ok)

	secreto, ok := (&ports.BoldConfig{Secret: "abc"}).SecretKey()
	assert.True(t, ok)
	assert.Equal(t, "abc", secreto)
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	negativo := n < 0
	if negativo {
		n = -n
	}
	var buf [24]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if negativo {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
