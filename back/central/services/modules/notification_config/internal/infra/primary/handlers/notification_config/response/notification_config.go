package response

// NotificationConfig es el DTO de respuesta HTTP con tags JSON
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas + datos relacionados
type NotificationConfig struct {
	ID                      uint   `json:"id"`
	BusinessID              *uint  `json:"business_id,omitempty"`
	IntegrationID           uint   `json:"integration_id"`
	NotificationTypeID      uint   `json:"notification_type_id"`
	NotificationEventTypeID uint   `json:"notification_event_type_id"`
	Enabled                 bool   `json:"enabled"`
	Description             string `json:"description"`
	OrderStatusIDs          []uint `json:"order_status_ids"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`

	// Campos adicionales para el frontend (incluyen datos de relaciones)
	EventType             *string  `json:"event_type,omitempty"`              // Ej: "order.created" (desde NotificationEventType.EventCode)
	Channels              []string `json:"channels,omitempty"`                // Ej: ["sse", "whatsapp"] (campo deprecated temporal)
	NotificationTypeName  *string  `json:"notification_type_name,omitempty"`  // Ej: "WhatsApp" (desde NotificationType.Name)
	NotificationEventName *string  `json:"notification_event_name,omitempty"` // Ej: "Confirmación de Pedido" (desde NotificationEventType.EventName)
}
