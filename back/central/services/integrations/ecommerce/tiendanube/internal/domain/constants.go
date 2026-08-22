package domain

const (
	TypeID = uint(17)

	OAuthAuthorizeURLTemplate = "https://www.tiendanube.com/apps/%s/authorize"
	OAuthTokenURL             = "https://www.tiendanube.com/apps/authorize/token"

	AuthMethodOAuth = "oauth"

	ConfigAuthMethod        = "auth_method"
	ConfigScope             = "scope"
	ConfigStatusSyncEnabled = "status_sync_enabled"
)

const (
	FulfillmentUnpacked   = "UNPACKED"
	FulfillmentPacked     = "PACKED"
	FulfillmentDispatched = "DISPATCHED"
	FulfillmentDelivered  = "DELIVERED"
)

const (
	TrackingDispatched     = "dispatched"
	TrackingInTransit      = "in_transit"
	TrackingOutForDelivery = "out_for_delivery"
	TrackingDelivered      = "delivered"
	TrackingAttemptFailed  = "delivery_attempt_failed"
	TrackingReturned       = "returned_to_sender"
)

const (
	EventOrderCreated   = "order/created"
	EventOrderPaid      = "order/paid"
	EventOrderUpdated   = "order/updated"
	EventOrderCancelled = "order/cancelled"
	EventOrderFulfilled = "order/fulfilled"
	EventProductUpdated = "product/updated"
	EventAppUninstalled = "app/uninstalled"
)

var WebhookEvents = []string{
	EventOrderCreated,
	EventOrderPaid,
	EventOrderUpdated,
	EventOrderCancelled,
	EventOrderFulfilled,
	EventProductUpdated,
	EventAppUninstalled,
}
