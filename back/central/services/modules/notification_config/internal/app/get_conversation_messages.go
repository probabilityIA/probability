package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

// GetConversationMessages obtiene los mensajes de una conversación para la vista de chat
func (uc *useCase) GetConversationMessages(ctx context.Context, conversationID string, businessID uint) (*dtos.ConversationDetailResponseDTO, error) {
	conv, messages, err := uc.messageAuditQuerier.GetConversationMessages(ctx, conversationID, businessID)
	if err != nil {
		uc.logger.Error().Err(err).Str("conversation_id", conversationID).Msg("Error getting conversation messages")
		return nil, err
	}

	// Map messages to response DTOs
	msgDTOs := make([]dtos.ConversationMessageResponseDTO, len(messages))
	for i, msg := range messages {
		dto := dtos.ConversationMessageResponseDTO{
			ID:           msg.ID,
			Direction:    msg.Direction,
			MessageID:    msg.MessageID,
			TemplateName: msg.TemplateName,
			Content:      msg.Content,
			Status:       msg.Status,
			CreatedAt:    msg.CreatedAt.Format(time.RFC3339),
		}
		if msg.DeliveredAt != nil {
			formatted := msg.DeliveredAt.Format(time.RFC3339)
			dto.DeliveredAt = &formatted
		}
		if msg.ReadAt != nil {
			formatted := msg.ReadAt.Format(time.RFC3339)
			dto.ReadAt = &formatted
		}
		msgDTOs[i] = dto
	}

	return &dtos.ConversationDetailResponseDTO{
		ConversationID: conv.ID,
		PhoneNumber:    conv.PhoneNumber,
		OrderNumber:    conv.OrderNumber,
		CurrentState:   conv.CurrentState,
		Messages:       msgDTOs,
	}, nil
}
