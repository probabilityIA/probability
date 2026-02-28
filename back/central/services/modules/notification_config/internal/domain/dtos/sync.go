package dtos

// SyncRuleDTO representa una regla individual dentro de un sync request
type SyncRuleDTO struct {
	ID                      *uint  // nil = crear nueva, valor = actualizar existente
	NotificationTypeID      uint   // Canal de salida (WhatsApp, Email, SMS, SSE)
	NotificationEventTypeID uint   // Tipo de evento (order.created, order.shipped, etc.)
	Enabled                 bool   // Estado de la regla
	Description             string // Descripción opcional
	OrderStatusIDs          []uint // Filtro de estados de orden
}

// SyncNotificationConfigsDTO representa el request de sync batch
type SyncNotificationConfigsDTO struct {
	BusinessID    uint          // Negocio
	IntegrationID uint          // Integración origen
	Rules         []SyncRuleDTO // Reglas a sincronizar
}

// SyncNotificationConfigsResponseDTO representa la respuesta del sync
type SyncNotificationConfigsResponseDTO struct {
	Created int                           `json:"created"`
	Updated int                           `json:"updated"`
	Deleted int                           `json:"deleted"`
	Configs []NotificationConfigResponseDTO `json:"configs"`
}
