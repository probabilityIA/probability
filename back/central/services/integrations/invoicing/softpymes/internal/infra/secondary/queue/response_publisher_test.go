package queue

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/testkit"
)

func nuevoPublisher(q *testkit.QueueMock) *ResponsePublisher {
	if q == nil {
		return New(nil, testkit.NewSilentLogger())
	}
	return New(q, testkit.NewSilentLogger())
}

func respuestaFactura() *InvoiceResponseMessage {
	return &InvoiceResponseMessage{
		InvoiceID:     501,
		Provider:      "softpymes",
		Status:        "success",
		Operation:     "create",
		InvoiceNumber: "FE-1001",
		CorrelationID: "corr-1",
	}
}

func decodificar(t *testing.T, q *testkit.QueueMock, destino any) {
	t.Helper()
	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	require.NoError(t, json.Unmarshal(pubs[0].Body, destino))
}

func TestPublishResponse_PublicaEnLaColaDeRespuestasDeFacturacion(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), respuestaFactura()))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	assert.Equal(t, rabbitmq.QueueInvoicingResponses, pubs[0].Queue,
		"todos los proveedores de facturacion responden a la misma cola; el campo provider distingue quien contesto")
}

func TestPublishResponse_ConservaLosDatosDeLaFacturaEmitida(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := respuestaFactura()
	resp.CUFE = "CUFE-ABC"
	resp.PDFURL = "https://s3/factura.pdf"
	resp.ExternalID = "SP-9"

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), resp))

	var msg InvoiceResponseMessage
	decodificar(t, q, &msg)
	assert.EqualValues(t, 501, msg.InvoiceID,
		"sin invoice_id el modulo de facturacion no sabe que factura marcar como emitida")
	assert.Equal(t, "FE-1001", msg.InvoiceNumber)
	assert.Equal(t, "CUFE-ABC", msg.CUFE,
		"el CUFE es el identificador fiscal de la DIAN: si se pierde, la factura no se puede validar despues")
	assert.Equal(t, "https://s3/factura.pdf", msg.PDFURL)
	assert.Equal(t, "corr-1", msg.CorrelationID)
}

func TestPublishResponse_LaRespuestaDeErrorLlevaCodigoYDetalle(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := respuestaFactura()
	resp.Status = "error"
	resp.Error = "NIT del cliente invalido"
	resp.ErrorCode = "validation_error"
	resp.ErrorDetails = map[string]interface{}{"campo": "nit"}

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), resp))

	var msg InvoiceResponseMessage
	decodificar(t, q, &msg)
	assert.Equal(t, "error", msg.Status)
	assert.Equal(t, "validation_error", msg.ErrorCode)
	assert.Equal(t, "nit", msg.ErrorDetails["campo"])
}

func TestPublishResponse_LlevaLaAuditoriaDelHTTPAlProveedor(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := respuestaFactura()
	resp.AuditRequestURL = "https://api.softpymes.com/invoices"
	resp.AuditRequestPayload = map[string]interface{}{"total": 150000}
	resp.AuditResponseStatus = 201
	resp.AuditResponseBody = `{"id":9}`
	resp.CashReceiptRequestURL = "https://api.softpymes.com/receipts"
	resp.CashReceiptResponseStatus = 200

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), resp))

	var msg InvoiceResponseMessage
	decodificar(t, q, &msg)
	assert.Equal(t, "https://api.softpymes.com/invoices", msg.AuditRequestURL,
		"la auditoria del HTTP crudo es lo que permite reclamarle al proveedor cuando una factura sale mal")
	assert.Equal(t, 201, msg.AuditResponseStatus)
	assert.Equal(t, 200, msg.CashReceiptResponseStatus,
		"la factura y el recibo de caja son dos llamadas distintas y cada una se audita aparte")
}

func TestPublishResponse_LeAsignaFechaSiLlegaSinElla(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := respuestaFactura()

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), resp))

	assert.False(t, resp.Timestamp.IsZero())
}

func TestPublishResponse_RespetaLaFechaQueYaTrae(t *testing.T) {
	fijo := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	q := &testkit.QueueMock{}
	resp := respuestaFactura()
	resp.Timestamp = fijo

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), resp))

	assert.True(t, resp.Timestamp.Equal(fijo))
}

func TestPublishResponse_SinBrokerDevuelveNilYLaFacturaQuedaColgada(t *testing.T) {
	err := nuevoPublisher(nil).PublishResponse(context.Background(), respuestaFactura())

	assert.NoError(t, err,
		"con rabbit nil devuelve nil a proposito para no romper el flujo. La factura YA se emitio en Softpymes pero el modulo de facturacion nunca se entera: queda en processing para siempre. Contrasta con el publisher de transporte, que si devuelve error")
}

func TestPublishResponse_PropagaElErrorDelBroker(t *testing.T) {
	q := &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("canal cerrado")
		},
	}

	err := nuevoPublisher(q).PublishResponse(context.Background(), respuestaFactura())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish response",
		"un fallo del broker si sube; solo el caso de rabbit nil se traga")
}

func TestPublishResponse_UnPayloadNoSerializableNoSePublicaAMedias(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := respuestaFactura()
	resp.DocumentJSON = map[string]interface{}{"malo": math.Inf(1)}

	err := nuevoPublisher(q).PublishResponse(context.Background(), resp)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal response")
	assert.Empty(t, q.Publicaciones())
}

func TestPublishCompareResponse_LlevaElRangoDeFechasYLosDocumentosDelProveedor(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &CompareResponseMessage{
		Operation:     "compare",
		CorrelationID: "corr-1",
		BusinessID:    26,
		DateFrom:      "2026-08-01",
		DateTo:        "2026-08-05",
		ProviderDocuments: []CompareDocument{{
			DocumentNumber: "FE-1001",
			Total:          "150000",
			CustomerNit:    "900123456",
			Annuled:        false,
			Details:        []CompareDocumentDetail{{ItemCode: "SKU-1", Quantity: "2", Value: "50000"}},
		}},
	}

	require.NoError(t, nuevoPublisher(q).PublishCompareResponse(context.Background(), resp))

	var msg CompareResponseMessage
	decodificar(t, q, &msg)
	assert.Equal(t, "compare", msg.Operation)
	assert.EqualValues(t, 26, msg.BusinessID,
		"la comparacion cruza documentos de un negocio contra los del proveedor: sin business_id se cruzarian catalogos ajenos")
	assert.Equal(t, "2026-08-01", msg.DateFrom)
	require.Len(t, msg.ProviderDocuments, 1)
	assert.Equal(t, "FE-1001", msg.ProviderDocuments[0].DocumentNumber)
	require.Len(t, msg.ProviderDocuments[0].Details, 1)
	assert.Equal(t, "SKU-1", msg.ProviderDocuments[0].Details[0].ItemCode)
}

func TestPublishCompareResponse_MarcaLosDocumentosAnulados(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &CompareResponseMessage{
		CorrelationID:     "corr-1",
		ProviderDocuments: []CompareDocument{{DocumentNumber: "FE-1", Annuled: true, ElectronicDocument: true}},
	}

	require.NoError(t, nuevoPublisher(q).PublishCompareResponse(context.Background(), resp))

	var msg CompareResponseMessage
	decodificar(t, q, &msg)
	assert.True(t, msg.ProviderDocuments[0].Annuled,
		"la bandera de anulada es lo que dispara liberar la orden para re-facturar")
	assert.True(t, msg.ProviderDocuments[0].ElectronicDocument)
}

func TestPublishCompareResponse_UnaComparacionSinDocumentosSePublicaIgual(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPublisher(q).PublishCompareResponse(context.Background(),
		&CompareResponseMessage{CorrelationID: "corr-1"}))

	assert.Len(t, q.Publicaciones(), 1,
		"cero documentos es un resultado valido de la comparacion, no un error")
}

func TestPublishCompareResponse_SinBrokerTambienDevuelveNil(t *testing.T) {
	err := nuevoPublisher(nil).PublishCompareResponse(context.Background(),
		&CompareResponseMessage{CorrelationID: "corr-1"})

	assert.NoError(t, err)
}

func TestPublishCompareResponse_PropagaElErrorDelBroker(t *testing.T) {
	q := &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("canal cerrado")
		},
	}

	err := nuevoPublisher(q).PublishCompareResponse(context.Background(),
		&CompareResponseMessage{CorrelationID: "corr-1"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish compare response")
}

func TestPublishListItemsResponse_LlevaElCatalogoDelProveedor(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &ListItemsResponseMessage{
		Operation:     "list_items",
		CorrelationID: "corr-1",
		BusinessID:    26,
		Items: []ListItemsItem{
			{ItemCode: "SKU-1", ItemName: "Camiseta", ItemPrice: 50000, UnitCost: 20000},
		},
	}

	require.NoError(t, nuevoPublisher(q).PublishListItemsResponse(context.Background(), resp))

	var msg ListItemsResponseMessage
	decodificar(t, q, &msg)
	assert.Equal(t, "list_items", msg.Operation)
	require.Len(t, msg.Items, 1)
	assert.Equal(t, "SKU-1", msg.Items[0].ItemCode,
		"el item_code del proveedor es la llave con la que se hace match contra el SKU de Probability")
	assert.InDelta(t, 50000, msg.Items[0].ItemPrice, 0.001)
}

func TestPublishListItemsResponse_UnaConsultaFallidaViajaConSuError(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPublisher(q).PublishListItemsResponse(context.Background(),
		&ListItemsResponseMessage{CorrelationID: "corr-1", Error: "credenciales vencidas"}))

	var msg ListItemsResponseMessage
	decodificar(t, q, &msg)
	assert.Equal(t, "credenciales vencidas", msg.Error,
		"el error viaja en el mensaje, no como fallo de publicacion: el front que espera por SSE tiene que enterarse")
	assert.Empty(t, msg.Items)
}

func TestPublishListItemsResponse_SinBrokerDevuelveNil(t *testing.T) {
	err := nuevoPublisher(nil).PublishListItemsResponse(context.Background(),
		&ListItemsResponseMessage{CorrelationID: "corr-1"})

	assert.NoError(t, err)
}

func TestPublishListItemsResponse_PropagaElErrorDelBroker(t *testing.T) {
	q := &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("canal cerrado")
		},
	}

	err := nuevoPublisher(q).PublishListItemsResponse(context.Background(),
		&ListItemsResponseMessage{CorrelationID: "corr-1"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish list_items response")
}

func TestPublishListBankAccountsResponse_LlevaLasCuentasDelProveedor(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &ListBankAccountsResponseMessage{
		Operation:     "list_bank_accounts",
		CorrelationID: "corr-1",
		BusinessID:    26,
		Items: []BankAccountItem{
			{AccountNumber: "1234", Name: "Bancolombia Ahorros", NameType: "savings"},
		},
	}

	require.NoError(t, nuevoPublisher(q).PublishListBankAccountsResponse(context.Background(), resp))

	var msg ListBankAccountsResponseMessage
	decodificar(t, q, &msg)
	assert.Equal(t, "list_bank_accounts", msg.Operation)
	require.Len(t, msg.Items, 1)
	assert.Equal(t, "1234", msg.Items[0].AccountNumber,
		"la cuenta es a donde se imputa el recibo de caja; el usuario la elige de esta lista")
	assert.EqualValues(t, 26, msg.BusinessID)
}

func TestPublishListBankAccountsResponse_SinBrokerDevuelveNil(t *testing.T) {
	err := nuevoPublisher(nil).PublishListBankAccountsResponse(context.Background(),
		&ListBankAccountsResponseMessage{CorrelationID: "corr-1"})

	assert.NoError(t, err)
}

func TestPublishListBankAccountsResponse_PropagaElErrorDelBroker(t *testing.T) {
	q := &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("canal cerrado")
		},
	}

	err := nuevoPublisher(q).PublishListBankAccountsResponse(context.Background(),
		&ListBankAccountsResponseMessage{CorrelationID: "corr-1"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish list_bank_accounts response")
}

func TestTodasLasRespuestas_CompartenLaMismaColaYNoDeclaranNada(t *testing.T) {
	q := &testkit.QueueMock{}
	p := nuevoPublisher(q)
	ctx := context.Background()

	require.NoError(t, p.PublishResponse(ctx, respuestaFactura()))
	require.NoError(t, p.PublishCompareResponse(ctx, &CompareResponseMessage{CorrelationID: "c"}))
	require.NoError(t, p.PublishListItemsResponse(ctx, &ListItemsResponseMessage{CorrelationID: "c"}))
	require.NoError(t, p.PublishListBankAccountsResponse(ctx, &ListBankAccountsResponseMessage{CorrelationID: "c"}))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 4)
	for _, pub := range pubs {
		assert.Equal(t, rabbitmq.QueueInvoicingResponses, pub.Queue,
			"el consumidor discrimina por el campo operation dentro del cuerpo, no por la cola")
	}
	assert.Empty(t, q.Declaradas,
		"ninguna declara la cola ni usa publisher confirms")
}
