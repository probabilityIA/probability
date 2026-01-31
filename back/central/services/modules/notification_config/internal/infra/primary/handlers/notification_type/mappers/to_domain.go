package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_type/request"
)

// CreateRequestToDomain convierte un CreateNotificationType HTTP a entidad de dominio
func CreateRequestToDomain(req *request.CreateNotificationType) *entities.NotificationType {
	return &entities.NotificationType{
		Name:         req.Name,
		Code:         req.Code,
		Description:  req.Description,
		Icon:         req.Icon,
		IsActive:     req.IsActive,
		ConfigSchema: req.ConfigSchema,
	}
}

// UpdateRequestToDomain convierte un UpdateNotificationType HTTP a entidad de dominio
func UpdateRequestToDomain(req *request.UpdateNotificationType, existing *entities.NotificationType) *entities.NotificationType {
	result := *existing // Copiar la entidad existente

	// Actualizar solo los campos que fueron enviados
	if req.Name != nil {
		result.Name = *req.Name
	}
	if req.Description != nil {
		result.Description = *req.Description
	}
	if req.Icon != nil {
		result.Icon = *req.Icon
	}
	if req.IsActive != nil {
		result.IsActive = *req.IsActive
	}
	if req.ConfigSchema != nil {
		result.ConfigSchema = *req.ConfigSchema
	}

	return &result
}
