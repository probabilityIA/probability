package rabbitmq

// ═══════════════════════════════════════════════════════════════════
// CONSTANTES CENTRALIZADAS DE COLAS Y EXCHANGES RABBITMQ
//
// Todas las colas y exchanges del sistema están definidos aquí.
// Los módulos importan estas constantes en lugar de definir string
// literals o constantes locales.
//
// Convención de nombres:
//   Exchange*  → exchanges
//   Queue*     → queues
//
// Documentación por constante: Publisher → Consumer
// ═══════════════════════════════════════════════════════════════════

// ─── Exchanges ───────────────────────────────────────────────────

const (
	// ExchangeEvents es el exchange topic del sistema unificado de eventos.
	// Cualquier módulo publica aquí usando PublishEvent().
	// Publisher: cualquier módulo | Consumer: events.unified (via binding #)
	ExchangeEvents = "events.exchange"

	// ExchangeOrderEvents es el exchange fanout de eventos de órdenes.
	// Distribuye eventos a invoicing, whatsapp, score e inventory.
	// Publisher: modules/orders | Consumer: (fanout a 4 colas)
	ExchangeOrderEvents = "orders.events"

	// ExchangeInventory es el exchange topic de inventario.
	// Routing key: sync.{integration_id}
	// Publisher: modules/inventory | Consumer: ecommerce integrations
	ExchangeInventory = "probability.inventory"
)

// ─── Queues: Events ──────────────────────────────────────────────

const (
	// QueueEventsUnified es la cola del sistema unificado de eventos.
	// Bindeada a ExchangeEvents con routing key "#" (recibe todo).
	// Publisher: (via exchange) | Consumer: events/rabbitmq_consumer
	QueueEventsUnified = "events.unified"
)

// ─── Queues: Orders ──────────────────────────────────────────────

const (
	// QueueOrdersCanonical recibe órdenes en formato canónico desde ecommerce.
	// Publisher: ecommerce integrations (shopify, meli, vtex, etc.)
	// Consumer: modules/orders/consumer
	QueueOrdersCanonical = "probability.orders.canonical"

	// QueueOrdersToInvoicing recibe eventos de órdenes para facturación automática.
	// Bindeada a ExchangeOrderEvents (fanout).
	// Publisher: (via exchange) | Consumer: modules/invoicing/order_consumer
	QueueOrdersToInvoicing = "orders.events.invoicing"

	// QueueOrdersToWhatsApp recibe eventos de órdenes para notificaciones WhatsApp.
	// Bindeada a ExchangeOrderEvents (fanout).
	// Publisher: (via exchange) | Consumer: events module
	QueueOrdersToWhatsApp = "orders.events.whatsapp"

	// QueueOrdersToScore recibe eventos de órdenes para cálculo de score.
	// Bindeada a ExchangeOrderEvents (fanout).
	// Publisher: (via exchange) | Consumer: modules/orders/score
	QueueOrdersToScore = "orders.events.score"

	// QueueOrdersToInventory recibe eventos de órdenes para gestión de inventario.
	// Bindeada a ExchangeOrderEvents (fanout).
	// Publisher: (via exchange) | Consumer: modules/inventory/order_consumer
	QueueOrdersToInventory = "orders.events.inventory"

	// QueueOrdersConfirmationRequested solicitudes de confirmación vía WhatsApp.
	// Publisher: modules/orders, events/channel_publisher
	// Consumer: integrations/messaging/whatsapp/consumerorder
	QueueOrdersConfirmationRequested = "orders.confirmation.requested"

	// QueueWhatsAppOrderConfirmed confirmaciones de orden desde WhatsApp.
	// Publisher: integrations/messaging/whatsapp/webhook_publisher
	// Consumer: modules/orders/whatsapp_consumer
	QueueWhatsAppOrderConfirmed = "orders.whatsapp.confirmed"

	// QueueWhatsAppOrderCancelled cancelaciones de orden desde WhatsApp.
	// Publisher: integrations/messaging/whatsapp/webhook_publisher
	// Consumer: modules/orders/whatsapp_consumer
	QueueWhatsAppOrderCancelled = "orders.whatsapp.cancelled"

	// QueueWhatsAppOrderNovelty novedades de orden desde WhatsApp.
	// Publisher: integrations/messaging/whatsapp/webhook_publisher
	// Consumer: modules/orders/whatsapp_consumer
	QueueWhatsAppOrderNovelty = "orders.whatsapp.novelty"
)

// ─── Routing Keys: Orders (fanout exchange, informational only) ──

const (
	// RoutingKeyOrderCreated routing key para eventos de orden creada.
	RoutingKeyOrderCreated = "orders.events.created"

	// RoutingKeyOrderUpdated routing key para eventos de orden actualizada.
	RoutingKeyOrderUpdated = "orders.events.updated"

	// RoutingKeyOrderCancelled routing key para eventos de orden cancelada.
	RoutingKeyOrderCancelled = "orders.events.cancelled"

	// RoutingKeyOrderStatusChanged routing key para eventos de cambio de estado.
	RoutingKeyOrderStatusChanged = "orders.events.status_changed"

	// RoutingKeyOrderGeneric routing key genérico para eventos de orden.
	RoutingKeyOrderGeneric = "orders.events.generic"
)

// ─── Queues: Invoicing ───────────────────────────────────────────

const (
	// QueueInvoicingRequests cola unificada de solicitudes de facturación.
	// El router decide a qué proveedor enrutar.
	// Publisher: modules/invoicing | Consumer: integrations/invoicing/router
	QueueInvoicingRequests = "invoicing.requests"

	// QueueInvoicingResponses respuestas de todos los proveedores de facturación.
	// Publisher: integrations/invoicing/* (cada proveedor)
	// Consumer: modules/invoicing/response_consumer
	QueueInvoicingResponses = "invoicing.responses"

	// QueueInvoicingEvents eventos de facturación (creada, cancelada, fallida).
	// Publisher: modules/invoicing/event_publisher | Consumer: (event subscribers)
	QueueInvoicingEvents = "invoicing.events"

	// QueueInvoicingBulkCreate trabajos de facturación masiva.
	// Publisher: modules/invoicing/handlers | Consumer: modules/invoicing/bulk_consumer
	QueueInvoicingBulkCreate = "invoicing.bulk.create"

	// QueueInvoicingSoftpymesRequests solicitudes para proveedor Softpymes.
	// Publisher: integrations/invoicing/router | Consumer: integrations/invoicing/softpymes
	QueueInvoicingSoftpymesRequests = "invoicing.softpymes.requests"

	// QueueInvoicingFactusRequests solicitudes para proveedor Factus.
	// Publisher: integrations/invoicing/router | Consumer: integrations/invoicing/factus
	QueueInvoicingFactusRequests = "invoicing.factus.requests"

	// QueueInvoicingSiigoRequests solicitudes para proveedor Siigo.
	// Publisher: integrations/invoicing/router | Consumer: integrations/invoicing/siigo
	QueueInvoicingSiigoRequests = "invoicing.siigo.requests"

	// QueueInvoicingAlegraRequests solicitudes para proveedor Alegra.
	// Publisher: integrations/invoicing/router | Consumer: integrations/invoicing/alegra
	QueueInvoicingAlegraRequests = "invoicing.alegra.requests"

	// QueueInvoicingWorldOfficeRequests solicitudes para proveedor World Office.
	// Publisher: integrations/invoicing/router | Consumer: integrations/invoicing/world_office
	QueueInvoicingWorldOfficeRequests = "invoicing.world_office.requests"

	// QueueInvoicingHelisaRequests solicitudes para proveedor Helisa.
	// Publisher: integrations/invoicing/router | Consumer: integrations/invoicing/helisa
	QueueInvoicingHelisaRequests = "invoicing.helisa.requests"
)

// ─── Queues: Pay ─────────────────────────────────────────────────

const (
	// QueuePayRequests cola unificada de solicitudes de pago.
	// El router decide a qué gateway enrutar.
	// Publisher: modules/pay | Consumer: integrations/pay/router
	QueuePayRequests = "pay.requests"

	// QueuePayResponses respuestas de todos los gateways de pago.
	// Publisher: integrations/pay/* (cada gateway)
	// Consumer: modules/pay/response_consumer
	QueuePayResponses = "pay.responses"

	// QueuePayNequiRequests solicitudes para gateway Nequi.
	QueuePayNequiRequests = "pay.nequi.requests"

	// QueuePayBoldRequests solicitudes para gateway Bold.
	QueuePayBoldRequests = "pay.bold.requests"

	// QueuePayWompiRequests solicitudes para gateway Wompi.
	QueuePayWompiRequests = "pay.wompi.requests"

	// QueuePayStripeRequests solicitudes para gateway Stripe.
	QueuePayStripeRequests = "pay.stripe.requests"

	// QueuePayPayURequests solicitudes para gateway PayU.
	QueuePayPayURequests = "pay.payu.requests"

	// QueuePayEPaycoRequests solicitudes para gateway EPayco.
	QueuePayEPaycoRequests = "pay.epayco.requests"

	// QueuePayMeliPagoRequests solicitudes para gateway MeliPago.
	QueuePayMeliPagoRequests = "pay.melipago.requests"
)

// ─── Queues: Transport ───────────────────────────────────────────

const (
	// QueueTransportRequests cola unificada de solicitudes de transporte.
	// El router decide a qué carrier enrutar.
	// Publisher: modules/shipments | Consumer: integrations/transport/router
	QueueTransportRequests = "transport.requests"

	// QueueTransportResponses respuestas de todos los carriers de transporte.
	// Publisher: integrations/transport/* (cada carrier)
	// Consumer: modules/shipments/response_consumer
	QueueTransportResponses = "transport.responses"

	// QueueTransportEnvioclickRequests solicitudes para carrier EnvioClick.
	QueueTransportEnvioclickRequests = "transport.envioclick.requests"

	// QueueTransportEnviameRequests solicitudes para carrier Enviame.
	QueueTransportEnviameRequests = "transport.enviame.requests"

	// QueueTransportTuRequests solicitudes para carrier Tu.
	QueueTransportTuRequests = "transport.tu.requests"

	// QueueTransportMiPaqueteRequests solicitudes para carrier MiPaquete.
	QueueTransportMiPaqueteRequests = "transport.mipaquete.requests"
)

// ─── Queues: Monitoring ──────────────────────────────────────────

const (
	// QueueMonitoringAlerts alertas del sistema de monitoreo.
	// Publisher: modules/monitoring | Consumer: integrations/messaging/whatsapp/consumeralert
	QueueMonitoringAlerts = "monitoring.alerts"
)

// ─── Queues: WhatsApp ────────────────────────────────────────────

const (
	// QueueWhatsAppCustomerHandoff solicitudes de atención humana desde WhatsApp.
	// Publisher: integrations/messaging/whatsapp/webhook_publisher
	// Consumer: (futuro customer service module)
	QueueWhatsAppCustomerHandoff = "customer.whatsapp.handoff"
)
