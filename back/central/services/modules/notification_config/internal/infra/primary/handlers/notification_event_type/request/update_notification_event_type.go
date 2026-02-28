package request

// UpdateNotificationEventType representa la petición HTTP para actualizar un tipo de evento de notificación
type UpdateNotificationEventType struct {
	EventName             *string                 `json:"event_name"`
	Description           *string                 `json:"description"`
	TemplateConfig        *map[string]interface{} `json:"template_config"`
	IsActive              *bool                   `json:"is_active"`
	AllowedOrderStatusIDs *[]uint                 `json:"allowed_order_status_ids"` // nil = no cambiar, [] = limpiar, [ids] = reemplazar
}
