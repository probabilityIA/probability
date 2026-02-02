package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// CachedNotificationConfig representa una configuración de notificación en Redis
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas
type CachedNotificationConfig struct {
	ID                      uint   `json:"id"`
	BusinessID              *uint  `json:"business_id,omitempty"`
	IntegrationID           uint   `json:"integration_id"`
	NotificationTypeID      uint   `json:"notification_type_id"`
	NotificationEventTypeID uint   `json:"notification_event_type_id"`
	Enabled                 bool   `json:"enabled"`
	Description             string `json:"description"`
	OrderStatusIDs          []uint `json:"order_status_ids"`
}

// ToCachedConfig convierte una entidad de dominio a estructura cacheada
// NUEVA ESTRUCTURA: Usa IDs en lugar de strings y nested objects
func ToCachedConfig(entity *entities.IntegrationNotificationConfig) *CachedNotificationConfig {
	return &CachedNotificationConfig{
		ID:                      entity.ID,
		BusinessID:              entity.BusinessID,
		IntegrationID:           entity.IntegrationID,
		NotificationTypeID:      entity.NotificationTypeID,
		NotificationEventTypeID: entity.NotificationEventTypeID,
		Enabled:                 entity.Enabled,
		Description:             entity.Description,
		OrderStatusIDs:          entity.OrderStatusIDs,
	}
}
