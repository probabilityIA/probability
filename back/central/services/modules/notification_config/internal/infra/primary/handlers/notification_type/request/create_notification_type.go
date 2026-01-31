package request

// CreateNotificationType representa la petición HTTP para crear un tipo de notificación
type CreateNotificationType struct {
	Name         string                 `json:"name" binding:"required"`
	Code         string                 `json:"code" binding:"required"`
	Description  string                 `json:"description"`
	Icon         string                 `json:"icon"`
	IsActive     bool                   `json:"is_active"`
	ConfigSchema map[string]interface{} `json:"config_schema"`
}
