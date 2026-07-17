package domain

const (
	OAuthAuthorizeURL = "https://accounts.jumpseller.com/oauth/authorize"
	OAuthTokenURL     = "https://accounts.jumpseller.com/oauth/token"

	AuthMethodOAuth = "oauth"

	ConfigAuthMethod     = "auth_method"
	ConfigTokenExpiresAt = "token_expires_at"

	DefaultScopes = "read_store read_orders write_orders read_products write_products read_customers write_hooks"
)

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
	EventOrderPaid           = "order_paid"
	EventOrderPendingPayment = "order_pending_payment"
	EventOrderShipped        = "order_shipped"
	EventOrderCanceled       = "order_canceled"
	EventOrderAbandoned      = "order_abandoned"
	EventOrderUpdated        = "order_updated"
	EventProductStockUpdate  = "product_stock_update"
)

var WebhookOrderEvents = []string{
	EventOrderPaid,
	EventOrderCanceled,
	EventOrderShipped,
}
