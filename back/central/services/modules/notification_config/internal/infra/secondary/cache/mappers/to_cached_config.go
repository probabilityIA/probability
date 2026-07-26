package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// CachedNotificationConfig representa una configuración de notificación en Redis
// Usa IDs de tablas normalizadas + campos resueltos para lookup rápido
type CachedNotificationConfig struct {
	ID                      uint     `json:"id"`
	BusinessID              *uint    `json:"business_id,omitempty"`
	IntegrationID           uint     `json:"integration_id"`
	NotificationTypeID      uint     `json:"notification_type_id"`
	NotificationEventTypeID uint     `json:"notification_event_type_id"`
	Enabled                 bool     `json:"enabled"`
	Description             string   `json:"description"`
	OrderStatusIDs          []uint   `json:"order_status_ids"`
	EventCode               string   `json:"event_code,omitempty"`
	OrderStatusCodes        []string `json:"order_status_codes,omitempty"`
}

// ToCachedConfig convierte una entidad de dominio a estructura cacheada
// EventCode se extrae del preload de NotificationEventType si está disponible
// OrderStatusCodes debe ser seteado externamente por el cache manager (requiere query)
func ToCachedConfig(entity *entities.IntegrationNotificationConfig) *CachedNotificationConfig {
	cached := &CachedNotificationConfig{
		ID:                      entity.ID,
		BusinessID:              entity.BusinessID,
		IntegrationID:           entity.IntegrationID,
		NotificationTypeID:      entity.NotificationTypeID,
		NotificationEventTypeID: entity.NotificationEventTypeID,
		Enabled:                 entity.Enabled,
		Description:             entity.Description,
		OrderStatusIDs:          entity.OrderStatusIDs,
	}

	// Extraer EventCode del preload si está disponible
	if entity.NotificationEventType != nil {
		cached.EventCode = entity.NotificationEventType.EventCode
	}

	return cached
}
