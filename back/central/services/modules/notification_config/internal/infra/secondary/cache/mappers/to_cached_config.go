package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// CachedNotificationConfig representa una configuración de notificación en Redis
// Estructura plana para facilitar serialización JSON
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

// ToCachedConfig convierte una entidad de dominio a estructura cacheada
func ToCachedConfig(entity *entities.IntegrationNotificationConfig) *CachedNotificationConfig {
	return &CachedNotificationConfig{
		ID:               entity.ID,
		IntegrationID:    entity.IntegrationID,
		NotificationType: entity.NotificationType,
		IsActive:         entity.IsActive,
		Priority:         entity.Priority,
		Description:      entity.Description,

		// Conditions
		Trigger:             entity.Conditions.Trigger,
		Statuses:            entity.Conditions.Statuses,
		PaymentMethods:      entity.Conditions.PaymentMethods,
		SourceIntegrationID: entity.Conditions.SourceIntegrationID,

		// Config
		TemplateName:  entity.Config.TemplateName,
		RecipientType: entity.Config.RecipientType,
		Language:      entity.Config.Language,
	}
}
