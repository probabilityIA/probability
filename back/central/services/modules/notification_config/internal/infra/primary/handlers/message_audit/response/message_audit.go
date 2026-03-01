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
