package dtos

// NotificationConfigResponseDTO representa la respuesta de una configuración
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas + datos relacionados
type NotificationConfigResponseDTO struct {
	ID                      uint     `json:"id"`
	BusinessID              *uint    `json:"business_id,omitempty"`
	IntegrationID           uint     `json:"integration_id"`
	NotificationTypeID      uint     `json:"notification_type_id"`
	NotificationEventTypeID uint     `json:"notification_event_type_id"`
	Enabled                 bool     `json:"enabled"`
	Description             string   `json:"description"`
	OrderStatusIDs          []uint   `json:"order_status_ids"`
	CreatedAt               string   `json:"created_at"`
	UpdatedAt               string   `json:"updated_at"`

	// CAMPOS ADICIONALES: Datos relacionados para el frontend
	EventType              *string  `json:"event_type,omitempty"`              // Ej: "order.created"
	Channels               []string `json:"channels,omitempty"`                // Ej: ["sse", "whatsapp"]
	NotificationTypeName   *string  `json:"notification_type_name,omitempty"`  // Ej: "WhatsApp"
	NotificationEventName  *string  `json:"notification_event_name,omitempty"` // Ej: "Confirmación de Pedido"
}
