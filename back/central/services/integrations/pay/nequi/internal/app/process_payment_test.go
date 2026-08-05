package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/pay/nequi/internal/mocks"
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
		GatewayCode:          "nequi",
		Amount:               150000,
		Currency:             "COP",
		Reference:            "ORD-1001",
		Description:          "Pedido 1001",
		CorrelationID:        "corr-abc",
	}
}

func TestProcessPayment_PublicaElQRYSuIdentificadorDeTransaccion(t *testing.T) {
	e := nuevoEntorno()

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	require.Len(t, e.pub.Publicados, 1)
	resp := e.pub.Publicados[0]
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "nequi", resp.GatewayCode)
	assert.Equal(t, "00020101021243650016", resp.GatewayResponse["qr_value"])
	assert.Equal(t, "TX-1", resp.GatewayResponse["transaction_id"])
	assert.Equal(t, "corr-abc", resp.CorrelationID)
}

func TestProcessPayment_ElIDExternoLlevaElPrefijoDeLaPasarela(t *testing.T) {
	e := nuevoEntorno()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	require.NotNil(t, e.pub.Publicados[0].ExternalID)
	assert.Equal(t, "nequi-TX-1", *e.pub.Publicados[0].ExternalID,
		"a diferencia de las otras pasarelas, aca el id externo se prefija; el id crudo queda en gateway_response.transaction_id")
}

func TestProcessPayment_NoLeMandaLaMonedaANequi(t *testing.T) {
	e := nuevoEntorno()
	msg := solicitud()
	msg.Currency = "USD"

	require.NoError(t, e.uc.ProcessPayment(context.Background(), msg))

	require.Len(t, e.client.Llamadas, 1)
	assert.InDelta(t, 150000, e.client.Llamadas[0].Amount, 0.001)
	assert.Equal(t, "ORD-1001", e.client.Llamadas[0].Reference,
		"Nequi solo opera en pesos, por eso el cliente no recibe moneda; una solicitud en USD se cobraria igual en COP")
}

func TestProcessPayment_NoLeMandaLaDescripcionANequi(t *testing.T) {
	e := nuevoEntorno()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.NotNil(t, e.client.Llamadas[0].Config,
		"las credenciales salen de integration_types, nunca del mensaje de la cola")
}

func TestProcessPayment_SinCredencialesPublicaErrorDeConfiguracion(t *testing.T) {
	e := nuevoEntorno()
	e.repo.GetNequiConfigFn = func(ctx context.Context) (*ports.NequiConfig, error) {
		return nil, errors.New("configuration not found")
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	assert.Equal(t, "config_error", e.pub.Publicados[0].ErrorCode)
	assert.Equal(t, "error", e.pub.Publicados[0].Status)
	assert.Empty(t, e.client.Llamadas)
}

func TestProcessPayment_SiNequiFallaPublicaErrorDeAPI(t *testing.T) {
	e := nuevoEntorno()
	e.client.GenerateQRFn = func(ctx context.Context, config *ports.NequiConfig, amount float64, reference string) (string, string, error) {
		return "", "", errors.New("nequi devolvio 401")
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	resp := e.pub.Publicados[0]
	assert.Equal(t, "api_error", resp.ErrorCode)
	assert.Contains(t, resp.Error, "401")
	assert.Nil(t, resp.ExternalID)
	assert.EqualValues(t, 501, resp.PaymentTransactionID)
	assert.Equal(t, "corr-abc", resp.CorrelationID)
}

func TestProcessPayment_SiFallaLaPublicacionElErrorSubeALaCola(t *testing.T) {
	fallo := errors.New("rabbitmq caido")
	e := nuevoEntorno()
	e.pub.PublishPaymentResponseFn = func(ctx context.Context, msg *ports.PaymentResponseMsg) error {
		return fallo
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	assert.ErrorIs(t, err, fallo)
}

func TestProcessPayment_NoEsIdempotente(t *testing.T) {
	e := nuevoEntorno()
	e.pub.PublishPaymentResponseFn = func(ctx context.Context, msg *ports.PaymentResponseMsg) error {
		return errors.New("rabbitmq caido")
	}

	require.Error(t, e.uc.ProcessPayment(context.Background(), solicitud()))
	require.Error(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.Len(t, e.client.Llamadas, 2,
		"cada reintento genera otro QR para la misma transaccion")
}

func TestProcessPayment_ReportaElTiempoDeProcesamiento(t *testing.T) {
	e := nuevoEntorno()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.GreaterOrEqual(t, e.pub.Publicados[0].ProcessingTimeMs, int64(0))
}
