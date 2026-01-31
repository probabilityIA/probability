package request

// UpdateNotificationConfig es el DTO de transporte HTTP para actualizar configuración
type UpdateNotificationConfig struct {
	NotificationType *string                 `json:"notification_type,omitempty" binding:"omitempty,oneof=whatsapp email sms"`
	IsActive         *bool                   `json:"is_active,omitempty"`
	Conditions       *NotificationConditions `json:"conditions,omitempty"`
	Config           *NotificationConfig     `json:"config,omitempty"`
	Description      *string                 `json:"description,omitempty"`
	Priority         *int                    `json:"priority,omitempty"`
}
