package rabbitmq

const (
	ExchangeEvents = "events.exchange"

	ExchangeOrderEvents = "orders.events"

	ExchangeInventory = "probability.inventory"
)

const (
	QueueEventsUnified = "events.unified"
)

const (
	QueueOrdersCanonical = "probability.orders.canonical"

	QueueOrdersToInvoicing = "orders.events.invoicing"

	QueueOrdersToScore = "orders.events.score"

	QueueOrdersToInventory = "orders.events.inventory"

	QueueOrdersToEvents = "orders.events.events"

	QueueOrdersToCustomers = "orders.events.customers"

	QueueOrdersConfirmationRequested = "orders.confirmation.requested"

	QueueWhatsAppOrderConfirmed = "orders.whatsapp.confirmed"

	QueueWhatsAppOrderCancelled = "orders.whatsapp.cancelled"

	QueueWhatsAppOrderNovelty = "orders.whatsapp.novelty"

	QueueSyncBatches = "integration.sync.batches"
)

const (
	RoutingKeyOrderCreated = "orders.events.created"

	RoutingKeyOrderUpdated = "orders.events.updated"

	RoutingKeyOrderCancelled = "orders.events.cancelled"

	RoutingKeyOrderStatusChanged = "orders.events.status_changed"

	RoutingKeyOrderGeneric = "orders.events.generic"
)

const (
	QueueInvoicingRequests = "invoicing.requests"

	QueueInvoicingResponses = "invoicing.responses"

	QueueInvoicingEvents = "invoicing.events"

	QueueInvoicingBulkCreate = "invoicing.bulk.create"

	QueueInvoicingSoftpymesRequests = "invoicing.softpymes.requests"

	QueueInvoicingFactusRequests = "invoicing.factus.requests"

	QueueInvoicingSiigoRequests = "invoicing.siigo.requests"

	QueueInvoicingAlegraRequests = "invoicing.alegra.requests"

	QueueInvoicingWorldOfficeRequests = "invoicing.world_office.requests"

	QueueInvoicingHelisaRequests = "invoicing.helisa.requests"
)

const (
	QueuePayRequests = "pay.requests"

	QueuePayResponses = "pay.responses"

	QueuePayNequiRequests = "pay.nequi.requests"

	QueuePayBoldRequests = "pay.bold.requests"

	QueuePayWompiRequests = "pay.wompi.requests"

	QueuePayStripeRequests = "pay.stripe.requests"

	QueuePayPayURequests = "pay.payu.requests"

	QueuePayEPaycoRequests = "pay.epayco.requests"

	QueuePayMeliPagoRequests = "pay.melipago.requests"
)

const (
	QueueTransportRequests = "transport.requests"

	QueueTransportResponses = "transport.responses"

	QueueTransportEnvioclickRequests = "transport.envioclick.requests"

	QueueTransportEnviameRequests = "transport.enviame.requests"

	QueueTransportTuRequests = "transport.tu.requests"

	QueueTransportMiPaqueteRequests = "transport.mipaquete.requests"
)

const (
	QueueMonitoringAlerts = "monitoring.alerts"
)

const (
	QueueShipmentsWhatsAppGuideNotification = "shipments.whatsapp.guide_notification"
)

const (
	QueueWhatsAppCustomerHandoff = "customer.whatsapp.handoff"

	QueueWhatsAppConversationEvents = "whatsapp.conversation.events"

	QueueWhatsAppMessageLogEvents = "whatsapp.messagelog.events"
)

const (
	QueueWhatsAppAIIncoming = "whatsapp.ai.incoming"

	QueueWhatsAppAIResponse = "whatsapp.ai.response"

	QueueAIOrderResult = "ai.order.result"
)

const (
	QueueInventoryBulkLoad = "inventory.bulk_load.requests"

	QueueInventoryOrderFeedback = "inventory.orders.feedback"
)

const (
	QueueMessagingEmailRequests = "messaging.email.requests"
)

const (
	QueueNotificationDeliveryResults = "notification.delivery.results"
)
