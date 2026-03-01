package dtos

// MessageAuditFilterDTO define los filtros para consultar logs de auditoría
type MessageAuditFilterDTO struct {
	BusinessID   uint
	Status       *string
	Direction    *string
	TemplateName *string
	DateFrom     *string
	DateTo       *string
	Page         int
	PageSize     int
}

// MessageAuditLogResponseDTO es la respuesta de un registro de auditoría
type MessageAuditLogResponseDTO struct {
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

// MessageAuditStatsResponseDTO contiene estadísticas de mensajes
type MessageAuditStatsResponseDTO struct {
	TotalSent      int64   `json:"total_sent"`
	TotalDelivered int64   `json:"total_delivered"`
	TotalRead      int64   `json:"total_read"`
	TotalFailed    int64   `json:"total_failed"`
	SuccessRate    float64 `json:"success_rate"`
}

// PaginatedMessageAuditResponseDTO es la respuesta paginada de logs de auditoría
type PaginatedMessageAuditResponseDTO struct {
	Data       []MessageAuditLogResponseDTO `json:"data"`
	Total      int64                        `json:"total"`
	Page       int                          `json:"page"`
	PageSize   int                          `json:"page_size"`
	TotalPages int                          `json:"total_pages"`
}
