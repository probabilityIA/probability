package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/request"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/response"
)

// DomainToResponse convierte DTO de dominio a response HTTP
func DomainToResponse(dto dtos.NotificationConfigResponseDTO) response.NotificationConfig {
	return response.NotificationConfig{
		ID:               dto.ID,
		IntegrationID:    dto.IntegrationID,
		NotificationType: dto.NotificationType,
		IsActive:         dto.IsActive,
		Conditions: request.NotificationConditions{
			Trigger:        dto.Conditions.Trigger,
			Statuses:       dto.Conditions.Statuses,
			PaymentMethods: dto.Conditions.PaymentMethods,
		},
		Config: request.NotificationConfig{
			TemplateName:  dto.Config.TemplateName,
			RecipientType: dto.Config.RecipientType,
			Language:      dto.Config.Language,
		},
		Description: dto.Description,
		Priority:    dto.Priority,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

// DomainListToResponse convierte lista de DTOs de dominio a lista de responses HTTP
func DomainListToResponse(dtos []dtos.NotificationConfigResponseDTO) []response.NotificationConfig {
	responses := make([]response.NotificationConfig, len(dtos))
	for i, dto := range dtos {
		responses[i] = DomainToResponse(dto)
	}
	return responses
}
