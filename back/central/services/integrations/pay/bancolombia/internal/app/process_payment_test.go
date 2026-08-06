package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/mocks"
)

type entornoPago struct {
	repo   *mocks.IntegrationRepositoryMock
	client *mocks.ClientMock
	pub    *mocks.ResponsePublisherMock
	uc     IUseCase
}

func nuevoEntornoPago() *entornoPago {
	e := &entornoPago{
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
		GatewayCode:          gatewayCode,
		Amount:               150000,
		Currency:             "COP",
		Reference:            "ORD-1001",
		Description:          "Pedido 1001",
		CorrelationID:        "corr-abc",
	}
}

func TestProcessPayment_GeneraElQRYPublicaLaRespuestaExitosa(t *testing.T) {
	e := nuevoEntornoPago()

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	require.Len(t, e.pub.Publicados, 1)
	resp := e.pub.Publicados[0]
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "bancolombia", resp.GatewayCode)
	require.NotNil(t, resp.ExternalID)
	assert.Equal(t, "QR-1", *resp.ExternalID)
	assert.Equal(t, "00020101021243650016", resp.GatewayResponse["qr_value"])
	assert.Equal(t, "iVBORw0KGgo=", resp.GatewayResponse["qr_image_base64"])
	assert.Equal(t, "sandbox", resp.GatewayResponse["environment"])
}

func TestProcessPayment_ElQRNaceEnPendingNoEnPagado(t *testing.T) {
	e := nuevoEntornoPago()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.Equal(t, "PENDING", e.pub.Publicados[0].GatewayResponse["status"],
		"generar el QR no es cobrar: el pago solo se confirma cuando llega el webhook de transferencia")
}

func TestProcessPayment_IncluyeLaRespuestaCrudaDelBancoCuandoLaHay(t *testing.T) {
	e := nuevoEntornoPago()
	e.client.GenerateQRFn = func(ctx context.Context, config *ports.BancolombiaConfig, amount float64, currency, reference, description string) (*ports.QRGenerationResult, error) {
		return &ports.QRGenerationResult{
			ExternalID:  "QR-9",
			RawResponse: map[string]interface{}{"qrId": "QR-9"},
		}, nil
	}

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.NotNil(t, e.pub.Publicados[0].GatewayResponse["raw"],
		"la respuesta cruda queda guardada para conciliar contra el extracto del banco")
}

func TestProcessPayment_SinRespuestaCrudaNoAgregaLaClave(t *testing.T) {
	e := nuevoEntornoPago()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	_, existe := e.pub.Publicados[0].GatewayResponse["raw"]
	assert.False(t, existe)
}

func TestProcessPayment_UsaLasCredencialesDelModoDelMensaje(t *testing.T) {
	for _, esPrueba := range []bool{true, false} {
		e := nuevoEntornoPago()
		msg := solicitud()
		msg.IsTest = esPrueba

		require.NoError(t, e.uc.ProcessPayment(context.Background(), msg))

		assert.Equal(t, []bool{esPrueba}, e.repo.ModosConsultados,
			"el modo lo trae el mensaje de la cola, no una variable global: sandbox y produccion conviven")
	}
}

func TestProcessPayment_SinMonedaAsumePesosColombianos(t *testing.T) {
	e := nuevoEntornoPago()
	msg := solicitud()
	msg.Currency = ""

	require.NoError(t, e.uc.ProcessPayment(context.Background(), msg))

	assert.Equal(t, "COP", e.client.Llamadas[0].Currency)
}

func TestProcessPayment_PasaMontoReferenciaYDescripcionAlBanco(t *testing.T) {
	e := nuevoEntornoPago()

	require.NoError(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	require.Len(t, e.client.Llamadas, 1)
	assert.InDelta(t, 150000, e.client.Llamadas[0].Amount, 0.001)
	assert.Equal(t, "ORD-1001", e.client.Llamadas[0].Reference,
		"la referencia es lo unico que ata el QR con la orden cuando vuelva el webhook")
	assert.Equal(t, "Pedido 1001", e.client.Llamadas[0].Description)
}

func TestProcessPayment_SinCredencialesPublicaErrorDeConfiguracion(t *testing.T) {
	e := nuevoEntornoPago()
	e.repo.GetConfigForModeFn = func(ctx context.Context, testMode bool) (*ports.BancolombiaConfig, error) {
		return nil, errors.New("configuration not found")
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	assert.Equal(t, "config_error", e.pub.Publicados[0].ErrorCode)
	assert.Equal(t, "error", e.pub.Publicados[0].Status)
	assert.Empty(t, e.client.Llamadas)
}

func TestProcessPayment_SiElBancoFallaPublicaErrorDeAPI(t *testing.T) {
	e := nuevoEntornoPago()
	e.client.GenerateQRFn = func(ctx context.Context, config *ports.BancolombiaConfig, amount float64, currency, reference, description string) (*ports.QRGenerationResult, error) {
		return nil, errors.New("bancolombia devolvio 500")
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	require.NoError(t, err)
	resp := e.pub.Publicados[0]
	assert.Equal(t, "api_error", resp.ErrorCode)
	assert.Contains(t, resp.Error, "500")
	assert.Nil(t, resp.ExternalID)
	assert.EqualValues(t, 501, resp.PaymentTransactionID)
	assert.Equal(t, "corr-abc", resp.CorrelationID)
}

func TestProcessPayment_SiFallaLaPublicacionElErrorSubeALaCola(t *testing.T) {
	fallo := errors.New("rabbitmq caido")
	e := nuevoEntornoPago()
	e.pub.PublishPaymentResponseFn = func(ctx context.Context, msg *ports.PaymentResponseMsg) error {
		return fallo
	}

	err := e.uc.ProcessPayment(context.Background(), solicitud())

	assert.ErrorIs(t, err, fallo,
		"sin respuesta publicada el cobro queda colgado, por eso si se reintenta")
}

func TestProcessPayment_NoEsIdempotente(t *testing.T) {
	e := nuevoEntornoPago()
	e.pub.PublishPaymentResponseFn = func(ctx context.Context, msg *ports.PaymentResponseMsg) error {
		return errors.New("rabbitmq caido")
	}

	require.Error(t, e.uc.ProcessPayment(context.Background(), solicitud()))
	require.Error(t, e.uc.ProcessPayment(context.Background(), solicitud()))

	assert.Len(t, e.client.Llamadas, 2,
		"cada reintento genera otro QR para la misma transaccion; el comprador podria pagar dos veces si le llegan los dos")
}
