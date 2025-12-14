package request

type WebhookHeaders struct {
	Topic      string `header:"X-Shopify-Topic" binding:"required"`
	Hmac       string `header:"X-Shopify-Hmac-Sha256" binding:"required"`
	ShopDomain string `header:"X-Shopify-Shop-Domain" binding:"required"`
	WebhookID  string `header:"X-Shopify-Webhook-Id"`
}

type WebhookQuery struct {
	IntegrationID string `form:"integration_id" binding:"required"`
}
