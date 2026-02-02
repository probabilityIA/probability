package domain

// IWebhookClient define la interfaz para enviar webhooks de WhatsApp
type IWebhookClient interface {
	SendWebhook(payload WebhookPayload) error
}
