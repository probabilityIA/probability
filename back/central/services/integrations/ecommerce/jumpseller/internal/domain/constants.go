package domain

const (
	HeaderHmac      = "Jumpseller-Hmac-Sha256"
	HeaderEvent     = "Jumpseller-Event"
	HeaderStoreCode = "Jumpseller-Store-Code"

	RateLimitHeader = "Jumpseller-PerMinuteRateLimit-Limit"
)

const (
	StatusPendingPayment = "Pending Payment"
	StatusPaid           = "Paid"
	StatusCanceled       = "Canceled"
	StatusAbandoned      = "Abandoned"
)

const (
	ShipmentRequested = "requested"
	ShipmentInTransit = "in_transit"
	ShipmentDelivered = "delivered"
	ShipmentFailed    = "failed"
)

const (
	EventOrderCreated        = "order_created"
	EventOrderPaid           = "order_paid"
	EventOrderPendingPayment = "order_pending_payment"
	EventOrderShipped        = "order_shipped"
	EventOrderCanceled       = "order_canceled"
	EventOrderAbandoned      = "order_abandoned"
	EventOrderUpdated        = "order_updated"
	EventProductStockUpdate  = "product_stock_update"
)

var WebhookOrderEvents = []string{
	EventOrderCreated,
	EventOrderPaid,
	EventOrderCanceled,
	EventOrderShipped,
}
