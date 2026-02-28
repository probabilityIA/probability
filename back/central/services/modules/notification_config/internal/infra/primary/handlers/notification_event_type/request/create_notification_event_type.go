package request

// CreateNotificationEventType representa la petición HTTP para crear un tipo de evento de notificación
type CreateNotificationEventType struct {
	NotificationTypeID    uint                   `json:"notification_type_id" binding:"required"`
	EventCode             string                 `json:"event_code" binding:"required"`
	EventName             string                 `json:"event_name" binding:"required"`
	Description           string                 `json:"description"`
	TemplateConfig        map[string]interface{} `json:"template_config"`
	IsActive              bool                   `json:"is_active"`
	AllowedOrderStatusIDs []uint                 `json:"allowed_order_status_ids"` // Estados de orden permitidos (vacío = todos)
}
