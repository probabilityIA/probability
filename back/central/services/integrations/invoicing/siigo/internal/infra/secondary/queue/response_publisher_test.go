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

func brokerCaido() *testkit.QueueMock {
	return &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("canal cerrado")
		},
	}
}

func decodificar(t *testing.T, q *testkit.QueueMock, destino any) {
	t.Helper()
	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	require.NoError(t, json.Unmarshal(pubs[0].Body, destino))
}

func TestPublishResponse_PublicaLaFacturaEnLaColaCompartidaDeFacturacion(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &InvoiceResponseMessage{
		InvoiceID: 501, Provider: "siigo", Status: "success",
		InvoiceNumber: "FV-1001", CUFE: "CUFE-ABC", CorrelationID: "corr-1",
	}

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), resp))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 1)
	assert.Equal(t, rabbitmq.QueueInvoicingResponses, pubs[0].Queue,
		"siigo y softpymes comparten cola: el consumidor distingue por el campo provider")

	var msg InvoiceResponseMessage
	decodificar(t, q, &msg)
	assert.EqualValues(t, 501, msg.InvoiceID)
	assert.Equal(t, "siigo", msg.Provider)
	assert.Equal(t, "CUFE-ABC", msg.CUFE,
		"el CUFE es el identificador fiscal de la DIAN: sin el la factura no se puede validar despues")
}

func TestPublishResponse_LaRespuestaDeErrorConservaCodigoYAuditoria(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &InvoiceResponseMessage{
		InvoiceID: 501, Status: "error", CorrelationID: "corr-1",
		Error: "cliente no existe en Siigo", ErrorCode: "not_found",
		AuditRequestURL: "https://api.siigo.com/v1/invoices", AuditResponseStatus: 404,
		AuditResponseBody: `{"Errors":[{"Code":"not_found"}]}`,
	}

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), resp))

	var msg InvoiceResponseMessage
	decodificar(t, q, &msg)
	assert.Equal(t, "error", msg.Status)
	assert.Equal(t, "not_found", msg.ErrorCode)
	assert.Equal(t, 404, msg.AuditResponseStatus,
		"guardar el HTTP crudo es lo que permite reclamarle a Siigo cuando una factura sale mal")
}

func TestPublishResponse_LeAsignaFechaSiLlegaSinElla(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &InvoiceResponseMessage{CorrelationID: "corr-1"}

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), resp))

	assert.False(t, resp.Timestamp.IsZero())
}

func TestPublishResponse_RespetaLaFechaQueYaTrae(t *testing.T) {
	fijo := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	q := &testkit.QueueMock{}
	resp := &InvoiceResponseMessage{CorrelationID: "corr-1", Timestamp: fijo}

	require.NoError(t, nuevoPublisher(q).PublishResponse(context.Background(), resp))

	assert.True(t, resp.Timestamp.Equal(fijo))
}

func TestPublishCompareResponse_LlevaLosDocumentosDelProveedor(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &CompareResponseMessage{
		Operation: "compare", CorrelationID: "corr-1", BusinessID: 26,
		DateFrom: "2026-08-01", DateTo: "2026-08-05",
	}

	require.NoError(t, nuevoPublisher(q).PublishCompareResponse(context.Background(), resp))

	var msg CompareResponseMessage
	decodificar(t, q, &msg)
	assert.Equal(t, "compare", msg.Operation)
	assert.EqualValues(t, 26, msg.BusinessID,
		"sin business_id se cruzarian documentos de negocios distintos")
	assert.Equal(t, "2026-08-01", msg.DateFrom)
}

func TestPublishListItemsResponse_LlevaElCatalogoDeSiigo(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &ListItemsResponseMessage{
		Operation: "list_items", CorrelationID: "corr-1", BusinessID: 26,
		Items: []ListItemsItem{{ItemCode: "SKU-1", ItemName: "Camiseta", ItemPrice: 50000}},
	}

	require.NoError(t, nuevoPublisher(q).PublishListItemsResponse(context.Background(), resp))

	var msg ListItemsResponseMessage
	decodificar(t, q, &msg)
	require.Len(t, msg.Items, 1)
	assert.Equal(t, "SKU-1", msg.Items[0].ItemCode,
		"el item_code de Siigo es la llave con la que se hace match contra el SKU de Probability")
	assert.InDelta(t, 50000, msg.Items[0].ItemPrice, 0.001)
}

func TestPublishListBankAccountsResponse_LlevaLasCuentas(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &ListBankAccountsResponseMessage{
		Operation: "list_bank_accounts", CorrelationID: "corr-1", BusinessID: 26,
		Items: []BankAccountItem{{AccountNumber: "1234", Name: "Bancolombia"}},
	}

	require.NoError(t, nuevoPublisher(q).PublishListBankAccountsResponse(context.Background(), resp))

	var msg ListBankAccountsResponseMessage
	decodificar(t, q, &msg)
	require.Len(t, msg.Items, 1)
	assert.Equal(t, "1234", msg.Items[0].AccountNumber)
}

func TestPublishListSiigoWarehousesResponse_LlevaLasBodegasDelProveedor(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &ListSiigoWarehousesResponseMessage{
		Operation: "list_siigo_warehouses", CorrelationID: "corr-1", BusinessID: 26,
		Items: []SiigoWarehouseItem{{ID: 9, Name: "Bodega Principal"}},
	}

	require.NoError(t, nuevoPublisher(q).PublishListSiigoWarehousesResponse(context.Background(), resp))

	var msg ListSiigoWarehousesResponseMessage
	decodificar(t, q, &msg)
	assert.Equal(t, "list_siigo_warehouses", msg.Operation,
		"es la unica operacion que siigo tiene y softpymes no")
	require.Len(t, msg.Items, 1)
	assert.Equal(t, 9, msg.Items[0].ID)
	assert.Equal(t, "Bodega Principal", msg.Items[0].Name)
}

func TestTodasLasOperaciones_ComparteLaMismaColaYFechanElMensaje(t *testing.T) {
	q := &testkit.QueueMock{}
	p := nuevoPublisher(q)
	ctx := context.Background()

	require.NoError(t, p.PublishResponse(ctx, &InvoiceResponseMessage{CorrelationID: "c"}))
	require.NoError(t, p.PublishCompareResponse(ctx, &CompareResponseMessage{CorrelationID: "c"}))
	require.NoError(t, p.PublishListItemsResponse(ctx, &ListItemsResponseMessage{CorrelationID: "c"}))
	require.NoError(t, p.PublishListBankAccountsResponse(ctx, &ListBankAccountsResponseMessage{CorrelationID: "c"}))
	require.NoError(t, p.PublishListSiigoWarehousesResponse(ctx, &ListSiigoWarehousesResponseMessage{CorrelationID: "c"}))

	pubs := q.Publicaciones()
	require.Len(t, pubs, 5)
	for _, pub := range pubs {
		assert.Equal(t, rabbitmq.QueueInvoicingResponses, pub.Queue)
	}
	assert.Empty(t, q.Declaradas,
		"ninguna declara la cola ni usa publisher confirms")
}

func TestTodasLasOperaciones_SinBrokerDevuelvenNilYLaOperacionQuedaColgada(t *testing.T) {
	p := nuevoPublisher(nil)
	ctx := context.Background()

	casos := map[string]error{
		"invoice":               p.PublishResponse(ctx, &InvoiceResponseMessage{CorrelationID: "c"}),
		"compare":               p.PublishCompareResponse(ctx, &CompareResponseMessage{CorrelationID: "c"}),
		"list_items":            p.PublishListItemsResponse(ctx, &ListItemsResponseMessage{CorrelationID: "c"}),
		"list_bank_accounts":    p.PublishListBankAccountsResponse(ctx, &ListBankAccountsResponseMessage{CorrelationID: "c"}),
		"list_siigo_warehouses": p.PublishListSiigoWarehousesResponse(ctx, &ListSiigoWarehousesResponseMessage{CorrelationID: "c"}),
	}

	for nombre, err := range casos {
		assert.NoError(t, err,
			"%s con rabbit nil devuelve nil: la factura ya se emitio en Siigo pero el modulo de facturacion nunca se entera y queda en processing", nombre)
	}
}

func TestTodasLasOperaciones_PropaganElErrorDelBrokerConSuTipo(t *testing.T) {
	ctx := context.Background()
	casos := map[string]func(*ResponsePublisher) error{
		"invoice":    func(p *ResponsePublisher) error { return p.PublishResponse(ctx, &InvoiceResponseMessage{}) },
		"compare":    func(p *ResponsePublisher) error { return p.PublishCompareResponse(ctx, &CompareResponseMessage{}) },
		"list_items": func(p *ResponsePublisher) error { return p.PublishListItemsResponse(ctx, &ListItemsResponseMessage{}) },
		"list_bank_accounts": func(p *ResponsePublisher) error {
			return p.PublishListBankAccountsResponse(ctx, &ListBankAccountsResponseMessage{})
		},
		"list_siigo_warehouses": func(p *ResponsePublisher) error {
			return p.PublishListSiigoWarehousesResponse(ctx, &ListSiigoWarehousesResponseMessage{})
		},
	}

	for tipo, publicar := range casos {
		t.Run(tipo, func(t *testing.T) {
			err := publicar(nuevoPublisher(brokerCaido()))

			require.Error(t, err)
			assert.Contains(t, err.Error(), "failed to publish "+tipo+" response",
				"el tipo va en el mensaje de error para poder rastrear que operacion se perdio")
		})
	}
}

func TestPublish_UnPayloadNoSerializableNoSePublicaAMedias(t *testing.T) {
	q := &testkit.QueueMock{}
	resp := &InvoiceResponseMessage{
		CorrelationID: "corr-1",
		DocumentJSON:  map[string]interface{}{"malo": math.Inf(1)},
	}

	err := nuevoPublisher(q).PublishResponse(context.Background(), resp)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal invoice response")
	assert.Empty(t, q.Publicaciones())
}
