package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
)

// whatsAppPersister persiste eventos de WhatsApp en DB
type whatsAppPersister struct {
	db     db.IDatabase
	logger log.ILogger
}

// NewWhatsAppPersister crea una nueva instancia del persister de WhatsApp
func NewWhatsAppPersister(database db.IDatabase, logger log.ILogger) ports.IWhatsAppPersister {
	return &whatsAppPersister{
		db:     database,
		logger: logger.WithModule("whatsapp_persister"),
	}
}

// CreateConversation persiste una nueva conversación de WhatsApp
func (r *whatsAppPersister) CreateConversation(ctx context.Context, conv *entities.WhatsAppConversation) error {
	convID, err := uuid.Parse(conv.ID)
	if err != nil {
		convID = uuid.New()
	}

	metadataJSON, _ := json.Marshal(conv.Metadata)

	model := &models.WhatsAppConversation{
		ID:             convID,
		PhoneNumber:    conv.PhoneNumber,
		OrderNumber:    conv.OrderNumber,
		BusinessID:     conv.BusinessID,
		CurrentState:   conv.CurrentState,
		LastMessageID:  conv.LastMessageID,
		LastTemplateID: conv.LastTemplateID,
		Metadata:       datatypes.JSON(metadataJSON),
		CreatedAt:      conv.CreatedAt,
		UpdatedAt:      conv.UpdatedAt,
		ExpiresAt:      conv.ExpiresAt,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("error creando whatsapp_conversation: %w", err)
	}

	return nil
}

// UpdateConversation actualiza una conversación existente de WhatsApp
func (r *whatsAppPersister) UpdateConversation(ctx context.Context, conv *entities.WhatsAppConversation) error {
	convID, err := uuid.Parse(conv.ID)
	if err != nil {
		return fmt.Errorf("conversation ID inválido: %s", conv.ID)
	}

	metadataJSON, _ := json.Marshal(conv.Metadata)

	updates := map[string]interface{}{
		"current_state":   conv.CurrentState,
		"last_message_id": conv.LastMessageID,
		"last_template_id": conv.LastTemplateID,
		"metadata":        datatypes.JSON(metadataJSON),
		"updated_at":      conv.UpdatedAt,
		"expires_at":      conv.ExpiresAt,
	}

	result := r.db.Conn(ctx).
		Model(&models.WhatsAppConversation{}).
		Where("id = ?", convID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("error actualizando whatsapp_conversation: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		// Conversación no existe en DB aún (puede llegar update antes de create)
		// Intentar crear
		r.logger.Warn(ctx).
			Str("conversation_id", conv.ID).
			Msg("Conversación no encontrada para update, intentando create")
		return r.CreateConversation(ctx, conv)
	}

	return nil
}

// ExpireConversation marca una conversación como expirada
func (r *whatsAppPersister) ExpireConversation(ctx context.Context, id string) error {
	convID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("conversation ID inválido: %s", id)
	}

	result := r.db.Conn(ctx).
		Model(&models.WhatsAppConversation{}).
		Where("id = ?", convID).
		Update("expires_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("error expirando whatsapp_conversation: %w", result.Error)
	}

	return nil
}

// CreateMessageLog persiste un nuevo message log de WhatsApp
func (r *whatsAppPersister) CreateMessageLog(ctx context.Context, entry *entities.WhatsAppMessageLogEntry) error {
	msgID := uuid.New()

	convID, err := uuid.Parse(entry.ConversationID)
	if err != nil {
		return fmt.Errorf("conversation_id inválido: %s", entry.ConversationID)
	}

	model := &models.WhatsAppMessageLog{
		ID:             msgID,
		ConversationID: convID,
		Direction:      entry.Direction,
		MessageID:      entry.MessageID,
		TemplateName:   entry.TemplateName,
		Content:        entry.Content,
		Status:         entry.Status,
		DeliveredAt:    entry.DeliveredAt,
		ReadAt:         entry.ReadAt,
		CreatedAt:      entry.CreatedAt,
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("error creando whatsapp_message_log: %w", err)
	}

	return nil
}

// UpdateMessageLogStatus actualiza el estado de un message log por WhatsApp message_id
func (r *whatsAppPersister) UpdateMessageLogStatus(ctx context.Context, messageID, status string, deliveredAt, readAt *string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if deliveredAt != nil {
		if t, err := time.Parse(time.RFC3339, *deliveredAt); err == nil {
			updates["delivered_at"] = t
		}
	}

	if readAt != nil {
		if t, err := time.Parse(time.RFC3339, *readAt); err == nil {
			updates["read_at"] = t
		}
	}

	result := r.db.Conn(ctx).
		Model(&models.WhatsAppMessageLog{}).
		Where("message_id = ?", messageID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("error actualizando whatsapp_message_log status: %w", result.Error)
	}

	return nil
}
