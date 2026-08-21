package domain

const (
	TypeID = uint(17)

	OAuthAuthorizeURLTemplate = "https://www.tiendanube.com/apps/%s/authorize"
	OAuthTokenURL             = "https://www.tiendanube.com/apps/authorize/token"

	AuthMethodOAuth = "oauth"

	ConfigAuthMethod = "auth_method"
	ConfigScope      = "scope"
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
