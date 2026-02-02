package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/notification_config/response"
)

// DomainToResponse convierte DTO de dominio a response HTTP
// NUEVA ESTRUCTURA: Usa IDs de tablas normalizadas + datos relacionados
func DomainToResponse(dto dtos.NotificationConfigResponseDTO) response.NotificationConfig {
	return response.NotificationConfig{
		ID:                      dto.ID,
		BusinessID:              dto.BusinessID,
		IntegrationID:           dto.IntegrationID,
		NotificationTypeID:      dto.NotificationTypeID,
		NotificationEventTypeID: dto.NotificationEventTypeID,
		Enabled:                 dto.Enabled,
		Description:             dto.Description,
		OrderStatusIDs:          dto.OrderStatusIDs,
		CreatedAt:               dto.CreatedAt,
		UpdatedAt:               dto.UpdatedAt,

		// Campos adicionales para frontend
		EventType:             dto.EventType,
		Channels:              dto.Channels,
		NotificationTypeName:  dto.NotificationTypeName,
		NotificationEventName: dto.NotificationEventName,
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
