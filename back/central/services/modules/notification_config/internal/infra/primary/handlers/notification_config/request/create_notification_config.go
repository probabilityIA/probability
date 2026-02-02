package request

// CreateNotificationConfig es el DTO de transporte HTTP para crear configuración
// NUEVA ESTRUCTURA: Usa IDs de tablas en lugar de strings hardcodeados
type CreateNotificationConfig struct {
	BusinessID               *uint  `json:"business_id"`                                    // Opcional - auto-asignado por middleware
	IntegrationID            uint   `json:"integration_id" binding:"required"`             // Integración origen (Shopify, WhatsApp, etc.)
	NotificationTypeID       uint   `json:"notification_type_id" binding:"required"`       // Canal de salida (WhatsApp, Email, SMS)
	NotificationEventTypeID  uint   `json:"notification_event_type_id" binding:"required"` // Tipo de evento (order_created, order_paid, etc.)
	Enabled                  bool   `json:"enabled"`                                        // Estado de la configuración
	Description              string `json:"description"`                                    // Descripción opcional
	OrderStatusIDs           []uint `json:"order_status_ids"`                              // Filtro opcional de estados
}

// ============================================
// LEGACY: Estructura anterior (DEPRECADA)
// Se mantiene temporalmente para compatibilidad
// ============================================

type CreateNotificationConfigLegacy struct {
	IntegrationID    uint                      `json:"integration_id" binding:"required"`
	NotificationType string                    `json:"notification_type" binding:"required,oneof=whatsapp email sms"`
	IsActive         bool                      `json:"is_active"`
	Conditions       NotificationConditions    `json:"conditions" binding:"required"`
	Config           NotificationConfig        `json:"config" binding:"required"`
	Description      string                    `json:"description"`
	Priority         int                       `json:"priority"`
}

type NotificationConditions struct {
	Trigger        string   `json:"trigger" binding:"required"`
	Statuses       []string `json:"statuses"`
	PaymentMethods []uint   `json:"payment_methods"`
}

type NotificationConfig struct {
	TemplateName  string `json:"template_name" binding:"required"`
	RecipientType string `json:"recipient_type" binding:"required"`
	Language      string `json:"language" binding:"required"`
}
