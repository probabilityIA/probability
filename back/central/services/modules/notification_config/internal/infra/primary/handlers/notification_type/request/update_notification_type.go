package request

// UpdateNotificationType representa la petición HTTP para actualizar un tipo de notificación
type UpdateNotificationType struct {
	Name         *string                 `json:"name"`
	Description  *string                 `json:"description"`
	Icon         *string                 `json:"icon"`
	IsActive     *bool                   `json:"is_active"`
	ConfigSchema *map[string]interface{} `json:"config_schema"`
}
