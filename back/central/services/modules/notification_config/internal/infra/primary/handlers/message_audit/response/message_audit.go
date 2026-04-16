package response

// MessageAuditLog es el DTO de respuesta HTTP para un log de auditoría
type MessageAuditLog struct {
	ID             string  `json:"id"`
	ConversationID string  `json:"conversation_id"`
	MessageID      string  `json:"message_id"`
	Direction      string  `json:"direction"`
	TemplateName   string  `json:"template_name"`
	Content        string  `json:"content"`
	Status         string  `json:"status"`
	DeliveredAt    *string `json:"delivered_at,omitempty"`
	ReadAt         *string `json:"read_at,omitempty"`
	CreatedAt      string  `json:"created_at"`
	PhoneNumber    string  `json:"phone_number"`
	OrderNumber    string  `json:"order_number"`
	BusinessID     uint    `json:"business_id"`
}

// MessageAuditStats es el DTO de respuesta HTTP para estadísticas
type MessageAuditStats struct {
	TotalSent      int64   `json:"total_sent"`
	TotalDelivered int64   `json:"total_delivered"`
	TotalRead      int64   `json:"total_read"`
	TotalFailed    int64   `json:"total_failed"`
	SuccessRate    float64 `json:"success_rate"`
}

// PaginatedMessageAuditResponse es la respuesta paginada de logs
type PaginatedMessageAuditResponse struct {
	Data       []MessageAuditLog `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

// --- Conversation View ---

// ConversationSummary es el resumen de una conversación para la vista de lista
type ConversationSummary struct {
	ID                   string `json:"id"`
	PhoneNumber          string `json:"phone_number"`
	OrderNumber          string `json:"order_number"`
	CurrentState         string `json:"current_state"`
	MessageCount         int    `json:"message_count"`
	LastMessageContent   string `json:"last_message_content"`
	LastMessageDirection string `json:"last_message_direction"`
	LastMessageStatus    string `json:"last_message_status"`
	LastActivity         string `json:"last_activity"`
	CreatedAt            string `json:"created_at"`
}

// PaginatedConversationListResponse es la respuesta paginada de conversaciones
type PaginatedConversationListResponse struct {
	Data       []ConversationSummary `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// ConversationMessage es un mensaje dentro de una conversación para la vista de chat
type ConversationMessage struct {
	ID           string  `json:"id"`
	Direction    string  `json:"direction"`
	MessageID    string  `json:"message_id"`
	TemplateName string  `json:"template_name"`
	Content      string  `json:"content"`
	Status       string  `json:"status"`
	DeliveredAt  *string `json:"delivered_at,omitempty"`
	ReadAt       *string `json:"read_at,omitempty"`
	CreatedAt    string  `json:"created_at"`
}

// ConversationDetailResponse es la respuesta con el detalle de una conversación
type ConversationDetailResponse struct {
	ConversationID string                `json:"conversation_id"`
	PhoneNumber    string                `json:"phone_number"`
	OrderNumber    string                `json:"order_number"`
	CurrentState   string                `json:"current_state"`
	AiPaused       bool                  `json:"ai_paused"`
	Messages       []ConversationMessage `json:"messages"`
}
