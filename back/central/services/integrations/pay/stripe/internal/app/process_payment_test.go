package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/pay/stripe/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/pay/stripe/internal/mocks"
)

type entorno struct {
	repo   *mocks.IntegrationRepositoryMock
	client *mocks.ClientMock
	pub    *mocks.ResponsePublisherMock
	uc     IUseCase
}

func nuevoEntorno() *entorno {
	e := &entorno{
		repo:   &mocks.IntegrationRepositoryMock{},
		client: &mocks.ClientMock{},
		pub:    &mocks.ResponsePublisherMock{},
	}
	e.uc = New(e.client, e.repo, e.pub, mocks.NewSilentLogger())
	return e
}

func solicitud() *PaymentRequestMsg {
	return &PaymentRequestMsg{
		PaymentTransactionID: 501,
		BusinessID:           26,
		GatewayCode:          "stripe",
		Amount:               150000,
		Currency:             "usd",
		Reference:            "ORD-1001",
		Description:          "Pedido 1001",
		CorrelationID:        "corr-abc",
	}
}

func TestProcessPayment_PublicaLaRespuestaExitosaConElIDExterno(t *testing.T) {
	e := nuevoEntorno()

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	require.Len(t, e.pub.Publicados, 1)
	resp := e.pub.Publicados[0]
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "stripe", resp.GatewayCode)
	assert.EqualValues(t, 501, resp.PaymentTransactionID)
	require.NotNil(t, resp.ExternalID)
	assert.Equal(t, "pi_123", *resp.ExternalID,
		"el id externo es lo que permite conciliar despues contra la pasarela")
	assert.Equal(t, "pi_123", resp.GatewayResponse["payment_intent_id"])
	assert.Equal(t, "pi_123_secret_abc", resp.GatewayResponse["client_secret"])
	assert.Equal(t, "corr-abc", resp.CorrelationID)
}

func TestProcessPayment_PasaMontoReferenciaYDescripcionALaPasarela(t *testing.T) {
	e := nuevoEntorno()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	require.Len(t, e.client.Llamadas, 1)
	llamada := e.client.Llamadas[0]
	assert.InDelta(t, 150000, llamada.Amount, 0.001)
	assert.Equal(t, "ORD-1001", llamada.Reference)
	assert.Equal(t, "Pedido 1001", llamada.Description)
	assert.NotNil(t, llamada.Config,
		"las credenciales salen de integration_types, nunca del mensaje de la cola")
}

func TestProcessPayment_SinMonedaUsaElDefaultDeLaPasarela(t *testing.T) {
	e := nuevoEntorno()
	msg := solicitud()
	msg.Currency = ""

	require.NoError(t, e.uc.ProcessPayment(context.Background(), msg))

	assert.Equal(t, "usd", e.client.Llamadas[0].Currency,
		"una moneda vacia no debe reventar el cobro")
}

func TestProcessPayment_RespetaLaMonedaDelMensaje(t *testing.T) {
	e := nuevoEntorno()
	msg := solicitud()
	msg.Currency = "EUR"

	require.NoError(t, e.uc.ProcessPayment(context.Background(), msg))

	assert.Equal(t, "EUR", e.client.Llamadas[0].Currency)
}

func TestProcessPayment_SinCredencialesPublicaErrorDeConfiguracion(t *testing.T) {
	e := nuevoEntorno()
	e.repo.GetStripeConfigFn = func(ctx context.Context) (*ports.StripeConfig, error) {
		return nil, errors.New("configuration not found")
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err,
		"la respuesta de fallo ya se publico; devolver error solo repetiria el mismo fallo de config en el reintento")
	require.Len(t, e.pub.Publicados, 1)
	assert.Equal(t, "error", e.pub.Publicados[0].Status)
	assert.Equal(t, "config_error", e.pub.Publicados[0].ErrorCode)
	assert.Empty(t, e.client.Llamadas, "sin credenciales ni se llama a la pasarela")
}

func TestProcessPayment_SiLaPasarelaFallaPublicaErrorDeAPI(t *testing.T) {
	e := nuevoEntorno()
	e.client.CreatePaymentIntentFn = func(ctx context.Context, config *ports.StripeConfig, amount float64, currency, reference, description string) (string, string, error) {
		return "", "", errors.New("la pasarela devolvio 401")
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	resp := e.pub.Publicados[0]
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "api_error", resp.ErrorCode)
	assert.Contains(t, resp.Error, "401")
	assert.Nil(t, resp.ExternalID)
	assert.EqualValues(t, 501, resp.PaymentTransactionID,
		"sin el id de transaccion el modulo de pagos no sabria que cobro marcar como fallido")
	assert.Equal(t, "corr-abc", resp.CorrelationID)
}

func TestProcessPayment_SiFallaLaPublicacionElErrorSubeALaCola(t *testing.T) {
	fallo := errors.New("rabbitmq caido")
	e := nuevoEntorno()
	e.pub.PublishPaymentResponseFn = func(ctx context.Context, msg *ports.PaymentResponseMsg) error {
		return fallo
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	assert.ErrorIs(t, err, fallo,
		"si la respuesta no se publica el cobro queda colgado en processing para siempre")
}

func TestProcessPayment_NoEsIdempotente(t *testing.T) {
	e := nuevoEntorno()
	e.pub.PublishPaymentResponseFn = func(ctx context.Context, msg *ports.PaymentResponseMsg) error {
		return errors.New("rabbitmq caido")
	}

	require.Error(t, e.uc.ProcessPayment(context.Background(), solicitud()))
	require.Error(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.Len(t, e.client.Llamadas, 2,
		"cada reintento abre otro cobro en la pasarela para la misma transaccion")
}

func TestProcessPayment_ReportaElTiempoDeProcesamiento(t *testing.T) {
	e := nuevoEntorno()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.GreaterOrEqual(t, e.pub.Publicados[0].ProcessingTimeMs, int64(0))
}
