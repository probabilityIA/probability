package mappers

import (
	"encoding/json"

	"github.com/google/uuid"
		"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
)

// ConversationDomainToModel convierte una entidad del dominio a un modelo GORM
func ConversationDomainToModel(entity *entities.Conversation) (*models.WhatsAppConversation, error) {
	if entity == nil {
		return nil, nil
	}

	// Parsear UUID
	id, err := uuid.Parse(entity.ID)
	if err != nil {
		// Si no tiene ID aún (creación nueva), generar uno nuevo
		id = uuid.New()
	}

	// Serializar metadata a JSONB
	metadataJSON, err := json.Marshal(entity.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	return &models.WhatsAppConversation{
		ID:             id,
		PhoneNumber:    entity.PhoneNumber,
		OrderNumber:    entity.OrderNumber,
		BusinessID:     entity.BusinessID,
		CurrentState:   string(entity.CurrentState),
		LastMessageID:  entity.LastMessageID,
		LastTemplateID: entity.LastTemplateID,
		Metadata:       datatypes.JSON(metadataJSON),
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		ExpiresAt:      entity.ExpiresAt,
	}, nil
}

// MessageLogDomainToModel convierte una entidad del dominio de mensaje a un modelo GORM
func MessageLogDomainToModel(entity *entities.MessageLog) (*models.WhatsAppMessageLog, error) {
	if entity == nil {
		return nil, nil
	}

	// Parsear UUIDs
	id, err := uuid.Parse(entity.ID)
	if err != nil {
		// Si no tiene ID aún (creación nueva), generar uno nuevo
		id = uuid.New()
	}

	conversationID, err := uuid.Parse(entity.ConversationID)
	if err != nil {
		return nil, err
	}

	return &models.WhatsAppMessageLog{
		ID:             id,
		ConversationID: conversationID,
		Direction:      string(entity.Direction),
		MessageID:      entity.MessageID,
		TemplateName:   entity.TemplateName,
		Content:        entity.Content,
		Status:         string(entity.Status),
		DeliveredAt:    entity.DeliveredAt,
		ReadAt:         entity.ReadAt,
		CreatedAt:      entity.CreatedAt,
	}, nil
}
