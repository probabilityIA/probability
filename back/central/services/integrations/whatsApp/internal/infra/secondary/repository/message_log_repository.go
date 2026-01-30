package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
		"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/infra/secondary/repository/mappers"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type MessageLogRepository struct {
	db  db.IDatabase
	log log.ILogger
}

// Create crea un nuevo log de mensaje en la base de datos
func (r *MessageLogRepository) Create(ctx context.Context, messageLog *entities.MessageLog) error {
	r.log.Info(ctx).
		Str("conversation_id", messageLog.ConversationID).
		Str("message_id", messageLog.MessageID).
		Str("direction", string(messageLog.Direction)).
		Str("template_name", messageLog.TemplateName).
		Msg("[WhatsApp Repository] - creando log de mensaje")

	// Convertir entidad del dominio a modelo GORM
	model, err := mappers.MessageLogDomainToModel(messageLog)
	if err != nil {
		r.log.Error(ctx).Err(err).Msg("[WhatsApp Repository] - error convirtiendo log de mensaje")
		return fmt.Errorf("error al convertir log de mensaje: %w", err)
	}

	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		r.log.Error(ctx).Err(err).
			Str("conversation_id", messageLog.ConversationID).
			Str("message_id", messageLog.MessageID).
			Msg("[WhatsApp Repository] - error creando log de mensaje")
		return fmt.Errorf("error al crear log de mensaje: %w", err)
	}

	// Actualizar la entidad con el ID generado
	messageLog.ID = model.ID.String()

	r.log.Info(ctx).
		Str("id", messageLog.ID).
		Str("message_id", messageLog.MessageID).
		Msg("[WhatsApp Repository] - log de mensaje creado exitosamente")

	return nil
}

// GetByID obtiene un log de mensaje por su ID
func (r *MessageLogRepository) GetByID(ctx context.Context, id string) (*entities.MessageLog, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("ID inválido: %w", err)
	}

	var model models.WhatsAppMessageLog

	if err := r.db.Conn(ctx).
		Where("id = ?", parsedID).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Debug(ctx).Str("id", id).Msg("[WhatsApp Repository] - log de mensaje no encontrado")
			return nil, fmt.Errorf("log de mensaje no encontrado")
		}
		r.log.Error(ctx).Err(err).
			Str("id", id).
			Msg("[WhatsApp Repository] - error obteniendo log de mensaje")
		return nil, fmt.Errorf("error al obtener log de mensaje: %w", err)
	}

	return mappers.MessageLogModelToDomain(&model), nil
}

// GetByMessageID obtiene un log de mensaje por el message_id de WhatsApp
func (r *MessageLogRepository) GetByMessageID(ctx context.Context, messageID string) (*entities.MessageLog, error) {
	var model models.WhatsAppMessageLog

	if err := r.db.Conn(ctx).
		Where("message_id = ?", messageID).
		First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.log.Debug(ctx).
				Str("message_id", messageID).
				Msg("[WhatsApp Repository] - log de mensaje no encontrado por message_id")
			return nil, fmt.Errorf("log de mensaje no encontrado")
		}
		return nil, fmt.Errorf("error al buscar log de mensaje: %w", err)
	}

	return mappers.MessageLogModelToDomain(&model), nil
}

// GetByConversation obtiene todos los logs de mensajes de una conversación
func (r *MessageLogRepository) GetByConversation(ctx context.Context, conversationID string) ([]entities.MessageLog, error) {
	parsedID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, fmt.Errorf("ID de conversación inválido: %w", err)
	}

	var models []models.WhatsAppMessageLog

	if err := r.db.Conn(ctx).
		Where("conversation_id = ?", parsedID).
		Order("created_at ASC").
		Find(&models).Error; err != nil {
		r.log.Error(ctx).Err(err).
			Str("conversation_id", conversationID).
			Msg("[WhatsApp Repository] - error obteniendo logs de conversación")
		return nil, fmt.Errorf("error al obtener logs de conversación: %w", err)
	}

	r.log.Debug(ctx).
		Str("conversation_id", conversationID).
		Int("count", len(models)).
		Msg("[WhatsApp Repository] - logs de conversación obtenidos")

	return mappers.MessageLogsModelToDomain(models), nil
}

// UpdateStatus actualiza el estado y timestamps de un mensaje
func (r *MessageLogRepository) UpdateStatus(ctx context.Context, messageID string, status entities.MessageStatus, timestamps map[string]time.Time) error {
	r.log.Info(ctx).
		Str("message_id", messageID).
		Str("status", string(status)).
		Msg("[WhatsApp Repository] - actualizando estado de mensaje")

	updates := map[string]interface{}{
		"status": string(status),
	}

	// Agregar timestamps si están presentes
	if deliveredAt, ok := timestamps["delivered_at"]; ok {
		updates["delivered_at"] = deliveredAt
	}
	if readAt, ok := timestamps["read_at"]; ok {
		updates["read_at"] = readAt
	}

	if err := r.db.Conn(ctx).
		Model(&models.WhatsAppMessageLog{}).
		Where("message_id = ?", messageID).
		Updates(updates).Error; err != nil {
		r.log.Error(ctx).Err(err).
			Str("message_id", messageID).
			Msg("[WhatsApp Repository] - error actualizando estado de mensaje")
		return fmt.Errorf("error al actualizar estado de mensaje: %w", err)
	}

	return nil
}

// Delete elimina un log de mensaje por su ID
func (r *MessageLogRepository) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("ID inválido: %w", err)
	}

	r.log.Info(ctx).
		Str("id", id).
		Msg("[WhatsApp Repository] - eliminando log de mensaje")

	if err := r.db.Conn(ctx).
		Where("id = ?", parsedID).
		Delete(&models.WhatsAppMessageLog{}).Error; err != nil {
		r.log.Error(ctx).Err(err).
			Str("id", id).
			Msg("[WhatsApp Repository] - error eliminando log de mensaje")
		return fmt.Errorf("error al eliminar log de mensaje: %w", err)
	}

	return nil
}
