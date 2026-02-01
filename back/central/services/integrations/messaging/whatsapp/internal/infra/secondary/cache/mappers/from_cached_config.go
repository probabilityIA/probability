package mappers

import (
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/secondary/repository"
)

// CachedNotificationConfig representa una configuración de notificación en Redis
// Debe coincidir con la estructura del cache manager de notification_config
type CachedNotificationConfig struct {
	ID               uint     `json:"id"`
	IntegrationID    uint     `json:"integration_id"`
	NotificationType string   `json:"notification_type"`
	IsActive         bool     `json:"is_active"`
	Priority         int      `json:"priority"`
	Description      string   `json:"description"`

	// Conditions (aplanadas)
	Trigger             string   `json:"trigger"`
	Statuses            []string `json:"statuses"`
	PaymentMethods      []uint   `json:"payment_methods"`
	SourceIntegrationID *uint    `json:"source_integration_id"`

	// Config (aplanadas)
	TemplateName  string `json:"template_name"`
	RecipientType string `json:"recipient_type"`
	Language      string `json:"language"`
}

// FromCachedConfig convierte una configuración cacheada a NotificationConfigData
func FromCachedConfig(cached *CachedNotificationConfig) repository.NotificationConfigData {
	return repository.NotificationConfigData{
		ID:                  cached.ID,
		IntegrationID:       cached.IntegrationID,
		NotificationType:    cached.NotificationType,
		IsActive:            cached.IsActive,
		TemplateName:        cached.TemplateName,
		Language:            cached.Language,
		RecipientType:       cached.RecipientType,
		Trigger:             cached.Trigger,
		Statuses:            cached.Statuses,
		PaymentMethods:      cached.PaymentMethods,
		SourceIntegrationID: cached.SourceIntegrationID,
		Priority:            cached.Priority,
		Description:         cached.Description,
	}
}
