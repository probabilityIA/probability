package queue

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/testkit"
)

func f64(v float64) *float64 { return &v }

func txt(v string) *string { return &v }

func u(v uint) *uint { return &v }

func orden() *entities.ProbabilityOrder {
	return &entities.ProbabilityOrder{
		ID:            "ord-1",
		BusinessID:    u(26),
		BusinessName:  "Tienda Probability",
		IntegrationID: 197,
		Platform:      "shopify",
		OrderNumber:   "1001",
		TotalAmount:   150000,
		Currency:      "COP",
		CustomerName:  "Ana Perez",
		CustomerPhone: "573001112233",
		CustomerEmail: "ana@correo.com",
	}
}

func nuevoPublisher(q *testkit.QueueMock) *OrderRabbitPublisher {
	return NewOrderRabbitPublisher(q, testkit.NewSilentLogger()).(*OrderRabbitPublisher)
}

func cuerpo(t *testing.T, p testkit.Publicacion) map[string]any {
	t.Helper()
	var out map[string]any
	require.NoError(t, json.Unmarshal(p.Body, &out))
	return out
}

func TestPublishOrderCreated_VaAlExchangeDeEventosDeOrden(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPublisher(q).PublishOrderCreated(context.Background(), orden()))

	require.Len(t, q.Publicados, 1)
	assert.Equal(t, rabbitmq.ExchangeOrderEvents, q.Publicados[0].Exchange)
	assert.Empty(t, q.Publicados[0].Queue,
		"los eventos de orden no van a una cola directa sino al exchange fanout")
	assert.Equal(t, "order.created", cuerpo(t, q.Publicados[0])["event_type"])
}

func TestPublishOrder_LosCuatroEventosUsanElMismoExchangeYSoloCambiaElTipo(t *testing.T) {
	casos := map[string]func(*OrderRabbitPublisher) error{
		"order.created":   func(p *OrderRabbitPublisher) error { return p.PublishOrderCreated(context.Background(), orden()) },
		"order.updated":   func(p *OrderRabbitPublisher) error { return p.PublishOrderUpdated(context.Background(), orden()) },
		"order.cancelled": func(p *OrderRabbitPublisher) error { return p.PublishOrderCancelled(context.Background(), orden()) },
		"order.status_changed": func(p *OrderRabbitPublisher) error {
			return p.PublishOrderStatusChanged(context.Background(), orden(), "pending", "confirmed")
		},
	}

	for tipo, publicar := range casos {
		t.Run(tipo, func(t *testing.T) {
			q := &testkit.QueueMock{}

			require.NoError(t, publicar(nuevoPublisher(q)))

			require.Len(t, q.Publicados, 1)
			assert.Equal(t, rabbitmq.ExchangeOrderEvents, q.Publicados[0].Exchange)
			assert.Equal(t, tipo, cuerpo(t, q.Publicados[0])["event_type"])
		})
	}
}

func TestPublishToQueue_LaRoutingKeyViajaVaciaPorqueElExchangeEsFanout(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPublisher(q).PublishOrderCreated(context.Background(), orden()))

	assert.Empty(t, q.Publicados[0].RoutingKey,
		"publishToQueue descarta la routing key que recibe y publica con cadena vacia: el fanout entrega a todos los bindings (invoicing, score, inventory, events, customers)")
}

func TestPublishOrderStatusChanged_LlevaElEstadoAnteriorYElNuevo(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPublisher(q).PublishOrderStatusChanged(context.Background(), orden(), "pending", "confirmed"))

	cambios := cuerpo(t, q.Publicados[0])["changes"].(map[string]any)
	assert.Equal(t, "pending", cambios["previous_status"])
	assert.Equal(t, "confirmed", cambios["current_status"],
		"el consumidor necesita ambos estados para decidir si la transicion dispara facturacion o notificacion")
}

func TestPublishOrderCreated_CadaEventoLlevaUnIDUnico(t *testing.T) {
	q := &testkit.QueueMock{}
	p := nuevoPublisher(q)

	require.NoError(t, p.PublishOrderCreated(context.Background(), orden()))
	require.NoError(t, p.PublishOrderCreated(context.Background(), orden()))

	primero := cuerpo(t, q.Publicados[0])["event_id"]
	segundo := cuerpo(t, q.Publicados[1])["event_id"]
	assert.NotEmpty(t, primero)
	assert.NotEqual(t, primero, segundo,
		"el event_id es lo unico que el consumidor tiene para deduplicar; si se repitiera, dos publicaciones de la misma orden se verian como el mismo evento")
}

func TestPublishOrderCreated_IdentificaLaOrdenElNegocioYLaIntegracion(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPublisher(q).PublishOrderCreated(context.Background(), orden()))

	msg := cuerpo(t, q.Publicados[0])
	assert.Equal(t, "ord-1", msg["order_id"])
	assert.EqualValues(t, 26, msg["business_id"],
		"sin business_id el evento no se puede atribuir a un negocio del otro lado del fanout")
	assert.EqualValues(t, 197, msg["integration_id"])
}

func TestPublishOrderCreated_PropagaElErrorDelBroker(t *testing.T) {
	fallo := errors.New("canal cerrado")
	q := &testkit.QueueMock{
		PublishToExchangeFn: func(ctx context.Context, exchangeName, routingKey string, message []byte) error {
			return fallo
		},
	}

	err := nuevoPublisher(q).PublishOrderCreated(context.Background(), orden())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error publishing to exchange")
}

func TestPublishOrderEvent_UsaElTimestampYElIDDelEventoNoLosGeneraDeNuevo(t *testing.T) {
	momento := time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)
	q := &testkit.QueueMock{}
	evento := &entities.OrderEvent{
		ID:            "evt-fijo",
		Type:          entities.OrderEventTypeCreated,
		OrderID:       "ord-1",
		BusinessID:    u(26),
		IntegrationID: u(197),
		Timestamp:     momento,
	}

	require.NoError(t, nuevoPublisher(q).PublishOrderEvent(context.Background(), evento, orden()))

	msg := cuerpo(t, q.Publicados[0])
	assert.Equal(t, "evt-fijo", msg["event_id"],
		"al reenviar un evento ya existente se conserva su id: generar uno nuevo romperia la deduplicacion del consumidor")
	assert.Equal(t, "order.created", msg["event_type"])
	assert.Contains(t, msg["timestamp"], "2026-08-05T12:00:00")
}

func TestPublishOrderEvent_TrasladaLaMetadataDelEvento(t *testing.T) {
	q := &testkit.QueueMock{}
	evento := &entities.OrderEvent{
		ID: "evt-1", Type: entities.OrderEventTypeShipped, OrderID: "ord-1",
		Metadata: map[string]interface{}{"origen": "webhook"},
	}

	require.NoError(t, nuevoPublisher(q).PublishOrderEvent(context.Background(), evento, orden()))

	meta := cuerpo(t, q.Publicados[0])["metadata"].(map[string]any)
	assert.Equal(t, "webhook", meta["origen"])
}

func TestGetQueueForEventType_MapeaCadaTipoASuRoutingKey(t *testing.T) {
	p := nuevoPublisher(&testkit.QueueMock{})
	casos := map[entities.OrderEventType]string{
		entities.OrderEventTypeCreated:               rabbitmq.RoutingKeyOrderCreated,
		entities.OrderEventTypeUpdated:               rabbitmq.RoutingKeyOrderUpdated,
		entities.OrderEventTypeCancelled:             rabbitmq.RoutingKeyOrderCancelled,
		entities.OrderEventTypeStatusChanged:         rabbitmq.RoutingKeyOrderStatusChanged,
		entities.OrderEventTypeConfirmationRequested: rabbitmq.QueueOrdersConfirmationRequested,
		entities.OrderEventTypeRefunded:              rabbitmq.RoutingKeyOrderGeneric,
		"tipo.inventado":                             rabbitmq.RoutingKeyOrderGeneric,
	}

	for tipo, esperado := range casos {
		assert.Equal(t, esperado, p.getQueueForEventType(tipo),
			"cualquier tipo no contemplado cae a la routing key generica en vez de perderse")
	}
}

func TestGetQueueForEventType_ElResultadoNoLlegaAAfectarLaPublicacion(t *testing.T) {
	q := &testkit.QueueMock{}
	evento := &entities.OrderEvent{ID: "evt-1", Type: entities.OrderEventTypeCancelled, OrderID: "ord-1"}

	require.NoError(t, nuevoPublisher(q).PublishOrderEvent(context.Background(), evento, orden()))

	assert.Empty(t, q.Publicados[0].RoutingKey,
		"la cola calculada por getQueueForEventType se pasa a publishToQueue y ahi se descarta: hoy es codigo muerto que solo sirve como documentacion del mapeo")
}

func TestPublishConfirmationRequested_VaALaColaDirectaDeConfirmacion(t *testing.T) {
	q := &testkit.QueueMock{}

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), orden()))

	require.Len(t, q.Publicados, 1)
	assert.Equal(t, rabbitmq.QueueOrdersConfirmationRequested, q.Publicados[0].Queue,
		"la confirmacion no va al fanout sino a una cola directa: solo el bot de WhatsApp la consume")
	assert.Empty(t, q.Publicados[0].Exchange)
}

func TestPublishConfirmationRequested_LlevaLosDatosQueElBotNecesitaParaEscribirle(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.ShippingStreet = "Calle 1 #2-3"
	o.ShippingCity = "Bogota"
	o.PaymentMethodName = "Contraentrega"

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

	msg := cuerpo(t, q.Publicados[0])
	assert.Equal(t, "573001112233", msg["customer_phone"],
		"sin telefono el bot no tiene a quien escribirle")
	assert.Equal(t, "Ana Perez", msg["customer_name"])
	assert.Equal(t, "1001", msg["order_number"])
	assert.Equal(t, "Calle 1 #2-3", msg["shipping_address"])
	assert.Equal(t, "Bogota", msg["shipping_city"])
	assert.Equal(t, "Contraentrega", msg["payment_method_name"])
}

func TestPublishConfirmationRequested_UnaOrdenQueNoEsCODNoLlevaValoresDeRecaudo(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.IsCod = false
	o.CodTotal = f64(99000)
	o.Shipments = []entities.ProbabilityShipment{{CodCarrierFee: f64(5000)}}

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

	msg := cuerpo(t, q.Publicados[0])
	assert.EqualValues(t, 0, msg["cod_total"],
		"aunque la orden traiga cod_total cargado, si no es COD se manda cero para que el bot no le pida plata al cliente")
	assert.EqualValues(t, 0, msg["cod_carrier_fee"])
	assert.Equal(t, false, msg["is_cod"])
}

func TestPublishConfirmationRequested_UnaOrdenCODLlevaElTotalYLaTarifaDeLaTransportadora(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.IsCod = true
	o.CodTotal = f64(150000)
	o.Shipments = []entities.ProbabilityShipment{{CodCarrierFee: f64(8000)}}

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

	msg := cuerpo(t, q.Publicados[0])
	assert.EqualValues(t, 150000, msg["cod_total"])
	assert.EqualValues(t, 8000, msg["cod_carrier_fee"])
}

func TestPublishConfirmationRequested_UnaOrdenCODSinTotalMandaCero(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.IsCod = true
	o.CodTotal = nil

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

	assert.EqualValues(t, 0, cuerpo(t, q.Publicados[0])["cod_total"],
		"un puntero nil no puede desreferenciarse: se manda cero en vez de reventar el publisher")
}

func TestCodCarrierFee_TomaLaPrimeraTarifaPositivaDeLosEnvios(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.IsCod = true
	o.CodTotal = f64(1000)
	o.Shipments = []entities.ProbabilityShipment{
		{CodCarrierFee: nil},
		{CodCarrierFee: f64(0)},
		{CodCarrierFee: f64(7000)},
		{CodCarrierFee: f64(9000)},
	}

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

	assert.EqualValues(t, 7000, cuerpo(t, q.Publicados[0])["cod_carrier_fee"],
		"se salta los nil y los ceros y toma la primera tarifa real; con varios envios cobrados, los siguientes se ignoran")
}

func TestCodCarrierFee_SinEnviosConTarifaMandaCero(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.IsCod = true
	o.CodTotal = f64(1000)
	o.Shipments = []entities.ProbabilityShipment{{CodCarrierFee: nil}}

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

	assert.EqualValues(t, 0, cuerpo(t, q.Publicados[0])["cod_carrier_fee"])
}

func TestBuildItemsSummary_ArmaElResumenLegibleParaElMensaje(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.OrderItems = []entities.ProbabilityOrderItem{
		{ProductName: "Camiseta", Quantity: 2},
		{ProductName: "Pantalon", Quantity: 1},
	}

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

	assert.Equal(t, "Camiseta x2, Pantalon x1", cuerpo(t, q.Publicados[0])["items_summary"])
}

func TestBuildItemsSummary_ColapsaLosEspaciosDelNombreDelProducto(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.OrderItems = []entities.ProbabilityOrderItem{
		{ProductName: "  Camiseta   azul \n  talla M ", Quantity: 1},
	}

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

	assert.Equal(t, "Camiseta azul talla M x1", cuerpo(t, q.Publicados[0])["items_summary"],
		"los saltos de linea y espacios repetidos del catalogo romperian el formato del mensaje de WhatsApp")
}

func TestBuildItemsSummary_DescartaLosItemsSinNombre(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.OrderItems = []entities.ProbabilityOrderItem{
		{ProductName: "   ", Quantity: 5},
		{ProductName: "Camiseta", Quantity: 1},
	}

	require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

	assert.Equal(t, "Camiseta x1", cuerpo(t, q.Publicados[0])["items_summary"],
		"un item sin nombre se omite en vez de mandar ' x5' al cliente")
}

func TestBuildItemsSummary_SinItemsUtilesDiceSinItems(t *testing.T) {
	casos := map[string][]entities.ProbabilityOrderItem{
		"sin items":        nil,
		"todos sin nombre": {{ProductName: "", Quantity: 1}, {ProductName: "  ", Quantity: 2}},
	}

	for nombre, items := range casos {
		t.Run(nombre, func(t *testing.T) {
			q := &testkit.QueueMock{}
			o := orden()
			o.OrderItems = items

			require.NoError(t, nuevoPublisher(q).PublishConfirmationRequested(context.Background(), o))

			assert.Equal(t, "Sin items", cuerpo(t, q.Publicados[0])["items_summary"])
		})
	}
}

func TestPublishConfirmationRequested_PropagaElErrorDelBroker(t *testing.T) {
	q := &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("broker caido")
		},
	}

	err := nuevoPublisher(q).PublishConfirmationRequested(context.Background(), orden())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error publishing to RabbitMQ")
}

func TestPublishGuideNotificationRequested_VaALaColaDeNotificacionDeGuia(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.TrackingNumber = txt("GUIA-123")
	o.Shipments = []entities.ProbabilityShipment{{Carrier: txt("SERVIENTREGA")}}

	require.NoError(t, nuevoPublisher(q).PublishGuideNotificationRequested(context.Background(), o))

	require.Len(t, q.Publicados, 1)
	assert.Equal(t, rabbitmq.QueueShipmentsWhatsAppGuideNotification, q.Publicados[0].Queue)
	msg := cuerpo(t, q.Publicados[0])
	assert.Equal(t, "GUIA-123", msg["tracking_number"])
	assert.Equal(t, "SERVIENTREGA", msg["carrier"],
		"la plantilla de WhatsApp usa el nombre real de la transportadora, nunca uno generico")
}

func TestPublishGuideNotificationRequested_SinGuiaNiTransportadoraMandaCadenasVacias(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.TrackingNumber = nil
	o.Shipments = nil

	require.NoError(t, nuevoPublisher(q).PublishGuideNotificationRequested(context.Background(), o))

	msg := cuerpo(t, q.Publicados[0])
	assert.Equal(t, "", msg["tracking_number"],
		"se publica igual con los campos vacios: el consumidor decide si puede armar el mensaje")
	assert.Equal(t, "", msg["carrier"])
}

func TestPublishGuideNotificationRequested_TomaLaTransportadoraDelPrimerEnvio(t *testing.T) {
	q := &testkit.QueueMock{}
	o := orden()
	o.Shipments = []entities.ProbabilityShipment{
		{Carrier: txt("SERVIENTREGA")},
		{Carrier: txt("INTERRAPIDISIMO")},
	}

	require.NoError(t, nuevoPublisher(q).PublishGuideNotificationRequested(context.Background(), o))

	assert.Equal(t, "SERVIENTREGA", cuerpo(t, q.Publicados[0])["carrier"],
		"con envios partidos solo se notifica la transportadora del primero")
}

func TestPublishGuideNotificationRequested_PropagaElErrorDelBroker(t *testing.T) {
	q := &testkit.QueueMock{
		PublishFn: func(ctx context.Context, queueName string, message []byte) error {
			return errors.New("broker caido")
		},
	}

	err := nuevoPublisher(q).PublishGuideNotificationRequested(context.Background(), orden())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error publishing to RabbitMQ")
}

func TestPublicaciones_NingunaEsperaConfirmacionDelBroker(t *testing.T) {
	q := &testkit.QueueMock{}
	p := nuevoPublisher(q)

	require.NoError(t, p.PublishOrderCreated(context.Background(), orden()))
	require.NoError(t, p.PublishConfirmationRequested(context.Background(), orden()))

	assert.Len(t, q.Publicados, 2,
		"ningun metodo declara la cola ni usa publisher confirms: un nil aqui significa 'se entrego al canal', no 'el broker lo persistio'")
	assert.Empty(t, q.Declaradas)
}
