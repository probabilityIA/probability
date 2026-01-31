package mappers

import (
	"encoding/json"

		"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

// ConversationModelToDomain convierte un modelo GORM a una entidad del dominio
func ConversationModelToDomain(model *models.WhatsAppConversation) *entities.Conversation {
	if model == nil {
		return nil
	}

	// Deserializar metadata JSONB
	metadata := make(map[string]interface{})
	if model.Metadata != nil {
		_ = json.Unmarshal(model.Metadata, &metadata)
	}

	return &entities.Conversation{
		ID:             model.ID.String(),
		PhoneNumber:    model.PhoneNumber,
		OrderNumber:    model.OrderNumber,
		BusinessID:     model.BusinessID,
		CurrentState:   entities.ConversationState(model.CurrentState),
		LastMessageID:  model.LastMessageID,
		LastTemplateID: model.LastTemplateID,
		Metadata:       metadata,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		ExpiresAt:      model.ExpiresAt,
	}
}

// MessageLogModelToDomain convierte un modelo GORM de mensaje a una entidad del dominio
func MessageLogModelToDomain(model *models.WhatsAppMessageLog) *entities.MessageLog {
	if model == nil {
		return nil
	}

	return &entities.MessageLog{
		ID:             model.ID.String(),
		ConversationID: model.ConversationID.String(),
		Direction:      entities.MessageDirection(model.Direction),
		MessageID:      model.MessageID,
		TemplateName:   model.TemplateName,
		Content:        model.Content,
		Status:         entities.MessageStatus(model.Status),
		DeliveredAt:    model.DeliveredAt,
		ReadAt:         model.ReadAt,
		CreatedAt:      model.CreatedAt,
	}
}

// MessageLogsModelToDomain convierte una lista de modelos GORM a entidades del dominio
func MessageLogsModelToDomain(models []models.WhatsAppMessageLog) []entities.MessageLog {
	if len(models) == 0 {
		return []entities.MessageLog{}
	}

	result := make([]entities.MessageLog, len(models))
	for i, model := range models {
		if log := MessageLogModelToDomain(&model); log != nil {
			result[i] = *log
		}
	}

	return result
}
