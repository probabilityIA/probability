package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

// ListConversations obtiene conversaciones agrupadas con resumen para la vista de lista
func (uc *useCase) ListConversations(ctx context.Context, filter dtos.ConversationListFilterDTO) (*dtos.PaginatedConversationListResponseDTO, error) {
	// Validar y aplicar defaults de paginación
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	conversations, total, err := uc.messageAuditQuerier.ListConversations(ctx, filter)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error listing conversations")
		return nil, err
	}

	// Map to response DTOs
	data := make([]dtos.ConversationSummaryResponseDTO, len(conversations))
	for i, conv := range conversations {
		data[i] = dtos.ConversationSummaryResponseDTO{
			ID:                   conv.ID,
			PhoneNumber:          conv.PhoneNumber,
			OrderNumber:          conv.OrderNumber,
			CurrentState:         conv.CurrentState,
			MessageCount:         conv.MessageCount,
			LastMessageContent:   conv.LastMessageContent,
			LastMessageDirection: conv.LastMessageDirection,
			LastMessageStatus:    conv.LastMessageStatus,
			LastActivity:         conv.LastActivity.Format(time.RFC3339),
			CreatedAt:            conv.CreatedAt.Format(time.RFC3339),
		}
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize != 0 {
		totalPages++
	}

	return &dtos.PaginatedConversationListResponseDTO{
		Data:       data,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}
