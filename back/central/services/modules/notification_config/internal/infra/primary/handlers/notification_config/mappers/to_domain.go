package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/request"
)

// CreateRequestToDomain convierte request HTTP a DTO de dominio (NUEVA ESTRUCTURA)
func CreateRequestToDomain(req *request.CreateNotificationConfig) dtos.CreateNotificationConfigDTO {
	return dtos.CreateNotificationConfigDTO{
		BusinessID:              req.BusinessID,
		IntegrationID:           req.IntegrationID,
		NotificationTypeID:      req.NotificationTypeID,
		NotificationEventTypeID: req.NotificationEventTypeID,
		Enabled:                 req.Enabled,
		Description:             req.Description,
		OrderStatusIDs:          req.OrderStatusIDs,
	}
}

// UpdateRequestToDomain convierte request de actualización a DTO de dominio
func UpdateRequestToDomain(req *request.UpdateNotificationConfig) dtos.UpdateNotificationConfigDTO {
	dto := dtos.UpdateNotificationConfigDTO{
		NotificationType: req.NotificationType,
		IsActive:         req.IsActive,
		Description:      req.Description,
		Priority:         req.Priority,
	}

	if req.Conditions != nil {
		dto.Conditions = &entities.NotificationConditions{
			Trigger:        req.Conditions.Trigger,
			Statuses:       req.Conditions.Statuses,
			PaymentMethods: req.Conditions.PaymentMethods,
		}
	}

	if req.Config != nil {
		dto.Config = &entities.NotificationConfig{
			TemplateName:  req.Config.TemplateName,
			RecipientType: req.Config.RecipientType,
			Language:      req.Config.Language,
		}
	}

	return dto
}

// FilterRequestToDomain convierte query params a DTO de dominio
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas
func FilterRequestToDomain(req *request.FilterNotificationConfig) dtos.FilterNotificationConfigDTO {
	return dtos.FilterNotificationConfigDTO{
		IntegrationID:           req.IntegrationID,
		NotificationTypeID:      req.NotificationTypeID,
		NotificationEventTypeID: req.NotificationEventTypeID,
		Enabled:                 req.Enabled,
	}
}
