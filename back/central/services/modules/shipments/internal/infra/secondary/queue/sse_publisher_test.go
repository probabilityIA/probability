package queue

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/testkit"
)

func f64(v float64) *float64 { return &v }

func nuevoSSE(q *testkit.QueueMock) domain.IShipmentSSEPublisher {
	return NewSSEPublisher(q, testkit.NewSilentLogger())
}

func sobre(t *testing.T, q *testkit.QueueMock) rabbitmq.EventEnvelope {
	t.Helper()
	pub := q.EsperarPublicacion(t)
	var env rabbitmq.EventEnvelope
	require.NoError(t, json.Unmarshal(pub.Body, &env))
	return env
}

func datos(t *testing.T, q *testkit.QueueMock) map[string]any {
	t.Helper()
	env := sobre(t, q)
	crudo, err := json.Marshal(env.Data)
	require.NoError(t, err)
	var out map[string]any
	require.NoError(t, json.Unmarshal(crudo, &out))
	return out
}

func TestPublish_TodosLosEventosVanAlExchangeDeEventosConSuTipo(t *testing.T) {
	casos := map[string]func(domain.IShipmentSSEPublisher){
		"shipment.quote_received": func(p domain.IShipmentSSEPublisher) {
			p.PublishQuoteReceived(context.Background(), 26, "corr-1", map[string]interface{}{"a": 1})
		},
		"shipment.quote_failed": func(p domain.IShipmentSSEPublisher) {
			p.PublishQuoteFailed(context.Background(), 26, "corr-1", "sin cobertura")
		},
		"shipment.guide_failed": func(p domain.IShipmentSSEPublisher) {
			p.PublishGuideFailed(context.Background(), 26, 5, "corr-1", "rechazada")
		},
		"shipment.tracking_updated": func(p domain.IShipmentSSEPublisher) {
			p.PublishTrackingUpdated(context.Background(), 26, "corr-1", map[string]interface{}{"estado": "en_ruta"})
		},
		"shipment.tracking_failed": func(p domain.IShipmentSSEPublisher) {
			p.PublishTrackingFailed(context.Background(), 26, "corr-1", "timeout")
		},
		"shipment.cancelled": func(p domain.IShipmentSSEPublisher) {
			p.PublishShipmentCancelled(context.Background(), 26, 5)
		},
		"shipment.cancel_failed": func(p domain.IShipmentSSEPublisher) {
			p.PublishCancelFailed(context.Background(), 26, 5, "corr-1", "ya despachada")
		},
		"shipment.sync_batch_completed": func(p domain.IShipmentSSEPublisher) {
			p.PublishSyncBatchCompleted(context.Background(), 26, "corr-1", map[string]interface{}{"total": 10})
		},
	}

	for tipo, publicar := range casos {
		t.Run(tipo, func(t *testing.T) {
			q := &testkit.QueueMock{}

			publicar(nuevoSSE(q))

			env := sobre(t, q)
			assert.Equal(t, tipo, env.Type)
			assert.Equal(t, "shipment", env.Category,
				"la categoria es lo que el gateway SSE usa para saber a que canal del front empujar el evento")
			assert.EqualValues(t, 26, env.BusinessID,
				"sin business_id el evento se le empujaria al front de otro negocio")
		})
	}
}

func TestPublish_LaRoutingKeyEsElTipoDeEvento(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishShipmentCancelled(context.Background(), 26, 5)

	pub := q.EsperarPublicacion(t)
	assert.Equal(t, rabbitmq.EventsExchangeName, pub.Exchange)
	assert.Equal(t, "shipment.cancelled", pub.RoutingKey,
		"a diferencia del fanout de ordenes, aca la routing key si se usa: los suscriptores filtran por tipo")
}

func TestPublish_CadaEventoLlevaIDYFechaGeneradosSiNoVienen(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishShipmentCancelled(context.Background(), 26, 5)

	env := sobre(t, q)
	assert.NotEmpty(t, env.ID)
	assert.False(t, env.Timestamp.IsZero())
}

func TestPublishQuoteReceived_LlevaLaCorrelacionYLasCotizaciones(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishQuoteReceived(context.Background(), 26, "corr-1",
		map[string]interface{}{"servientrega": 12000})

	d := datos(t, q)
	assert.Equal(t, "corr-1", d["correlation_id"],
		"el front abrio el SSE con ese correlation_id: sin el no sabe que cotizacion acaba de llegar")
	assert.NotNil(t, d["quotes"])
}

func TestPublishQuoteFailed_LlevaElMensajeDeError(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishQuoteFailed(context.Background(), 26, "corr-1", "sin cobertura en el destino")

	d := datos(t, q)
	assert.Equal(t, "sin cobertura en el destino", d["error_message"])
	assert.Equal(t, "corr-1", d["correlation_id"])
}

func TestPublishGuideGenerated_LlevaLaGuiaYLaEtiqueta(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishGuideGenerated(context.Background(), 26, 5, "corr-1",
		"GUIA-123", "https://s3/label.pdf", "SERVIENTREGA", nil)

	d := datos(t, q)
	assert.EqualValues(t, 5, d["shipment_id"])
	assert.Equal(t, "GUIA-123", d["tracking_number"])
	assert.Equal(t, "https://s3/label.pdf", d["label_url"])
	assert.Equal(t, "SERVIENTREGA", d["carrier"])
}

func TestPublishGuideGenerated_SinDatosDeNotificacionNoAgregaCamposDelCliente(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishGuideGenerated(context.Background(), 26, 5, "corr-1", "G", "L", "C", nil)

	d := datos(t, q)
	assert.NotContains(t, d, "customer_phone",
		"si no hay datos de notificacion no se emiten claves vacias: el consumidor distingue ausente de vacio")
	assert.NotContains(t, d, "cod_total")
}

func TestPublishGuideGenerated_ConDatosDeNotificacionLlevaAlClienteYElRecaudo(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishGuideGenerated(context.Background(), 26, 5, "corr-1", "G", "L", "SERVIENTREGA",
		&domain.GuideNotificationData{
			CustomerName:  "Ana Perez",
			CustomerPhone: "573001112233",
			OrderNumber:   "1001",
			BusinessName:  "Tienda",
			IntegrationID: 197,
			CodTotal:      f64(150000),
			CodCarrierFee: f64(8000),
			TrackingURL:   "https://track/G",
		})

	d := datos(t, q)
	assert.Equal(t, "Ana Perez", d["customer_name"])
	assert.Equal(t, "573001112233", d["customer_phone"])
	assert.Equal(t, "1001", d["order_number"])
	assert.EqualValues(t, 150000, d["cod_total"])
	assert.EqualValues(t, 8000, d["cod_carrier_fee"])
	assert.Equal(t, "https://track/G", d["tracking_url"])
}

func TestPublishGuideGenerated_LosPunterosNilDelRecaudoNoSeEmiten(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishGuideGenerated(context.Background(), 26, 5, "corr-1", "G", "L", "C",
		&domain.GuideNotificationData{CustomerName: "Ana", CodTotal: nil, CodCarrierFee: nil, TrackingURL: ""})

	d := datos(t, q)
	assert.Equal(t, "Ana", d["customer_name"])
	assert.NotContains(t, d, "cod_total",
		"una orden que no es contraentrega no emite cod_total, en vez de emitir cero y confundir al front")
	assert.NotContains(t, d, "cod_carrier_fee")
	assert.NotContains(t, d, "tracking_url", "una URL vacia tampoco se emite")
}

func TestPublishGuideGenerated_ElIntegrationIDDeLaNotificacionSubeAlSobre(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishGuideGenerated(context.Background(), 26, 5, "corr-1", "G", "L", "C",
		&domain.GuideNotificationData{IntegrationID: 197})

	env := sobre(t, q)
	assert.EqualValues(t, 197, env.IntegrationID,
		"el integration_id se extrae del mapa de datos y se copia al sobre para que el suscriptor pueda filtrar por canal")
}

func TestPublish_SinIntegrationIDEnLosDatosElSobreLlevaCero(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishShipmentCancelled(context.Background(), 26, 5)

	assert.EqualValues(t, 0, sobre(t, q).IntegrationID)
}

func TestPublishSyncBatchCompleted_FusionaElResumenConLaCorrelacion(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishSyncBatchCompleted(context.Background(), 26, "corr-1",
		map[string]interface{}{"total": 10, "fallidas": 2})

	d := datos(t, q)
	assert.Equal(t, "corr-1", d["correlation_id"])
	assert.EqualValues(t, 10, d["total"])
	assert.EqualValues(t, 2, d["fallidas"])
}

func TestPublishSyncBatchCompleted_UnResumenQueTraeCorrelationIDPisaElDelParametro(t *testing.T) {
	q := &testkit.QueueMock{}

	nuevoSSE(q).PublishSyncBatchCompleted(context.Background(), 26, "corr-parametro",
		map[string]interface{}{"correlation_id": "corr-resumen"})

	assert.Equal(t, "corr-resumen", datos(t, q)["correlation_id"],
		"el bucle que copia el resumen sobreescribe la clave: si el resumen trae correlation_id, gana el del resumen")
}

func TestPublish_ElEnvioEsAsincronoYNoDevuelveError(t *testing.T) {
	q := &testkit.QueueMock{
		PublishToExchangeFn: func(ctx context.Context, exchangeName, routingKey string, message []byte) error {
			return errors.New("broker caido")
		},
	}

	nuevoSSE(q).PublishShipmentCancelled(context.Background(), 26, 5)

	q.EsperarPublicacion(t)
	assert.True(t, true,
		"los metodos no devuelven error: si el broker falla el evento SSE se pierde en silencio y solo queda el log. Es aceptable para un aviso de UI, no lo seria para un cobro")
}

func TestPublish_UsaUnContextoPropioYSobreviveALaCancelacionDelRequest(t *testing.T) {
	q := &testkit.QueueMock{}
	ctx, cancelar := context.WithCancel(context.Background())

	nuevoSSE(q).PublishShipmentCancelled(ctx, 26, 5)
	cancelar()

	q.EsperarPublicacion(t)
	assert.Len(t, q.Publicaciones(), 1,
		"la goroutine usa context.Background(), asi que cancelar el request HTTP no aborta el envio del evento")
}

func TestNoopSSEPublisher_NoPublicaNadaYNoRevienta(t *testing.T) {
	p := NewNoopSSEPublisher()
	ctx := context.Background()

	p.PublishQuoteReceived(ctx, 26, "c", nil)
	p.PublishQuoteFailed(ctx, 26, "c", "e")
	p.PublishGuideGenerated(ctx, 26, 5, "c", "g", "l", "car", nil)
	p.PublishGuideFailed(ctx, 26, 5, "c", "e")
	p.PublishTrackingUpdated(ctx, 26, "c", nil)
	p.PublishTrackingFailed(ctx, 26, "c", "e")
	p.PublishShipmentCancelled(ctx, 26, 5)
	p.PublishCancelFailed(ctx, 26, 5, "c", "e")
	p.PublishSyncBatchCompleted(ctx, 26, "c", nil)

	assert.True(t, true,
		"la version noop deja correr los casos de uso sin broker (tests y arranques degradados) sin cambiar su codigo")
}
