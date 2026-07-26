package mappers

import (
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

// CachedNotificationConfigJSON es la representación JSON del cache de notification_config.
// Mirrors the cache struct from notification_config module.
type CachedNotificationConfigJSON struct {
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

// FromCachedConfig convierte la representación JSON a la entidad de dominio
func FromCachedConfig(cached *CachedNotificationConfigJSON) entities.CachedNotificationConfig {
	return entities.CachedNotificationConfig{
		ID:                      cached.ID,
		BusinessID:              cached.BusinessID,
		IntegrationID:           cached.IntegrationID,
		NotificationTypeID:      cached.NotificationTypeID,
		NotificationEventTypeID: cached.NotificationEventTypeID,
		Enabled:                 cached.Enabled,
		Description:             cached.Description,
		OrderStatusIDs:          cached.OrderStatusIDs,
		EventCode:               cached.EventCode,
		OrderStatusCodes:        cached.OrderStatusCodes,
	}
}
