package mappers

import (
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
)

// CachedNotificationConfig mirrors the cache struct from notification_config module
// Used for deserializing JSON from Redis secondary keys
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

// FromCachedConfig convierte una configuración cacheada a NotificationConfigData de dominio
func FromCachedConfig(cached *CachedNotificationConfig) dtos.NotificationConfigData {
	return dtos.NotificationConfigData{
		ID:            cached.ID,
		IntegrationID: cached.IntegrationID,
		IsActive:      cached.Enabled,
		Trigger:       cached.EventCode,
		Statuses:      cached.OrderStatusCodes,
		Description:   cached.Description,
	}
}
