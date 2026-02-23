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
)
