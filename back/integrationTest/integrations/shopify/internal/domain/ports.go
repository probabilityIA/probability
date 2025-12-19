package domain

// IWebhookClient define la interfaz para enviar webhooks
type IWebhookClient interface {
	SendWebhook(topic string, shopDomain string, payload interface{}) error
}



