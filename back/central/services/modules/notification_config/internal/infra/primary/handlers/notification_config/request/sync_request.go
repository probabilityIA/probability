package request

// SyncRule representa una regla individual en el request de sync
type SyncRule struct {
	ID                      *uint  `json:"id"`                                             // nil = crear, valor = actualizar
	NotificationTypeID      uint   `json:"notification_type_id" binding:"required"`       // Canal de salida
	NotificationEventTypeID uint   `json:"notification_event_type_id" binding:"required"` // Tipo de evento
	Enabled                 bool   `json:"enabled"`                                        // Estado
	Description             string `json:"description"`                                    // Descripción opcional
	OrderStatusIDs          []uint `json:"order_status_ids"`                              // Filtro de estados
}

// SyncNotificationConfigs es el DTO de transporte HTTP para sync batch
type SyncNotificationConfigs struct {
	IntegrationID uint       `json:"integration_id" binding:"required"` // Integración origen
	Rules         []SyncRule `json:"rules" binding:"required"`          // Reglas a sincronizar
}
