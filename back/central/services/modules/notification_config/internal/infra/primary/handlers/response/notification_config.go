package response

import "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/request"

// NotificationConfig es el DTO de respuesta HTTP con tags JSON
type NotificationConfig struct {
	ID               uint                           `json:"id"`
	IntegrationID    uint                           `json:"integration_id"`
	NotificationType string                         `json:"notification_type"`
	IsActive         bool                           `json:"is_active"`
	Conditions       request.NotificationConditions `json:"conditions"`
	Config           request.NotificationConfig     `json:"config"`
	Description      string                         `json:"description"`
	Priority         int                            `json:"priority"`
	CreatedAt        string                         `json:"created_at"`
	UpdatedAt        string                         `json:"updated_at"`
}
