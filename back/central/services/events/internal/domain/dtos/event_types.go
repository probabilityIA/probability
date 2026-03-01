package dtos

// Categorías de eventos
const (
	CategoryOrder       = "order"
	CategoryInvoice     = "invoice"
	CategoryShipment    = "shipment"
	CategoryIntegration = "integration"
	CategoryInventory   = "inventory"
	CategoryPay         = "pay"
)

// ───────────────────────────────────────────
// ORDER EVENT TYPES
// ───────────────────────────────────────────

const (
	OrderCreated         = "order.created"
	OrderUpdated         = "order.updated"
	OrderStatusChanged   = "order.status_changed"
	OrderCancelled       = "order.cancelled"
	OrderDelivered       = "order.delivered"
	OrderShipped         = "order.shipped"
	OrderPaymentReceived = "order.payment_received"
	OrderRefunded        = "order.refunded"
	OrderFailed          = "order.failed"
	OrderOnHold          = "order.on_hold"
	OrderProcessing      = "order.processing"

	OrderScoreCalculationRequested = "order.score_calculation_requested"
	OrderNotificationSent          = "order.notification_sent"
	OrderNotificationFailed        = "order.notification_failed"
)

// ───────────────────────────────────────────
// INVOICE EVENT TYPES
// ───────────────────────────────────────────

const (
	InvoiceCreated      = "invoice.created"
	InvoiceFailed       = "invoice.failed"
	InvoiceCancelled    = "invoice.cancelled"
	CreditNoteCreated   = "credit_note.created"
	BulkJobProgress     = "bulk_job.progress"
	BulkJobCompleted    = "bulk_job.completed"
	InvoiceCompareReady = "invoice.compare_ready"
)

// ───────────────────────────────────────────
// SHIPMENT EVENT TYPES
// ───────────────────────────────────────────

const (
	ShipmentQuoteReceived   = "shipment.quote_received"
	ShipmentQuoteFailed     = "shipment.quote_failed"
	ShipmentGuideGenerated  = "shipment.guide_generated"
	ShipmentGuideFailed     = "shipment.guide_failed"
	ShipmentTrackingUpdated = "shipment.tracking_updated"
	ShipmentTrackingFailed  = "shipment.tracking_failed"
	ShipmentCancelled       = "shipment.cancelled"
	ShipmentCancelFailed    = "shipment.cancel_failed"
)

// ───────────────────────────────────────────
// INTEGRATION / SYNC EVENT TYPES
// ───────────────────────────────────────────

const (
	IntegrationSyncOrderCreated  = "integration.sync.order.created"
	IntegrationSyncOrderUpdated  = "integration.sync.order.updated"
	IntegrationSyncOrderRejected = "integration.sync.order.rejected"
	IntegrationSyncStarted       = "integration.sync.started"
	IntegrationSyncCompleted     = "integration.sync.completed"
	IntegrationSyncFailed        = "integration.sync.failed"
)

// ───────────────────────────────────────────
// INVENTORY EVENT TYPES
// ───────────────────────────────────────────

const (
	InventorySyncStarted   = "inventory_sync_started"
	InventorySyncProgress  = "inventory_sync_progress"
	InventorySyncCompleted = "inventory_sync_completed"
	InventorySyncFailed    = "inventory_sync_failed"
	ProductSynced          = "product_synced"
	ProductFailed          = "product_failed"
	ConnectionEstablished  = "connection_established"
	BatchStarted           = "batch_started"
	BatchCompleted         = "batch_completed"
)

// NotificationTypeSSE es el ID del tipo de notificación SSE
const NotificationTypeSSE = 1

// NotificationTypeWhatsApp es el ID del tipo de notificación WhatsApp
const NotificationTypeWhatsApp = 2

// NotificationTypeEmail es el ID del tipo de notificación Email
const NotificationTypeEmail = 3
