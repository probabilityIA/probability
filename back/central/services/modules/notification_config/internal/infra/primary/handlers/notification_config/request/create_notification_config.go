package request

// CreateNotificationConfig es el DTO de transporte HTTP para crear configuración
type CreateNotificationConfig struct {
	IntegrationID    uint                      `json:"integration_id" binding:"required"`
	NotificationType string                    `json:"notification_type" binding:"required,oneof=whatsapp email sms"`
	IsActive         bool                      `json:"is_active"`
	Conditions       NotificationConditions    `json:"conditions" binding:"required"`
	Config           NotificationConfig        `json:"config" binding:"required"`
	Description      string                    `json:"description"`
	Priority         int                       `json:"priority"`
}

// NotificationConditions representa las condiciones en formato HTTP
type NotificationConditions struct {
	Trigger        string   `json:"trigger" binding:"required"`
	Statuses       []string `json:"statuses"`
	PaymentMethods []uint   `json:"payment_methods"`
}

// NotificationConfig representa la configuración en formato HTTP
type NotificationConfig struct {
	TemplateName  string `json:"template_name" binding:"required"`
	RecipientType string `json:"recipient_type" binding:"required"`
	Language      string `json:"language" binding:"required"`
}
