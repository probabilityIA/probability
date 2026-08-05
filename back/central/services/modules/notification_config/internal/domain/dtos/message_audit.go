package dtos

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

type MessageAuditStatsResponseDTO struct {
	TotalSent      int64   `json:"total_sent"`
	TotalDelivered int64   `json:"total_delivered"`
	TotalRead      int64   `json:"total_read"`
	TotalFailed    int64   `json:"total_failed"`
	SuccessRate    float64 `json:"success_rate"`
}

type PaginatedMessageAuditResponseDTO struct {
	Data       []MessageAuditLogResponseDTO `json:"data"`
	Total      int64                        `json:"total"`
	Page       int                          `json:"page"`
	PageSize   int                          `json:"page_size"`
	TotalPages int                          `json:"total_pages"`
}

type ConversationListFilterDTO struct {
	BusinessID uint
	DateFrom   *string
	DateTo     *string
	State      *string
	Phone      *string
	Page       int
	PageSize   int
}

type ConversationSummaryResponseDTO struct {
	ID                   string `json:"id"`
	PhoneNumber          string `json:"phone_number"`
	OrderNumber          string `json:"order_number"`
	ConversationType     string `json:"conversation_type"`
	CurrentState         string `json:"current_state"`
	MessageCount         int    `json:"message_count"`
	LastMessageContent   string `json:"last_message_content"`
	LastMessageDirection string `json:"last_message_direction"`
	LastMessageStatus    string `json:"last_message_status"`
	LastActivity         string `json:"last_activity"`
	CreatedAt            string `json:"created_at"`
}

type PaginatedConversationListResponseDTO struct {
	Data       []ConversationSummaryResponseDTO `json:"data"`
	Total      int64                            `json:"total"`
	Page       int                              `json:"page"`
	PageSize   int                              `json:"page_size"`
	TotalPages int                              `json:"total_pages"`
}

type ConversationMessageResponseDTO struct {
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

type ConversationDetailResponseDTO struct {
	ConversationID   string                           `json:"conversation_id"`
	PhoneNumber      string                           `json:"phone_number"`
	OrderNumber      string                           `json:"order_number"`
	ConversationType string                           `json:"conversation_type"`
	CurrentState     string                           `json:"current_state"`
	AiPaused         bool                             `json:"ai_paused"`
	Messages         []ConversationMessageResponseDTO `json:"messages"`
}
