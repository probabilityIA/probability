package response

import "time"

// NotificationEventType representa la respuesta HTTP de un tipo de evento de notificación
type NotificationEventType struct {
	ID                 uint                   `json:"id"`
	NotificationTypeID uint                   `json:"notification_type_id"`
	EventCode          string                 `json:"event_code"`
	EventName          string                 `json:"event_name"`
	Description        string                 `json:"description"`
	TemplateConfig     map[string]interface{} `json:"template_config,omitempty"`
	IsActive           bool                   `json:"is_active"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`

	// Relación expandida (opcional)
	NotificationType *NotificationTypeBasic `json:"notification_type,omitempty"`
}

// NotificationTypeBasic representa información básica del tipo de notificación (para evitar anidación profunda)
type NotificationTypeBasic struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
