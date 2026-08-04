package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/response"
)

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

func DomainToStatsResponse(dto *dtos.MessageAuditStatsResponseDTO) response.MessageAuditStats {
	return response.MessageAuditStats{
		TotalSent:      dto.TotalSent,
		TotalDelivered: dto.TotalDelivered,
		TotalRead:      dto.TotalRead,
		TotalFailed:    dto.TotalFailed,
		SuccessRate:    dto.SuccessRate,
	}
}

func DomainToConversationListResponse(dto *dtos.PaginatedConversationListResponseDTO) response.PaginatedConversationListResponse {
	conversations := make([]response.ConversationSummary, len(dto.Data))
	for i, d := range dto.Data {
		conversations[i] = response.ConversationSummary{
			ID:                   d.ID,
			PhoneNumber:          d.PhoneNumber,
			OrderNumber:          d.OrderNumber,
			ConversationType:     d.ConversationType,
			CurrentState:         d.CurrentState,
			MessageCount:         d.MessageCount,
			LastMessageContent:   d.LastMessageContent,
			LastMessageDirection: d.LastMessageDirection,
			LastMessageStatus:    d.LastMessageStatus,
			LastActivity:         d.LastActivity,
			CreatedAt:            d.CreatedAt,
		}
	}

	return response.PaginatedConversationListResponse{
		Data:       conversations,
		Total:      dto.Total,
		Page:       dto.Page,
		PageSize:   dto.PageSize,
		TotalPages: dto.TotalPages,
	}
}

func DomainToConversationDetailResponse(dto *dtos.ConversationDetailResponseDTO) response.ConversationDetailResponse {
	messages := make([]response.ConversationMessage, len(dto.Messages))
	for i, m := range dto.Messages {
		messages[i] = response.ConversationMessage{
			ID:           m.ID,
			Direction:    m.Direction,
			MessageID:    m.MessageID,
			TemplateName: m.TemplateName,
			Content:      m.Content,
			Status:       m.Status,
			DeliveredAt:  m.DeliveredAt,
			ReadAt:       m.ReadAt,
			CreatedAt:    m.CreatedAt,
		}
	}

	return response.ConversationDetailResponse{
		ConversationID:   dto.ConversationID,
		PhoneNumber:      dto.PhoneNumber,
		OrderNumber:      dto.OrderNumber,
		ConversationType: dto.ConversationType,
		CurrentState:     dto.CurrentState,
		AiPaused:         dto.AiPaused,
		Messages:         messages,
	}
}
