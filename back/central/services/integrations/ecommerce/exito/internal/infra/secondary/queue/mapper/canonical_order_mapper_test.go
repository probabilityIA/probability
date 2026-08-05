package mapper

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
)

func momento() time.Time {
	return time.Date(2026, 8, 5, 12, 30, 0, 0, time.UTC)
}

func ptrTiempo(t time.Time) *time.Time { return &t }

func ptrTexto(s string) *string { return &s }

func ordenMinima() *canonical.ProbabilityOrderDTO {
	return &canonical.ProbabilityOrderDTO{
		OccurredAt: momento(),
		ImportedAt: momento(),
	}
}

func TestMapDomainToSerializable_FormateaLasFechasObligatoriasEnRFC3339(t *testing.T) {
	out := MapDomainToSerializable(ordenMinima())

	require.NotNil(t, out)
	assert.Equal(t, "2026-08-05T12:30:00Z", out.OccurredAt)
	assert.Equal(t, "2026-08-05T12:30:00Z", out.ImportedAt,
		"el consumidor parsea RFC3339; cualquier otro formato rompe la ingesta de la orden")
}

func TestMapDomainToSerializable_UnaOrdenSinColeccionesNoProduceNil(t *testing.T) {
	out := MapDomainToSerializable(ordenMinima())

	assert.NotNil(t, out.OrderItems, "make con len 0 devuelve slice vacio, no nil: el JSON sale [] y no null")
	assert.Empty(t, out.OrderItems)
	assert.Empty(t, out.Addresses)
	assert.Empty(t, out.Payments)
	assert.Empty(t, out.Shipments)
	assert.Nil(t, out.ChannelMetadata, "el channel metadata si es puntero y queda nil cuando no viene")
}

func TestMapDomainToSerializable_ConservaLosIdentificadoresDelNegocioYElCanal(t *testing.T) {
	businessID := uint(26)
	orden := ordenMinima()
	orden.BusinessID = &businessID
	orden.IntegrationID = 197
	orden.ExternalID = "EXT-1"
	orden.OrderNumber = "1001"

	out := MapDomainToSerializable(orden)

	require.NotNil(t, out.BusinessID)
	assert.EqualValues(t, 26, *out.BusinessID,
		"sin business_id la orden no se puede atribuir a ningun negocio al otro lado de la cola")
	assert.EqualValues(t, 197, out.IntegrationID)
	assert.Equal(t, "EXT-1", out.ExternalID)
	assert.Equal(t, "1001", out.OrderNumber)
}

func TestMapDomainToSerializable_TrasladaLosItemsConSusPrecios(t *testing.T) {
	orden := ordenMinima()
	orden.OrderItems = []canonical.ProbabilityOrderItemDTO{
		{ProductSKU: "SKU-1", ProductName: "Camiseta", Quantity: 2, UnitPrice: 25000, TotalPrice: 50000, Currency: "COP"},
		{ProductSKU: "SKU-2", Quantity: 1, UnitPrice: 10000, TotalPrice: 10000},
	}

	out := MapDomainToSerializable(orden)

	require.Len(t, out.OrderItems, 2)
	assert.Equal(t, "SKU-1", out.OrderItems[0].ProductSKU)
	assert.Equal(t, "Camiseta", out.OrderItems[0].ProductName)
	assert.EqualValues(t, 2, out.OrderItems[0].Quantity)
	assert.InDelta(t, 25000, out.OrderItems[0].UnitPrice, 0.001)
	assert.InDelta(t, 50000, out.OrderItems[0].TotalPrice, 0.001,
		"el total por linea viaja tal cual: el mapper no recalcula precio por cantidad")
	assert.Equal(t, "SKU-2", out.OrderItems[1].ProductSKU)
}

func TestMapDomainToSerializable_TrasladaLasDirecciones(t *testing.T) {
	orden := ordenMinima()
	orden.Addresses = []canonical.ProbabilityAddressDTO{
		{Type: "shipping", FirstName: "Ana", City: "Bogota", Country: "CO", Street: "Calle 1"},
	}

	out := MapDomainToSerializable(orden)

	require.Len(t, out.Addresses, 1)
	assert.Equal(t, "shipping", out.Addresses[0].Type)
	assert.Equal(t, "Ana", out.Addresses[0].FirstName)
	assert.Equal(t, "Bogota", out.Addresses[0].City)
	assert.Equal(t, "Calle 1", out.Addresses[0].Street)
}

func TestMapDomainToSerializable_LasFechasOpcionalesDePagoQuedanNilSiNoVienen(t *testing.T) {
	orden := ordenMinima()
	orden.Payments = []canonical.ProbabilityPaymentDTO{{Amount: 50000, Status: "paid"}}

	out := MapDomainToSerializable(orden)

	require.Len(t, out.Payments, 1)
	assert.Nil(t, out.Payments[0].PaidAt,
		"un pago sin fecha manda null, no la cadena vacia ni la fecha cero de Go")
	assert.Nil(t, out.Payments[0].ProcessedAt)
	assert.Nil(t, out.Payments[0].RefundedAt)
}

func TestMapDomainToSerializable_FormateaLasFechasDePagoQueSiVienen(t *testing.T) {
	orden := ordenMinima()
	orden.Payments = []canonical.ProbabilityPaymentDTO{{
		Amount: 50000, Status: "paid",
		PaidAt: ptrTiempo(momento()), ProcessedAt: ptrTiempo(momento()), RefundedAt: ptrTiempo(momento()),
	}}

	out := MapDomainToSerializable(orden)

	require.NotNil(t, out.Payments[0].PaidAt)
	assert.Equal(t, "2026-08-05T12:30:00Z", *out.Payments[0].PaidAt)
	assert.Equal(t, "2026-08-05T12:30:00Z", *out.Payments[0].ProcessedAt)
	assert.Equal(t, "2026-08-05T12:30:00Z", *out.Payments[0].RefundedAt)
	assert.InDelta(t, 50000, out.Payments[0].Amount, 0.001)
}

func TestMapDomainToSerializable_TrasladaLosEnviosConSusFechasOpcionales(t *testing.T) {
	orden := ordenMinima()
	orden.Shipments = []canonical.ProbabilityShipmentDTO{{
		TrackingNumber: ptrTexto("GUIA-1"), Carrier: ptrTexto("SERVIENTREGA"),
		ShippedAt: ptrTiempo(momento()),
	}}

	out := MapDomainToSerializable(orden)

	require.Len(t, out.Shipments, 1)
	require.NotNil(t, out.Shipments[0].TrackingNumber)
	assert.Equal(t, "GUIA-1", *out.Shipments[0].TrackingNumber)
	require.NotNil(t, out.Shipments[0].Carrier)
	assert.Equal(t, "SERVIENTREGA", *out.Shipments[0].Carrier)
	require.NotNil(t, out.Shipments[0].ShippedAt)
	assert.Equal(t, "2026-08-05T12:30:00Z", *out.Shipments[0].ShippedAt)
	assert.Nil(t, out.Shipments[0].DeliveredAt)
	assert.Nil(t, out.Shipments[0].EstimatedDelivery)
}

func TestMapDomainToSerializable_MapeaElChannelMetadataCuandoViene(t *testing.T) {
	orden := ordenMinima()
	orden.ChannelMetadata = &canonical.ProbabilityChannelMetadataDTO{
		ChannelSource: "canal",
		RawData:       []byte(`{"id":123}`),
		ReceivedAt:    momento(),
		IsLatest:      true,
		SyncStatus:    "synced",
	}

	out := MapDomainToSerializable(orden)

	require.NotNil(t, out.ChannelMetadata)
	assert.Equal(t, "canal", out.ChannelMetadata.ChannelSource)
	assert.JSONEq(t, `{"id":123}`, string(out.ChannelMetadata.RawData),
		"el payload crudo del canal viaja como JSON embebido, no como cadena escapada")
	assert.Equal(t, "2026-08-05T12:30:00Z", out.ChannelMetadata.ReceivedAt)
	assert.True(t, out.ChannelMetadata.IsLatest)
	assert.Nil(t, out.ChannelMetadata.ProcessedAt)
	assert.Nil(t, out.ChannelMetadata.LastSyncedAt)
}

func TestMapDomainToSerializable_LosBloquesCrudosViajanComoJSONEmbebido(t *testing.T) {
	orden := ordenMinima()
	orden.Items = []byte(`{"a":1}`)
	orden.Metadata = []byte(`{"b":2}`)
	orden.FinancialDetails = []byte(`{"c":3}`)
	orden.ShippingDetails = []byte(`{"d":4}`)
	orden.PaymentDetails = []byte(`{"e":5}`)
	orden.FulfillmentDetails = []byte(`{"f":6}`)

	out := MapDomainToSerializable(orden)

	assert.JSONEq(t, `{"a":1}`, string(out.Items))
	assert.JSONEq(t, `{"b":2}`, string(out.Metadata))
	assert.JSONEq(t, `{"c":3}`, string(out.FinancialDetails))
	assert.JSONEq(t, `{"d":4}`, string(out.ShippingDetails))
	assert.JSONEq(t, `{"e":5}`, string(out.PaymentDetails))
	assert.JSONEq(t, `{"f":6}`, string(out.FulfillmentDetails))
}

func TestMapDomainToSerializable_ElResultadoEsSerializableSinPerderLaOrden(t *testing.T) {
	businessID := uint(26)
	orden := ordenMinima()
	orden.BusinessID = &businessID
	orden.ExternalID = "EXT-1"
	orden.TotalAmount = 60000
	orden.Currency = "COP"
	orden.OrderItems = []canonical.ProbabilityOrderItemDTO{{ProductSKU: "SKU-1", Quantity: 1}}

	crudo, err := json.Marshal(MapDomainToSerializable(orden))

	require.NoError(t, err,
		"si el mapeo produce algo no serializable la orden nunca llega a la cola")
	var vuelta map[string]any
	require.NoError(t, json.Unmarshal(crudo, &vuelta))
	assert.Equal(t, "EXT-1", vuelta["external_id"])
	assert.EqualValues(t, 60000, vuelta["total_amount"])
	assert.Equal(t, "COP", vuelta["currency"])
}

func TestMapDomainToSerializable_UnaOrdenSinBusinessIDIgualSeSerializa(t *testing.T) {
	out := MapDomainToSerializable(ordenMinima())

	assert.Nil(t, out.BusinessID,
		"el mapper no valida: una orden sin negocio llega igual a la cola y el fallo aparece recien en el consumidor")
	_, err := json.Marshal(out)
	assert.NoError(t, err)
}
