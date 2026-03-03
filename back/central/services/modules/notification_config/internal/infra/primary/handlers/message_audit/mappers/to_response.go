package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/response"
)

// DomainToListResponse convierte la respuesta paginada de dominio a response HTTP
func DomainToListResponse(dto *dtos.PaginatedMessageAuditResponseDTO) response.PaginatedMessageAuditResponse {
	logs := make([]response.MessageAuditLog, len(dto.Data))
	for i, d := range dto.Data {
		logs[i] = response.MessageAuditLog{
			ID:             d.ID,
			ConversationID: d.ConversationID,
			MessageID:      d.MessageID,
			Direction:      d.Direction,
			TemplateName:   d.TemplateName,
			Content:        d.Content,
			Status:         d.Status,
			DeliveredAt:    d.DeliveredAt,
			ReadAt:         d.ReadAt,
			CreatedAt:      d.CreatedAt,
			PhoneNumber:    d.PhoneNumber,
			OrderNumber:    d.OrderNumber,
			BusinessID:     d.BusinessID,
		}
	}

	return response.PaginatedMessageAuditResponse{
		Data:       logs,
		Total:      dto.Total,
		Page:       dto.Page,
		PageSize:   dto.PageSize,
		TotalPages: dto.TotalPages,
	}
}

// DomainToStatsResponse convierte estadísticas de dominio a response HTTP
func DomainToStatsResponse(dto *dtos.MessageAuditStatsResponseDTO) response.MessageAuditStats {
	return response.MessageAuditStats{
		TotalSent:      dto.TotalSent,
		TotalDelivered: dto.TotalDelivered,
		TotalRead:      dto.TotalRead,
		TotalFailed:    dto.TotalFailed,
		SuccessRate:    dto.SuccessRate,
	}
}
