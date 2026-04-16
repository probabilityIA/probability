package response

// SendTemplateResponse define la estructura de respuesta del envío de plantilla
type SendTemplateResponse struct {
	MessageID      string `json:"message_id"`
	Status         string `json:"status"`
	ConversationID string `json:"conversation_id,omitempty"`
}
