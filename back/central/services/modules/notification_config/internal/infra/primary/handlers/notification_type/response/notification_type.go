package response

import "time"

// NotificationType representa la respuesta HTTP de un tipo de notificación
type NotificationType struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name"`
	Code         string                 `json:"code"`
	Description  string                 `json:"description"`
	Icon         string                 `json:"icon"`
	IsActive     bool                   `json:"is_active"`
	ConfigSchema map[string]interface{} `json:"config_schema,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
