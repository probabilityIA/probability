package response

type WebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
