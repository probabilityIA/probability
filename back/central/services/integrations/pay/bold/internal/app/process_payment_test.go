package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bold/internal/mocks"
)

type entornoPago struct {
	repo   *mocks.IntegrationRepositoryMock
	client *mocks.BoldClientMock
	pub    *mocks.ResponsePublisherMock
	uc     IUseCase
}

func nuevoEntornoPago() *entornoPago {
	e := &entornoPago{
		repo:   &mocks.IntegrationRepositoryMock{},
		client: &mocks.BoldClientMock{},
		pub:    &mocks.ResponsePublisherMock{},
	}
	e.uc = New(e.client, e.repo, e.pub, mocks.NewSilentLogger())
	return e
}

func solicitud() *PaymentRequestMsg {
	return &PaymentRequestMsg{
		PaymentTransactionID: 501,
		BusinessID:           26,
		GatewayCode:          "bold",
		Amount:               150000,
		Currency:             "COP",
		Reference:            "ORD-1001",
		Description:          "Pedido 1001",
		CorrelationID:        "corr-abc",
	}
}

func TestProcessPayment_CreaElLinkYPublicaLaRespuestaExitosa(t *testing.T) {
	e := nuevoEntornoPago()

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	require.Len(t, e.pub.Publicados, 1)
	resp := e.pub.Publicados[0]
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "bold", resp.GatewayCode)
	assert.EqualValues(t, 501, resp.PaymentTransactionID)
	require.NotNil(t, resp.ExternalID)
	assert.Equal(t, "LNK-1", *resp.ExternalID)
	assert.Equal(t, "https://checkout.bold.co/LNK-1", resp.GatewayResponse["checkout_url"],
		"el checkout_url es lo que el front le muestra al comprador para pagar")
	assert.Equal(t, "corr-abc", resp.CorrelationID,
		"el correlation_id ata la respuesta con la solicitud original en la cola")
}

func TestProcessPayment_PasaMontoReferenciaYDescripcionAlClienteDeBold(t *testing.T) {
	e := nuevoEntornoPago()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	require.Len(t, e.client.Llamadas, 1)
	llamada := e.client.Llamadas[0]
	assert.InDelta(t, 150000, llamada.Amount, 0.001)
	assert.Equal(t, "ORD-1001", llamada.Reference)
	assert.Equal(t, "Pedido 1001", llamada.Description)
	assert.NotNil(t, llamada.Config, "las credenciales salen de integration_types, no del mensaje de la cola")
}

func TestProcessPayment_SinMonedaAsumePesosColombianos(t *testing.T) {
	e := nuevoEntornoPago()
	msg := solicitud()
	msg.Currency = ""

	require.NoError(t, e.uc.ProcessPayment(context.Background(), msg))

	assert.Equal(t, "COP", e.client.Llamadas[0].Currency,
		"el default es COP porque Bold es una pasarela colombiana; una moneda vacia no debe reventar el cobro")
}

func TestProcessPayment_RespetaLaMonedaQueLlegaEnElMensaje(t *testing.T) {
	e := nuevoEntornoPago()
	msg := solicitud()
	msg.Currency = "USD"

	require.NoError(t, e.uc.ProcessPayment(context.Background(), msg))

	assert.Equal(t, "USD", e.client.Llamadas[0].Currency)
}

func TestProcessPayment_SiNoHayCredencialesPublicaUnErrorDeConfiguracion(t *testing.T) {
	e := nuevoEntornoPago()
	e.repo.GetBoldConfigFn = func(ctx context.Context) (*ports.BoldConfig, error) {
		return nil, errors.New("bold integration type configuration not found")
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err,
		"no devuelve error a la cola: la respuesta de fallo ya se publico, reintentar solo repetiria el mismo fallo de config")
	require.Len(t, e.pub.Publicados, 1)
	assert.Equal(t, "error", e.pub.Publicados[0].Status)
	assert.Equal(t, "config_error", e.pub.Publicados[0].ErrorCode)
	assert.Empty(t, e.client.Llamadas, "sin credenciales ni se llama a Bold")
}

func TestProcessPayment_SiBoldFallaPublicaUnErrorDeAPI(t *testing.T) {
	e := nuevoEntornoPago()
	e.client.CreatePaymentLinkFn = func(ctx context.Context, config *ports.BoldConfig, amount float64, currency, reference, description string) (string, string, error) {
		return "", "", errors.New("bold devolvio 401")
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	require.Len(t, e.pub.Publicados, 1)
	resp := e.pub.Publicados[0]
	assert.Equal(t, "error", resp.Status)
	assert.Equal(t, "api_error", resp.ErrorCode)
	assert.Contains(t, resp.Error, "401")
	assert.Nil(t, resp.ExternalID, "sin link no hay id externo que guardar")
}

func TestProcessPayment_LaRespuestaDeErrorConservaLaTransaccionYElCorrelationID(t *testing.T) {
	e := nuevoEntornoPago()
	e.repo.GetBoldConfigFn = func(ctx context.Context) (*ports.BoldConfig, error) {
		return nil, errors.New("x")
	}

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	resp := e.pub.Publicados[0]
	assert.EqualValues(t, 501, resp.PaymentTransactionID,
		"sin el id de transaccion el modulo de pagos no sabria que cobro marcar como fallido")
	assert.Equal(t, "corr-abc", resp.CorrelationID)
	assert.Equal(t, "bold", resp.GatewayCode)
}

func TestProcessPayment_SiFallaLaPublicacionElErrorSiSubeALaCola(t *testing.T) {
	fallo := errors.New("rabbitmq caido")
	e := nuevoEntornoPago()
	e.pub.PublishPaymentResponseFn = func(ctx context.Context, msg *ports.PaymentResponseMsg) error {
		return fallo
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	assert.ErrorIs(t, err, fallo,
		"si la respuesta no se pudo publicar hay que reintentar: si no, el cobro queda colgado en processing para siempre")
}

func TestProcessPayment_UnLinkCreadoQueNoSePudoResponderSeVuelveACrearAlReintentar(t *testing.T) {
	e := nuevoEntornoPago()
	e.pub.PublishPaymentResponseFn = func(ctx context.Context, msg *ports.PaymentResponseMsg) error {
		return errors.New("rabbitmq caido")
	}

	require.Error(t, e.uc.ProcessPayment(context.Background(), solicitud()))
	require.Error(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.Len(t, e.client.Llamadas, 2,
		"no hay idempotencia: cada reintento crea otro link de pago en Bold para la misma transaccion")
}

func TestProcessPayment_ReportaElTiempoDeProcesamiento(t *testing.T) {
	e := nuevoEntornoPago()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.GreaterOrEqual(t, e.pub.Publicados[0].ProcessingTimeMs, int64(0))
}
