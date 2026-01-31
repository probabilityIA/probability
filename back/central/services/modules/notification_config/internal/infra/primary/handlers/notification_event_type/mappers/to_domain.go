package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_event_type/request"
)

// CreateRequestToDomain convierte un CreateNotificationEventType HTTP a entidad de dominio
func CreateRequestToDomain(req *request.CreateNotificationEventType) *entities.NotificationEventType {
	return &entities.NotificationEventType{
		NotificationTypeID: req.NotificationTypeID,
		EventCode:          req.EventCode,
		EventName:          req.EventName,
		Description:        req.Description,
		TemplateConfig:     req.TemplateConfig,
		IsActive:           req.IsActive,
	}
}

// UpdateRequestToDomain convierte un UpdateNotificationEventType HTTP a entidad de dominio
func UpdateRequestToDomain(req *request.UpdateNotificationEventType, existing *entities.NotificationEventType) *entities.NotificationEventType {
	result := *existing // Copiar la entidad existente

	// Actualizar solo los campos que fueron enviados
	if req.EventName != nil {
		result.EventName = *req.EventName
	}
	if req.Description != nil {
		result.Description = *req.Description
	}
	if req.TemplateConfig != nil {
		result.TemplateConfig = *req.TemplateConfig
	}
	if req.IsActive != nil {
		result.IsActive = *req.IsActive
	}

	return &result
}
