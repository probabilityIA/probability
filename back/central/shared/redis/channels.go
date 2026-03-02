package redis

// Canales Redis Pub/Sub del proyecto Probability.
//
// Usar siempre estas constantes — nunca strings literales —
// para garantizar que publishers y subscribers usen exactamente el mismo nombre.
const (
	// ChannelOrdersEvents publica cambios de estado en órdenes internas del sistema.
	// Publisher : modules/orders
	// Consumers : modules/events (SSE), modules/orders (score), integrations/messaging/whatsapp
	ChannelOrdersEvents = "probability:orders:state:events"

	// ChannelInvoicingEvents publica resultados de facturación electrónica (creada, fallida, cancelada).
	// Publisher : services/invoicing (factus, siigo, softpymes)
	// Consumers : modules/events (SSE)
	ChannelInvoicingEvents = "probability:invoicing:state:events"

	// ChannelIntegrationsSyncOrders publica resultados de sincronización de órdenes
	// desde plataformas externas (Shopify, WooCommerce, etc.).
	// Publisher : services/integrations/events
	// Consumers : modules/events (SSE)
	ChannelIntegrationsSyncOrders = "probability:integrations:orders:sync:events"

	// DEPRECATED: Shipments ahora publica al dispatcher central via RabbitMQ (ExchangeEvents).
	// Mantener temporalmente por backward compatibility.
	// Publisher : (ninguno — reemplazado por queue/sse_publisher.go)
	// Consumers : (ninguno — frontend ya usa SSE del events module)
	ChannelShipmentsEvents = "probability:shipments:state:events"

	// ChannelPayEvents publica resultados de transacciones de pago (completada, fallida, cancelada).
	// Publisher : modules/pay (response consumer)
	// Consumers : modules/events (SSE)
	ChannelPayEvents = "probability:pay:state:events"

	// ChannelInventoryEvents publica resultados de operaciones de inventario
	// (reservas, confirmaciones, liberaciones, devoluciones por órdenes).
	// Publisher : modules/inventory (order consumer)
	// Consumers : modules/events (SSE)
	ChannelInventoryEvents = "probability:inventory:state:events"
)
