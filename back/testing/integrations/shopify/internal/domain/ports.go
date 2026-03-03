package domain

import sharedtypes "github.com/secamc93/probability/back/testing/shared/types"

// IWebhookClient define la interfaz para enviar webhooks
type IWebhookClient interface {
	SendWebhook(topic string, shopDomain string, payload interface{}) error
	BuildWebhook(topic string, shopDomain string, payload interface{}, baseURL string) (*sharedtypes.WebhookPayload, error)
}
